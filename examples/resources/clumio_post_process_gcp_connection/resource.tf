resource "clumio_post_process_gcp_connection" "test" {
  project_id            = "project_id"
  project_number        = "project_number"
  project_name          = "project_name"
  token                 = "token"
  service_account_email = "service_account_email_changed"
  wif_pool_id           = "wif_pool_id"
  wif_provider_id       = "wif_provider_id"
  config_version        = "1.2"
  protect_gcs_version   = "1.1"
  regions               = ["us-west1"]
  region_configuration = [
    {
      region                       = "us-west1"
      inventory_bridge_bucket_name = "clumio-inventory-bridge-us-west1-project_id"
    }
  ]
  properties = {
    key = "value"
  }
}
