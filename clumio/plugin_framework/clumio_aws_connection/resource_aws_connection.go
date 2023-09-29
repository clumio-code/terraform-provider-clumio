// Copyright 2023. Clumio, Inc.

// clumio_aws_connection definition and CRUD implementation.

package clumio_aws_connection

import (
	"context"
	"fmt"
	"strings"

	aws_connections "github.com/clumio-code/clumio-go-sdk/controllers/aws_connections"
	awsEnvs "github.com/clumio-code/clumio-go-sdk/controllers/aws_environments"
	orgUnits "github.com/clumio-code/clumio-go-sdk/controllers/organizational_units"
	"github.com/clumio-code/clumio-go-sdk/models"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &clumioAWSConnectionResource{}
	_ resource.ResourceWithConfigure   = &clumioAWSConnectionResource{}
	_ resource.ResourceWithImportState = &clumioAWSConnectionResource{}
)

// NewClumioAWSConnectionResource is a helper function to simplify the provider implementation.
func NewClumioAWSConnectionResource() resource.Resource {
	return &clumioAWSConnectionResource{}
}

// clumioAWSConnectionResource is the resource implementation.
type clumioAWSConnectionResource struct {
	client *common.ApiClient
}

// clumioAWSConnectionResource model
type clumioAWSConnectionResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	AccountNativeID      types.String `tfsdk:"account_native_id"`
	AWSRegion            types.String `tfsdk:"aws_region"`
	Description          types.String `tfsdk:"description"`
	OrganizationalUnitID types.String `tfsdk:"organizational_unit_id"`
	ConnectionStatus     types.String `tfsdk:"connection_status"`
	Token                types.String `tfsdk:"token"`
	Namespace            types.String `tfsdk:"namespace"`
	ClumioAWSAccountID   types.String `tfsdk:"clumio_aws_account_id"`
	ClumioAWSRegion      types.String `tfsdk:"clumio_aws_region"`
	ExternalID           types.String `tfsdk:"role_external_id"`
	DataPlaneAccountID   types.String `tfsdk:"data_plane_account_id"`
}

// Metadata returns the resource type name.
func (r *clumioAWSConnectionResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_aws_connection"
}

// Schema defines the schema for the resource.
func (r *clumioAWSConnectionResource) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Clumio AWS Connection Resource used to connect" +
			" AWS accounts to Clumio.",
		Attributes: map[string]schema.Attribute{
			schemaId: schema.StringAttribute{
				Description: "Clumio AWS Connection Id.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			schemaAccountNativeId: schema.StringAttribute{
				Description: "AWS Account Id to connect to Clumio.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			schemaAwsRegion: schema.StringAttribute{
				Description: "AWS Region of account.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			schemaDescription: schema.StringAttribute{
				Description: "Clumio AWS Connection Description.",
				Optional:    true,
			},
			schemaOrganizationalUnitId: schema.StringAttribute{
				Description: "Clumio Organizational Unit Id.",
				Optional:    true,
				Computed:    true,
			},
			schemaConnectionStatus: schema.StringAttribute{
				Description: "The status of the connection. Possible values include " +
					"connecting, connected and unlinked.",
				Computed: true,
			},
			schemaToken: schema.StringAttribute{
				Description: "The 36-character Clumio AWS integration ID token used to" +
					" identify the installation of the Terraform template on the account.",
				Computed: true,
			},
			schemaNamespace: schema.StringAttribute{
				Description: "K8S Namespace.",
				Computed:    true,
			},
			schemaClumioAwsAccountId: schema.StringAttribute{
				Description: "Clumio AWS AccountId.",
				Computed:    true,
			},
			schemaClumioAwsRegion: schema.StringAttribute{
				Description: "Clumio AWS Region.",
				Computed:    true,
			},
			schemaExternalId: schema.StringAttribute{
				Description: "A key used by Clumio to assume the service role in your account.",
				Computed: true,
			},
			schemaDataPlaneAccountId: schema.StringAttribute{
				Description: "The internal representation to uniquely identify a given data plane.",
				Computed: true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *clumioAWSConnectionResource) Configure(
	_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*common.ApiClient)
}

