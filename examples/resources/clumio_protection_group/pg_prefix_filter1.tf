resource "clumio_protection_group" "prefix_filter_1" {
  name        = "example-protection_group-3"
  description = "example protection group-3"
  object_filter {
    latest_version_only = false
    prefix_filters {
      excluded_sub_prefixes = ["abc", "xyz"] # exludes /abc, /xyz only /123, /foo are backed up
      prefix                = "/"            # /abc, /xyz, /123, /foo
    }
    storage_classes = ["S3 Standard"]
  }
}
