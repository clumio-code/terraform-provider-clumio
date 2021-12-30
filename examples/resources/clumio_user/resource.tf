resource "clumio_user" "example" {
  full_name               = "full name"
  email                   = "example@someorg.com"
  organizational_unit_ids = ["organizational_unit_id1"]
}