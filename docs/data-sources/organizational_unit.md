---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "clumio_organizational_unit Data Source - terraform-provider-clumio"
subcategory: ""
description: |-
  clumio_organizational_unit data source is used to retrieve details of an organizational unit for use in other resources.
---

# clumio_organizational_unit (Data Source)

clumio_organizational_unit data source is used to retrieve details of an organizational unit for use in other resources.

## Example Usage

```terraform
data "clumio_organizational_unit" "example" {
  name = "organizational-unit-name"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the organizational unit.

### Read-Only

- `organizational_units` (Attributes Set) OrganizationalUnits which match the given name. (see [below for nested schema](#nestedatt--organizational_units))

<a id="nestedatt--organizational_units"></a>
### Nested Schema for `organizational_units`

Read-Only:

- `descendant_ids` (Set of String) List of all recursive descendent organizational units.
- `description` (String) Brief description to denote details of the organizational unit.
- `id` (String) Unique identifier of the organizational unit.
- `name` (String) The name of the organizational unit.
- `parent_id` (String) The identifier of the parent organizational unit under which the organizational unit was created.
- `users_with_role` (Attributes Set) List of user ids, with role assigned to this organizational unit. (see [below for nested schema](#nestedatt--organizational_units--users_with_role))

<a id="nestedatt--organizational_units--users_with_role"></a>
### Nested Schema for `organizational_units.users_with_role`

Read-Only:

- `assigned_role` (String) Identifier of the role associated with the user assigned to the organizational unit.
- `user_id` (String) Identifier of the user assigned to the organizational unit.
