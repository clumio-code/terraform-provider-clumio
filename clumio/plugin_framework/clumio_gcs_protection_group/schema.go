// Copyright 2026. Clumio, Inc.

// This file holds the type definition and Schema resource function used by the resource model for
// the clumio_gcs_protection_group Terraform resource.

package clumio_gcs_protection_group

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// clumioGCSProtectionGroupResourceModel is the resource model for the clumio_gcs_protection_group
// Terraform resource.
type clumioGCSProtectionGroupResourceModel struct {
	ID                       types.String     `tfsdk:"id"`
	Name                     types.String     `tfsdk:"name"`
	BucketRule               *bucketRuleModel `tfsdk:"bucket_rule"`
	IncludePrefixes          types.Set        `tfsdk:"include_prefixes"`
	ExcludePrefixes          types.Set        `tfsdk:"exclude_prefixes"`
	LatestVersionOnly        types.Bool       `tfsdk:"latest_version_only"`
	OrganizationalUnitID     types.String     `tfsdk:"organizational_unit_id"`
	ProtectionStatus         types.String     `tfsdk:"protection_status"`
	ProtectionInfo           types.List       `tfsdk:"protection_info"`
	BucketCount              types.Int64      `tfsdk:"bucket_count"`
	BucketUuids              types.Set        `tfsdk:"bucket_uuids"`
	BucketRuleMatchedCount   types.Int64      `tfsdk:"bucket_rule_matched_bucket_count"`
	BucketRuleMatchedUuids   types.Set        `tfsdk:"bucket_rule_matched_bucket_uuids"`
	ManualAddedBucketCount   types.Int64      `tfsdk:"manual_added_bucket_count"`
	Labels                   types.Set        `tfsdk:"labels"`
	Location                 types.String     `tfsdk:"location"`
	LocationType             types.String     `tfsdk:"location_type"`
	LastBackupTimestamp      types.String     `tfsdk:"last_backup_timestamp"`
	TotalBackedUpObjectCount types.Int64      `tfsdk:"total_backed_up_object_count"`
	TotalBackedUpSizeBytes   types.Int64      `tfsdk:"total_backed_up_size_bytes"`
	BackupStatusStats        types.Object     `tfsdk:"backup_status_stats"`
	CreatedTimestamp         types.String     `tfsdk:"created_timestamp"`
	ModifiedTimestamp        types.String     `tfsdk:"modified_timestamp"`
	Version                  types.Int64      `tfsdk:"version"`
}

// bucketRuleModel maps to the bucket_rule block. Each condition field is optional and itself a
// single nested block.
type bucketRuleModel struct {
	GcpLabel     *gcpLabelOperatorModel  `tfsdk:"gcp_label"`
	GcpLocation  *gcpStringOperatorModel `tfsdk:"gcp_location"`
	GcpProjectId *gcpStringOperatorModel `tfsdk:"gcp_project_id"`
}

// gcpStringOperatorModel maps to a string-valued condition field (gcp_location, gcp_project_id).
type gcpStringOperatorModel struct {
	Eq    types.String `tfsdk:"eq"`
	In    types.Set    `tfsdk:"in"`
	NotEq types.String `tfsdk:"not_eq"`
	NotIn types.Set    `tfsdk:"not_in"`
}

// gcpLabelOperatorModel maps to the gcp_label condition field. The single-label and
// match-all operators are expressed as maps of key => value (GCP label keys are unique), while
// in/not_in are repeated blocks because they may list multiple values for the same key.
type gcpLabelOperatorModel struct {
	Eq          types.Map     `tfsdk:"eq"`
	Contains    types.Map     `tfsdk:"contains"`
	All         types.Map     `tfsdk:"all"`
	NotEq       types.Map     `tfsdk:"not_eq"`
	NotContains types.Map     `tfsdk:"not_contains"`
	NotAll      types.Map     `tfsdk:"not_all"`
	In          []*labelModel `tfsdk:"in"`
	NotIn       []*labelModel `tfsdk:"not_in"`
}

// labelModel maps to a single {key, value} GCP label.
type labelModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

// labelBlockAttrs returns the attributes of a single GCP label block.
func labelBlockAttrs() map[string]schema.Attribute {
	// key and value are Optional rather than Required: a Required attribute under an optional
	// SingleNestedBlock is validated even when the parent block is absent (a Plugin Framework
	// limitation). The API requires both when a label operator is set.
	return map[string]schema.Attribute{
		schemaKey: schema.StringAttribute{
			Optional: true, Description: "The GCP label key (required when the operator is set)."},
		schemaValue: schema.StringAttribute{
			Optional: true, Description: "The GCP label value (required when the operator is set)."},
	}
}

