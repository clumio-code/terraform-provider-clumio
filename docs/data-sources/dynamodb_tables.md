---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "clumio_dynamodb_tables Data Source - terraform-provider-clumio"
subcategory: ""
description: |-
  clumio_dynamo_db_tables data source is used to retrieve details of the DynamoDB tables for use in other resources.
---

# clumio_dynamodb_tables (Data Source)

clumio_dynamo_db_tables data source is used to retrieve details of the DynamoDB tables for use in other resources.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `account_native_id` (String) The identifier of the AWS account under which the DynamoDB bucket was created.
- `aws_region` (String) The AWS region associated with the DynamoDB tables.

### Optional

- `name` (String) The DynamoDB table name to be queried.
- `table_native_id` (String) Native identifier of the DynamoDB table to be queried.

### Read-Only

- `dynamodb_tables` (Attributes Set) List of DynamoDB tables which matched the query criteria. (see [below for nested schema](#nestedatt--dynamodb_tables))

<a id="nestedatt--dynamodb_tables"></a>
### Nested Schema for `dynamodb_tables`

Read-Only:

- `id` (String) Unique identifier of the DynamoDB table.
- `name` (String) Name of the DynamoDB table.
- `table_native_id` (String) Native identifier of the DynamoDB table.