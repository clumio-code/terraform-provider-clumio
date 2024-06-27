// Copyright 2024. Clumio, Inc.

// This file holds the type definition and Schema resource function used by the resource model for
// the clumio_protection_group_bucket Terraform resource.

package clumio_protection_group_bucket

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// clumioProtectionGroupBucketResourceModel is the resource model for the
// clumio_protection_group_bucket Terraform resource. It represents the schema of the resource and
// the data it holds. This schema is used by customers to configure the resource and by the Clumio
// provider to read and write the resource.
type clumioProtectionGroupBucketResourceModel struct {
	ID                types.String `tfsdk:"id"`
	BucketID          types.String `tfsdk:"bucket_id"`
	ProtectionGroupID types.String `tfsdk:"protection_group_id"`
}

// Schema defines the structure and constraints of the clumio_protection_group_bucket Terraform
// resource. Schema is a method on the clumioProtectionGroupBucketResource struct. It sets the
// schema for the clumio_protection_group_bucket Terraform resource, which is used to manage bucket
// assignment to Protection Group. The schema defines various attributes such as the ID, bucket ID
// and Protection Group ID. ID attribute is computed, meaning it is determined by Clumio at runtime,
// while others are required or optional inputs from the user.
func (r *clumioProtectionGroupBucketResource) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Clumio S3 Protection Group Bucket Resource used to assign a bucket to a " +
			"Protection Group.",
		Attributes: map[string]schema.Attribute{
			schemaId: schema.StringAttribute{
				Description: "Unique identifier for the Clumio Protection Group bucket association.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			schemaBucketId: schema.StringAttribute{
				Description: "Clumio assigned unique identifier of the AWS S3 bucket.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			schemaProtectionGroupId: schema.StringAttribute{
				Description: "Unique identifier of the Protection Group.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
	}
}
