#### [ Pre-requisites for using clumio provider](#pre-requisites)
#### [Example 1: Basic Install](#example1)
#### [Example 2: Cross-Region Install](#example2)
#### [Example 3: Cross-Account Install with Assume Role](#example3)
#### [Example 4: Cross-Account Install with Profiles](#example4)
#### [Example 5: Cross-Account Install with Static Credentials](#example5)
#### [Example 6: Importing Resources](#example6)
#### [Example 7: Instantiate the Clumio Provider in the Context of an OU](#example7)
#### [Example 8: Instantiate an S3 Protection Group and Assign a Policy](#example8)
<a name="pre-requisites"></a>
## Pre-requisites for using Clumio Provider
### Obtain a Clumio API Key:
API Token(personal or service token) can be generated from Settings->Access Management->API Tokens section in ClumioUI.
### Prepare your environment
- These credentials will be used by the AWS Terraform provider to provision resources required by Clumio in the target AWS account (see Example 1 or Example 2 below). If resources must be created in a cross-AWS-account, these credentials will need the ability to assume a role there (see Example 3 below).
- NOTE: For Example 4 and Example 5 credentials are provided in a different manner to the AWS Terraform provider and thus this step should not be required.
```shell
$ export AWS_ACCESS_KEY_ID=<AWS_ACCESS_KEY_ID>
$ export AWS_SECRET_ACCESS_KEY=<AWS_SECRET_ACCESS_KEY>

# If a session token is required ...
$ export AWS_SESSION_TOKEN=<AWS_SESSION_TOKEN>
```
### Instantiate Custom Clumio Resources and the Clumio Module for Protect
- NOTE: In all examples, replace <clumio_api_token> and <clumio_api_base_url> with the appropriate values. The clumio_api_token should be set to the token obtained in the first step. <clumio_api_base_url>  should be one of the following depending on the Clumio portal in-use:
    - https://us-west-2.api.clumio.com
        - portal: https://west.portal.clumio.com/
    - https://us-east-1.api.clumio.com
        - portal: https://east.portal.clumio.com/
    - https://ca-central-1.ca.api.clumio.com
        - portal: https://canada.portal.clumio.com/

### Set up data group if required
- When a new organization is being used, a data group needs to be created from the Clumio UI before clumio_organizational_unit resource could be used.
    - DataGroup can be created from Settings->Organizational Units->Set up data groups

<a name="example1"></a>
## Example 1: Basic Install
The following:
- Downloads the latest Clumio Provider and instantiates it.
- Creates a clumio_aws_connection and installs the Clumio module for Protect onto it.
- Creates several additional, custom Clumio resources including:
    - clumio_policy
    - clumio_policy_rule
    - clumio_organizational_unit
    - clumio_user
