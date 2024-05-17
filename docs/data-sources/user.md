---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "clumio_user Data Source - terraform-provider-clumio"
subcategory: ""
description: |-
  clumio_user data source is used to retrieve details of a user for use in other resources.
---

# clumio_user (Data Source)

clumio_user data source is used to retrieve details of a user for use in other resources.

## Example Usage

```terraform
data "clumio_user" "ds_user" {
  name    = "user-name"
  role_id = "role-id"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `name` (String) The name of the user.
- `role_id` (String) Unique identifier of the role assigned to the user.

### Read-Only

- `users` (Attributes Set) Users that match the given name and/or role_id. (see [below for nested schema](#nestedatt--users))

<a id="nestedatt--users"></a>
### Nested Schema for `users`

Read-Only:

- `access_control_configuration` (Attributes Set) Identifiers of the organizational units, along with the identifier of the role assigned to the user. (see [below for nested schema](#nestedatt--users--access_control_configuration))
- `full_name` (String) The name of the user.
- `id` (String) Unique identifier of the user.

<a id="nestedatt--users--access_control_configuration"></a>
### Nested Schema for `users.access_control_configuration`

Read-Only:

- `organizational_unit_ids` (Set of String) Identifiers of the organizational units assigned to the user.
- `role_id` (String) Identifier of the role assigned to the user.