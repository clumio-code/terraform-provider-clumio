resource "clumio_protection_group" "basic_example" {
  name        = "example-protection_group-1"
  description = "example protection group-1"
  bucket_rule = jsonencode({
    "aws_tag" = {
      "$eq" = {
        "key"   = "Key1"
        "value" = "Value1"
      }
    }
  })
  object_filter {
    storage_classes = [
      "S3 Standard"
    ]
  }
}
