// Copyright 2023. Clumio, Inc.

// Package clumio_protection_group contains the resource definition and CRUD implementation.
package clumio_protection_group

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/clumio-code/clumio-go-sdk/config"
	protectionGroups "github.com/clumio-code/clumio-go-sdk/controllers/protection_groups"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &protectionGroupResource{}
	_ resource.ResourceWithConfigure   = &protectionGroupResource{}
	_ resource.ResourceWithImportState = &protectionGroupResource{}
)

type protectionGroupResource struct {
	client *common.ApiClient
}

// NewProtectionGroupResource is a helper function to simplify the provider implementation.
func NewProtectionGroupResource() resource.Resource {
	return &protectionGroupResource{}
}

type prefixFilterModel struct {
	ExcludedSubPrefixes []types.String `tfsdk:"excluded_sub_prefixes"`
	Prefix              types.String   `tfsdk:"prefix"`
}

type objectFilterModel struct {
	LatestVersionOnly types.Bool           `tfsdk:"latest_version_only"`
	PrefixFilters     []*prefixFilterModel `tfsdk:"prefix_filters"`
	StorageClasses    []types.String       `tfsdk:"storage_classes"`
}

type protectionGroupResourceModel struct {
	ID                   types.String         `tfsdk:"id"`
	Name                 types.String         `tfsdk:"name"`
	Description          types.String         `tfsdk:"description"`
	BucketRule           types.String         `tfsdk:"bucket_rule"`
	ObjectFilter         []*objectFilterModel `tfsdk:"object_filter"`
	ProtectionStatus     types.String         `tfsdk:"protection_status"`
	ProtectionInfo       types.List           `tfsdk:"protection_info"`
	OrganizationalUnitID types.String         `tfsdk:"organizational_unit_id"`
}

// Schema defines the schema for the data source.
func (r *protectionGroupResource) Schema(
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
			Description: "Storage class to include in the backup. If not specified," +
				" then all objects across all storage classes will be backed up." +
				" Valid values are: S3 Standard, S3 Standard-IA," +
				" S3 Intelligent-Tiering, and S3 One Zone-IA.",
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
				Description: "Protection Group Id.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			schemaDescription: schema.StringAttribute{
				Description: "The user-assigned description of the protection group.",
				Optional:    true,
			},
			schemaName: schema.StringAttribute{
				Description: "The user-assigned name of the protection group.",
				Required:    true,
			},
			schemaBucketRule: schema.StringAttribute{
				Description: "Describes the possible conditions for a bucket to be " +
					"automatically added to a protection group. For example: " +
					"{\"aws_tag\":{\"$eq\":{\"key\":\"Environment\", \"value\":\"Prod\"}}}",
				Optional: true,
			},
			schemaOrganizationalUnitId: schema.StringAttribute{
				Description: "The Clumio-assigned ID of the organizational unit" +
					" associated with the protection group.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			schemaProtectionStatus: schema.StringAttribute{
				Description: "The protection status of the protection group. Possible" +
					" values include \"protected\", \"unprotected\", and" +
					" \"unsupported\". If the protection group does not support backups," +
					" then this field has a value of unsupported.",
				Computed: true,
			},
			schemaProtectionInfo: schema.ListNestedAttribute{
				Description: "The protection policy applied to this resource.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						schemaInheritingEntityId: schema.StringAttribute{
							Description: "The ID of the entity from which protection" +
								" was inherited.",
							Computed: true,
						},
						schemaInheritingEntityType: schema.StringAttribute{
							Description: "The type of the entity from which" +
								" protection was inherited.",
							Computed: true,
						},
						schemaPolicyId: schema.StringAttribute{
							Description: "ID of policy to apply on Protection Group",
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
					setvalidator.SizeAtMost(1),
				},
			},
		},
	}
}

// Metadata returns the data source type name.
func (r *protectionGroupResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_protection_group"
}

// Configure adds the provider configured client to the data source.
func (r *protectionGroupResource) Configure(
	_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*common.ApiClient)
}

