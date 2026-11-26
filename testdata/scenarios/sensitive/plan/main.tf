variable "db_password" {
  type      = string
  sensitive = true
  default   = "SYNTHETIC-SECRET-PASSWORD"
}

resource "terraform_data" "database" {
  input = {
    username = "synthetic"
    password = var.db_password
  }
}

output "password" {
  value     = var.db_password
  sensitive = true
}
