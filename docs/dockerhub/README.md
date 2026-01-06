# Quick reference

* **Maintained by**:<br>
  [Germain LEFEBVRE](https://github.com/germainlefebvre4)

* **Where to get help**:<br>
  [Github Discussions](https://github.com/germainlefebvre4/tf2report/discussions)

# Supported tags and respective Dockerfile links

* [`latest`, `v0`, `v0.0`, `v0.0.0`](https://github.com/germainlefebvre4/tf2report/blob/v0.0.0/Dockerfile)

# Quick reference (cont.)

* **Where to file issues**:<br>
  https://github.com/germainlefebvre4/tf2report/issues⁠

* Supported architectures: ([more info⁠]())<br>
  amd64, arm64

* **Source of this description**:<br>
  [tf2report repo's `docs/dockerhub/` directory](https://github.com/germainlefebvre4/tf2report/tree/main/docs/dockerhub/) ([history](https://github.com/docker-library/docs/commits/master/nginx))

# What is TF2Report?

TF2Report is a CLI tool that analyzes Terraform plan files and generates human-readable resource change reports.

## Features

* **Parse Terraform plan files:** Analyze Terraform plan JSON files to extract resource changes
* **Multiple output formats:** Generate reports in Markdown, plain text, or JSON format
* **Flexible filtering:** Filter by resource type and change action (create, update, delete, replace)
* **Configuration file support:** Use a `tf2report.yaml` file to specify default settings for your projects
* **CI/CD integration:** Lightweight CLI tool that can be easily integrated into your automation workflows
* **Detailed summaries:** Get clear overviews of infrastructure changes before applying them

# How to use this image

TF2Report helper.

```raw
tf2report analyzes Terraform plan files and generates resource change reports.

It parses Terraform plan JSON files, extracts resource changes, and generates
summaries in multiple output formats including Markdown, plain text, and JSON.

Usage:
  tf2report [flags]

Flags:
  -a, --action strings   filter by action (create, update, delete, replace)
      --config string    config file path (default: ./tf2report.yaml)
  -f, --format string    output format (markdown, text, json)
  -h, --help             help for tf2report
  -p, --plan string      path to Terraform plan JSON file
  -t, --type strings     filter by resource type (can be specified multiple times)
  -v, --verbose          verbose output
      --version          version for tf2report
```

# License

View [license information⁠](https://github.com/germainlefebvre4/tf2report/blob/main/LICENSE) for the software contained in this image.

As with all Docker images, these likely also contain other software which may be under other licenses (such as Bash, etc from the base distribution, along with any direct or indirect dependencies of the primary software being contained).

As for any pre-built image usage, it is the image user's responsibility to ensure that any use of this image complies with any relevant licenses for all software contained within.
