terraform {
  required_providers {
    random = { source = "hashicorp/random", version = "= 3.7.2" }
  }
}

# Provider instances created with for_each are specific to OpenTofu.
provider "random" {
  alias    = "region"
  for_each = toset(["eu", "us"])
}

resource "random_pet" "region" {
  for_each = toset(["eu", "us"])
  provider = random.region[each.key]
}

resource "terraform_data" "anchor" { input = "x" }
