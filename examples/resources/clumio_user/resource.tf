resource "clumio_user" "example" {
  full_name = "full name"
  email     = "example@someorg.com"
  access_control_configuration = [
    {
      role_id                 = "role_id1"
      organizational_unit_ids = ["organizational_unit_id1"]
    },
    {
      role_id                 = "role_id2"
      organizational_unit_ids = ["organizational_unit_id2"]
    }
  ]
}
