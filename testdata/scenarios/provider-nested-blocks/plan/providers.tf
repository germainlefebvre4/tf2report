terraform {
  required_providers {
    aws        = { source = "hashicorp/aws", version = "= 5.100.0" }
    google     = { source = "hashicorp/google", version = "= 6.50.0" }
    kubernetes = { source = "hashicorp/kubernetes", version = "= 2.38.0" }
    github     = { source = "integrations/github", version = "= 6.6.0" }
  }
}

# Synthetic credentials and no remote call: every provider is configured so that planning works offline.
provider "aws" {
  region                      = "eu-west-1"
  access_key                  = "SYNTHETIC"
  secret_key                  = "SYNTHETIC"
  skip_credentials_validation = true
  skip_requesting_account_id  = true
  skip_metadata_api_check     = true
  skip_region_validation      = true
}

provider "google" {
  project      = "synthetic-project"
  region       = "europe-west1"
  access_token = "SYNTHETIC"
}

provider "kubernetes" {
  host     = "https://127.0.0.1:1"
  token    = "SYNTHETIC"
  insecure = true
}

provider "github" {
  owner = "synthetic-org"
  token = "SYNTHETIC"
}

