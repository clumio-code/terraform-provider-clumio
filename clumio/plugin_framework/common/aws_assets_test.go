// Copyright 2026. Clumio, Inc.

// This file contains the unit tests for the functions in aws_assets.go.

//go:build unit

package common

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
)

// Unit test for the following cases:
//   - All assets enabled returns every asset type
//   - No assets enabled returns an empty list
//   - Null Iceberg attributes, as written before Iceberg support, return only the other assets
//   - Only the enabled assets are returned
func TestListEnabledAwsAssetTypes(t *testing.T) {

	enabled := basetypes.NewBoolValue(true)
	disabled := basetypes.NewBoolValue(false)
	unset := basetypes.NewBoolNull()

	tests := []struct {
		name   string
		assets *AwsAssetsEnabledModel
		want   []string
	}{
		{
			name: "all assets enabled",
			assets: &AwsAssetsEnabledModel{
				EBS: enabled, S3: enabled, RDS: enabled, DynamoDB: enabled,
				EC2MSSQL: enabled, IcebergOnGlue: enabled, IcebergOnS3Tables: enabled,
			},
			want: []string{AssetTypeEBS, AssetTypeS3, AssetTypeRDS, AssetTypeDynamoDB,
				AssetTypeEC2MSSQL, AssetTypeIcebergOnGlue, AssetTypeIcebergOnS3Tables},
		},
		{
			name: "no assets enabled",
			assets: &AwsAssetsEnabledModel{
				EBS: disabled, S3: disabled, RDS: disabled, DynamoDB: disabled,
				EC2MSSQL: disabled, IcebergOnGlue: disabled, IcebergOnS3Tables: disabled,
			},
			want: []string{},
		},
		{
			name: "Iceberg attributes left unset",
			assets: &AwsAssetsEnabledModel{
				EBS: enabled, S3: enabled, RDS: enabled, DynamoDB: enabled,
				EC2MSSQL: disabled, IcebergOnGlue: unset, IcebergOnS3Tables: unset,
			},
			want: []string{AssetTypeEBS, AssetTypeS3, AssetTypeRDS, AssetTypeDynamoDB},
		},
		{
			name: "only Iceberg on S3 Tables enabled",
			assets: &AwsAssetsEnabledModel{
				EBS: disabled, S3: disabled, RDS: disabled, DynamoDB: disabled,
				EC2MSSQL: disabled, IcebergOnGlue: disabled, IcebergOnS3Tables: enabled,
			},
			want: []string{AssetTypeIcebergOnS3Tables},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			got := []string{}
			for _, assetType := range ListEnabledAwsAssetTypes(test.assets) {
				got = append(got, *assetType)
			}
			assert.Equal(t, test.want, got, "assets: %+v", test.assets)
		})
	}
}

// Unit test for the following cases:
//   - Every asset attribute name is mapped
//   - The Iceberg attributes take the optional value and the rest take the required value
func TestBuildAwsAssetAttributes(t *testing.T) {

	attributes := BuildAwsAssetAttributes("required", "optional")

	t.Run("Maps every asset attribute name", func(t *testing.T) {

		assert.Len(t, attributes, 7)
	})

	t.Run("Only the Iceberg attributes are optional", func(t *testing.T) {

		assert.Equal(t, map[string]string{
			SchemaIsEbsEnabled:               "required",
			SchemaIsDynamoDBEnabled:          "required",
			SchemaIsRDSEnabled:               "required",
			SchemaIsS3Enabled:                "required",
			SchemaIsMssqlEnabled:             "required",
			SchemaIsIcebergOnGlueEnabled:     "optional",
			SchemaIsIcebergOnS3TablesEnabled: "optional",
		}, attributes)
	})
}
