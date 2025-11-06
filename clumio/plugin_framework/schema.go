// Copyright 2024. Clumio, Inc.

// This file holds the type definition and Schema resource function used by the provider model for
// the Clumio Terraform provider.

package clumio_pf

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// clumioProviderModel is the provider model for the Clumio Provider for Terraform. It holds the
// values required to make API calls to the Clumio backend. This includes the context for the Clumio
// organizational unit from which calls will be made.
type clumioProviderModel struct {
	ClumioApiToken                  types.String `tfsdk:"clumio_api_token"`
	ClumioApiBaseUrl                types.String `tfsdk:"clumio_api_base_url"`
	ClumioOrganizationalUnitContext types.String `tfsdk:"clumio_organizational_unit_context"`
}

// Schema defines the structure and constraints of the provider block for the Clumio Provider for
// Terraform. The atributes of the schema are mainly used to initialize an API client to the Clumio
// backend. However, all such attributes are marked "Optional" as the attributes can also be
// initialized using environment variables.
func (p *clumioProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "",
		Attributes: map[string]schema.Attribute{
			"clumio_api_token": schema.StringAttribute{
				MarkdownDescription: "The API token required to invoke Clumio APIs. " +
					"Informations for generating this token are available here: " +
					"https://documentation.commvault.com/clumio/api_tokens.html#manage-tokens",
				Optional: true,
			},
			"clumio_api_base_url": schema.StringAttribute{
				MarkdownDescription: "The base URL for Clumio APIs. The following are the valid " +
					"values for clumio_api_base_url. Use the appropriate value depending" +
					" on the region for which your credentials were created. " +
					"Below are the URLs to access the Clumio portal for each region and the corresponding API Base URLs:\n\n\t\t" +
					"Portal: https://west.portal.clumio.com/\n\n\t\t" +
					"API Base URL: https://us-west-2.api.clumio.com\n\n\t\t" +
					"Portal: https://east.portal.clumio.com/\n\n\t\t" +
					"API Base URL: https://us-east-1.api.clumio.com\n\n\t\t" +
					"Portal: https://canada.portal.clumio.com/\n\n\t\t" +
					"API Base URL:  https://ca-central-1.ca.api.clumio.com\n\n\t\t" +
					"Portal: https://eu1.portal.clumio.com/\n\n\t\t" +
					"API Base URL:  https://eu-central-1.de.api.clumio.com\n\n\t\t" +
					"Portal: https://au.portal.clumio.com/\n\n\t\t" +
					"API Base URL:  https://ap-southeast-2.au.api.clumio.com\n\n\t\t",
				Optional: true,
			},
			"clumio_organizational_unit_context": schema.StringAttribute{
				MarkdownDescription: "Organizational Unit context in which to create the" +
					" clumio resources. If not set, the resources will be created in" +
					" the context of the Global Organizational Unit. The value should" +
					" be the id of the Organizational Unit and not the name.",
				Optional: true,
			},
		},
	}
}
