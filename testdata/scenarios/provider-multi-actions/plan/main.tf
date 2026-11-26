# bucket name is immutable: replace
resource "aws_s3_bucket" "data" { bucket = "synthetic-data-bucket-v2" }

# aws_s3_bucket.logs is in the prior state only: delete

# in-place update
resource "google_storage_bucket" "data" {
  name     = "synthetic-data-bucket"
  location = "EU"
  labels   = { env = "prod" }
}

# unchanged: no-op
resource "kubernetes_namespace" "app" {
  metadata {
    name = "synthetic-app"
  }
}

# new: create
resource "github_repository" "app" {
  name       = "synthetic-app"
  visibility = "private"
}
