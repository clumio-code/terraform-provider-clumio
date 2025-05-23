---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "clumio_aws_manual_connection_resources Data Source - terraform-provider-clumio"
subcategory: ""
description: |-
  Clumio AWS Manual Connection Resources Datasource to get resources for manual connections.
---

# clumio_aws_manual_connection_resources (Data Source)

Clumio AWS Manual Connection Resources Datasource to get resources for manual connections.

## Example Usage

```terraform
data "clumio_aws_manual_connection_resources" "test_get_resources" {
  account_native_id = "aws_account_id"
  aws_region        = "aws_region"
  asset_types_enabled = {
    ebs   = true
    rds   = true
    ddb   = true
    s3    = true
    mssql = false # Note that "mssql" is only available on legacy connections.
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `account_native_id` (String) AWS Account ID to be connected to Clumio.
- `asset_types_enabled` (Object) Assets to be connected to Clumio. Note that `mssql` is only available for legacy connections. (see [below for nested schema](#nestedatt--asset_types_enabled))
- `aws_region` (String) AWS Region to be connected to Clumio.

### Read-Only

- `id` (String) Combination of provided Account Native ID and AWS Region.
- `resources` (String) Generated manual resources for provided configuration.

<a id="nestedatt--asset_types_enabled"></a>
### Nested Schema for `asset_types_enabled`

Required:

- `ddb` (Boolean)
- `ebs` (Boolean)
- `mssql` (Boolean)
- `rds` (Boolean)
- `s3` (Boolean)