// Create creates the resource and sets the initial Terraform state.
func (r *clumioAWSConnectionResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client := r.client
	var plan clumioAWSConnectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	awsConnection := aws_connections.NewAwsConnectionsV1(client.ClumioConfig)
	accountNativeId := plan.AccountNativeID.ValueString()
	awsRegion := plan.AWSRegion.ValueString()
	description := plan.Description.ValueString()
	organizationalUnitId := plan.OrganizationalUnitID.ValueString()
	res, apiErr := awsConnection.CreateAwsConnection(&models.CreateAwsConnectionV1Request{
		AccountNativeId:      &accountNativeId,
		AwsRegion:            &awsRegion,
		Description:          &description,
		OrganizationalUnitId: &organizationalUnitId,
	})
	if apiErr != nil {
		resp.Diagnostics.AddError(
			"Error creating Clumio AWS Connection.",
			fmt.Sprintf(errorFmt, string(apiErr.Response)))
		return
	}
	plan.ID = types.StringValue(*res.Id)
	plan.Token = types.StringValue(*res.Token)
	if res.Namespace != nil {
		plan.Namespace = types.StringValue(*res.Namespace)
	}
	plan.ClumioAWSAccountID = types.StringValue(*res.ClumioAwsAccountId)
	plan.ClumioAWSRegion = types.StringValue(*res.ClumioAwsRegion) // Set state to fully populated data
	plan.OrganizationalUnitID = types.StringValue(*res.OrganizationalUnitId)
	plan.ConnectionStatus = types.StringValue(*res.ConnectionStatus)
	plan.ExternalID = types.StringValue(*res.ExternalId)
	plan.DataPlaneAccountID = types.StringValue(*res.DataPlaneAccountId)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *clumioAWSConnectionResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state clumioAWSConnectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	awsConnection := aws_connections.NewAwsConnectionsV1(r.client.ClumioConfig)
	res, apiErr := awsConnection.ReadAwsConnection(state.ID.ValueString())
	if apiErr != nil {
		if strings.Contains(apiErr.Error(), "The resource is not found.") {
			state.ID = types.StringValue("")
		}
		resp.Diagnostics.AddError(
			"Error retrieving Clumio AWS Connection.",
			fmt.Sprintf(errorFmt, string(apiErr.Response)))
		return
	}
	state.ConnectionStatus = types.StringValue(*res.ConnectionStatus)
	state.Token = types.StringValue(*res.Token)
	if res.Namespace != nil {
		state.Namespace = types.StringValue(*res.Namespace)
	}
	state.ClumioAWSAccountID = types.StringValue(*res.ClumioAwsAccountId)
	state.ClumioAWSRegion = types.StringValue(*res.ClumioAwsRegion)
	state.AccountNativeID = types.StringValue(*res.AccountNativeId)
	state.AWSRegion = types.StringValue(*res.AwsRegion)
	state.ExternalID = types.StringValue(*res.ExternalId)
	state.DataPlaneAccountID = types.StringValue(*res.DataPlaneAccountId)
	if res.Description != nil {
		state.Description = types.StringValue(*res.Description)
	}
	if res.OrganizationalUnitId != nil {
		state.OrganizationalUnitID = types.StringValue(*res.OrganizationalUnitId)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *clumioAWSConnectionResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan clumioAWSConnectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state clumioAWSConnectionResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updated := updateOUForConnectionIfNeeded(ctx, r.client, req, resp)
	if resp.Diagnostics.HasError() {
		return
	}
	if plan.Description == state.Description {
		if updated {
			state.OrganizationalUnitID = plan.OrganizationalUnitID
			diags = resp.State.Set(ctx, plan)
			resp.Diagnostics.Append(diags...)
		}
		return
	}
	awsConnection := aws_connections.NewAwsConnectionsV1(r.client.ClumioConfig)
	description := plan.Description.ValueString()
	res, apiErr := awsConnection.UpdateAwsConnection(plan.ID.ValueString(),
		models.UpdateAwsConnectionV1Request{
			Description: &description,
		})
	if apiErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf(
				"Error updating description of Clumio AWS Connection %v.",
				plan.ID.ValueString()),
			fmt.Sprintf(errorFmt, string(apiErr.Response)))
		return
	}
	plan.Token = types.StringValue(*res.Token)
	if res.Namespace != nil {
		plan.Namespace = types.StringValue(*res.Namespace)
	}
	plan.ClumioAWSAccountID = types.StringValue(*res.ClumioAwsAccountId)
	plan.ClumioAWSRegion = types.StringValue(*res.ClumioAwsRegion)
	plan.OrganizationalUnitID = types.StringValue(*res.OrganizationalUnitId)
	plan.ConnectionStatus = types.StringValue(*res.ConnectionStatus)
	plan.ExternalID = types.StringValue(*res.ExternalId)
	plan.DataPlaneAccountID = types.StringValue(*res.DataPlaneAccountId)
	plan.ID = types.StringValue(*res.Id)
	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// updateOUForConnectionIfNeeded updates the OU for the connection if the new OU provided
