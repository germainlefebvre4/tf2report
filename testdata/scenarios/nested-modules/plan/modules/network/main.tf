resource "terraform_data" "vpc" { input = "vpc" }

module "subnets" {
  source = "./subnets"
}