// stringOperatorBlock returns a single nested block for a string-valued condition field.
func stringOperatorBlock(desc string) schema.SingleNestedBlock {
	return schema.SingleNestedBlock{
		Description: desc,
		Attributes: map[string]schema.Attribute{
			schemaEq: schema.StringAttribute{
				Optional: true, Description: "Match values equal to this value."},
			schemaIn: schema.SetAttribute{
				Optional: true, ElementType: types.StringType,
				Description: "Match values in this set."},
			schemaNotEq: schema.StringAttribute{
				Optional: true, Description: "Match values not equal to this value."},
			schemaNotIn: schema.SetAttribute{
				Optional: true, ElementType: types.StringType,
				Description: "Match values not in this set."},
		},
	}
}

// bucketRuleNotEmptyValidator ensures that, when the bucket_rule block is present, it specifies at
// least one condition field (gcp_label, gcp_location, or gcp_project_id). An empty bucket_rule is
// meaningless and is rejected at plan time rather than sending an empty rule to the API.
type bucketRuleNotEmptyValidator struct{}

func (bucketRuleNotEmptyValidator) Description(_ context.Context) string {
	return "bucket_rule must specify at least one of gcp_label, gcp_location, or gcp_project_id."
}

func (v bucketRuleNotEmptyValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (bucketRuleNotEmptyValidator) ValidateObject(
	_ context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {

	// An absent (or not-yet-known) bucket_rule is allowed; the block is optional.
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	attrs := req.ConfigValue.Attributes()
	for _, name := range []string{schemaGcpLabel, schemaGcpLocation, schemaGcpProjectId} {
		if a, ok := attrs[name]; ok && !a.IsNull() {
			return
		}
	}
	resp.Diagnostics.AddAttributeError(req.Path, "Empty bucket_rule",
		"bucket_rule must specify at least one condition: gcp_label, gcp_location, or"+
			" gcp_project_id.")
}

// Schema defines the structure and constraints of the clumio_gcs_protection_group resource.
func (r *clumioGCSProtectionGroupResource) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {

	// labelMapSingle is for single-label operators (eq/contains/not_eq/not_contains): a map of
	// key => value restricted to one entry.
	labelMapSingle := func(desc string) schema.MapAttribute {
		return schema.MapAttribute{
			Description: desc + " Provide a single label as {key = value}.",
			ElementType: types.StringType,
			Optional:    true,
			Validators:  []validator.Map{mapvalidator.SizeAtMost(1)},
		}
	}
	// labelMapMulti is for match-all operators (all/not_all): a map of key => value with one or
	// more entries, all of which must match.
	labelMapMulti := func(desc string) schema.MapAttribute {
		return schema.MapAttribute{
			Description: desc + " Provide labels as {key = value, ...}.",
			ElementType: types.StringType,
			Optional:    true,
			Validators:  []validator.Map{mapvalidator.SizeAtLeast(1)},
		}
	}
	// labelSet is for in/not_in: repeated {key, value} blocks, allowing the same key to appear
	// with multiple values.
	labelSet := func(desc string) schema.SetNestedBlock {
		return schema.SetNestedBlock{
			Description:  desc,
			NestedObject: schema.NestedBlockObject{Attributes: labelBlockAttrs()},
		}
	}

	resp.Schema = schema.Schema{
		Description: "Clumio GCP Protection Group Resource used to create and manage GCS" +
			" protection groups.",
		Attributes: map[string]schema.Attribute{
			schemaId: schema.StringAttribute{
				Description: "Unique identifier for the Clumio GCP protection group.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			schemaName: schema.StringAttribute{
				Description: "The user-assigned name of the protection group.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			schemaIncludePrefixes: schema.SetAttribute{
				Description: "A list of prefixes to include in the backup.",
				ElementType: types.StringType,
				Optional:    true,
				Validators:  []validator.Set{setvalidator.SizeAtLeast(1)},
			},
			schemaExcludePrefixes: schema.SetAttribute{
				Description: "A list of prefixes to exclude from the backup.",
				ElementType: types.StringType,
				Optional:    true,
				Validators:  []validator.Set{setvalidator.SizeAtLeast(1)},
			},
			schemaLatestVersionOnly: schema.BoolAttribute{
				Description: "Whether to back up only the latest object version. Defaults to true.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			schemaOrganizationalUnitId: schema.StringAttribute{
				Description: "The Clumio-assigned ID of the organizational unit associated with" +
					" the protection group.",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			schemaProtectionStatus: schema.StringAttribute{
				Description: "The protection status of the protection group.",
				Computed:    true,
			},
			schemaProtectionInfo: schema.ListNestedAttribute{
				Description: "The protection policy applied to this resource.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						schemaPolicyId: schema.StringAttribute{
							Description: "Identifier of the policy applied to the protection group.",
							Computed:    true,
						},
					},
				},
			},
			schemaBucketCount: schema.Int64Attribute{
				Description: "The total number of unique buckets in this protection group.",
				Computed:    true,
			},
			schemaBucketUuids: schema.SetAttribute{
				Description: "The set of bucket UUIDs directly assigned to this protection group.",
				ElementType: types.StringType,
				Computed:    true,
			},
			schemaBucketRuleMatchedCount: schema.Int64Attribute{
				Description: "The number of buckets matched by the bucket rule.",
				Computed:    true,
			},
			schemaBucketRuleMatchedUuids: schema.SetAttribute{
				Description: "The set of bucket UUIDs matched by the bucket rule.",
				ElementType: types.StringType,
				Computed:    true,
			},
			schemaManualAddedBucketCount: schema.Int64Attribute{
				Description: "The number of buckets manually assigned to this protection group.",
				Computed:    true,
			},
			schemaLabels: schema.SetNestedAttribute{
				Description: "The GCP labels associated with the protection group.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						schemaKey: schema.StringAttribute{
							Description: "The GCP label key.", Computed: true},
						schemaValue: schema.StringAttribute{
							Description: "The GCP label value.", Computed: true},
					},
				},
			},
			schemaLocation: schema.StringAttribute{
				Description: "The location of the protection group.",
				Computed:    true,
			},
			schemaLocationType: schema.StringAttribute{
				Description: "The location type of the protection group.",
				Computed:    true,
			},
			schemaLastBackupTimestamp: schema.StringAttribute{
				Description: "Time of the last backup in RFC-3339 format.",
				Computed:    true,
			},
			schemaTotalBackedUpObjectCount: schema.Int64Attribute{
				Description: "The total number of objects backed up in this protection group.",
				Computed:    true,
			},
			schemaTotalBackedUpSizeBytes: schema.Int64Attribute{
				Description: "The total size in bytes backed up in this protection group.",
				Computed:    true,
			},
			schemaBackupStatusStats: schema.SingleNestedAttribute{
				Description: "Aggregated backup status statistics for the protection group.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					schemaFailureCount: schema.Int64Attribute{
						Description: "Number of entities with a backup status of failure.",
						Computed:    true},
					schemaNoBackupCount: schema.Int64Attribute{
						Description: "Number of entities with a backup status of no_backup.",
						Computed:    true},
					schemaPartialSuccessCount: schema.Int64Attribute{
						Description: "Number of entities with a backup status of partial_success.",
						Computed:    true},
					schemaSuccessCount: schema.Int64Attribute{
						Description: "Number of entities with a backup status of success.",
						Computed:    true},
				},
			},
			schemaCreatedTimestamp: schema.StringAttribute{
				Description: "Creation time of the protection group in RFC-3339 format.",
				Computed:    true,
			},
			schemaModifiedTimestamp: schema.StringAttribute{
				Description: "Modified time of the protection group in RFC-3339 format.",
				Computed:    true,
			},
			schemaVersion: schema.Int64Attribute{
				Description: "Version of the protection group. Incremented on every change.",
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
			schemaBucketRule: schema.SingleNestedBlock{
				Description: "Rule that automatically adds matching GCS buckets to the protection" +
					" group. Within each condition field at most one include operator" +
					" (eq/contains/in/all) and one exclude operator" +
					" (not_eq/not_contains/not_in/not_all) may be set; the API rejects invalid" +
					" combinations.",
				Blocks: map[string]schema.Block{
					schemaGcpLabel: schema.SingleNestedBlock{
						Description: "Label-based bucket conditions.",
						Attributes: map[string]schema.Attribute{
							schemaEq:          labelMapSingle("Match buckets with this exact label."),
							schemaContains:    labelMapSingle("Match buckets containing this label."),
							schemaAll:         labelMapMulti("Match buckets with all of these labels."),
							schemaNotEq:       labelMapSingle("Exclude buckets with this exact label."),
							schemaNotContains: labelMapSingle("Exclude buckets containing this label."),
							schemaNotAll:      labelMapMulti("Exclude buckets with all of these labels."),
						},
						Blocks: map[string]schema.Block{
							schemaIn: labelSet("Match buckets with any of these labels. Repeat the" +
								" block to allow multiple values for the same key."),
							schemaNotIn: labelSet("Exclude buckets with any of these labels. Repeat" +
								" the block to exclude multiple values for the same key."),
						},
					},
					schemaGcpLocation:  stringOperatorBlock("GCP location-based bucket conditions."),
					schemaGcpProjectId: stringOperatorBlock("GCP project-id-based bucket conditions."),
				},
				Validators: []validator.Object{bucketRuleNotEmptyValidator{}},
			},
		},
	}
}
