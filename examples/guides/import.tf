terraform {
  required_providers {
    clumio = {
      source  = "clumio.com/providers/clumio"
      version = "~>0.4.0"
    }
  }
}

# Instantiate the Clumio provider
provider "clumio" {
  clumio_api_token    = "<clumio_api_token>"
  clumio_api_base_url = "<clumio_api_base_url>"
}

# Instantiate an empty policy
resource "clumio_policy" "policy" {
}
