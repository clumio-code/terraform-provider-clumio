// Copyright 2024. Clumio, Inc.

// This file holds the type definition and Schema resource function used by the resource model for
// the clumio_protection_group Terraform resource.

package clumio_protection_group

import (
	"context"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// clumioProtectionGroupResourceModel is the resource model for the clumio_protection_group
// Terraform resource. It represents the schema of the resource and the data it holds. This schema
// is used by customers to configure the resource and by the Clumio provider to read and write the
// resource.
type clumioProtectionGroupResourceModel struct {
	ID               types.String         `tfsdk:"id"`
	Name             types.String         `tfsdk:"name"`
	Description      types.String         `tfsdk:"description"`
	BucketRule       types.String         `tfsdk:"bucket_rule"`
	ObjectFilter     []*objectFilterModel `tfsdk:"object_filter"`
	ProtectionStatus types.String         `tfsdk:"protection_status"`
	ProtectionInfo   types.List           `tfsdk:"protection_info"`
}

// objectFilterModel maps to the 'object_filter' field in clumioProtectionGroupResourceModel and
// refers to the list of object filters in a protection group
type objectFilterModel struct {
	LatestVersionOnly             types.Bool           `tfsdk:"latest_version_only"`
	PrefixFilters                 []*prefixFilterModel `tfsdk:"prefix_filters"`
	StorageClasses                []types.String       `tfsdk:"storage_classes"`
	EarliestLastModifiedTimestamp types.String         `tfsdk:"earliest_last_modified_timestamp"`
}

// prefixFilterModel maps to 'prefix_filters' field in objectFilterModel and refers to list of
// prefix filters inside an object filter
type prefixFilterModel struct {
	ExcludedSubPrefixes []types.String `tfsdk:"excluded_sub_prefixes"`
	Prefix              types.String   `tfsdk:"prefix"`
}

// Schema defines the structure and constraints of the clumio_protection_group Terraform resource.
// Schema is a method on the clumioProtectionGroupResource struct. It sets the schema for the
// clumio_protection_group Terraform resource, which is used to manage protection groups within
// Clumio. The schema defines various attributes such as the protection group ID, name, description,
// etc. Some of these attributes are computed, meaning they are determined by Clumio at runtime,
// while others are required or optional inputs from the user.
func (r *clumioProtectionGroupResource) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	prefixFilterSchemaAttributes := map[string]schema.Attribute{
		schemaExcludedSubPrefixes: schema.SetAttribute{
			Description: "List of subprefixes to exclude from the prefix.",
			ElementType: types.StringType,
			Optional:    true,
			Validators:  []validator.Set{setvalidator.SizeAtLeast(1)},
		},
		schemaPrefix: schema.StringAttribute{
			Required: true,
			Description: "Prefix to include. To include all objects in the bucket specify empty " +
				"string \"\".",
		},
	}

	objectFilterSchemaAttributes := map[string]schema.Attribute{
		schemaLatestVersionOnly: schema.BoolAttribute{
			Optional:    true,
			Computed:    true,
			Description: "Whether to back up only the latest object version.",
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		},
		schemaStorageClasses: schema.SetAttribute{
			Description: "Storage class to include in the backup. Valid values are: S3 Standard," +
				" S3 Standard-IA, S3 Intelligent-Tiering, and S3 One Zone-IA.",
			ElementType: types.StringType,
			Required:    true,
		},
		schemaEarliestLastModifiedTimestamp: schema.StringAttribute{
			Description: "The cutoff date for inclusion objects from the backup. Any object with" +
				" a last modified date after or equal than this value will be included in the " +
				"backup. This is useful for filtering out old or irrelevant objects based on " +
				"their modification timestamps. This supports RFC-3339 format.",
			Optional: true,
			Computed: true,
		},
	}

	objectFilterSchemaBlocks := map[string]schema.Block{
		schemaPrefixFilters: schema.SetNestedBlock{
			Description: "Prefix Filters.",
			NestedObject: schema.NestedBlockObject{
				Attributes: prefixFilterSchemaAttributes,
			},
		},
	}

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Clumio S3 Protection Group Resource used to create and manage Protection Groups.",
		Attributes: map[string]schema.Attribute{
			schemaId: schema.StringAttribute{
				Description: "Unique identifier for the Clumio Protection Group.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			schemaDescription: schema.StringAttribute{
				Description: "Brief description to denote details of the protection group.",
				Optional:    true,
			},
			schemaName: schema.StringAttribute{
				Description: "The user-assigned name of the protection group. Must be globally-unique.",
				Required:    true,
			},
			schemaBucketRule: schema.StringAttribute{
				Description: "The following table describes the possible conditions for a bucket" +
					" to be automatically added to a protection group. <br><table><tr><th>Field" +
					"</th><th>Rule Condition</th><th>Description</th></tr><tr><td>aws_tag</td>" +
					"<td>$eq, $not_eq, $contains, $not_contains, $all, $not_all, $in, $not_in" +
					"</td><td>Denotes the AWS tag(s) to conditionalize on<code>{\"aws_tag\":" +
					"{\"$eq\":{\"key\":\"Environment\", \"value\":\"Prod\"}}}</code></td></tr>" +
					"<tr><td>aws_account_native_id</td><td>$eq, $in</td><td>Denotes the AWS " +
					"account to conditionalize on<code>{\"aws_account_native_id\":{\"$eq\":" +
					"\"111111111111\"}}</code></td></tr><tr><td>account_native_id<br><b>" +
					"Deprecated</b></td><td>$eq, $in</td><td>This will be deprecated and use " +
					"aws_account_native_id instead.<br>Denotes the AWS account to conditionalize" +
					" on<code>{\"account_native_id\":{\"$in\":[\"111111111111\"]}}</code></td>" +
					"</tr><tr><td>aws_region</td><td>$eq, $in</td><td>Denotes the AWS region to " +
					"conditionalize on<code>{\"aws_region\":{\"$eq\":\"us-west-2\"}}</code></td>" +
					"</tr></table>",
				Optional: true,
			},
			schemaProtectionStatus: schema.StringAttribute{
				Description: "The protection status of the protection group. Possible values include" +
					"\"protected\", \"unprotected\", and \"unsupported\". If the protection group does not" +
					"support backups, then this field has a value of unsupported.",
				Computed: true,
			},
			schemaProtectionInfo: schema.ListNestedAttribute{
				Description: "The protection policy applied to this resource.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						schemaInheritingEntityId: schema.StringAttribute{
							Description: "The identifier of the entity from which protection was inherited.",
							Computed:    true,
						},
						schemaInheritingEntityType: schema.StringAttribute{
							Description: "The type of the entity from which protection was inherited.",
							Computed:    true,
						},
						schemaPolicyId: schema.StringAttribute{
							Description: "Identifier of the policy to apply on Protection Group",
							Computed:    true,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			schemaObjectFilter: schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: objectFilterSchemaAttributes,
					Blocks:     objectFilterSchemaBlocks,
				},
				Validators: []validator.Set{
					common.WrapSetValidator(setvalidator.SizeAtMost(1)),
				},
			},
		},
	}
}
