resource "clumio_protection_group" "example" {
  name                   = "example-protection_group"
  description            = "example protection group"
  organizational_unit_id = "organizational_unit_id"
  object_filter {
    latest_version_only = false
    prefix_filters {
      excluded_sub_prefixes = ["prefix1", "prefix2"]
      prefix                = "prefix"
    }
    storage_classes = [
      "S3 Intelligent-Tiering", "S3 One Zone-IA", "S3 Standard", "S3 Standard-IA",
      "S3 Reduced Redundancy"
    ]
  }
}
