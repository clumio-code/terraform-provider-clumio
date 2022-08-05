package clumio_wallet

import (
	"context"

	"github.com/clumio-code/clumio-go-sdk/models"

	wallets "github.com/clumio-code/clumio-go-sdk/controllers/wallets"
	"github.com/clumio-code/terraform-provider-clumio/clumio/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ClumioWallet returns the resource for Clumio Wallet.

func ClumioWallet() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Clumio Wallet Resource to create and manage wallets in Clumio. " +
			"Wallets should be created \"after\" connecting an AWS account to Clumio.<br>" +
			"**NOTE:** To protect against accidental deletion, wallets cannot be destroyed. " +
			"To remove a wallet, contact Clumio support.",

		CreateContext: clumioWalletCreate,
		ReadContext:   clumioWalletRead,
		DeleteContext: clumioWalletDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			schemaId: {
				Description: "Wallet Id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			schemaAccountNativeId: {
				Description: "The AWS account id to be associated with the wallet.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			schemaToken: {
				Description: "Token is used to identify and authenticate the CloudFormation stack creation.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			schemaState: {
				Description: "State describes the state of the wallet. Valid states are:" +
					" Waiting: The wallet has been created, but a stack hasn't been" +
					" created. The wallet can't be used in this state. Enabled: The" +
					" wallet has been created and a stack has been created for the" +
					" wallet. This is the normal expected state of a wallet in use." +
					" Error: The wallet is inaccessible. See ErrorCode and ErrorMessage" +
					" fields for additional details.",
				Type:     schema.TypeString,
				Computed: true,
			},
			schemaInstalledRegions: {
				Type:     schema.TypeSet,
				Set:      common.SchemaSetHashString,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The regions where the wallet is installed.",
			},
			schemaSupportedRegions: {
				Type:     schema.TypeSet,
				Set:      common.SchemaSetHashString,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The supported regions for the wallet.",
			},
			schemaClumioAccountId: {
				Type:        schema.TypeString,
				Description: "Clumio Account ID.",
				Computed:    true,
			},
			schemaClumioControlPlaneAccountId: {
				Type:        schema.TypeString,
				Description: "Clumio Control Plane Account ID.",
				Computed:    true,
			},
		},
	}
}

// clumioWalletCreate handles the Create action for the ClumioWallet resource.
func clumioWalletCreate(
	ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.ApiClient)
	walletsAPI := wallets.NewWalletsV1(client.ClumioConfig)
	accountNativeId := common.GetStringValue(d, schemaAccountNativeId)
	apiRes, apiErr := walletsAPI.CreateWallet(&models.CreateWalletV1Request{
		AccountNativeId: &accountNativeId,
	})
	if apiErr != nil {
		return diag.Errorf("Error creating wallet. Error: %v", apiErr)
	}
	d.SetId(*apiRes.Id)
	return clumioWalletRead(ctx, d, meta)
}

// clumioWalletRead handles the Read action for the ClumioWallet resource.
func clumioWalletRead(
	_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.ApiClient)
	walletsAPI := wallets.NewWalletsV1(client.ClumioConfig)
	apiRes, apiErr := walletsAPI.ReadWallet(d.Id())
	if apiErr != nil {
		return diag.Errorf("Error reading wallet: %v. Error: %v", d.Id(), apiErr)
	}
	err := d.Set(schemaToken, *apiRes.Token)
	if err != nil {
		return diag.Errorf(common.SchemaAttributeSetError, schemaToken, err)
	}
	err = d.Set(schemaState, *apiRes.State)
	if err != nil {
		return diag.Errorf(common.SchemaAttributeSetError, schemaState, err)
	}
	err = d.Set(schemaClumioAccountId, *apiRes.ClumioAwsAccountId)
	if err != nil {
		return diag.Errorf(
			common.SchemaAttributeSetError, schemaClumioAccountId, err)
	}
	err = d.Set(schemaClumioControlPlaneAccountId, *apiRes.ClumioControlPlaneAwsAccountId)
	if err != nil {
		return diag.Errorf(
			common.SchemaAttributeSetError, schemaClumioControlPlaneAccountId, err)
	}
	if apiRes.InstalledRegions != nil && len(apiRes.InstalledRegions) > 0 {
		installedRegions := make([]string, 0)
		for _, installedRegion := range apiRes.InstalledRegions {
			installedRegions = append(installedRegions, *installedRegion)
		}
		err = d.Set(schemaInstalledRegions, installedRegions)
		if err != nil {
			return diag.Errorf(common.SchemaAttributeSetError, schemaInstalledRegions, err)
		}
	}
	if apiRes.SupportedRegions != nil && len(apiRes.SupportedRegions) > 0 {
		SupportedRegions := make([]string, 0)
		for _, installedRegion := range apiRes.SupportedRegions {
			SupportedRegions = append(SupportedRegions, *installedRegion)
		}
		err = d.Set(schemaSupportedRegions, SupportedRegions)
		if err != nil {
			return diag.Errorf(common.SchemaAttributeSetError, schemaSupportedRegions, err)
		}
	}
	return nil
}

// clumioWalletDelete handles the Delete action for the ClumioWallet resource.
func clumioWalletDelete(
	_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.ApiClient)
	walletsAPI := wallets.NewWalletsV1(client.ClumioConfig)
	_, apiErr := walletsAPI.DeleteWallet(d.Id())
	if apiErr != nil {
		return diag.Errorf("To protect against accidental deletion, wallets cannot be destroyed. "+
			"To remove a wallet, contact Clumio support: %v.", d.Id())
	}
	return nil
}
