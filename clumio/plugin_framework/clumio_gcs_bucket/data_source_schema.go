// Copyright 2026. Clumio, Inc.

// This file holds the type definition and Schema data source function used by the data source model
// for the clumio_gcs_bucket Terraform data source.

package clumio_gcs_bucket

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// clumioGCSBucketDataSourceModel is the data source model for the clumio_gcs_bucket Terraform
// data source. It represents the schema of the data source and the data it holds. The bucket_names
// attribute is used to query the Clumio GCS buckets, whereas the gcs_buckets attribute is computed
// and holds the buckets returned by the Clumio API.
type clumioGCSBucketDataSourceModel struct {
	BucketNames types.Set `tfsdk:"bucket_names"`
	GCSBuckets  types.Set `tfsdk:"gcs_buckets"`
}

// Schema defines the structure and constraints of the clumio_gcs_bucket Terraform data source.
// Schema is a method on the clumioGCSBucketDataSource struct. It sets the schema for the
// clumio_gcs_bucket Terraform data source, which retrieves the Clumio GCS buckets matching the
// given 'bucket_names'.
func (r *clumioGCSBucketDataSource) Schema(
	_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			schemaBucketNames: schema.SetAttribute{
				Description: "The list of GCS bucket names to be queried.",
				Required:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			schemaGCSBuckets: schema.SetNestedAttribute{
				Description: "The GCS buckets matching the given bucket_names.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						schemaId: schema.StringAttribute{
							Description: "The Clumio-assigned ID of the GCS bucket.",
							Computed:    true,
						},
						schemaBucketName: schema.StringAttribute{
							Description: "The name of the GCS bucket.",
							Computed:    true,
						},
						schemaProjectId: schema.StringAttribute{
							Description: "The GCP project ID associated with the bucket.",
							Computed:    true,
						},
						schemaProjectUuid: schema.StringAttribute{
							Description: "The Clumio-assigned UUID of the GCP project associated" +
								" with the bucket.",
							Computed: true,
						},
						schemaLocation: schema.StringAttribute{
							Description: "The GCP location associated with the bucket.",
							Computed:    true,
						},
						schemaLocationType: schema.StringAttribute{
							Description: "The location type of the bucket (e.g., \"Region\"," +
								" \"Dual-region\", \"Multi-region\").",
							Computed: true,
						},
						schemaLocationUuid: schema.StringAttribute{
							Description: "The Clumio-assigned UUID of the GCP location associated" +
								" with the bucket.",
							Computed: true,
						},
						schemaOrganizationalUnitId: schema.StringAttribute{
							Description: "The Clumio-assigned ID of the organizational unit" +
								" associated with the bucket.",
							Computed: true,
						},
						schemaObjectCount: schema.Int64Attribute{
							Description: "The number of objects in the bucket.",
							Computed:    true,
						},
						schemaSizeBytes: schema.Int64Attribute{
							Description: "Total size in bytes of all objects in the bucket.",
							Computed:    true,
						},
						schemaProtectionGroupCount: schema.Int64Attribute{
							Description: "The number of protection groups associated with the" +
								" bucket.",
							Computed: true,
						},
						schemaIsDeleted: schema.BoolAttribute{
							Description: "Determines whether the bucket has been deleted.",
							Computed:    true,
						},
						schemaIsVersioningEnabled: schema.BoolAttribute{
							Description: "Determines whether versioning is enabled for the bucket.",
							Computed:    true,
						},
						schemaCreatedTimestamp: schema.StringAttribute{
							Description: "Creation time of the bucket in RFC-3339 format.",
							Computed:    true,
						},
						schemaLastBackupTimestamp: schema.StringAttribute{
							Description: "Time of the last backup in RFC-3339 format.",
							Computed:    true,
						},
						schemaLabels: schema.SetNestedAttribute{
							Description: "The GCP labels associated with the bucket.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									schemaKey: schema.StringAttribute{
										Description: "The GCP label key.",
										Computed:    true,
									},
									schemaValue: schema.StringAttribute{
										Description: "The GCP label value.",
										Computed:    true,
									},
								},
							},
						},
					},
				},
				Computed: true,
			},
		},
		Description: "clumio_gcs_bucket data source is used to retrieve details of the GCS" +
			" buckets for use in other resources.",
	}
}