```terraform
terraform {
  required_providers {
    clumio = {
      source  = "clumio-code/clumio"
      version = "~>0.2.2"
    }
    aws = {}
  }
}

# Instantiate the Clumio provider
provider "clumio" {
  clumio_api_token    = "<clumio_api_token>"
  clumio_api_base_url = "<clumio_api_base_url>"
}

# Instantiate the AWS provider
provider "aws" {
  region = "us-west-2"
}

# Retrieve the effective AWS account ID and region
data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

# Register a new Clumio connection for the effective AWS account ID and region
resource "clumio_aws_connection" "connection" {
  account_native_id           = data.aws_caller_identity.current.account_id
  aws_region                  = data.aws_region.current.name
  services_enabled            = ["discover", "protect"]
  protect_asset_types_enabled = ["EBS", "RDS", "DynamoDB", "EC2MSSQL", "S3"]
  description                 = "My Clumio Connection"
}

# Install the Clumio Protect template onto the registered connection
module "clumio_protect" {
  providers = {
    clumio = clumio
    aws    = aws
  }
  source                = "clumio-code/aws-template/clumio"
  clumio_token          = clumio_aws_connection.connection.token
  role_external_id      = "my_external_id"
  aws_account_id        = clumio_aws_connection.connection.account_native_id
  aws_region            = clumio_aws_connection.connection.aws_region
  clumio_aws_account_id = clumio_aws_connection.connection.clumio_aws_account_id

  # Enablement of datasources in the module are based on the registered connection
  is_ebs_enabled               = contains(clumio_aws_connection.connection.protect_asset_types_enabled, "EBS")
  is_rds_enabled               = contains(clumio_aws_connection.connection.protect_asset_types_enabled, "RDS")
  is_ec2_mssql_enabled         = contains(clumio_aws_connection.connection.protect_asset_types_enabled, "EC2MSSQL")
  is_warmtier_enabled          = contains(clumio_aws_connection.connection.protect_asset_types_enabled, "DynamoDB")
  is_warmtier_dynamodb_enabled = contains(clumio_aws_connection.connection.protect_asset_types_enabled, "DynamoDB")
  is_s3_enabled                = contains(clumio_aws_connection.connection.protect_asset_types_enabled, "S3")
}

# Create a Clumio policy with a 7-day RPO and 14-day retention
resource "clumio_policy" "policy" {
  name = "Gold"
  operations {
    action_setting = "immediate"
    type           = "aws_ebs_volume_backup"
    slas {
      retention_duration {
        unit  = "days"
        value = 14
      }
      rpo_frequency {
        unit  = "days"
        value = 7
      }
    }
  }
}

# Create a Clumio policy rule and associate it with the policy
resource "clumio_policy_rule" "rule" {
  name           = "Tag-Based Rule"
  policy_id      = clumio_policy.policy.id
  condition      = "{\"entity_type\":{\"$eq\":\"aws_ebs_volume\"}, \"aws_tag\":{\"$eq\":{\"key\":\"random-test-123\", \"value\":\"random-test-123\"}}}"
  before_rule_id = ""
}

# Retrive the role for OU Admin
data "clumio_role" "ou_admin" {
  name = "Organizational Unit Admin"
}

# Create a new OU
resource "clumio_organizational_unit" "ou" {
  name = "My OU"
}

# Create a user for the OU
resource "clumio_user" "user" {
  full_name               = "Foo Bar"
  email                   = "foobar@clumio.com"
  assigned_role           = data.clumio_role.ou_admin.id
  organizational_unit_ids = [clumio_organizational_unit.ou.id]
}
```
<a name="example2"></a>
## Example 2: Cross-Region Install
The following:
- Downloads the latest Clumio Provider and instantiates it.
- Instantiates two AWS Providers where the attributes for each denote two different AWS regions.
- Creates two instances of clumio_aws_connection and installs the Clumio module for Protect onto both.
```terraform
terraform {
  required_providers {
    clumio = {
      source  = "clumio-code/clumio"
      version = "~> 0.2.2"
    }
    aws = {}
  }
}

# Instantiate the Clumio provider
provider "clumio" {
  clumio_api_token    = "<clumio_api_token>"
  clumio_api_base_url = "<clumio_api_base_url>"
}

# Instantiate two AWS providers with different regions
provider "aws" {
  region = "us-west-2"
}
provider "aws" {
  alias  = "east"
  region = "us-east-1"
}

# Retrieve the effective AWS account ID
data "aws_caller_identity" "current" {}

# Register a new Clumio connection on us-west-2 for the effective AWS account ID
resource "clumio_aws_connection" "connection_west" {
  account_native_id           = data.aws_caller_identity.current.account_id
  aws_region                  = "us-west-2"
  services_enabled            = ["discover", "protect"]
  protect_asset_types_enabled = ["EBS", "RDS", "DynamoDB", "EC2MSSQL", "S3"]
  description                 = "My Clumio Connection West"
}

# Register a new Clumio connection on us-east-1 for the effective AWS account ID
resource "clumio_aws_connection" "connection_east" {
  account_native_id           = data.aws_caller_identity.current.account_id
  aws_region                  = "us-east-1"
  services_enabled            = ["discover", "protect"]
  protect_asset_types_enabled = ["EBS", "RDS", "DynamoDB", "EC2MSSQL", "S3"]
  description                 = "My Clumio Connection East"
}

# Install the Clumio Protect template onto the registered connection for West
module "clumio_protect_west" {
  providers = {
    clumio = clumio
    aws    = aws
  }
  source                = "clumio-code/aws-template/clumio"
  clumio_token          = clumio_aws_connection.connection_west.token
  role_external_id      = "my_external_id_west"
  aws_account_id        = clumio_aws_connection.connection_west.account_native_id
  aws_region            = clumio_aws_connection.connection_west.aws_region
  clumio_aws_account_id = clumio_aws_connection.connection_west.clumio_aws_account_id

  # Enablement of datasources in the module are based on the registered connection
  is_ebs_enabled               = contains(clumio_aws_connection.connection_west.protect_asset_types_enabled, "EBS")
  is_rds_enabled               = contains(clumio_aws_connection.connection_west.protect_asset_types_enabled, "RDS")
  is_ec2_mssql_enabled         = contains(clumio_aws_connection.connection_west.protect_asset_types_enabled, "EC2MSSQL")
  is_warmtier_enabled          = contains(clumio_aws_connection.connection_west.protect_asset_types_enabled, "DynamoDB")
  is_warmtier_dynamodb_enabled = contains(clumio_aws_connection.connection_west.protect_asset_types_enabled, "DynamoDB")
  is_s3_enabled                = contains(clumio_aws_connection.connection_west.protect_asset_types_enabled, "S3")
}

# Install the Clumio Protect template onto the registered connection for East
module "clumio_protect_east" {
  providers = {
    clumio = clumio
    aws    = aws.east
  }
  source                = "clumio-code/aws-template/clumio"
  clumio_token          = clumio_aws_connection.connection_east.token
  role_external_id      = "my_external_id_east"
  aws_account_id        = clumio_aws_connection.connection_east.account_native_id
  aws_region            = clumio_aws_connection.connection_east.aws_region
  clumio_aws_account_id = clumio_aws_connection.connection_east.clumio_aws_account_id

  # Enablement of datasources in the module are based on the registered connection
  is_ebs_enabled               = contains(clumio_aws_connection.connection_east.protect_asset_types_enabled, "EBS")
  is_rds_enabled               = contains(clumio_aws_connection.connection_east.protect_asset_types_enabled, "RDS")
  is_ec2_mssql_enabled         = contains(clumio_aws_connection.connection_east.protect_asset_types_enabled, "EC2MSSQL")
  is_warmtier_enabled          = contains(clumio_aws_connection.connection_east.protect_asset_types_enabled, "DynamoDB")
  is_warmtier_dynamodb_enabled = contains(clumio_aws_connection.connection_east.protect_asset_types_enabled, "DynamoDB")
  is_s3_enabled                = contains(clumio_aws_connection.connection_east.protect_asset_types_enabled, "S3")
}
```
<a name="example3"></a>
## Example 3: Cross-Account Install with Assume Role
The following:
- Downloads the latest Clumio Provider and instantiates it.
- Instantiates two AWS Providers where the attributes for each denote the cross-account role to assume on two different AWS accounts.
- Creates two instances of clumio_aws_connection and installs the Clumio module for Protect onto both.
- NOTE: Replace <assume_role_arn_1>, <assume_role_external_id_1>, <assume_role_arn_2>, and <assume_role_external_id_2> with the appropriate values to assume a cross-account role on the two different AWS accounts.
```terraform
terraform {
  required_providers {
    clumio = {
      source  = "clumio-code/clumio"
      version = "~> 0.2.2"
    }
    aws = {}
  }
}

# Instantiate the Clumio provider
provider "clumio" {
  clumio_api_token    = "<clumio_api_token>"
  clumio_api_base_url = "<clumio_api_base_url>"
}

# Instantiate two AWS providers with different regions and roles to assume in
# order to access different AWS accounts
provider "aws" {
  region = "us-west-2"
  assume_role {
    role_arn    = "<assume_role_arn_1>"
    external_id = "<assume_role_external_id_1>"
  }
}
provider "aws" {
  alias  = "account_2_east"
  region = "us-east-1"
  assume_role {
    role_arn    = "<assume_role_arn_2>"
    external_id = "<assume_role_external_id_2>"
  }
}

# Retrieve the effective AWS account IDs for the different AWS accounts
data "aws_caller_identity" "account_1" {}
data "aws_caller_identity" "account_2" {
  provider = aws.account_2_east
}

# Register a new Clumio connection on us-west-2 for the first AWS account ID
resource "clumio_aws_connection" "connection_account_1_west" {
  account_native_id           = data.aws_caller_identity.account_1.account_id
  aws_region                  = "us-west-2"
  services_enabled            = ["discover", "protect"]
  protect_asset_types_enabled = ["EBS", "RDS", "DynamoDB", "EC2MSSQL", "S3"]
  description                 = "My Clumio Connection Account 1 West"
}

# Register a new Clumio connection on us-east-1 for the second AWS account ID
resource "clumio_aws_connection" "connection_account_2_east" {
  account_native_id           = data.aws_caller_identity.account_2.account_id
  aws_region                  = "us-east-1"
  services_enabled            = ["discover", "protect"]
  protect_asset_types_enabled = ["EBS", "RDS", "DynamoDB", "EC2MSSQL", "S3"]
  description                 = "My Clumio Connection Account 2 East"
}

# Install the Clumio Protect template onto the registered connection for the
# first AWS account ID on West
module "clumio_protect_account_1_west" {
  providers = {
    clumio = clumio
    aws    = aws
  }
  source                = "clumio-code/aws-template/clumio"
  clumio_token          = clumio_aws_connection.connection_account_1_west.token
  role_external_id      = "my_external_id_account_1_west"
  aws_account_id        = clumio_aws_connection.connection_account_1_west.account_native_id
  aws_region            = clumio_aws_connection.connection_account_1_west.aws_region
  clumio_aws_account_id = clumio_aws_connection.connection_account_1_west.clumio_aws_account_id

  # Enablement of datasources in the module are based on the registered connection
  is_ebs_enabled               = contains(clumio_aws_connection.connection_account_1_west.protect_asset_types_enabled, "EBS")
  is_rds_enabled               = contains(clumio_aws_connection.connection_account_1_west.protect_asset_types_enabled, "RDS")
  is_ec2_mssql_enabled         = contains(clumio_aws_connection.connection_account_1_west.protect_asset_types_enabled, "EC2MSSQL")
  is_warmtier_enabled          = contains(clumio_aws_connection.connection_account_1_west.protect_asset_types_enabled, "DynamoDB")
  is_warmtier_dynamodb_enabled = contains(clumio_aws_connection.connection_account_1_west.protect_asset_types_enabled, "DynamoDB")
  is_s3_enabled                = contains(clumio_aws_connection.connection_account_1_west.protect_asset_types_enabled, "S3")
}

# Install the Clumio Protect template onto the registered connection for the
# second AWS account ID on East
module "clumio_protect_account_2_east" {
  providers = {
    clumio = clumio
    aws    = aws.account_2_east
  }
  source                = "clumio-code/aws-template/clumio"
  clumio_token          = clumio_aws_connection.connection_account_2_east.token
  role_external_id      = "my_external_id_account_2_east"
  aws_account_id        = clumio_aws_connection.connection_account_2_east.account_native_id
  aws_region            = clumio_aws_connection.connection_account_2_east.aws_region
  clumio_aws_account_id = clumio_aws_connection.connection_account_2_east.clumio_aws_account_id

  # Enablement of datasources in the module are based on the registered connection
  is_ebs_enabled               = contains(clumio_aws_connection.connection_account_2_east.protect_asset_types_enabled, "EBS")
  is_rds_enabled               = contains(clumio_aws_connection.connection_account_2_east.protect_asset_types_enabled, "RDS")
  is_ec2_mssql_enabled         = contains(clumio_aws_connection.connection_account_2_east.protect_asset_types_enabled, "EC2MSSQL")
  is_warmtier_enabled          = contains(clumio_aws_connection.connection_account_2_east.protect_asset_types_enabled, "DynamoDB")
  is_warmtier_dynamodb_enabled = contains(clumio_aws_connection.connection_account_2_east.protect_asset_types_enabled, "DynamoDB")
  is_s3_enabled                = contains(clumio_aws_connection.connection_account_2_east.protect_asset_types_enabled, "S3")
}
```
<a name="example4"></a>
## Example 4: Cross-Account Install with Profiles
- https://registry.terraform.io/providers/hashicorp/aws/latest/docs#shared-credentials-file
- Similar to Example 3, but appropriate values for the profile attribute must be given (it is assumed that default values for shared_config_files and shared_credentials_files are okay). For example:
```terraform
provider "aws" {
  region  = "us-west-2"
  profile = "custom-profile"
}
```
<a name="example5"></a>
## Example 5: Cross-Account Install with Static Credentials
- NOTE: Hard-coded credentials are not recommended in any Terraform configuration.
- https://registry.terraform.io/providers/hashicorp/aws/latest/docs#static-credentials.
- Similar to Example 3, but appropriate values for the access_key and secret_key attributes must be given. For example:
```terraform
provider "aws" {
  region     = "us-west-2"
  access_key = "my-access-key"
  secret_key = "my-secret-key"

  # If a session token is required ...
  token = "my-token"
}
```
<a name="example6"></a>
## Example 6: Importing Resources
- Specify an empty custom resource in the Terraform file for the Clumio resource to be imported:
```terraform
terraform {
  required_providers {
    clumio = {
      source  = "clumio.com/providers/clumio"
      version = "~>0.2.2"
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
```
- Then run the terraform import command. For example:
```shell
$ terraform import clumio_policy.policy <clumio_policy_id>
```
- where <clumio_policy_id> is replaced with the appropriate Clumio policy ID to import
- Once complete, the Terraform file should be populated with the details of the imported resource. This can be obtained using the terraform state show command:
```shell
$ terraform state show clumio_policy.policy

# clumio_policy.policy:
resource "clumio_policy" "policy" {
    activation_status      = "activated"
    id                     = "a44a4d91-c50c-4a61-bd4d-c46cf1c4427a"
    lock_status            = "unlocked"
    name                   = "demo"
    organizational_unit_id = "00000000-0000-0000-0000-000000000000"

    operations {
        action_setting = "immediate"
        type           = "aws_ebs_volume_backup"

        slas {
            retention_duration {
                unit  = "days"
                value = 14
            }

            rpo_frequency {
                unit  = "days"
                value = 1
            }
        }
        slas {
            retention_duration {
                unit  = "days"
                value = 7
            }

            rpo_frequency {
                unit  = "on_demand"
                value = 0
            }
        }
    }
}
```
<a name="example7"></a>
## Example7: Instantiate the Clumio Provider in the Context of an OU
```terraform
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
```
where <clumio_ou_id> is replaced with the ID of an existing organizational unit.
<a name="example8"></a>
## Instantiate an S3 Protection Group and Assign a Policy
```terraform
# Create a Clumio policy with a 7-day RPO and 6-month retention
terraform {
  required_providers {
    clumio = {
      source  = "clumio.com/providers/clumio"
      version = "~>0.2.2"
    }
  }
}

resource "clumio_policy" "policy" {
  name = "Gold"
  operations {
    action_setting = "immediate"
    type           = "protection_group_backup"
    advanced_settings {
      protection_group_backup {
        backup_tier = "cold"
      }
    }
    slas {
      retention_duration {
        unit  = "months"
        value = 6
      }
      rpo_frequency {
        unit  = "days"
        value = 7
      }
    }
  }
}

# Create a Clumio protection group for S3
resource "clumio_protection_group" "protection_group" {
  name = "S3 Protection Group"
  object_filter {
    storage_classes = ["S3 Standard", "S3 Standard-IA"]
  }
}

# Assign a policy to the protection group
resource "clumio_policy_assignment" "protection_group_policy_assignment" {
  entity_type = "protection_group"
  entity_id   = clumio_protection_group.protection_group.id
  policy_id   = clumio_policy.policy.id
}
```