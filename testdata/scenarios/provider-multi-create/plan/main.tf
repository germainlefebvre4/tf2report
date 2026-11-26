resource "aws_s3_bucket" "data" { bucket = "synthetic-data-bucket" }

resource "google_storage_bucket" "data" {
  name     = "synthetic-data-bucket"
  location = "EU"
}

resource "kubernetes_namespace" "app" {
  metadata {
    name = "synthetic-app"
  }
}

resource "github_repository" "app" {
  name       = "synthetic-app"
  visibility = "private"
}
