// Copyright 2024. Clumio, Inc.

// This file holds the type definition and Schema resource function used by the resource model for
// the clumio_wallet Terraform resource.

package clumio_wallet

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// clumioWalletResourceModel is the resource model for the clumio_wallet Terraform
// resource. It represents the schema of the resource and the data it holds. This schema is used by
// customers to configure the resource and by the Clumio provider to read and write the resource.
type clumioWalletResourceModel struct {
	Id              types.String `tfsdk:"id"`
	AccountNativeId types.String `tfsdk:"account_native_id"`
	Token           types.String `tfsdk:"token"`
	State           types.String `tfsdk:"state"`
	ClumioAccountId types.String `tfsdk:"clumio_account_id"`
}

// Schema defines the structure and constraints of the clumio_wallet Terraform resource.
// Schema is a method on the clumioWalletResource struct. It sets the schema for the
// clumio_wallet Terraform resource, which is used to create a Clumio wallet.
func (r *clumioWalletResource) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{
		Description: "Clumio Wallet Resource to create and manage wallets in Clumio. " +
			"Wallets should be created \"after\" connecting an AWS account to Clumio.<br>" +
			"**NOTE:** To protect against accidental deletion, wallets cannot be destroyed once the" +
			" byok-template module has been installed. To remove a wallet, contact Clumio support.",
		Attributes: map[string]schema.Attribute{
			schemaId: schema.StringAttribute{
				Description: "Unique identifier for the Clumio Wallet.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			schemaAccountNativeId: schema.StringAttribute{
				Description: "Identifier of the AWS account to be setup with BYOK and associated" +
					" with the wallet.",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			schemaToken: schema.StringAttribute{
				Description: "Token used to identify resources set up by the BYOK template" +
					" installation on the account being connected.",
				Computed: true,
			},
			schemaState: schema.StringAttribute{
				Description: "State describes the state of the wallet. Valid states are:\n" +
					"\tWaiting: The wallet has been created, but a stack hasn't been created. The" +
					" wallet can't be used in this state.\n" +
					"\tEnabled: The wallet has been created and a stack has been created for the" +
					" wallet. This is the normal expected state of a wallet in use.\n" +
					"\tError: The wallet is inaccessible.",
				Computed: true,
			},
			schemaClumioAccountId: schema.StringAttribute{
				Description: "Identifier of the AWS account associated with Clumio. This " +
					"identifier is provided so that access to the service role for Clumio can be " +
					"restricted to just this account.",
				Computed: true,
			},
		},
	}
}