func (r *protectionGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest,
	resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Create creates the resource and sets the initial Terraform state.
func (r *protectionGroupResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan protectionGroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.OrganizationalUnitID.ValueString() != "" {
		r.client.ClumioConfig.OrganizationalUnitContext =
			plan.OrganizationalUnitID.ValueString()
		defer r.clearOUContext()
	}
	protectionGroup := protectionGroups.NewProtectionGroupsV1(r.client.ClumioConfig)
	name := plan.Name.ValueString()
	description := plan.Description.ValueString()
	bucketRule := plan.BucketRule.ValueString()
	objectFilter := mapSchemaObjectFilterToClumioObjectFilter(plan.ObjectFilter)
	response, apiErr := protectionGroup.CreateProtectionGroup(
		models.CreateProtectionGroupV1Request{
			BucketRule:   &bucketRule,
			Description:  &description,
			Name:         &name,
			ObjectFilter: objectFilter,
		})
	if apiErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error creating Protection Group %v.", name),
			fmt.Sprintf(errorFmt, apiErr.Response))
		return
	}
	err := pollForProtectionGroup(ctx, *response.Id, r.client.ClumioConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading the created Protection Group: %v", name),
			fmt.Sprintf(errorFmt, apiErr.Response))
		return
	}

	plan.ID = types.StringValue(*response.Id)
	readResponse, apiErr := protectionGroup.ReadProtectionGroup(plan.ID.ValueString())
	if apiErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf(errorProtectionGroupReadFmt, plan.Name.ValueString()),
			fmt.Sprintf(errorFmt, apiErr.Response))
		return
	}
	if !plan.Description.IsNull() || *readResponse.Description != "" {
		plan.Description = types.StringValue(*readResponse.Description)
	}
	if !plan.BucketRule.IsNull() || *readResponse.BucketRule != "" {
		plan.BucketRule = types.StringValue(*readResponse.BucketRule)
	}
	plan.Name = types.StringValue(*readResponse.Name)
	plan.OrganizationalUnitID = types.StringValue(*readResponse.OrganizationalUnitId)
	plan.ObjectFilter = mapClumioObjectFilterToSchemaObjectFilter(
		readResponse.ObjectFilter)
	plan.ProtectionStatus = types.StringValue(*readResponse.ProtectionStatus)
	plan.ProtectionInfo, diags = mapClumioProtectionInfoToSchemaProtectionInfo(
		readResponse.ProtectionInfo)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *protectionGroupResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan protectionGroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if plan.OrganizationalUnitID.ValueString() != "" {
		r.client.ClumioConfig.OrganizationalUnitContext =
			plan.OrganizationalUnitID.ValueString()
		defer r.clearOUContext()
	}
	protectionGroup := protectionGroups.NewProtectionGroupsV1(r.client.ClumioConfig)
	name := plan.Name.ValueString()
	description := plan.Description.ValueString()
	bucketRule := plan.BucketRule.ValueString()
	objectFilter := mapSchemaObjectFilterToClumioObjectFilter(plan.ObjectFilter)
	response, apiErr := protectionGroup.UpdateProtectionGroup(plan.ID.ValueString(),
		&models.UpdateProtectionGroupV1Request{
			BucketRule:   &bucketRule,
			Description:  &description,
			Name:         &name,
			ObjectFilter: objectFilter,
		})
	if apiErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error updating Protection Group %v.", name),
			fmt.Sprintf(errorFmt, apiErr.Response))
		return
	}
	err := pollForProtectionGroup(ctx, *response.Id, r.client.ClumioConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading the updated Protection Group: %v", name),
			fmt.Sprintf(errorFmt, apiErr.Response))
		return
	}
	readResponse, apiErr := protectionGroup.ReadProtectionGroup(plan.ID.ValueString())
	if apiErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf(errorProtectionGroupReadFmt, plan.Name.ValueString()),
			fmt.Sprintf(errorFmt, apiErr.Response))
		return
	}
	if !plan.Description.IsNull() || *readResponse.Description != "" {
		plan.Description = types.StringValue(*readResponse.Description)
	}
	if !plan.BucketRule.IsNull() || *readResponse.BucketRule != "" {
		plan.BucketRule = types.StringValue(*readResponse.BucketRule)
	}
	plan.Name = types.StringValue(*readResponse.Name)
	plan.OrganizationalUnitID = types.StringValue(*readResponse.OrganizationalUnitId)
	plan.ObjectFilter = mapClumioObjectFilterToSchemaObjectFilter(
		readResponse.ObjectFilter)
	plan.ProtectionStatus = types.StringValue(*readResponse.ProtectionStatus)
	plan.ProtectionInfo, diags = mapClumioProtectionInfoToSchemaProtectionInfo(
		readResponse.ProtectionInfo)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *protectionGroupResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state protectionGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if state.OrganizationalUnitID.ValueString() != "" {
		r.client.ClumioConfig.OrganizationalUnitContext =
			state.OrganizationalUnitID.ValueString()
		defer r.clearOUContext()
	}
	protectionGroup := protectionGroups.NewProtectionGroupsV1(r.client.ClumioConfig)
	readResponse, apiErr := protectionGroup.ReadProtectionGroup(state.ID.ValueString())
	if apiErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf(errorProtectionGroupReadFmt, state.Name.ValueString()),
			fmt.Sprintf(errorFmt, apiErr.Response))
		return
	}
	if !state.Description.IsNull() || *readResponse.Description != "" {
		state.Description = types.StringValue(*readResponse.Description)
	}
	if !state.BucketRule.IsNull() || *readResponse.BucketRule != "" {
		state.BucketRule = types.StringValue(*readResponse.BucketRule)
	}
	state.Name = types.StringValue(*readResponse.Name)
	state.OrganizationalUnitID = types.StringValue(*readResponse.OrganizationalUnitId)
	state.ObjectFilter = mapClumioObjectFilterToSchemaObjectFilter(
		readResponse.ObjectFilter)
	state.ProtectionStatus = types.StringValue(*readResponse.ProtectionStatus)
	state.ProtectionInfo, diags = mapClumioProtectionInfoToSchemaProtectionInfo(
		readResponse.ProtectionInfo)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *protectionGroupResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state protectionGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if state.OrganizationalUnitID.ValueString() != "" {
		r.client.ClumioConfig.OrganizationalUnitContext =
			state.OrganizationalUnitID.ValueString()
		defer r.clearOUContext()
	}
	protectionGroup := protectionGroups.NewProtectionGroupsV1(r.client.ClumioConfig)
	_, apiErr := protectionGroup.DeleteProtectionGroup(state.ID.ValueString())
	if apiErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error deleting Protection Group %v.", state.Name.ValueString()),
			fmt.Sprintf(errorFmt, apiErr.Response))
		return
	}
}

