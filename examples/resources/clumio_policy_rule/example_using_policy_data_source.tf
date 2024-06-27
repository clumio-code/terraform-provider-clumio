data  "clumio_policy" "test_policy" {
  name = "test-policy"
}

resource "clumio_policy_rule" "example_2" {
  name           = "example-policy-rule-2"
  # Using the clumio_policy data source to get the id of the policy. For more information on the
  # clumio_policy data source, please refer to the data source documentation.
  policy_id      = tolist(data.clumio_policy.test_policy.policies)[0].id
  before_rule_id = ""
  condition      = "{\"entity_type\":{\"$eq\":\"aws_ec2_instance\"}, \"aws_account_native_id\":{\"$eq\":\"aws_account_id_1\"}, \"aws_region\":{\"$eq\":\"us-west-2\"}, \"aws_tag\":{\"$contains\":{\"key\":\"aws_tag_key_substr\", \"value\":\"aws_tag_value_substr\"}}}"
}
