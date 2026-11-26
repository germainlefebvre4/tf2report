terraform {
  required_providers {
    random = { source = "hashicorp/random", version = "= 3.7.2" }
    local  = { source = "hashicorp/local", version = "= 2.5.3" }
  }
}

resource "random_pet" "name" { length = 2 }

resource "random_password" "db" {
  length  = 16
  special = true
}

resource "local_file" "greeting" {
  filename = "greeting.txt"
  content  = "SYNTHETIC-content"
}

resource "local_sensitive_file" "secret" {
  filename = "secret.txt"
  content  = "SYNTHETIC-SECRET-CONTENT"
}
