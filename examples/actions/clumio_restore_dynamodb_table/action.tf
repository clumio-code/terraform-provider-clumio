# Example restoring the latest SecureVault backup of a DynamoDB table to a new table. The action
# initiates the restore and reports the Clumio task tracking it without waiting for completion.
data "clumio_aws_environment" "example_aws_environment" {
  account_native_id = "AWS Account ID"
  aws_region        = "AWS Region"
}

data "clumio_dynamodb_backups" "example_backups" {
  table_id = "Clumio-assigned ID of the DynamoDB table"
  type     = "clumio_backup"
}

action "clumio_restore_dynamodb_table" "restore_from_backup" {
  config {
    source = {
      securevault_backup = {
        backup_id = data.clumio_dynamodb_backups.example_backups.backups[0].id
      }
    }
    target = {
      environment_id = data.clumio_aws_environment.example_aws_environment.id
      table_name     = "Name of the restored DynamoDB table"
    }
  }
}

# Example restoring a DynamoDB table to its latest restorable time using continuous backup
# (point-in-time restore)
action "clumio_restore_dynamodb_table" "restore_point_in_time" {
  config {
    source = {
      continuous_backup = {
        table_id                   = "Clumio-assigned ID of the DynamoDB table"
        use_latest_restorable_time = true
      }
    }
    target = {
      environment_id = data.clumio_aws_environment.example_aws_environment.id
      table_name     = "Name of the restored DynamoDB table"
    }
  }
}
