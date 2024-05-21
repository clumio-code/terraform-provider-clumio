## 0.7.2
Documentation updates for clumio_policy resource to resolve [#65](https://github.com/clumio-code/terraform-provider-clumio/issues/65).

## 0.7.1
Added documentation for clumio_s3_bucket datasource.

## 0.7.0
This update contains the following changes:
* Added clumio_s3_bucket_properties resource for enabling continuous backup for S3.
* Added clumio_s3_bucket and clumio_dynamo_db_tables data sources.
* Deprecated organizational_unit_id schema attribute in the following resources:
    - clumio_aws_connection
    - clumio_policy
    - clumio_policy_rule
    - clumio_policy_assignment
    - clumio_protection_group
* Deprecated the global timezone attribute in clumio_policy resource.
* Updates to documentation.

## 0.6.2
This update contains the following changes:
* Data sources added for the following:
    - clumio_aws_connection
    - clumio_organizational_unit
    - clumio_policy
    - clumio_policy_rule
    - clumio_protection_group
    - clumio_user
* Additional units tests were added to improve code coverage.
* Improvements to resolve a race condition in acceptance tests.

## 0.6.1
This update covers several improvements to enhance usability and testability. This includes:
* Restructured resources and data sources for enhanced unit test compatibility.
* Introduced unit tests for the provider and its associated resources and data sources.
* Added Sonar badges to surface codebase metrics in README.md.
* Fixes for minor bugs.
* Updates to documentation.

## 0.6.0
Multiple updates were applied across the provider, including, but not limited to:
  - Resolved Issue [#57](https://github.com/clumio-code/terraform-provider-clumio/issues/57) found in the GitHub public repository.
  - Implemented additional safeguards to further prevent provider failures.
  - Management of resources when externally altered, ensuring Terraform will regenerate a resource if it's deleted externally.
  - Broadened the scope of acceptance tests to cover a wider array of standard situations, including updates and recreations.
  - Adjusted the code structure, with such gradual refinements likely to continue for the next couple releases.
  - Revisions and updates to documentation.

## 0.5.9
Upgraded go dependencies to fix dependabot security alerts.

## 0.5.8
Changes to the clumio_policy resource for RDS Compliance Tier.

## 0.5.7
Documentation update to Getting Started guide.

## 0.5.6
Updates to clumio_policy resource and documentation changes.

## 0.5.5
Bug fix release to fix clumio_aws_connection resource.

## 0.5.4
New resource clumio_aws_manual_connection added and additional output attribute added to clumio_aws_connection resource.

## 0.5.3
Updates to clumio_policy resource.

## 0.5.2
Documentation updates and minor changes to the clumio_aws_connection and clumio_organizational_unit resources.

## 0.5.1
Updated clumio_user resource and removed deprecated attributes from clumio_aws_connection resource.

## 0.5.0
Migrated resources to terraform plugin framework. Also added new resources for user auto provisioning.

## 0.4.3
Minor updates to documentation of policy and organizatonal_unit resource.

## 0.4.2
Updates to documentation guides.

## 0.4.1
Minor update to clumio_post_process_kms resource and documentation changes.

## 0.4.0
Changes to the Clumio Wallet and Policy resources

## 0.3.0

Changes to include support for creating Clumio Wallet and related resources.

## 0.2.5

Changes to allow updating organizational_unit_id in clumio_aws_connection resource.

## 0.2.4

Validations and bug fixes added for resources.

## 0.2.3

Added support for bucket_rules in clumio_protection_group resource.

## 0.2.2

Added support to specify organizational_unit_context in the provider and added
clumio_protection_group and clumio_policy_assignment resources.

## 0.2.1

Bug fix release.

## 0.2.0

New resources added for clumio users, organizational-units, policy, policy-rules and
aws_connection.

## 0.1.4

Schema changes to clumio_callback_resource

## 0.1.3

Added support for AWS sso and shared_credentials_file in the provider.

## 0.1.2

Added support for AWS assume_role in the provider.

## 0.1.1

Updated implementation of the clumio_callback_resource.

## 0.1.0

Initial version of terraform-provider-clumio released
