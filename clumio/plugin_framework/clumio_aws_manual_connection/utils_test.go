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
func TestIsAssetConfigDowngraded(t *testing.T) {

	state := &clumioAWSManualConnectionResourceModel{
		AssetsEnabled: &AssetsEnabledModel{
			EBS:      basetypes.NewBoolValue(true),
			RDS:      basetypes.NewBoolValue(true),
			DynamoDB: basetypes.NewBoolValue(true),
			S3:       basetypes.NewBoolValue(true),
			EC2MSSQL: basetypes.NewBoolValue(true),
		},
	}

	plan := &clumioAWSManualConnectionResourceModel{
		AssetsEnabled: &AssetsEnabledModel{
			EBS:      basetypes.NewBoolValue(true),
			RDS:      basetypes.NewBoolValue(true),
			DynamoDB: basetypes.NewBoolValue(true),
			S3:       basetypes.NewBoolValue(true),
			EC2MSSQL: basetypes.NewBoolValue(true),
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
}
