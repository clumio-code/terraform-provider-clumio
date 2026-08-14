# Example retrieving the Clumio AWS environment of an AWS account and region
data "clumio_aws_environment" "example_aws_environment" {
  account_native_id = "AWS Account ID"
  aws_region        = "AWS Region"
}
