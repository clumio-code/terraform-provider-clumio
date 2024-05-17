// Copyright 2024. Clumio, Inc.

// This file holds the type definition and Schema resource function used by the resource model for
// the clumio_protection_group Terraform resource.

package clumio_protection_group

import (
	"context"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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
	ID                   types.String         `tfsdk:"id"`
	Name                 types.String         `tfsdk:"name"`
	Description          types.String         `tfsdk:"description"`
	BucketRule           types.String         `tfsdk:"bucket_rule"`
	ObjectFilter         []*objectFilterModel `tfsdk:"object_filter"`
	ProtectionStatus     types.String         `tfsdk:"protection_status"`
	ProtectionInfo       types.List           `tfsdk:"protection_info"`
	OrganizationalUnitID types.String         `tfsdk:"organizational_unit_id"`
}

// objectFilterModel maps to the 'object_filter' field in clumioProtectionGroupResourceModel and
// refers to the list of object filters in a protection group
type objectFilterModel struct {
	LatestVersionOnly types.Bool           `tfsdk:"latest_version_only"`
	PrefixFilters     []*prefixFilterModel `tfsdk:"prefix_filters"`
	StorageClasses    []types.String       `tfsdk:"storage_classes"`
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
		},
		schemaPrefix: schema.StringAttribute{
			Optional:    true,
			Description: "Prefix to include.",
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
				Description: "Describes the possible conditions for a bucket to be  automatically added" +
					" to a protection group. For example: " +
					"{\"aws_tag\":{\"$eq\":{\"key\":\"Environment\", \"value\":\"Prod\"}}}",
				Optional: true,
			},
			schemaOrganizationalUnitId: schema.StringAttribute{
				Description: "Identifier of the Clumio organizational unit associated with the " +
					"protection group. If not provided, the protection group will be associated " +
					"with the default organizational unit associated with the credentials used " +
					"to create the protection group.",
				Optional: true,
				Computed: true,
				DeprecationMessage: "Use the provider schema attribute " +
					"clumio_organizational_unit_context to create the resource in the context of " +
					"an Organizational Unit.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
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
