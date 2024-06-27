resource "clumio_protection_group" "prefix_filter_2" {
  name                   = "example-protection_group-4"
  description            = "example protection group-4"
  object_filter {
    latest_version_only = false
    prefix_filters {
      excluded_sub_prefixes = ["abcd", "/xyz"] # exludes all objects that start with path abcd or /xyz
      prefix                = "" # include all objects in the bucket
    }
    storage_classes = ["S3 Standard"]
  }
}
