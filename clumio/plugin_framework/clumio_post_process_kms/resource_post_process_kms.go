// Copyright 2023. Clumio, Inc.
//
// clumio_post_process_kms definition and CRUD implementation.

package clumio_post_process_kms

import (
	"context"
	"fmt"

	kms "github.com/clumio-code/clumio-go-sdk/controllers/post_process_kms"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &clumioPostProcessKmsResource{}
	_ resource.ResourceWithConfigure = &clumioPostProcessKmsResource{}
)

// NewClumioPostProcessKmsResource is a helper function to simplify the provider implementation.
func NewClumioPostProcessKmsResource() resource.Resource {
	return &clumioPostProcessKmsResource{}
}

// clumioPostProcessKmsResource is the resource implementation.
type clumioPostProcessKmsResource struct {
	client *common.ApiClient
}

// clumioPostProcessKmsResourceModel model
type clumioPostProcessKmsResourceModel struct {
	Id                    types.String `tfsdk:"id"`
	Token                 types.String `tfsdk:"token"`
	AccountId             types.String `tfsdk:"account_id"`
	Region                types.String `tfsdk:"region"`
	RoleId                types.String `tfsdk:"role_id"`
	RoleArn               types.String `tfsdk:"role_arn"`
	RoleExternalId        types.String `tfsdk:"role_external_id"`
	CreatedMultiRegionCMK types.Bool   `tfsdk:"created_multi_region_cmk"`
	MultiRegionCMKKeyId   types.String `tfsdk:"multi_region_cmk_key_id"`
	TemplateVersion       types.Int64  `tfsdk:"template_version"`
}

// Metadata returns the resource type name.
func (r *clumioPostProcessKmsResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_post_process_kms"
}

// Schema defines the schema for the resource.
func (r *clumioPostProcessKmsResource) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{
		Description: "Post-Process Clumio KMS Resource used to post-process KMS in Clumio.",
		Attributes: map[string]schema.Attribute{
			schemaId: schema.StringAttribute{
				Description: "The ID of this resource.",
				Computed:    true,
			},
			schemaToken: schema.StringAttribute{
				Description: "The AWS integration ID token.",
				Required:    true,
			},
			schemaAccountId: schema.StringAttribute{
				Description: "The AWS Customer Account ID associated with the connection.",
				Required:    true,
			},
			schemaRegion: schema.StringAttribute{
				Description: "The AWS Region.",
				Required:    true,
			},
			schemaRoleId: schema.StringAttribute{
				Description: "The ID of the IAM role to manage the CMK.",
				Required:    true,
			},
			schemaRoleArn: schema.StringAttribute{
				Description: "The ARN of the IAM role to manage the CMK.",
				Required:    true,
			},
			schemaRoleExternalId: schema.StringAttribute{
				Description: "The external ID to use when assuming the IAM role.",
				Required:    true,
			},
			schemaCreatedMultiRegionCMK: schema.BoolAttribute{
				Description: "Whether a new CMK was created.",
				Optional:    true,
			},
			schemaMultiRegionCMKKeyID: schema.StringAttribute{
				Description: "Multi Region CMK Key ID.",
				Optional:    true,
			},
			schemaTemplateVersion: schema.Int64Attribute{
				Description: "Template version",
				Optional:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *clumioPostProcessKmsResource) Configure(
	_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*common.ApiClient)
}

// Create creates the resource and sets the initial Terraform state.
func (r *clumioPostProcessKmsResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var plan clumioPostProcessKmsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = r.clumioPostProcessKmsCommon(ctx, plan, "Create")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	accountId := plan.AccountId.ValueString()
	awsRegion := plan.Region.ValueString()
	token := plan.Token.ValueString()
	plan.Id = types.StringValue(fmt.Sprintf("%v/%v/%v", accountId, awsRegion, token))
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *clumioPostProcessKmsResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *clumioPostProcessKmsResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan clumioPostProcessKmsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state clumioPostProcessKmsResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.Id = state.Id

	diags = r.clumioPostProcessKmsCommon(ctx, plan, "Update")
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

// Delete deletes the resource and removes the Terraform state on success.
func (r *clumioPostProcessKmsResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var state clumioPostProcessKmsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = r.clumioPostProcessKmsCommon(ctx, state, "Delete")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *clumioPostProcessKmsResource) clumioPostProcessKmsCommon(
	_ context.Context, state clumioPostProcessKmsResourceModel, eventType string) diag.Diagnostics {

	postProcessAwsKMS := kms.NewPostProcessKmsV1(r.client.ClumioConfig)
	accountId := state.AccountId.ValueString()
	awsRegion := state.Region.ValueString()
	multiRegionCMKKeyId := state.MultiRegionCMKKeyId.ValueString()
	token := state.Token.ValueString()
	roleId := state.RoleId.ValueString()
	roleArn := state.RoleArn.ValueString()
	roleExternalId := state.RoleExternalId.ValueString()
	createdMrCmk := state.CreatedMultiRegionCMK.ValueBool()
	templateVersion := uint64(state.TemplateVersion.ValueInt64())

	_, apiErr := postProcessAwsKMS.PostProcessKms(
		&models.PostProcessKmsV1Request{
			AccountNativeId:       &accountId,
			AwsRegion:             &awsRegion,
			RequestType:           &eventType,
			Token:                 &token,
			MultiRegionCmkKeyId:   &multiRegionCMKKeyId,
			RoleId:                &roleId,
			RoleArn:               &roleArn,
			RoleExternalId:        &roleExternalId,
			CreatedMultiRegionCmk: &createdMrCmk,
			Version:               &templateVersion,
		})
	if apiErr != nil {
		diagnostics := diag.Diagnostics{}
		diagnostics.AddError("Error in invoking Post-process Clumio KMS.",
			fmt.Sprintf("Error: %v", apiErr))
		return diagnostics
	}
	return nil
}