// mapSchemaObjectFilterToClumioObjectFilter converts the schema object_filter
// to the model Object Filter
func mapSchemaObjectFilterToClumioObjectFilter(objectFilterSlice []*objectFilterModel) *models.ObjectFilter {
	objectFilter := objectFilterSlice[0]
	latestVersionOnly := objectFilter.LatestVersionOnly.ValueBool()
	storageClasses := make([]*string, 0)
	if objectFilter.StorageClasses != nil {
		for _, storageClass := range objectFilter.StorageClasses {
			storageClassStr := storageClass.ValueString()
			storageClasses = append(storageClasses, &storageClassStr)
		}
	}
	modelPrefixFilters := make([]*models.PrefixFilter, 0)
	if objectFilter.PrefixFilters != nil {
		for _, prefixFilter := range objectFilter.PrefixFilters {
			excludedSubPrefixesList := make([]*string, 0)
			for _, excludedSubPrefix := range prefixFilter.ExcludedSubPrefixes {
				excludedSubPrefixStr := excludedSubPrefix.ValueString()
				excludedSubPrefixesList = append(
					excludedSubPrefixesList, &excludedSubPrefixStr)
			}
			prefix := prefixFilter.Prefix.ValueString()
			modelPrefixFilter := &models.PrefixFilter{
				ExcludedSubPrefixes: excludedSubPrefixesList,
				Prefix:              &prefix,
			}
			modelPrefixFilters = append(modelPrefixFilters, modelPrefixFilter)
		}
	}
	return &models.ObjectFilter{
		LatestVersionOnly: &latestVersionOnly,
		PrefixFilters:     modelPrefixFilters,
		StorageClasses:    storageClasses,
	}
}

