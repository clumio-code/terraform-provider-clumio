---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "clumio_callback_resource Resource - terraform-provider-clumio"
subcategory: ""
description: |-
  Clumio Callback Resource used while on-boarding AWS clients. The purpose of this resource is to send a SNS event with the necessary details of the AWS connection configuration done on the client AWS account so that necessary connection post processing can be done in Clumio.
---

# clumio_callback_resource (Resource)

Clumio Callback Resource used while on-boarding AWS clients. The purpose of this resource is to send a SNS event with the necessary details of the AWS connection configuration done on the client AWS account so that necessary connection post processing can be done in Clumio.

## Example Usage

```terraform
resource "clumio_callback_resource" "example" {
  # example configuration here
  topic               = "mytopic"
  token               = "mytoken"
  role_external_id    = "role_external_id"
  account_id          = "account_id"
  region              = "region"
  role_id             = "role_id"
  role_arn            = "role_arn"
  clumio_event_pub_id = "clumio_event_pub_id"
  type                = "type"
  properties = {
    "prop1" : {
      "key1" : val1
    }
  }
  config_version                     = "1"
  discover_version                   = "3"
  protect_config_version             = "18"
  protect_ebs_version                = "19"
  protect_rds_version                = "18"
  protect_ec2_mssql_version          = "1"
  protect_warm_tier_version          = "2"
  protect_warm_tier_dynamodb_version = "2"
  protect_s3_version                 = "1"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `account_id` (String) The AWS Customer Account ID.
- `bucket_name` (String) S3 bucket name where the status file is written.
- `canonical_user` (String) Canonical User ID of the account.
- `clumio_event_pub_id` (String) Clumio Event Pub SNS topic ID.
- `config_version` (String) Clumio Config version.
- `discover_version` (String) Clumio Discover version.
- `region` (String) The AWS Region.
- `role_arn` (String) Clumio IAM Role Arn.
- `role_external_id` (String) A key that must be used by Clumio to assume the service role in your account. This should be a secure string, like a password, but it does not need to be remembered (random characters are best).
- `role_id` (String) Clumio IAM Role ID.
- `sns_topic` (String) SNS Topic to publish event.
- `token` (String) The AWS integration ID token.
- `type` (String) Registration Type.

### Optional

- `discover_enabled` (Boolean) Is Clumio Discover enabled.
- `properties` (Map of String) Properties to be passed in the SNS event.
- `protect_config_version` (String) Clumio Protect Config version.
- `protect_dynamodb_version` (String) Clumio DynamoDB Protect version.
- `protect_ebs_version` (String) Clumio EBS Protect version.
- `protect_ec2_mssql_version` (String) Clumio EC2 MSSQL Protect version.
- `protect_enabled` (Boolean) Is Clumio Protect enabled.
- `protect_rds_version` (String) Clumio RDS Protect version.
- `protect_s3_version` (String) Clumio S3 Protect version.
- `protect_warm_tier_dynamodb_version` (String) Clumio DynamoDB Warmtier Protect version.
- `protect_warm_tier_version` (String) Clumio Warmtier Protect version.

### Read-Only

- `id` (String) The ID of this resource.

