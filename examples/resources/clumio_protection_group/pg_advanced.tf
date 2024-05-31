resource "clumio_protection_group" "advanced_example" {
  name                   = "example-protection_group-2"
  description            = "example protection group-2"
  bucket_rule = jsonencode(
    {
      "aws_tag": {
        "$eq": {
          "key": "Environment",
          "value": "Prod"
        }
      },
      "account_native_id": {
        "$eq": "AWS Account ID"
      },
      "aws_region": {
        "$eq": "AWS Region"
      }
    }
  )
  object_filter {
    latest_version_only = false
    prefix_filters {
      excluded_sub_prefixes = ["prefix1", "prefix2"]
      prefix                = "prefix"
    }
    storage_classes = [
      "S3 Standard", "S3 Standard-IA", "S3 Intelligent-Tiering", "S3 One Zone-IA",
      "S3 Reduced Redundancy"
    ]
  }
}
