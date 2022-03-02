# Create a Clumio policy with a 7-day RPO and 6-month retention
terraform {
  required_providers {
    clumio = {
      source  = "clumio.com/providers/clumio"
      version = "~>0.2.2"
    }
  }
}

resource "clumio_policy" "policy" {
  name = "Gold"
  operations {
    action_setting = "immediate"
    type           = "protection_group_backup"
    advanced_settings {
      protection_group_backup {
        backup_tier = "cold"
      }
    }
    slas {
      retention_duration {
        unit  = "months"
        value = 6
      }
      rpo_frequency {
        unit  = "days"
        value = 7
      }
    }
  }
}

# Create a Clumio protection group for S3
resource "clumio_protection_group" "protection_group" {
  name = "S3 Protection Group"
  object_filter {
    storage_classes = ["S3 Standard", "S3 Standard-IA"]
  }
}

# Assign a policy to the protection group
resource "clumio_policy_assignment" "protection_group_policy_assignment" {
  entity_type = "protection_group"
  entity_id   = clumio_protection_group.protection_group.id
  policy_id   = clumio_policy.policy.id
}
