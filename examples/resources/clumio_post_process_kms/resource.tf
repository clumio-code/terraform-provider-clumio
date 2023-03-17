resource "clumio_post_process_kms" "example" {
  # example configuration here
  token                   = "mytoken"
  account_id              = "account_id"
  region                  = "region"
  multi_region_cmk_key_id = "multi_region_cmk_key_id"
  role_external_id        = "role_external_id"
  role_arn                = "role_arn"
  role_id                 = "role_id"
}