// is either the parent of the current OU or one of its immediate descendant.
func updateOUForConnectionIfNeeded(ctx context.Context, client *common.ApiClient,
	req resource.UpdateRequest, resp *resource.UpdateResponse) bool {
	ouUpdated := false

	// Retrieve values from plan
	var plan clumioAWSConnectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return ouUpdated
	}

	var state clumioAWSConnectionResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return ouUpdated
	}

	if !plan.OrganizationalUnitID.IsUnknown() && plan.OrganizationalUnitID != state.OrganizationalUnitID {
		connectionStatus := plan.ConnectionStatus.ValueString()
		if connectionStatus != "" && connectionStatus != statusConnected {
			resp.Diagnostics.AddError(fmt.Sprintf("Connection status is %v.", connectionStatus),
				"Updating organizational_unit_id"+
					" is allowed only if the connection status in \"connected\". To make the"+
					" connection status as connected, install the clumio terraform aws"+
					" template module.")
			return ouUpdated
		}
		envId := GetEnvironmentId(ctx, client, req, resp)
		if resp.Diagnostics.HasError() {
			return ouUpdated
		}
		ouIdStr, isNewOUCurrentOUParent :=
			validateAndGetOUIDToPatch(ctx, client, req, resp)
		if resp.Diagnostics.HasError() {
			return ouUpdated
		}
		var removeEntityModels []*models.EntityModel
		var addEntityModels []*models.EntityModel
		awsEnv := awsEnvironment
		entityModels := []*models.EntityModel{
			{
				PrimaryEntity: &models.OrganizationalUnitPrimaryEntity{
					Id:         &envId,
					ClumioType: &awsEnv,
				},
			},
		}
		if isNewOUCurrentOUParent {
			removeEntityModels = entityModels
		} else {
			addEntityModels = entityModels
		}
		ouUpdateRequest := &models.PatchOrganizationalUnitV1Request{
			Entities: &models.UpdateEntities{
				Add:    addEntityModels,
				Remove: removeEntityModels,
			},
		}
		orgUnitsAPI := orgUnits.NewOrganizationalUnitsV1(client.ClumioConfig)
		res, apiErr := orgUnitsAPI.PatchOrganizationalUnit(ouIdStr, nil, ouUpdateRequest)
		if apiErr != nil {
			resp.Diagnostics.AddError(
				"Error updating the Organizational Unit for the connection.",
				fmt.Sprintf(errorFmt, apiErr))
			return ouUpdated
		}
		if res.StatusCode == http202 {
			err := common.PollTask(ctx, client, *res.Http202.TaskId, pollTimeoutInSec, pollIntervalInSec)
			if err != nil {
				resp.Diagnostics.AddError("Error while polling for the Update OU task for the connection",
					fmt.Sprintf(errorFmt, err))
				return ouUpdated
			}
		}
		ouUpdated = true
	}
	return ouUpdated
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *clumioAWSConnectionResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state clumioAWSConnectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	awsConnection := aws_connections.NewAwsConnectionsV1(r.client.ClumioConfig)
	_, apiErr := awsConnection.DeleteAwsConnection(state.ID.ValueString())
	if apiErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf(
				"Error deleting Clumio AWS Connection %v.", state.ID.ValueString()),
			fmt.Sprintf(errorFmt, string(apiErr.Response)))
	}

}

