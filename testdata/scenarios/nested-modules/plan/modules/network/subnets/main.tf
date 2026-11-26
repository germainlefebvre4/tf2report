resource "terraform_data" "subnet" {
  count = 2
  input = "subnet-${count.index}"
}
