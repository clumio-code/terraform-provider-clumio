terraform {
  required_providers {
    clumio = {
      source  = "clumio.com/providers/clumio"
      version = "~>0.2.2"
    }
  }
}

provider "clumio" {
  alias                              = "child_ou"
  clumio_api_token                   = "<clumio_api_token>"
  clumio_api_base_url                = "<clumio_api_base_url>"
  clumio_organizational_unit_context = "<clumio_ou_id>"
}

resource "clumio_policy_rule" "child_ou_rule" {
  provider = clumio.child_ou
  # other arguments
}
