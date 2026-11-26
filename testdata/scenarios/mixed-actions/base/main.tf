resource "terraform_data" "keep"    { input = "stable" }
resource "terraform_data" "change"  { input = "before" }
resource "terraform_data" "replace" { triggers_replace = "before" }
resource "terraform_data" "gone"    { input = "to-delete" }
