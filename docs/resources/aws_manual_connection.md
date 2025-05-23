---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "clumio_aws_manual_connection Resource - terraform-provider-clumio"
subcategory: ""
description: |-
  Clumio AWS Manual Connection Resource used to setup manual resources for connections.
---

# clumio_aws_manual_connection (Resource)

Clumio AWS Manual Connection Resource used to setup manual resources for connections.

## Example Usage

```terraform
resource "clumio_aws_manual_connection" "test_update_resources" {

  account_id = "aws_account_id"
  aws_region = "aws_region"
  assets_enabled = {
    ebs   = true
    rds   = true
    ddb   = true
    s3    = true
    mssql = false # Note that "mssql" is only available on legacy connections.
  }
  resources = {
    clumio_iam_role_arn     = "clumio_iam_role_arn"
    clumio_event_pub_arn    = "clumio_event_pub_arn"
    clumio_support_role_arn = "clumio_support_role_arn"
    event_rules = {
      cloudtrail_rule_arn = "cloudtrail_rule_arn"
      cloudwatch_rule_arn = "cloudwatch_rule_arn"
    }

    service_roles = {
      s3 = {
        continuous_backups_role_arn = "continuous_backups_role_arn"
      }
      mssql = {
        ssm_notification_role_arn    = "ssm_notification_role_arn"
        ec2_ssm_instance_profile_arn = "ec2_ssm_instance_profile_arn"
      }
    }
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `account_id` (String) Identifier of the AWS account to be linked with Clumio.
- `assets_enabled` (Object) Assets enabled for the connection. Note that `mssql` is only available for legacy connections. (see [below for nested schema](#nestedatt--assets_enabled))
- `aws_region` (String) Region of the AWS account to be linked with Clumio.
- `resources` (Object) An object containing the ARNs of the resources created for the manual AWS connection. Please refer to this guide for instructions on how to create them. - https://help.clumio.com/docs/manual-setup-for-aws-account-integration. If any of the ARNs are not applicable to the manual connection, provide an empty string "". (see [below for nested schema](#nestedatt--resources))

### Read-Only

- `id` (String) Unique identifier for the Clumio AWS manual connection.

<a id="nestedatt--assets_enabled"></a>
### Nested Schema for `assets_enabled`

Required:

- `ddb` (Boolean)
- `ebs` (Boolean)
- `mssql` (Boolean)
- `rds` (Boolean)
- `s3` (Boolean)


<a id="nestedatt--resources"></a>
### Nested Schema for `resources`

Required:

- `clumio_event_pub_arn` (String)
- `clumio_iam_role_arn` (String)
- `clumio_support_role_arn` (String)
- `event_rules` (Object) (see [below for nested schema](#nestedobjatt--resources--event_rules))
- `service_roles` (Object) (see [below for nested schema](#nestedobjatt--resources--service_roles))

<a id="nestedobjatt--resources--event_rules"></a>
### Nested Schema for `resources.event_rules`

Required:

- `cloudtrail_rule_arn` (String)
- `cloudwatch_rule_arn` (String)


<a id="nestedobjatt--resources--service_roles"></a>
### Nested Schema for `resources.service_roles`

Required:

- `mssql` (Object) (see [below for nested schema](#nestedobjatt--resources--service_roles--mssql))
- `s3` (Object) (see [below for nested schema](#nestedobjatt--resources--service_roles--s3))

<a id="nestedobjatt--resources--service_roles--mssql"></a>
### Nested Schema for `resources.service_roles.mssql`

Required:

- `ec2_ssm_instance_profile_arn` (String)
- `ssm_notification_role_arn` (String)


<a id="nestedobjatt--resources--service_roles--s3"></a>
### Nested Schema for `resources.service_roles.s3`

Required:

- `continuous_backups_role_arn` (String)
