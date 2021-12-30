resource "clumio_organizational_unit" "example" {
  name        = "Organizatonal Unit Name"
  description = "Description of Organizational Unit"
  parent_id   = "parent organizational unit id"
  users       = ["user-id1", "user-id2"]
}