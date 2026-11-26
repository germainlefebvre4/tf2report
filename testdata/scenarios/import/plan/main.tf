resource "terraform_data" "adopted" { input = "adopted" }

import {
  to = terraform_data.adopted
  id = "00000000-0000-0000-0000-000000000001"
}
