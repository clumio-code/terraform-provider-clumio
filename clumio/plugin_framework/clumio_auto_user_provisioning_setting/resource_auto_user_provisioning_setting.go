// Copyright 2023. Clumio, Inc.

package clumio_auto_user_provisioning_setting

import (
	"context"
	"fmt"

	aupSettings "github.com/clumio-code/clumio-go-sdk/controllers/auto_user_provisioning_settings"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &autoUserProvisioningSettingResource{}
	_ resource.ResourceWithConfigure = &autoUserProvisioningSettingResource{}
)

// autoUserProvisioningSettingResource is the resource implementation.
type autoUserProvisioningSettingResource struct {
	client *common.ApiClient
}

// NewAutoUserProvisioningSettingResource is a helper function to simplify the provider implementation.
func NewAutoUserProvisioningSettingResource() resource.Resource {
	return &autoUserProvisioningSettingResource{}
}

// autoUserProvisioningSettingResource model
type autoUserProvisioningSettingResourceModel struct {
	ID        types.String `tfsdk:"id"`
	IsEnabled types.Bool   `tfsdk:"is_enabled"`
}

// Metadata returns the resource type name.
func (r *autoUserProvisioningSettingResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_auto_user_provisioning_setting"
}

// Schema defines the schema for the data source.
func (r *autoUserProvisioningSettingResource) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Clumio Auto User Provisioning Setting Resource used to determine if Auto User Provisioning " +
			"feature is enabled or not.",
		Attributes: map[string]schema.Attribute{
			schemaId: schema.StringAttribute{
				Description: "Auto User Provisioning Setting Id.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			schemaIsEnabled: schema.BoolAttribute{
				Description: "Whether auto user provisioning is enabled or not.",
				Required:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *autoUserProvisioningSettingResource) Configure(
	_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*common.ApiClient)
}

// Create creates the resource and sets the initial Terraform state.
func (r *autoUserProvisioningSettingResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var plan autoUserProvisioningSettingResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	aups := aupSettings.NewAutoUserProvisioningSettingsV1(r.client.ClumioConfig)
	isEnabled := plan.IsEnabled.ValueBool()
	aupsRequest := &models.UpdateAutoUserProvisioningSettingV1Request{
		IsEnabled: &isEnabled,
	}
	_, apiErr := aups.UpdateAutoUserProvisioningSetting(aupsRequest)
	if apiErr != nil {
		resp.Diagnostics.AddError(
			"Error creating auto user provisioning setting.",
			fmt.Sprintf(errorFmt, string(apiErr.Response)),
		)
		return
	}

	plan.ID = types.StringValue(uuid.New().String())

	// Set the state.
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *autoUserProvisioningSettingResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state autoUserProvisioningSettingResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	aups := aupSettings.NewAutoUserProvisioningSettingsV1(r.client.ClumioConfig)
	res, apiErr := aups.ReadAutoUserProvisioningSetting()
	if apiErr != nil {
		resp.Diagnostics.AddError(
			"Error retrieving auto user provisioning setting.",
			fmt.Sprintf(errorFmt, string(apiErr.Response)),
		)
		return
	}

	state.IsEnabled = types.BoolValue(*res.IsEnabled)

	// Set refreshed state.
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *autoUserProvisioningSettingResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan autoUserProvisioningSettingResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	aups := aupSettings.NewAutoUserProvisioningSettingsV1(r.client.ClumioConfig)
	isEnabled := plan.IsEnabled.ValueBool()
	aupsRequest := &models.UpdateAutoUserProvisioningSettingV1Request{
		IsEnabled: &isEnabled,
	}

	res, apiErr := aups.UpdateAutoUserProvisioningSetting(aupsRequest)
	if apiErr != nil {
		resp.Diagnostics.AddError(
			"Error updating auto user provisioning setting.",
			fmt.Sprintf(errorFmt, string(apiErr.Response)),
		)
		return
	}

	plan.IsEnabled = types.BoolValue(*res.IsEnabled)

	// Set state to fully populated data.
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *autoUserProvisioningSettingResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state autoUserProvisioningSettingResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	aups := aupSettings.NewAutoUserProvisioningSettingsV1(r.client.ClumioConfig)
	isEnabled := false
	aupsRequest := &models.UpdateAutoUserProvisioningSettingV1Request{
		IsEnabled: &isEnabled,
	}
	_, apiErr := aups.UpdateAutoUserProvisioningSetting(aupsRequest)
	if apiErr != nil {
		resp.Diagnostics.AddError(
			"Error deleting auto user provisioning setting.",
			fmt.Sprintf(errorFmt, string(apiErr.Response)),
		)
	}
}
