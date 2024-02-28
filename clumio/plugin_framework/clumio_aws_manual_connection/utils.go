// Copyright 2024. Clumio, Inc.

// This file hold various utility functions used by the clumio_aws_manual_connection Terraform
// resource.

package clumio_aws_manual_connection

// isAssetConfigDowngraded checks if any previously added assets in the config are no being removed,
// as asset downgrades are not supoorted during update.
func isAssetConfigDowngraded(
	plan *clumioAWSManualConnectionResourceModel, state *clumioAWSManualConnectionResourceModel) bool {
	// If EBS was removed now
	if !plan.AssetsEnabled.EBS.ValueBool() && state.AssetsEnabled.EBS.ValueBool() {
		return true
	}
	// If S3 was removed now
	if !plan.AssetsEnabled.S3.ValueBool() && state.AssetsEnabled.S3.ValueBool() {
		return true
	}
	// If RDS was removed now
	if !plan.AssetsEnabled.RDS.ValueBool() && state.AssetsEnabled.RDS.ValueBool() {
		return true
	}
	// If DynamoDB was removed now
	if !plan.AssetsEnabled.DynamoDB.ValueBool() && state.AssetsEnabled.DynamoDB.ValueBool() {
		return true
	}
	// If EC2MSSQL was removed now
	if !plan.AssetsEnabled.EC2MSSQL.ValueBool() && state.AssetsEnabled.EC2MSSQL.ValueBool() {
		return true
	}
	return false
}