// GetEnvironmentId returns the Environment ID corresponding to the AWS connection
func GetEnvironmentId(ctx context.Context, client *common.ApiClient,
	req resource.UpdateRequest, resp *resource.UpdateResponse) string {
	// Retrieve values from plan
	var state clumioAWSConnectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return ""
	}

	awsEnvironmentsAPI := awsEnvs.NewAwsEnvironmentsV1(client.ClumioConfig)
	accountNativeId := state.AccountNativeID.ValueString()
	awsRegion := state.AWSRegion.ValueString()
	limit := int64(1)
	filterStr := fmt.Sprintf(
		"{\"account_native_id\":{\"$eq\":\"%v\"}, \"aws_region\":{\"$eq\":\"%v\"}}",
		accountNativeId, awsRegion)
	embed := "read-organizational-unit"
	envs, apiErr := awsEnvironmentsAPI.ListAwsEnvironments(
		&limit, nil, &filterStr, &embed)
	if apiErr != nil {
		resp.Diagnostics.AddError(
			"Error retrieving AWS Environment.",
			fmt.Sprintf("Error retrieving AWS Environment corresponding to "+
				"Account ID: %v and AWS Region: %v. Error: %v",
				accountNativeId, awsRegion, apiErr))
		return ""
	}
	if *envs.CurrentCount > 1 {
		resp.Diagnostics.AddError(
			"More than one environment found.",
			fmt.Sprintf("Expected only one environment but found %v", *envs.CurrentCount))
		return ""
	}
	envId := *envs.Embedded.Items[0].Id
	return envId
}

// validateAndGetOUIDToPatch validates the new organizational_unit_id and returns the
// organizational_unit_id to update.
func validateAndGetOUIDToPatch(ctx context.Context, client *common.ApiClient,
	req resource.UpdateRequest, resp *resource.UpdateResponse) (
	string, bool) {
	var state clumioAWSConnectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return "", false
	}
	var plan clumioAWSConnectionResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return "", false
	}
	orgUnitsAPI := orgUnits.NewOrganizationalUnitsV1(client.ClumioConfig)
	isValidNewOU := false
	isNewOUCurrentOUParent := false
	oldOUIdStr := state.OrganizationalUnitID.ValueString()
	oldOU, apiErr := orgUnitsAPI.ReadOrganizationalUnit(oldOUIdStr, nil)
	if apiErr != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error retrieving current OU %v.",
			state.OrganizationalUnitID.ValueString()),
			fmt.Sprintf("Error: %v.", apiErr))
		return "", false
	}
	newOUIdStr := plan.OrganizationalUnitID.ValueString()
	ouIdStr := newOUIdStr
	if oldOU.ParentId != nil && *oldOU.ParentId == newOUIdStr {
		isValidNewOU = true
		isNewOUCurrentOUParent = true
		ouIdStr = oldOUIdStr
	}
	filterStr := fmt.Sprintf("{\"parent_id\": {\"$eq\": \"%v\"}}", oldOUIdStr)
	listRes, apiErr := orgUnitsAPI.ListOrganizationalUnits(nil, nil, &filterStr)
	if apiErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error retrieving child OUs of current OU: %v.",
				state.OrganizationalUnitID.ValueString()),
			fmt.Sprintf("Error: %v.", apiErr))
		return "", false
	}
	if listRes != nil && listRes.Embedded != nil && len(listRes.Embedded.Items) > 0 {
		for _, ouObj := range listRes.Embedded.Items {
			if *ouObj.Id == newOUIdStr {
				isValidNewOU = true
				break
			}
		}
	}
	if !isValidNewOU {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Invalid Organizational Unit ID: %v specified.", newOUIdStr),
			fmt.Sprintf("Invalid Organizational Unit ID: %v specified."+
				" The Organizational Unit should either be a parent of the current"+
				" Organizational Unit or its immediate descendant.",
				newOUIdStr))
		return "", false
	}
	return ouIdStr, isNewOUCurrentOUParent
}

func (r *clumioAWSConnectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
