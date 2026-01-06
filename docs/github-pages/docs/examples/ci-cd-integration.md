---
sidebar_position: 3
---

# CI/CD Integration

Integrate tf2report into your CI/CD pipelines for automated Terraform plan analysis.

## GitHub Actions

### Basic Integration

```yaml title=".github/workflows/terraform.yml"
name: Terraform Plan

on:
  pull_request:
    branches: [main]

jobs:
  plan:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Setup Terraform
        uses: hashicorp/setup-terraform@v3

      - name: Install tf2report
        run: go install github.com/germainlefebvre4/tf2report/cmd/tf2report@latest

      - name: Terraform Init
        run: terraform init

      - name: Terraform Plan
        run: |
          terraform plan -out=tfplan
          terraform show -json tfplan > terraform.tfplan.json

      - name: Generate Report
        run: tf2report --plan terraform.tfplan.json --format markdown > report.md

      - name: Comment PR
        uses: actions/github-script@v7
        with:
          script: |
            const fs = require('fs');
            const report = fs.readFileSync('report.md', 'utf8');
            github.rest.issues.createComment({
              issue_number: context.issue.number,
              owner: context.repo.owner,
              repo: context.repo.repo,
              body: report
            });
```

### Advanced GitHub Actions with Checks

```yaml title=".github/workflows/terraform-advanced.yml"
name: Terraform Plan Analysis

on:
  pull_request:
    branches: [main]

jobs:
  plan-and-analyze:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: hashicorp/setup-terraform@v3

      - name: Install tf2report
        run: go install github.com/germainlefebvre4/tf2report/cmd/tf2report@latest

      - name: Terraform Init
        run: terraform init

      - name: Terraform Plan
        run: |
          terraform plan -out=tfplan
          terraform show -json tfplan > terraform.tfplan.json

      - name: Generate JSON Report
        run: tf2report --plan terraform.tfplan.json --format json > report.json

      - name: Check for Destructive Changes
        run: |
          DESTRUCTIVE=$(jq '.summary.to_destroy + .summary.to_replace' report.json)
          if [ "$DESTRUCTIVE" -gt 0 ]; then
            echo "::warning::Destructive changes detected: $DESTRUCTIVE resources"
          fi

      - name: Generate Markdown Report
        run: tf2report --plan terraform.tfplan.json > report.md

      - name: Generate Security Review
        run: |
          tf2report --plan terraform.tfplan.json \
            --type aws_iam_role \
            --type aws_iam_policy \
            --type aws_security_group \
            --type aws_kms_key \
            > security-review.md

      - name: Comment PR with Reports
        uses: actions/github-script@v7
        with:
          script: |
            const fs = require('fs');
            const report = fs.readFileSync('report.md', 'utf8');
            const security = fs.readFileSync('security-review.md', 'utf8');

            const comment = `## Terraform Plan Report\n\n${report}\n\n## Security Review\n\n${security}`;

            github.rest.issues.createComment({
              issue_number: context.issue.number,
              owner: context.repo.owner,
              repo: context.repo.repo,
              body: comment
            });

      - name: Upload Reports as Artifacts
        uses: actions/upload-artifact@v4
        with:
          name: terraform-reports
          path: |
            report.md
            report.json
            security-review.md
```

## GitLab CI/CD

### Basic Integration

```yaml title=".gitlab-ci.yml"
stages:
  - plan
  - report

terraform:plan:
  stage: plan
  image: hashicorp/terraform:latest
  script:
    - terraform init
    - terraform plan -out=tfplan
    - terraform show -json tfplan > terraform.tfplan.json
  artifacts:
    paths:
      - terraform.tfplan.json
    expire_in: 1 week

terraform:report:
  stage: report
  image: golang:1.25
  before_script:
    - go install github.com/germainlefebvre4/tf2report/cmd/tf2report@latest
  script:
    - tf2report --plan terraform.tfplan.json > report.md
    - tf2report --plan terraform.tfplan.json --format json > report.json
  artifacts:
    paths:
      - report.md
      - report.json
    expire_in: 1 week
```

### Advanced GitLab CI with MR Comments

```yaml title=".gitlab-ci.yml"
stages:
  - plan
  - analyze
  - comment

