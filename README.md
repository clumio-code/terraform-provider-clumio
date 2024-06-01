# Clumio Provider for Terraform

![Clumio](.github/logo.svg)

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=clumio_terraform-provider-clumio-internal&metric=alert_status&token=0a73b8177110fee6d39be9057997d2d666c0d662)](https://sonarcloud.io/summary/new_code?id=clumio_terraform-provider-clumio-internal)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=clumio_terraform-provider-clumio-internal&metric=security_rating&token=0a73b8177110fee6d39be9057997d2d666c0d662)](https://sonarcloud.io/summary/new_code?id=clumio_terraform-provider-clumio-internal)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=clumio_terraform-provider-clumio-internal&metric=reliability_rating&token=0a73b8177110fee6d39be9057997d2d666c0d662)](https://sonarcloud.io/summary/new_code?id=clumio_terraform-provider-clumio-internal)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=clumio_terraform-provider-clumio-internal&metric=sqale_rating&token=0a73b8177110fee6d39be9057997d2d666c0d662)](https://sonarcloud.io/summary/new_code?id=clumio_terraform-provider-clumio-internal)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=clumio_terraform-provider-clumio-internal&metric=coverage&token=0a73b8177110fee6d39be9057997d2d666c0d662)](https://sonarcloud.io/summary/new_code?id=clumio_terraform-provider-clumio-internal)
[![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=clumio_terraform-provider-clumio-internal&metric=ncloc&token=0a73b8177110fee6d39be9057997d2d666c0d662)](https://sonarcloud.io/summary/new_code?id=clumio_terraform-provider-clumio-internal)

The Clumio provider for Terraform is a plugin designed to facilitate the lifecycle management of
[Clumio](https://clumio.com/) resources. It simplifies the integration of Clumio's backup-as-a-service
for AWS, streamlining various tasks. These tasks range from linking multiple AWS accounts and
regions, establishing data protection policies and rules, to managing users and setting up
organizational units. The Clumio provider offers an easy to define, reproducible approach to
creating a data protection environment.

This provider makes use of the [Clumio REST API](https://help.clumio.com/reference) along with the
functionalities of the [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework).
The most current version of the Clumio provider is available in the
[Terraform Registry](https://registry.terraform.io/providers/clumio-code/clumio/latest).

## Getting Started

Initiating the use of the Clumio provider requires a valid API token from Clumio. Detailed
instructions for generating this token are available here:
[Creating an API Token](https://help.clumio.com/docs/api-tokens#managing-tokens)

Below is a basic Terraform configuration aimed at configuring the Clumio provider. Substitute
`<clumio_api_token>` with your specific API token and `<clumio_api_base_url>` with the corresponding
API base URL that aligns with the Clumio portal you are utilizing. For specifics about the API base
URLs, refer to the [Clumio provider documentation](https://registry.terraform.io/providers/clumio-code/clumio/latest/docs).
To install and initialize the provider, run `terraform init` followed by `terraform apply`:

```
terraform {
  required_providers {
    clumio = {
      source  = "clumio-code/clumio"
      version = "~>0.8.0"
    }
    aws = {}
  }
}

# Instantiate the Clumio provider
provider "clumio" {
  clumio_api_token    = "<clumio_api_token>"
  clumio_api_base_url = "<clumio_api_base_url>"
}
```

For more detailed guidance on how to provision resources in the Clumio cloud and integrate an AWS
environment for data protection, please refer to the
[Getting Started Guide](https://registry.terraform.io/providers/clumio-code/clumio/latest/docs/guides/getting_started).

NOTE: 0.1.x versions have been deprecated.

## Terraform CLI Compatibility

The provider is built with the terraform-plugin-go module and is compatible with
Terraform v1.0 and later.


## Building the Provider

To build the provider, ensure that [Go](https://go.dev/) is installed on your machine.

The provider can be built and installed directly with:
```shell
go install github.com/clumio-code/terraform-provider-clumio@latest
```

This will download the source code, dependencies, compile the provider binary
and install it under your local `$GOPATH/bin` directory.

To create or update the documentation, use the command `go generate`.

We follow the [Go support policy](https://golang.org/doc/devel/release.html#policy): the two most recent major releases of Go are supported to compile.

Currently, that means Go **1.21** or later must be used when compiling the provider from source.


## Running Tests

In order to run the full suite of unit tests, run `make testunit`.

To execute acceptance tests, it's necessary to set the following environment variables:

- `CLUMIO_API_TOKEN`
- `CLUMIO_API_BASE_URL`
- `CLUMIO_TEST_AWS_ACCOUNT_ID`
- `CLUMIO_TEST_AWS_ACCOUNT_ID2`
- `AWS_REGION`

Please ensure that the `CLUMIO_API_TOKEN` and `CLUMIO_API_BASE_URL` are linked to a functioning
Clumio organization, in line with the configurations detailed in the [Getting Started](#getting-started))
section. In addition, `CLUMIO_TEST_AWS_ACCOUNT_ID`, `CLUMIO_TEST_AWS_ACCOUNT_ID2` should correspond to two different, 
actual AWS accounts and `AWS_REGION` should point to a valid AWS region, 
as real resources will be provisioned during the testing process. Be aware that
conducting these tests could incur costs. 

Since account_native_id cannot be updated for the clumio_aws_connection resource, `CLUMIO_TEST_AWS_ACCOUNT_ID2` is used to test the scenario 
where updating the account_native_id will cause the resource to be destroyed and recreated with the new account_id.

Additionally, before executing the tests, the "data groups" selection for Organizational Units must be finalized through the UI.
This step is required only once.

Furthermore, some acceptance tests necessitate the configuration of Single Sign-On (SSO) prior to execution. Therefore, these tests are not run
by default. To include these tests in the acceptance test suite, SSO must be set up in advance, and the CLUMIO_TEST_IS_SSO_CONFIGURED environment
variable should be set to true. For guidance on configuring SSO, please consult the following documentation: [Authentication and Access](https://support.clumio.com/hc/en-us/sections/13440186425364-Authentication-and-Access)

In order to run the full suite of acceptance tests, run `make testacc`.
