resource "terraform_data" "keep"    { input = "stable" }
resource "terraform_data" "change"  { input = "after" }
resource "terraform_data" "replace" { triggers_replace = "after" }
resource "terraform_data" "added"   { input = "new" }

output "added_id" { value = terraform_data.added.id }