variables:
  TF_ROOT: ${CI_PROJECT_DIR}

terraform:plan:
  stage: plan
  image: hashicorp/terraform:latest
  script:
    - cd ${TF_ROOT}
    - terraform init
    - terraform plan -out=tfplan
    - terraform show -json tfplan > terraform.tfplan.json
  artifacts:
    paths:
      - ${TF_ROOT}/terraform.tfplan.json

terraform:analyze:
  stage: analyze
  image: golang:1.25
  before_script:
    - go install github.com/germainlefebvre4/tf2report/cmd/tf2report@latest
  script:
    - cd ${TF_ROOT}
    - tf2report --plan terraform.tfplan.json > report.md
    - tf2report --plan terraform.tfplan.json --format json > report.json

    # Check for destructive changes
    - |
      DESTRUCTIVE=$(jq '.summary.to_destroy + .summary.to_replace' report.json)
      if [ "$DESTRUCTIVE" -gt 0 ]; then
        echo "WARNING: $DESTRUCTIVE destructive changes detected"
        exit 1
      fi
  artifacts:
    paths:
      - ${TF_ROOT}/report.md
      - ${TF_ROOT}/report.json
  allow_failure: true

terraform:comment:
  stage: comment
  image: curlimages/curl:latest
  script:
    - |
      REPORT=$(cat ${TF_ROOT}/report.md)
      curl --request POST --header "PRIVATE-TOKEN: ${CI_JOB_TOKEN}" \
        --data "body=${REPORT}" \
        "${CI_API_V4_URL}/projects/${CI_PROJECT_ID}/merge_requests/${CI_MERGE_REQUEST_IID}/notes"
  only:
    - merge_requests
```

## Jenkins Pipeline

### Declarative Pipeline

```groovy title="Jenkinsfile"
pipeline {
    agent any

    environment {
        GO_VERSION = '1.25'
    }

    stages {
        stage('Setup') {
            steps {
                sh 'go install github.com/germainlefebvre4/tf2report/cmd/tf2report@latest'
            }
        }

        stage('Terraform Plan') {
            steps {
                sh '''
                    terraform init
                    terraform plan -out=tfplan
                    terraform show -json tfplan > terraform.tfplan.json
                '''
            }
        }

        stage('Generate Reports') {
            steps {
                sh '''
                    tf2report --plan terraform.tfplan.json > report.md
                    tf2report --plan terraform.tfplan.json --format json > report.json
                '''
            }
        }

        stage('Check Destructive Changes') {
            steps {
                script {
                    def destructive = sh(
                        script: 'jq ".summary.to_destroy + .summary.to_replace" report.json',
                        returnStdout: true
                    ).trim().toInteger()

                    if (destructive > 0) {
                        echo "WARNING: ${destructive} destructive changes detected"
                        currentBuild.result = 'UNSTABLE'
                    }
                }
            }
        }

        stage('Archive Reports') {
            steps {
                archiveArtifacts artifacts: '*.md,*.json', fingerprint: true
            }
        }
    }
}
```

## Azure DevOps

```yaml title="azure-pipelines.yml"
trigger:
  - main

pool:
  vmImage: 'ubuntu-latest'

steps:
  - task: GoTool@0
    inputs:
      version: '1.25'

  - script: |
      go install github.com/germainlefebvre4/tf2report/cmd/tf2report@latest
    displayName: 'Install tf2report'

  - script: |
      terraform init
      terraform plan -out=tfplan
      terraform show -json tfplan > terraform.tfplan.json
    displayName: 'Terraform Plan'

  - script: |
      $(go env GOPATH)/bin/tf2report --plan terraform.tfplan.json > report.md
      $(go env GOPATH)/bin/tf2report --plan terraform.tfplan.json --format json > report.json
    displayName: 'Generate Reports'

  - task: PublishBuildArtifacts@1
    inputs:
      pathToPublish: '$(Build.SourcesDirectory)'
      artifactName: 'terraform-reports'
```

## CircleCI

```yaml title=".circleci/config.yml"
version: 2.1

orbs:
  terraform: circleci/terraform@3.2

