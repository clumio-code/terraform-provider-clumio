// Copyright 2025. Clumio, Inc.

// This file holds the type definition and Schema resource function used by the resource model for
// the clumio_report_configuration Terraform resource.

package clumio_report_configuration

import (
	"context"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// reportConfigurationResourceModel is the resource model for the clumio_report_configuration
// Terraform resource. It represents the schema of the resource and the data it holds. This schema
// is used by customers to configure the resource and by the Clumio provider to read and write the
// resource.
type reportConfigurationResourceModel struct {
	ID           types.String         `tfsdk:"id"`
	Description  types.String         `tfsdk:"description"`
	Name         types.String         `tfsdk:"name"`
	CreatedAt    types.String         `tfsdk:"created"`
	Notification []*notificationModel `tfsdk:"notification"`
	Parameter    []*parameterModel    `tfsdk:"parameter"`
	Schedule     []*scheduleModel     `tfsdk:"schedule"`
}

// notificationModel maps to the 'notification' field in reportConfigurationResourceModel and
// refers to the list of notification targets in a report configuration.
type notificationModel struct {
	EmailList []types.String `tfsdk:"email_list"`
}

// parameterModel maps to the 'parameter' field in reportConfigurationResourceModel and refers to
// the list of parameters in a report configuration.
type parameterModel struct {
	Controls []*controlsModel `tfsdk:"controls"`
	Filters  []*filtersModel  `tfsdk:"filters"`
}

// controlsModel maps to the 'controls' field in parameterModel and refers to the set of controls
// supported in a report configuration.
type controlsModel struct {
	AssetBackupControl     []*assetBackupControl     `tfsdk:"asset_backup"`
	AssetProtectionControl []*assetProtectionControl `tfsdk:"asset_protection"`
	PolicyControl          []*policyControl          `tfsdk:"policy"`
}

// timeUnitModel defines the time unit for controls in a report configuration.
type timeUnitModel struct {
	Unit  types.String `tfsdk:"unit"`
	Value types.Int32  `tfsdk:"value"`
}

// assetBackupControl maps to the 'asset_backup' field in controlsModel and refers to the asset
// backup control in a report configuration.
type assetBackupControl struct {
	LookBackPeriod           []*timeUnitModel `tfsdk:"look_back_period"`
	MinimumRetentionDuration []*timeUnitModel `tfsdk:"minimum_retention_duration"`
	WindowSize               []*timeUnitModel `tfsdk:"window_size"`
}

// assetProtectionControl maps to the 'asset_protection' field in controlsModel and refers to the
// asset protection control in a report configuration.
type assetProtectionControl struct {
	ShouldIgnoreDeactivatedPolicy types.Bool `tfsdk:"should_ignore_deactivated_policy"`
}

// policyControl maps to the 'policy' field in controlsModel and refers to the policy control
// in a report configuration.
type policyControl struct {
	MinimumRetentionDuration []*timeUnitModel `tfsdk:"minimum_retention_duration"`
	MinimumRpoFrequency      []*timeUnitModel `tfsdk:"minimum_rpo_frequency"`
}

// filtersModel maps to the 'filters' field in parameterModel and refers to the set of filters
// supported in a report configuration.
type filtersModel struct {
	AssetFilter  []*assetFilter  `tfsdk:"asset"`
	CommonFilter []*commonFilter `tfsdk:"common"`
}

// assetFilter maps to the 'asset' field in filtersModel and refers to the asset filter
// in a report configuration.
type assetFilter struct {
	GroupsFilter []*assetGroupFilter `tfsdk:"groups"`
	TagOpMode    types.String        `tfsdk:"tag_op_mode"`
	TagsFilter   []*assetTagFilter   `tfsdk:"tags"`
}

// assetGroupFilter maps to the 'groups' field in assetFilter and refers to the group to be
// filtered in a report configuration.
type assetGroupFilter struct {
	ID     types.String `tfsdk:"id"`
	Region types.String `tfsdk:"region"`
	Type   types.String `tfsdk:"type"`
}

// assetTagFilter maps to the 'tags' field in assetFilter and refers to the tag to be filtered
// in a report configuration.
type assetTagFilter struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

// commonFilter maps to the 'common' field in filtersModel and refers to the common filter
// in a report configuration.
type commonFilter struct {
	AssetTypes          []types.String `tfsdk:"asset_types"`
	DataSources         []types.String `tfsdk:"data_sources"`
	OrganizationalUnits []types.String `tfsdk:"organizational_units"`
}

// scheduleModel maps to the 'schedule' field in reportConfigurationResourceModel and refers to
// the detail schedules of a report configuration.
type scheduleModel struct {
	DayOfMonth types.Int64  `tfsdk:"day_of_month"`
	DayOfWeek  types.String `tfsdk:"day_of_week"`
	Frequency  types.String `tfsdk:"frequency"`
	StartTime  types.String `tfsdk:"start_time"`
	Timezone   types.String `tfsdk:"timezone"`
}

// ifTagExistValidator is a custom validator for the 'tag_op_mode' field in
// assetFilter. It checks if the 'tag_op_mode' is set, then it requires that the 'tags' field
// is also set. If 'tag_op_mode' is set but 'tags' is not, it adds an error to the response.
type ifTagExistValidator struct{}

func (v ifTagExistValidator) Description(ctx context.Context) string {
	return "tag_op_mode must be set only if tags are provided."
}

func (v ifTagExistValidator) MarkdownDescription(ctx context.Context) string {
	return "tag_op_mode must be set only if tags are provided."
}

// Validate runs the main validation logic of the validator, reading configuration data out of `req` and updating `resp` with diagnostics.
func (v ifTagExistValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// If the value is unknown or null, there is nothing to validate.
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	tags := types.SetNull(types.StringType)
	diags := req.Config.GetAttribute(ctx, req.Path.ParentPath().AtName(schemaTagFilter), &tags)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if tags.IsNull() {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"`tag_op_mode` is only allowed when `tags` are provided.",
			"You must provide one or more `tags` to use `tag_op_mode`, otherwise remove `tag_op_mode`.",
		)
	}
}

