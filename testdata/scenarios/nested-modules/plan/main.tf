module "network" {
  source = "./modules/network"
}

module "app" {
  source = "./modules/app"
  count  = 2
}

module "region" {
  source   = "./modules/app"
  for_each = toset(["eu", "us"])
}
