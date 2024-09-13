# Using Multiple Clumio Providers

In this guide, we will explore how to use two Clumio providers in Terraform, distinguished by the `clumio_organizational_unit_context` variable.

## Prerequisites

Before we begin, make sure you have the following:

- Terraform installed on your machine
- Clumio API credentials for both Clumio providers

## Step 1: Set up Provider Configurations

First, we need to define the provider configurations in our Terraform code. Open your Terraform configuration file and add the following code:

```hcl
provider "clumio" {
    alias = "provider1"
    clumio_organizational_unit_context = "org_unit_1"
    clumio_api_token = "<provider1_api_token>"
    clumio_api_base_url = "<clumio_api_base_url>"
}

provider "clumio" {
    alias = "provider2"
    clumio_organizational_unit_context = "org_unit_2"
    clumio_api_token = "<provider2_api_token>"
    clumio_api_base_url = "<clumio_api_base_url>"
}
```

Replace `<provider1_api_token>`, `<provider2_api_token>` and `<clumio_api_base_url>` with the respective API tokens and provider base URL for each Clumio provider.

## Step 2: Define Resources

Next, we can define resources that will be managed by each Clumio provider. For example, let's create two Clumio policies:

```hcl
resource "clumio_policy" "policy1" {
    provider = clumio.provider1
    # Specify policy configuration for provider 1
}

resource "clumio_policy" "policy2" {
    provider = clumio.provider2
    # Specify policy configuration for provider 2
}
```

Make sure to customize the policy configurations according to your requirements.

## Step 3: Apply Changes

Finally, we can apply the changes to create the Clumio policies using both providers. Run the following command in your terminal:

```bash
terraform init
terraform apply
```

Terraform will detect the provider configurations and create the policies accordingly.

## Conclusion

By following this guide, you have learned how to use multiple Clumio providers in Terraform, distinguished by the `clumio_organizational_unit_context` variable. This allows you to manage resources across different Clumio organizational units efficiently.
