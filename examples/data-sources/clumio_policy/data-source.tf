data "clumio_policy" "example" {
  name = "example_policy"
  activation_status = "activated"
  operation_types = ["protection_group_backup", "aws_ebs_volume_backup"]
}
