# Migration guide

This document is meant to help you migrate your Clumio Terraform config to the newer version.
In migration guides, we will only describe deprecations or breaking changes and help you to change your configuration to keep the same (or similar) behavior across different versions.

## v0.6.2 ➞ v0.7.0

### Deprecated attribute organizational_unit_id
- [Immediate action](#actionou)
- [Sample Warning](#warningou)
- [Configs](#configsou)
- [Existing Sample Config with organizational_unit_id](#sampleou1)
- [Sample Config with organizational_unit_id removed from un-supported resources](#sampleou2)

Release 0.7.0 onwards of the Terraform provider
organizational_unit_id attribute is being deprecated as an attribute from the following resources below and will be un-supported in a minor release after 30 days from the 0.7.0 release.
* `clumio_aws_connection`
* `clumio_protection_group`
* `clumio_policy`
* `clumio_policy_assignment`
* `clumio_policy_rule`

<a name="Immediate action"></a>
#### actionou
There is no immediate action needed on the user part. If your Terraform config files specify organizational_unit_id attribute in the above mentioned resources, then you may notice the following warning when Terraform operations like apply or plan are triggered. However we recommend updating the config in preperation of organizational_unit_id being un-supported.

<a name="Sample Warning"></a>
#### warningou
Type of warning noticed, when organizational_unit_id is specified for the above resources in the Terraform config file:

```
Terraform has compared your real infrastructure against your configuration and found no differences, so no changes are needed.
╷
│ Warning: Attribute Deprecated
│ 
│   with clumio_policy.test_policy_1,
│   on main.tf line 33, in resource "clumio_policy" "test_policy_1":
│   33:   organizational_unit_id = clumio_organizational_unit.test_ou_0.id
│ 
│ Use the provider schema attribute clumio_organizational_unit_context to create the resource in the context of  an Organizational Unit.
```
<a name="Configs"></a>
#### configsou

<a name="Existing Sample Config with organizational_unit_id"></a>
#### sampleou1
Sample configuration using organizational_unit_id attribute in resources:

```
provider "clumio" {
  clumio_api_token    = ""
  clumio_api_base_url = ""
}

resource "clumio_organizational_unit" "test_ou_0" {
  name = "test_ou_0"
}


resource "clumio_policy" "test_policy_1" {
  name              = "test_policy_1"
  organizational_unit_id = clumio_organizational_unit.test_ou_0.id
  operations {
    action_setting = "immediate"
    type           = "protection_group_backup"
    slas {
      retention_duration {
        unit  = "months"
        value = 3
      }
      rpo_frequency {
        unit  = "days"
        value = 1
      }
    }
    advanced_settings {
      protection_group_backup {
        backup_tier = "cold"
      }
    }
  }
}
```

<a name="Sample Config with organizational_unit_id removed from un-supported resources"></a>
#### sampleou2
Suggested configuration change allowing to allow managing the resources created on other organizational units:
Please note the initialization of multiple clumio providers with and without clumio_organizational_unit_context.

```
provider "clumio" {
  clumio_api_token    = ""
  clumio_api_base_url = ""
}

resource "clumio_organizational_unit" "test_ou_0" {
  name = "test_ou_0"
}

provider "clumio" {
  alias = "clumio_ou_0"
  clumio_api_token    = ""
  clumio_api_base_url = ""
  clumio_organizational_unit_context = "" #populate OU ID
}

resource "clumio_ou_0_policy" "test_policy_1" { ### context aware clumio provider
  provider          = clumio.clumio_ou_0
  name              = "test_policy_1"
  # organizational_unit_id = clumio_organizational_unit.test_ou_0.id ### remove this attribute
  operations {
    action_setting = "immediate"
    type           = "protection_group_backup"
    slas {
      retention_duration {
        unit  = "months"
        value = 3
      }
      rpo_frequency {
        unit  = "days"
        value = 1
      }
    }
    advanced_settings {
      protection_group_backup {
        backup_tier = "cold"
      }
    }
  }
}

```
### Deprecated attribute timezone in resource clumio_policy
- [Immediate action](#actiontz)
- [Sample Warning](#warningtz)
- [Configs](#configstz)
- [Existing Sample Config with timezone at the toplevel of policy specification](#sampletz1)
- [Sample Config using timezone attribute at the operations specification](#sampletz2)
Release 0.7.0 onwards of the Terraform provider
timezone attribute in clumio_policy at the toplevel of policy specification is being deprecated as an attribute.
timezone attribute is now supported at the operations part of the resource config, which provides
more flexbility in allowing to specify timezone for each operation type. timezone at the top level of the policy specification will be un-supported in a minor release after 30 days from the 0.7.0 release

<a name="Immediate action"></a>
#### actiontz
There is no immediate action needed on the user part. If your Terraform config files specify timezone at the toplevel, then you may notice the following warning when Terraform operations like apply or plan are triggered.
 However we recommend updating the config in preperation of timezone being moved to the operations level of the policy specification.

<a name="Sample Warning"></a>
#### warningtz
Type of warning noticed, when timezone is specified at the top level of policy specification in the Terraform config file:

```
Terraform has compared your real infrastructure against your configuration and found no differences, so no changes are needed.
╷
│ Warning: Attribute Deprecated
│ 
│   with clumio_policy.example_backup_windown_timezone,
│   on main.tf line 15, in resource "clumio_policy" "example_backup_windown_timezone":
│   15:   timezone    = "America/Los_Angeles"
│ 
```
<a name="Configs"></a>
#### configstz

<a name="Existing Sample Config with timezone at the toplevel of policy specification"></a>
#### sampletz1
Sample configuration using timezone attribute at the toplevel of policy specification:

```
resource "clumio_policy" "example_backup_windown_timezone" {
  name              = "example-policy-Backup-Window-Timezone"
  activation_status = "activated"
  timezone          = "America/Los_Angeles" ### being deprecated
  operations {
    action_setting = "immediate"
    type           = "aws_ebs_volume_backup"
    slas {
      retention_duration {
        unit  = "days"
        value = 30
      }
      rpo_frequency {
        unit  = "days"
        value = 1
      }
    }
    backup_window_tz {
      start_time = "America/Los_Angeles"
    }
  }
}
```

<a name="Sample Config using timezone attribute at the operations specification"></a>
#### sampletz2
Suggested configuration change:
Please note timezone is specified ast the part of the operations part of the specification.

```
resource "clumio_policy" "example_backup_windown_timezone" {
  name              = "example-policy-Backup-Window-Timezone"
  activation_status = "activated"
  operations {
    action_setting = "immediate"
    type           = "aws_ebs_volume_backup"
    slas {
      retention_duration {
        unit  = "days"
        value = 30
      }
      rpo_frequency {
        unit  = "days"
        value = 1
      }
    }
    timezone          = "America/Los_Angeles" ### use this attribute instead
    backup_window_tz {
      start_time = "America/Los_Angeles"
    }
  }
}
```
