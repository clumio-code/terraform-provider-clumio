// Copyright 2023. Clumio, Inc.
//
// clumio_wallet resource definition and CRUD implementation.

package clumio_wallet

import (
	"context"
	"fmt"

	"github.com/clumio-code/clumio-go-sdk/controllers/wallets"
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
	_ resource.Resource                = &clumioWalletResource{}
	_ resource.ResourceWithConfigure   = &clumioWalletResource{}
	_ resource.ResourceWithImportState = &clumioWalletResource{}
)

// NewClumioWalletResource is a helper function to simplify the provider implementation.
func NewClumioWalletResource() resource.Resource {
	return &clumioWalletResource{}
}

// clumioWalletResource is the resource implementation.
type clumioWalletResource struct {
	client *common.ApiClient
}

// clumioWalletResource model
type clumioWalletResourceModel struct {
	Id              types.String `tfsdk:"id"`
	AccountNativeId types.String `tfsdk:"account_native_id"`
	Token           types.String `tfsdk:"token"`
	State           types.String `tfsdk:"state"`
	ClumioAccountId types.String `tfsdk:"clumio_account_id"`
}

// Metadata returns the resource type name.
func (r *clumioWalletResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_wallet"
}

// Schema defines the schema for the resource.
func (r *clumioWalletResource) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{
		Description: "Clumio Wallet Resource to create and manage wallets in Clumio. " +
			"Wallets should be created \"after\" connecting an AWS account to Clumio.<br>" +
			"**NOTE:** To protect against accidental deletion, wallets cannot be destroyed. " +
			"To remove a wallet, contact Clumio support.",
		Attributes: map[string]schema.Attribute{
			schemaId: schema.StringAttribute{
				Description: "Wallet Id.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			schemaAccountNativeId: schema.StringAttribute{
				Description: "The AWS account id to be associated with the wallet.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			schemaToken: schema.StringAttribute{
				Description: "Token is used to identify and authenticate the CloudFormation stack creation.",
				Computed:    true,
			},
			schemaState: schema.StringAttribute{
				Description: "State describes the state of the wallet. Valid states are:" +
					" Waiting: The wallet has been created, but a stack hasn't been" +
					" created. The wallet can't be used in this state. Enabled: The" +
					" wallet has been created and a stack has been created for the" +
					" wallet. This is the normal expected state of a wallet in use." +
					" Error: The wallet is inaccessible. See ErrorCode and ErrorMessage" +
					" fields for additional details.",
				Computed: true,
			},
			schemaClumioAccountId: schema.StringAttribute{
				Description: "Clumio Account ID.",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *clumioWalletResource) Configure(
	_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*common.ApiClient)
}

// Create creates the resource and sets the initial Terraform state.
func (r *clumioWalletResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var plan clumioWalletResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	walletsAPI := wallets.NewWalletsV1(r.client.ClumioConfig)
	accountNativeId := plan.AccountNativeId.ValueString()
	res, apiErr := walletsAPI.CreateWallet(&models.CreateWalletV1Request{
		AccountNativeId: &accountNativeId,
	})
	if apiErr != nil {
		resp.Diagnostics.AddError(
			"Error creating Clumio wallet.",
			fmt.Sprintf(errorFmt, string(apiErr.Response)))
		return
	}

	plan.Id = types.StringValue(*res.Id)
	plan.State = types.StringValue(*res.State)
	plan.AccountNativeId = types.StringValue(*res.AccountNativeId)
	plan.Token = types.StringValue(*res.Token)
	plan.ClumioAccountId = types.StringValue(*res.ClumioAwsAccountId)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *clumioWalletResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var state clumioWalletResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	walletsAPI := wallets.NewWalletsV1(r.client.ClumioConfig)
	res, apiErr := walletsAPI.ReadWallet(state.Id.ValueString())
	if apiErr != nil {
		resp.Diagnostics.AddError(
			"Error reading Clumio wallet.",
			fmt.Sprintf(errorFmt, string(apiErr.Response)))
		return
	}

	state.State = types.StringValue(*res.State)
	state.AccountNativeId = types.StringValue(*res.AccountNativeId)
	state.Token = types.StringValue(*res.Token)
	state.ClumioAccountId = types.StringValue(*res.ClumioAwsAccountId)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *clumioWalletResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

}

// Delete deletes the resource and removes the Terraform state on success.
func (r *clumioWalletResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var state clumioWalletResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	walletsAPI := wallets.NewWalletsV1(r.client.ClumioConfig)
	_, apiErr := walletsAPI.DeleteWallet(state.Id.ValueString())
	if apiErr != nil {
		resp.Diagnostics.AddError(
			"Error reading Clumio wallet.",
			fmt.Sprintf(errorFmt, string(apiErr.Response)))
	}
}

func (r *clumioWalletResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