jobs:
  plan-and-report:
    docker:
      - image: cimg/go:1.25
    steps:
      - checkout

      - run:
          name: Install Terraform
          command: |
            wget https://releases.hashicorp.com/terraform/1.10.5/terraform_1.10.5_linux_amd64.zip
            unzip terraform_1.10.5_linux_amd64.zip
            sudo mv terraform /usr/local/bin/

      - run:
          name: Install tf2report
          command: go install github.com/germainlefebvre4/tf2report/cmd/tf2report@latest

      - run:
          name: Terraform Plan
          command: |
            terraform init
            terraform plan -out=tfplan
            terraform show -json tfplan > terraform.tfplan.json

      - run:
          name: Generate Report
          command: |
            ~/go/bin/tf2report --plan terraform.tfplan.json > report.md
            ~/go/bin/tf2report --plan terraform.tfplan.json --format json > report.json

      - store_artifacts:
          path: report.md

      - store_artifacts:
          path: report.json

workflows:
  terraform:
    jobs:
      - plan-and-report
```

## Bitbucket Pipelines

```yaml title="bitbucket-pipelines.yml"
pipelines:
  pull-requests:
    '**':
      - step:
          name: Terraform Plan and Report
          image: golang:1.25
          script:
            - apt-get update && apt-get install -y wget unzip
            - wget https://releases.hashicorp.com/terraform/1.10.5/terraform_1.10.5_linux_amd64.zip
            - unzip terraform_1.10.5_linux_amd64.zip
            - mv terraform /usr/local/bin/
            - go install github.com/germainlefebvre4/tf2report/cmd/tf2report@latest
            - terraform init
            - terraform plan -out=tfplan
            - terraform show -json tfplan > terraform.tfplan.json
            - tf2report --plan terraform.tfplan.json > report.md
            - tf2report --plan terraform.tfplan.json --format json > report.json
          artifacts:
            - report.md
            - report.json
```

## Docker-Based CI/CD

### Dockerfile

```dockerfile title="Dockerfile"
FROM hashicorp/terraform:latest AS terraform

FROM golang:1.25 AS builder
RUN go install github.com/germainlefebvre4/tf2report/cmd/tf2report@latest

FROM alpine:latest
COPY --from=terraform /bin/terraform /usr/local/bin/terraform
COPY --from=builder /go/bin/tf2report /usr/local/bin/tf2report

WORKDIR /workspace

ENTRYPOINT ["/bin/sh", "-c"]
CMD ["terraform plan -out=tfplan && terraform show -json tfplan > terraform.tfplan.json && tf2report --plan terraform.tfplan.json"]
```

### Usage

```bash
docker build -t terraform-reporter .
docker run -v $(pwd):/workspace terraform-reporter
```

## Common CI/CD Patterns

### Pattern: Fail on Destructive Changes

```bash
#!/bin/bash
tf2report --plan terraform.tfplan.json --format json > report.json

DESTRUCTIVE=$(jq '.summary.to_destroy + .summary.to_replace' report.json)

if [ "$DESTRUCTIVE" -gt 0 ]; then
  echo "❌ FAILURE: $DESTRUCTIVE destructive changes detected"
  tf2report --plan terraform.tfplan.json --action delete --action replace
  exit 1
fi
```

### Pattern: Require Approval for Large Changes

```bash
#!/bin/bash
tf2report --plan terraform.tfplan.json --format json > report.json

TOTAL=$(jq '.summary.total_changes' report.json)

if [ "$TOTAL" -gt 10 ]; then
  echo "⚠️  WARNING: $TOTAL changes detected - manual approval required"
  exit 1
fi
```

### Pattern: Multiple Environment Reports

```bash
#!/bin/bash
for ENV in dev staging prod; do
  terraform workspace select $ENV
  terraform plan -out=tfplan
  terraform show -json tfplan > ${ENV}.tfplan.json
  tf2report --plan ${ENV}.tfplan.json > reports/${ENV}-report.md
done
```

## Next Steps

- **[Basic Usage](./basic-usage.md)** - Simple examples
- **[Filtering](./filtering.md)** - Filter reports in CI/CD