func ifTagExist() validator.String {
	return ifTagExistValidator{}
}

// Schema defines the structure and constraints of the clumio_report_configuration Terraform
// resource. Schema is a method on the clumioReportConfigurationResource struct. It sets the schema
// for the clumio_report_configuration Terraform resource, which is used to configure a report.
func (r *clumioReportConfigurationResource) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {

	TimeUnitSchemaBlock := schema.SetNestedBlock{
		Description: "The time unit used in control definition.",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				schemaUnit: schema.StringAttribute{
					Description: "Enum: `minutes` `hours` `days` `weeks` `months` `years`<br>Unit" +
						" indicates the unit for time unit param.",
					Required: true,
				},
				schemaValue: schema.Int32Attribute{
					Description: "Value indicates the value for time unit param.",
					Required:    true,
				},
			},
		},
		Validators: []validator.Set{
			common.WrapSetValidator(setvalidator.SizeAtMost(1)),
		},
	}

	AssetBackupControlSchemaBlock := schema.SetNestedBlock{
		Description: "The control for asset backup.",
		NestedObject: schema.NestedBlockObject{
			Blocks: map[string]schema.Block{
				schemaLookBackPeriod:           TimeUnitSchemaBlock,
				schemaMinimumRetentionDuration: TimeUnitSchemaBlock,
				schemaWindowSize:               TimeUnitSchemaBlock,
			},
		},
		Validators: []validator.Set{
			common.WrapSetValidator(setvalidator.SizeAtMost(1)),
		},
	}

	AssetProtectionControlSchemaBlock := schema.SetNestedBlock{
		Description: "The control for asset protection.",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				schemaShouldIgnoreDeactivatedPolicy: schema.BoolAttribute{
					Description: "Whether the report should ignore deactivated policy or not.",
					Required:    true,
				},
			},
		},
		Validators: []validator.Set{
			common.WrapSetValidator(setvalidator.SizeAtMost(1)),
		},
	}

	PolicyControlSchemaBlock := schema.SetNestedBlock{
		Description: "The control for policy.",
		NestedObject: schema.NestedBlockObject{
			Blocks: map[string]schema.Block{
				schemaMinimumRetentionDuration: TimeUnitSchemaBlock,
				schemaMinimumRpoFrequency:      TimeUnitSchemaBlock,
			},
		},
		Validators: []validator.Set{
			common.WrapSetValidator(setvalidator.SizeAtMost(1)),
		},
	}

	ControlsSchemaBlock := schema.SetNestedBlock{
		Description: "The set of controls supported in compliance report.",
		NestedObject: schema.NestedBlockObject{
			Blocks: map[string]schema.Block{
				schemaAssetBackupControl:     AssetBackupControlSchemaBlock,
				schemaAssetProtectionControl: AssetProtectionControlSchemaBlock,
				schemaPolicyControl:          PolicyControlSchemaBlock,
			},
		},
		Validators: []validator.Set{
			common.WrapSetValidator(setvalidator.IsRequired()),
			common.WrapSetValidator(setvalidator.SizeAtMost(1)),
		},
	}

	AssetGroupFilterSchemaBlock := schema.SetNestedBlock{
		Description: "The asset groups to be filtered.",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				schemaAssetGroupID: schema.StringAttribute{
					Description: "The id of asset group.",
					Optional:    true,
				},
				schemaAssetGroupRegion: schema.StringAttribute{
					Description: "The region of asset group. For example, `us-west-2`. This is " +
						"supported for AWS asset groups only.",
					Optional: true,
				},
				schemaAssetGroupType: schema.StringAttribute{
					Description: "Enum: `aws` `microsoft365` `vmware`<br>The type of asset group.",
					Required:    true,
				},
			},
		},
	}

	AssetTagFilterSchemaBlock := schema.SetNestedBlock{
		Description: "The asset tags to be filtered. This is supported for AWS assets only.",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				schemaTagKey: schema.StringAttribute{
					Description: "The key of tag to filter.",
					Required:    true,
				},
				schemaTagValue: schema.StringAttribute{
					Description: "The value of tag to filter.",
					Required:    true,
				},
			},
		},
	}

	AssetFilterSchemaBlock := schema.SetNestedBlock{
		Description: "The filter for asset. This will be applied to asset backup and asset " +
			"protection controls.",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				schemaTagOpMode: schema.StringAttribute{
					Description: "Enum: `equal` `or` `and`<br>The tag filter operation to be " +
						"applied to the given tags. This is supported for AWS assets only.",
					Optional:   true,
					Validators: []validator.String{ifTagExist()},
				},
			},
			Blocks: map[string]schema.Block{
				schemaAssetGroupFilter: AssetGroupFilterSchemaBlock,
				schemaTagFilter:        AssetTagFilterSchemaBlock,
			},
		},
		Validators: []validator.Set{
			common.WrapSetValidator(setvalidator.SizeAtMost(1)),
		},
	}

	CommonFilterSchemaBlock := schema.SetNestedBlock{
		Description: "The common filter which will be applied to all controls.",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				schemaAssetTypes: schema.ListAttribute{
					Description: "The asset types to be included in the report. For example, " +
						"[`aws_ec2_instance`, `microsoft365_drive`].",
					Optional:    true,
					ElementType: types.StringType,
				},
				schemaDataSources: schema.ListAttribute{
					Description: "The data sources to be included in the report. Possible values " +
						"include `aws`, `microsoft365` or `vmware`.",
					Optional:    true,
					ElementType: types.StringType,
				},
				schemaOrganizationalUnits: schema.ListAttribute{
					Description: "The organizational units to be included in the report.",
					Optional:    true,
					ElementType: types.StringType,
				},
			},
		},
		Validators: []validator.Set{
			common.WrapSetValidator(setvalidator.SizeAtMost(1)),
		},
	}

	FiltersSchemaBlock := schema.SetNestedBlock{
		Description: "The set of filters supported in compliance report.",
		NestedObject: schema.NestedBlockObject{
			Blocks: map[string]schema.Block{
				schemaAssetFilter:  AssetFilterSchemaBlock,
				schemaCommonFilter: CommonFilterSchemaBlock,
			},
		},
		Validators: []validator.Set{
			common.WrapSetValidator(setvalidator.SizeAtMost(1)),
		},
	}

	ScheduleSchemaAttributes := map[string]schema.Attribute{
		schemaDayOfMonth: schema.Int64Attribute{
			Description: "The day of the month when the report will be sent out. This is required" +
				" for the 'monthly' report frequency. It has to be >= 1 and <= 28, or '-1', which" +
				" signifies end of month. If the day_of_month is set to -1 then the report will " +
				"be sent out at the end of every month.",
			Optional: true,
		},
		schemaDayOfWeek: schema.StringAttribute{
			Description: "Enum: `sunday` `monday` `tuesday` `wednesday` `thursday` `friday` " +
				"`saturday`<br>Which day the report will be sent out. This is required for " +
				"'weekly' report frequency.",
			Optional: true,
		},
		schemaFrequency: schema.StringAttribute{
			Description: "Enum: `daily` `weekly` `monthly`<br>The unit of frequency in which the " +
				"report is generated.",
			Optional: true,
		},
		schemaStartTime: schema.StringAttribute{
			Description: "When the report will be send out. This field should follow the format " +
				"`HH:MM` based on a 24-hour clock. Only values where HH ranges from 0 to 23 and " +
				"MM ranges from 0 to 59 are allowed.",
			Required: true,
		},
		schemaTimezone: schema.StringAttribute{
			Description: "The timezone for the report schedule. The timezone must be a valid " +
				"location name from the IANA Time Zone database. For instance, it can be " +
				"`America/New_York`, `US/Central`, `UTC`, or similar. If empty, then the timezone" +
				" is considered as UTC.",
			Optional: true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
	}

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Clumio Report Configuration Resource used to manage reports.",
		Attributes: map[string]schema.Attribute{
			schemaCreatedAt: schema.StringAttribute{
				Description: "The RFC3339 format time when the report configuration was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			schemaDescription: schema.StringAttribute{
				Description: "The user-provided description of the compliance report " +
					"configuration.",
				Optional: true,
			},
			schemaID: schema.StringAttribute{
				Description: "The unique identifier of the report configuration.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			schemaName: schema.StringAttribute{
				Description: "The user-provided name of the compliance report configuration.",
				Required:    true,
			},
		},
		Blocks: map[string]schema.Block{
			schemaNotification: schema.SetNestedBlock{
				Description: "Notification channels to send the generated report runs.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						schemaEmailList: schema.SetAttribute{
							Optional:    true,
							Description: "Email list to send a generated report run.",
							ElementType: types.StringType,
						},
					},
				},
				Validators: []validator.Set{
					common.WrapSetValidator(setvalidator.SizeAtMost(1)),
				},
			},
			schemaParameter: schema.SetNestedBlock{
				Description: "Parameters for the report configuration.",
				NestedObject: schema.NestedBlockObject{
					Blocks: map[string]schema.Block{
						schemaControls: ControlsSchemaBlock,
						schemaFilters:  FiltersSchemaBlock,
					},
				},
				Validators: []validator.Set{
					common.WrapSetValidator(setvalidator.IsRequired()),
					common.WrapSetValidator(setvalidator.SizeAtMost(1)),
				},
			},
			schemaSchedule: schema.SetNestedBlock{
				Description: "When the report will be generated and sent. If the schedule is " +
					"not provided then a default value will be used.",
				NestedObject: schema.NestedBlockObject{
					Attributes: ScheduleSchemaAttributes,
				},
				Validators: []validator.Set{
					common.WrapSetValidator(setvalidator.SizeAtMost(1)),
				},
			},
		},
	}
}
