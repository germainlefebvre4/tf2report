# Example plans

Plans to try `tf2report` on. They contain no real infrastructure data, and no example output is stored
here: run the tool to see the report.

```bash
tf2report --plan examples/sample-plan.json
tf2report --plan examples/sample-plan.json --format json --type aws_instance
```

| File | What it shows |
|------|---------------|
| `sample-plan.json` | Add, change, replace and destroy across AWS, Azure, Google, Kubernetes, Datadog and Cloudflare, plus one no-op |
| `replace-create-before-destroy.json` | A replacement ordered `create` then `delete` (`create_before_destroy`) |
| `replace-forced.json` | Replacements requested with `-replace` and caused by a tainted resource, plus one untouched resource |
| `opentofu-replace-forced.json` | The same scenario planned by OpenTofu 1.9.0 instead of Terraform |
| `moved-resources.json` | A resource moved and changed in the same plan, and a pure move |
| `imported-resource.json` | A resource adopted through an `import` block |
| `no-changes.json` | A plan where every resource is up to date |
| `empty-config.json` | A plan for a configuration without any resource |
| `legacy-format-0.1.json` | The older plan format (`format_version` 0.1, Terraform 0.12) |

Every file here is a copy of a fixture under `testdata/plans/`, and a test keeps both in sync and checks
that `tf2report` reads all of them. The files named after a scenario were produced with Terraform by
`scripts/gen-plans.sh`; see `testdata/README.md`.

Only plans the tool can read today are listed. Plans that create or delete a resource in the exact form
Terraform exports are currently rejected (see the Known Gaps in
`specs/001-baseline-cli-behavior/spec.md`), so they live in `testdata/plans/terraform/` only.
