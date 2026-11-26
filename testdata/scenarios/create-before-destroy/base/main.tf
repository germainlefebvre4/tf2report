resource "terraform_data" "cbd" {
  triggers_replace = "before"
  lifecycle { create_before_destroy = true }
}
