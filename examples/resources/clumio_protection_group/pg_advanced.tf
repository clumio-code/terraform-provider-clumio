resource "clumio_protection_group" "advanced_example" {
  name        = "example-protection_group-2"
  description = "example protection group-2"
  bucket_rule = jsonencode(
    {
      "aws_tag" : {
        "$eq" : {
          "key" : "Key1",
          "value" : "Value1"
        }
      },
      "account_native_id" : {
        "$eq" : "123456789012"
      },
      "aws_region" : {
        "$eq" : "us-west-2"
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