// pollForProtectionGroup polls till the protection group becomes available after create
// or update protection group as they are asynchronous operations.
func pollForProtectionGroup(ctx context.Context, id string, config config.Config) error {
	protectionGroup := protectionGroups.NewProtectionGroupsV1(config)
	interval := time.Duration(intervalInSec) * time.Second
	ticker := time.NewTicker(interval)
	timeout := time.After(time.Duration(timeoutInSec) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return errors.New("context done")
		case <-ticker.C:
			_, err := protectionGroup.ReadProtectionGroup(id)
			if err != nil {
				errResponseStr := string(err.Response)
				if !strings.Contains(errResponseStr,
					"A resource with the requested ID could not be found.") {
					return errors.New(
						"error reading protection-group which was created")
				}
				continue
			}
			return nil
		case <-timeout:
			return errors.New("polling timeout")
		}
	}
}

// mapClumioObjectFilterToSchemaObjectFilter converts the Object Filter from the
// API to the schema object_filter
func mapClumioObjectFilterToSchemaObjectFilter(
	modelObjectFilter *models.ObjectFilter) []*objectFilterModel {
	schemaObjFilter := &objectFilterModel{}
	if modelObjectFilter.LatestVersionOnly != nil {
		schemaObjFilter.LatestVersionOnly = types.BoolValue(
			*modelObjectFilter.LatestVersionOnly)
	}
	if modelObjectFilter.PrefixFilters != nil {
		prefixFilters := make([]*prefixFilterModel, 0)
		for _, modelPrefixFilter := range modelObjectFilter.PrefixFilters {
			prefixFilter := &prefixFilterModel{}
			prefixFilter.Prefix = types.StringValue(*modelPrefixFilter.Prefix)
			if modelPrefixFilter.ExcludedSubPrefixes != nil {
				excludedSubPrefixes := make([]types.String, 0)
				for _, excludeSubPrefix := range modelPrefixFilter.ExcludedSubPrefixes {
					excludeSubPrefixStr := types.StringValue(*excludeSubPrefix)
					excludedSubPrefixes = append(excludedSubPrefixes, excludeSubPrefixStr)
				}
				prefixFilter.ExcludedSubPrefixes = excludedSubPrefixes
			}
			prefixFilters = append(prefixFilters, prefixFilter)
		}
		schemaObjFilter.PrefixFilters = prefixFilters
	}
	storageClasses := make([]types.String, 0)
	for _, storageClass := range modelObjectFilter.StorageClasses {
		storageClassStrType := types.StringValue(*storageClass)
		storageClasses = append(storageClasses, storageClassStrType)
	}
	schemaObjFilter.StorageClasses = storageClasses
	return []*objectFilterModel{schemaObjFilter}
}

// mapClumioProtectionInfoToSchemaProtectionInfo converts the Protection Info from the
// API to the schema protection_info.
func mapClumioProtectionInfoToSchemaProtectionInfo(
	modelProtectionInfo *models.ProtectionInfoWithRule) (types.List, diag.Diagnostics) {

	objtype := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			schemaPolicyId:             types.StringType,
			schemaInheritingEntityId:   types.StringType,
			schemaInheritingEntityType: types.StringType,
		},
	}
	if modelProtectionInfo == nil {
		return types.ListValue(objtype, []attr.Value{})
	}
	attrTypes := make(map[string]attr.Type)
	attrTypes[schemaPolicyId] = types.StringType
	attrTypes[schemaInheritingEntityType] = types.StringType
	attrTypes[schemaInheritingEntityId] = types.StringType

	attrValues := make(map[string]attr.Value)
	attrValues[schemaPolicyId] = types.StringValue("")
	attrValues[schemaInheritingEntityType] = types.StringValue("")
	attrValues[schemaInheritingEntityId] = types.StringValue("")
	if modelProtectionInfo != nil {
		if modelProtectionInfo.PolicyId != nil {
			attrValues[schemaPolicyId] = types.StringValue(*modelProtectionInfo.PolicyId)
		}
		if modelProtectionInfo.InheritingEntityType != nil {
			attrValues[schemaInheritingEntityType] =
				types.StringValue(*modelProtectionInfo.InheritingEntityType)
		}
		if modelProtectionInfo.InheritingEntityId != nil {
			attrValues[schemaInheritingEntityId] =
				types.StringValue(*modelProtectionInfo.InheritingEntityId)
		}
	}
	obj, diags := types.ObjectValue(attrTypes, attrValues)

	listobj, listdiag := types.ListValue(objtype, []attr.Value{obj})
	listdiag.Append(diags...)
	return listobj, listdiag
}

func (r *protectionGroupResource) clearOUContext() {
	r.client.ClumioConfig.OrganizationalUnitContext = ""
}
