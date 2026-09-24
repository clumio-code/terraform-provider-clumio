// Copyright 2026. Clumio, Inc.

// This file holds the asset types shared by the manual AWS connection resource and its resources
// data source.

package common

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	// Schema attribute names for the assets that can be enabled on a manual AWS connection.
	SchemaIsEbsEnabled               = "ebs"
	SchemaIsRDSEnabled               = "rds"
	SchemaIsDynamoDBEnabled          = "ddb"
	SchemaIsS3Enabled                = "s3"
	SchemaIsMssqlEnabled             = "mssql"
	SchemaIsIcebergOnGlueEnabled     = "iceberg_on_glue"
	SchemaIsIcebergOnS3TablesEnabled = "iceberg_on_s3_tables"

	// Asset type values accepted by the Clumio connection API.
	AssetTypeEBS               = "EBS"
	AssetTypeS3                = "S3"
	AssetTypeRDS               = "RDS"
	AssetTypeDynamoDB          = "DynamoDB"
	AssetTypeEC2MSSQL          = "EC2MSSQL"
	AssetTypeIcebergOnGlue     = "IcebergOnGlue"
	AssetTypeIcebergOnS3Tables = "IcebergOnS3Tables"
)

// AwsAssetsEnabledModel denotes which assets are enabled for a manual AWS connection.
type AwsAssetsEnabledModel struct {
	EBS               types.Bool `tfsdk:"ebs"`
	RDS               types.Bool `tfsdk:"rds"`
	DynamoDB          types.Bool `tfsdk:"ddb"`
	S3                types.Bool `tfsdk:"s3"`
	EC2MSSQL          types.Bool `tfsdk:"mssql"`
	IcebergOnGlue     types.Bool `tfsdk:"iceberg_on_glue"`
	IcebergOnS3Tables types.Bool `tfsdk:"iceberg_on_s3_tables"`
}

// BuildAwsAssetAttributes maps each asset attribute name to one of the two given attributes. The
// resource and the data source schemas take different attribute types, so the caller supplies them.
func BuildAwsAssetAttributes[T any](required T, optional T) map[string]T {
	return map[string]T{
		SchemaIsEbsEnabled:      required,
		SchemaIsDynamoDBEnabled: required,
		SchemaIsRDSEnabled:      required,
		SchemaIsS3Enabled:       required,
		SchemaIsMssqlEnabled:    required,
		// Optional so that configurations written before Iceberg support keep working.
		SchemaIsIcebergOnGlueEnabled:     optional,
		SchemaIsIcebergOnS3TablesEnabled: optional,
	}
}

// ListEnabledAwsAssetTypes returns the API asset type value of every enabled asset.
func ListEnabledAwsAssetTypes(assets *AwsAssetsEnabledModel) []*string {
	assetTypeByFlag := []struct {
		isEnabled types.Bool
		assetType string
	}{
		{assets.EBS, AssetTypeEBS},
		{assets.S3, AssetTypeS3},
		{assets.RDS, AssetTypeRDS},
		{assets.DynamoDB, AssetTypeDynamoDB},
		{assets.EC2MSSQL, AssetTypeEC2MSSQL},
		{assets.IcebergOnGlue, AssetTypeIcebergOnGlue},
		{assets.IcebergOnS3Tables, AssetTypeIcebergOnS3Tables},
	}

	enabled := []*string{}
	for _, entry := range assetTypeByFlag {
		if entry.isEnabled.ValueBool() {
			assetType := entry.assetType
			enabled = append(enabled, &assetType)
		}
	}
	return enabled
}
