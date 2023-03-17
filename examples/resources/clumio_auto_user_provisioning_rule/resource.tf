resource "clumio_auto_user_provisioning_rule" "example_1" {
  name                    = "example-auto-user-provisioning-rule-1"
  condition               = "{\"user.groups\":{\"$eq\":\"Group1\"}}"
  role_id                 = "role_id"
  organizational_unit_ids = ["organizational_unit_id1"]
}

resource "clumio_auto_user_provisioning_rule" "example_2" {
  name                    = "example-auto-user-provisioning-rule-2"
  condition               = "{\"user.groups\":{\"$in\":[\"Group1\",\"Group2\"]}}"
  role_id                 = "role_id"
  organizational_unit_ids = ["organizational_unit_id1"]
}

resource "clumio_auto_user_provisioning_rule" "example_3" {
  name                    = "example-auto-user-provisioning-rule-3"
  condition               = "{\"user.groups\":{\"$all\":[\"Group1\",\"Group2\"]}}"
  role_id                 = "role_id"
  organizational_unit_ids = ["organizational_unit_id1"]
}

resource "clumio_auto_user_provisioning_rule" "example_4" {
  name                    = "example-auto-user-provisioning-rule-4"
  condition               = "{\"user.groups\":{\"$contains\":{\"$in\":[\"Group1\",\"Group2\"]}}}"
  role_id                 = "role_id"
  organizational_unit_ids = ["organizational_unit_id1"]
}

resource "clumio_auto_user_provisioning_rule" "example_5" {
  name                    = "example-auto-user-provisioning-rule-5"
  condition               = "{\"user.groups\":{\"$contains\":{\"$all\":[\"Group1\",\"Group2\"]}}}"
  role_id                 = "role_id"
  organizational_unit_ids = ["organizational_unit_id1"]
}
