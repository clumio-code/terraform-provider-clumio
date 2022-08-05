resource "clumio_post_process_kms" "example" {
  # example configuration here
  token                   = "mytoken"
  account_id              = "account_id"
  region                  = "region"
  multi_region_cmk_key_id = "multi_region_cmk_key_id"
  stack_set_id            = "stack_set_id"
  other_regions           = "other_regions"
}
