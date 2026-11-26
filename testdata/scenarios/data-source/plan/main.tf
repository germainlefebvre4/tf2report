resource "terraform_data" "first" { input = "first" }

data "terraform_remote_state" "other" {
  backend    = "local"
  config     = { path = "${path.module}/other.tfstate" }
  depends_on = [terraform_data.first]
}
