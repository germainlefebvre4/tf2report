resource "terraform_data" "new_name" { input = "same" }
resource "terraform_data" "renamed_and_changed" { input = "after" }

moved {
  from = terraform_data.old_name
  to   = terraform_data.new_name
}
moved {
  from = terraform_data.moved_and_changed
  to   = terraform_data.renamed_and_changed
}
