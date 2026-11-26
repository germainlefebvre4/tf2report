resource "terraform_data" "a" { input = "a" }
resource "terraform_data" "b" { input = "b" }
resource "terraform_data" "c" { triggers_replace = "c" }
