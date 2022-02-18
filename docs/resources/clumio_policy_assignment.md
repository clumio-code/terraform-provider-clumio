---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "clumio_policy_assignment Resource - terraform-provider-clumio-internal"
subcategory: ""
description: |-
  Clumio Policy Assignment Resource used to assign (or unassign) policies.
  NOTE: Currently policy assignment is supported only for entity type "protection_group".
---

# clumio_policy_assignment (Resource)

Clumio Policy Assignment Resource used to assign (or unassign) policies.

 NOTE: Currently policy assignment is supported only for entity type "protection_group".

## Example Usage

```terraform
resource "clumio_policy_assignment" "example" {
  entity_id   = "entity_id"
  entity_type = "protection_group"
  policy_id   = "policy_id"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **entity_id** (String) The entity id.
- **entity_type** (String) The entity type. The supported entity type is"protection_group".
- **policy_id** (String) The Clumio-assigned ID of the policy.

### Optional

- **id** (String) The ID of this resource.
- **organizational_unit_id** (String) The Clumio-assigned ID of the organizational unit to use as the context for assigning the policy.

