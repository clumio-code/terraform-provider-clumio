// Copyright 2024. Clumio, Inc.

// This file holds the type definition and Schema datasource function used by the datasource model
// for the clumio_s3_bucket Terraform datasource.

package clumio_s3_bucket

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// clumioS3BucketDataSourceModel is the datasource model for the clumio_s3_bucket
// Terraform datasource. It represents the schema of the datasource and the data it holds. This
// schema is used by customers to configure the datasource and by the Clumio provider to read and
// write the datasource.
type clumioS3BucketDataSourceModel struct {
	BucketNames types.Set `tfsdk:"bucket_names"`
	S3Buckets   types.Set `tfsdk:"s3_buckets"`
}

// Schema defines the structure and constraints of the clumio_s3_bucket Terraform datasource. Schema
// is a method on the clumioS3BucketDataSource struct. It sets the schema for the clumio_s3_bucket
// Terraform datasource. The schema defines various attributes such as the name and s3_buckets where
// 's3_buckets' is computed, meaning is determined by Clumio at runtime, whereas the 'names'
// attribute us used to determine the Clumio s3 buckets to retrieve.
func (r *clumioS3BucketDataSource) Schema(
	_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			schemaBucketNames: schema.SetAttribute{
				Description: "The list of s3 bucket names to be queried.",
				Required:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			schemaS3Buckets: schema.SetNestedAttribute{
				Description: "S3Buckets which match the given name.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						schemaId: schema.StringAttribute{
							Description: "Unique identifier of the s3 bucket.",
							Computed:    true,
						},
						schemaName: schema.StringAttribute{
							Description: "The name of the s3 bucket.",
							Computed:    true,
						},
						schemaAccountNativeId: schema.StringAttribute{
							Description: "The identifier of the AWS account under which the S3 " +
								"bucket was created.",
							Computed: true,
						},
						schemaRegion: schema.StringAttribute{
							Description: "The AWS region associated with the S3 bucket.",
							Computed:    true,
						},
						schemaProtectionGroupCount: schema.Int64Attribute{
							Description: "Protection group count reflects how many protection groups" +
								" are linked to this bucket.",
							Computed: true,
						},
						schemaEventBridgeEnabled: schema.BoolAttribute{
							Description: "Determines if continuous backup is enabled for the S3" +
								" bucket.",
							Computed: true,
						},
						schemaLastBackupTimestamp: schema.StringAttribute{
							Description: "Time of the last backup in RFC-3339 format.",
							Computed:    true,
						},
						schemaLastContinuousBackupTimestamp: schema.StringAttribute{
							Description: "Time of the last continuous backup in RFC-3339 format.",
							Computed:    true,
						},
					},
				},
				Computed: true,
			},
		},
		Description: "clumio_s3_bucket data source is used to retrieve details of an" +
			" s3 bucket for use in other resources.",
	}
}
