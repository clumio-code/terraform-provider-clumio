# Terraform Provider Clumio

The Terraform Clumio provider is a plugin for Terraform that allows for the full lifecycle
management of Clumio resources.

NOTE: 0.1.x versions have been deprecated.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 0.14.x
- [Go](https://golang.org/doc/install) >= 1.17

## Installing the provider

The Clumio Terraform Provider is available in
the [Terraform Registry](https://registry.terraform.io/providers/clumio-code/clumio/latest)
.

The following terraform block will install the provider upon "terraform init":

```
terraform {
  required_providers {
    clumio = {
      source = "clumio-code/clumio"
      version = "~>0.2.1"
    }
  }
}
```
