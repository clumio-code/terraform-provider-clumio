resource "clumio_policy_rule" "example_1" {
  name           = "example-policy-rule-1"
  policy_id      = "policy_id"
  before_rule_id = clumio_policy_rule.example_2.id
  condition      = "{\"entity_type\":{\"$eq\":\"aws_ebs_volume\"}, \"aws_account_native_id\":{\"$in\":[\"aws_account_id_1\", \"aws_account_id_2\"]}, \"aws_tag\":{\"$eq\":{\"key\":\"aws_tag_key\", \"value\":\"aws_tag_value\"}}}"
}

resource "clumio_policy_rule" "example_2" {
  name           = "example-policy-rule-2"
  policy_id      = "policy_id"
  before_rule_id = ""
  condition      = "{\"entity_type\":{\"$eq\":\"aws_ec2_instance\"}, \"aws_account_native_id\":{\"$eq\":\"aws_account_id_1\"}, \"aws_region\":{\"$eq\":\"us-west-2\"}, \"aws_tag\":{\"$contains\":{\"key\":\"aws_tag_key_substr\", \"value\":\"aws_tag_value_substr\"}}}"
}
