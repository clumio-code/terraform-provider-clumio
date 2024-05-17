// Copyright 2024. Clumio, Inc.

// This file holds the type definition and Schema resource function used by the resource model for
// the clumio_s3_bucket_properties Terraform resource.

package clumio_s3_bucket_properties

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// clumioS3BucketResourceModel is the resource model for the clumio_s3_bucket Terraform resource. It
// represents the schema of the resource and the data it holds. This schema is used by customers to
// configure the resource and by the Clumio provider to read and write the resource.
type clumioS3BucketPropertiesResourceModel struct {
	ID                              types.String `tfsdk:"id"`
	BucketID                        types.String `tfsdk:"bucket_id"`
	EventBridgeEnabled              types.Bool   `tfsdk:"event_bridge_enabled"`
	EvendBridgeNotificationDisabled types.Bool   `tfsdk:"event_bridge_notification_disabled"`
}

// Schema defines the structure and constraints of the clumio_s3_bucket_properties Terraform
// resource. Schema is a method on the clumioS3BucketPropertiesResource struct. It sets the schema
// for the clumio_s3_bucket_properties Terraform resource, which is used to set the Clumio S3
// bucket properties. The schema defines various attributes such as the bucket ID,
// event_bridge_enabled, etc.
func (r *clumioS3BucketPropertiesResource) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{
		Description: "Clumio S3 Bucket Properties Resource to set and read S3 bucket properties" +
			" in Clumio.",
		Attributes: map[string]schema.Attribute{
			schemaId: schema.StringAttribute{
				Description: "Unique identifier for the S3 bucket properties.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			schemaBucketId: schema.StringAttribute{
				Description: "Unique Clumio assigned identifier for the S3 bucket.",
				Required:    true,
			},
			schemaEventBridgeEnabled: schema.BoolAttribute{
				Description: "Determines if continuous backup is enabled for the S3 bucket.",
				Required:    true,
			},
			schemaEventBridgeNotificationDisabled: schema.BoolAttribute{
				Description: "If true, tries to disable EventBridge notification for the given " +
					"bucket. This may override the existing bucket notification configuration in " +
					"the AWS account. Defaults to true if not specified.",
				Optional: true,
			},
		},
	}
}
