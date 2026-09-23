// Copyright 2024. Clumio, Inc.

// This file contains the unit tests for the functions in resource.go

//go:build unit

package clumio_aws_manual_connection

import (
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Unit test for the following cases:
//   - No change should return false
//   - Removing EBS asset should return true
//   - Removing RDS asset should return true
//   - Removing DynamoDB asset should return true
//   - Removing S3 asset should return true
//   - Removing EC2MSSQL asset should return true
//   - Removing Iceberg on Glue asset should return true
//   - Removing Iceberg on S3 Tables asset should return true
//   - Leaving both Iceberg assets unset should return false
func TestIsAssetConfigDowngraded(t *testing.T) {

	state := &clumioAWSManualConnectionResourceModel{
		AssetsEnabled: &AssetsEnabledModel{
			EBS:      basetypes.NewBoolValue(true),
			RDS:      basetypes.NewBoolValue(true),
			DynamoDB: basetypes.NewBoolValue(true),
			S3:       basetypes.NewBoolValue(true),
			EC2MSSQL: basetypes.NewBoolValue(true),

			IcebergOnGlue:     basetypes.NewBoolValue(true),
			IcebergOnS3Tables: basetypes.NewBoolValue(true),
		},
	}

	plan := &clumioAWSManualConnectionResourceModel{
		AssetsEnabled: &AssetsEnabledModel{
			EBS:      basetypes.NewBoolValue(true),
			RDS:      basetypes.NewBoolValue(true),
			DynamoDB: basetypes.NewBoolValue(true),
			S3:       basetypes.NewBoolValue(true),
			EC2MSSQL: basetypes.NewBoolValue(true),

			IcebergOnGlue:     basetypes.NewBoolValue(true),
			IcebergOnS3Tables: basetypes.NewBoolValue(true),
		},
	}

	// If there is no change return false.
	t.Run("Returns true if no change", func(t *testing.T) {

		downgrade := isAssetConfigDowngraded(plan, state)
		assert.False(t, downgrade)
	})

	// If plan has EBS as disabled and state has it as enabled, the return true.
	t.Run("Returns true if EBS is false", func(t *testing.T) {

		plan.AssetsEnabled.EBS = basetypes.NewBoolValue(false)
		downgrade := isAssetConfigDowngraded(plan, state)
		assert.True(t, downgrade)
		plan.AssetsEnabled.EBS = basetypes.NewBoolValue(true)
	})

	// If plan has RDS as disabled and state has it as enabled, the return true.
	t.Run("Returns true if RDS is false", func(t *testing.T) {

		plan.AssetsEnabled.RDS = basetypes.NewBoolValue(false)
		downgrade := isAssetConfigDowngraded(plan, state)
		assert.True(t, downgrade)
		plan.AssetsEnabled.RDS = basetypes.NewBoolValue(true)
	})

	// If plan has DynamoDB as disabled and state has it as enabled, the return true.
	t.Run("Returns true if DynamoDB is false", func(t *testing.T) {

		plan.AssetsEnabled.DynamoDB = basetypes.NewBoolValue(false)
		downgrade := isAssetConfigDowngraded(plan, state)
		assert.True(t, downgrade)
		plan.AssetsEnabled.DynamoDB = basetypes.NewBoolValue(true)
	})

	// If plan has S3 as disabled and state has it as enabled, the return true.
	t.Run("Returns true if S3 is false", func(t *testing.T) {

		plan.AssetsEnabled.S3 = basetypes.NewBoolValue(false)
		downgrade := isAssetConfigDowngraded(plan, state)
		assert.True(t, downgrade)
		plan.AssetsEnabled.S3 = basetypes.NewBoolValue(true)
	})

	// If plan has EC2MSSQL as disabled and state has it as enabled, the return true.
	t.Run("Returns true if EC2MSSQL is false", func(t *testing.T) {

		plan.AssetsEnabled.EC2MSSQL = basetypes.NewBoolValue(false)
		downgrade := isAssetConfigDowngraded(plan, state)
		assert.True(t, downgrade)
		plan.AssetsEnabled.EC2MSSQL = basetypes.NewBoolValue(true)
	})

	// If plan has Iceberg on Glue as disabled and state has it as enabled, then return true.
	t.Run("Returns true if Iceberg on Glue is false", func(t *testing.T) {

		plan.AssetsEnabled.IcebergOnGlue = basetypes.NewBoolValue(false)
		downgrade := isAssetConfigDowngraded(plan, state)
		assert.True(t, downgrade)
		plan.AssetsEnabled.IcebergOnGlue = basetypes.NewBoolValue(true)
	})

	// If plan has Iceberg on S3 Tables as disabled and state has it as enabled, then return true.
	t.Run("Returns true if Iceberg on S3 Tables is false", func(t *testing.T) {

		plan.AssetsEnabled.IcebergOnS3Tables = basetypes.NewBoolValue(false)
		downgrade := isAssetConfigDowngraded(plan, state)
		assert.True(t, downgrade)
		plan.AssetsEnabled.IcebergOnS3Tables = basetypes.NewBoolValue(true)
	})

	// A configuration written before Iceberg support leaves both Iceberg attributes null, and that
	// is not a downgrade.
	t.Run("Returns false if both Iceberg assets are null", func(t *testing.T) {

		nullIceberg := &clumioAWSManualConnectionResourceModel{
			AssetsEnabled: &AssetsEnabledModel{
				EBS:      basetypes.NewBoolValue(true),
				RDS:      basetypes.NewBoolValue(true),
				DynamoDB: basetypes.NewBoolValue(true),
				S3:       basetypes.NewBoolValue(true),
				EC2MSSQL: basetypes.NewBoolValue(true),

				IcebergOnGlue:     basetypes.NewBoolNull(),
				IcebergOnS3Tables: basetypes.NewBoolNull(),
			},
		}

		downgrade := isAssetConfigDowngraded(nullIceberg, nullIceberg)
		assert.False(t, downgrade)
	})
}
