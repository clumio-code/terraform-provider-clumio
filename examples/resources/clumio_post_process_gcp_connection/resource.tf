resource "clumio_post_process_gcp_connection" "test" {
  project_id            = "project_id"
  project_number        = "project_number"
  project_name          = "project_name"
  token                 = "token"
  service_account_email = "service_account_email_changed"
  wif_pool_id           = "wif_pool_id"
  wif_provider_id       = "wif_provider_id"
  clumio_aws_iam_role   = "clumio_aws_iam_role"
  config_version        = "1.2"
  protect_gcs_version   = "1.1"
  properties = {
    key = "value"
  }
}
