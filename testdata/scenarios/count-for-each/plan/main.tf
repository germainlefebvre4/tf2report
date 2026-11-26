resource "terraform_data" "counted" {
  count = 3
  input = "item-${count.index}"
}

resource "terraform_data" "keyed" {
  for_each = toset(["eu-west-1", "us-east-1", "with space", "quote\"d", "ünïcode-é"])
  input    = each.key
}
