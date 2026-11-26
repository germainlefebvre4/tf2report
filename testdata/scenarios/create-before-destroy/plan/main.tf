resource "terraform_data" "cbd" {
  triggers_replace = "after"
  lifecycle { create_before_destroy = true }
}
