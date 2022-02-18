# Terraform Provider Clumio

The Terraform Clumio provider is a plugin for Terraform that allows for the full lifecycle
management of Clumio resources.

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

For versions older than v0.1.3, install the provider using the following command:

```sh
$ curl https://raw.githubusercontent.com/clumio-code/terraform-provider-clumio/main/installer.sh | bash -s -- v0.1.2
```