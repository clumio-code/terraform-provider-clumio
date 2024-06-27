// Copyright 2024. Clumio, Inc.

// This file contains the unit test for the Schema function in schema.go.

//go:build unit

package clumio_pf

import (
	"context"
	"testing"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
)

// Unit test for the provider Metadata function.
func TestProviderMetadata(t *testing.T) {

	ctx := context.Background()
	clumioProvider := New()
	providerName := "clumio"

	// Tests that the provider name is set as part of the Metadata().
	metadataResp := &provider.MetadataResponse{}
	clumioProvider.Metadata(ctx, provider.MetadataRequest{}, metadataResp)

	assert.Equal(t, providerName, metadataResp.TypeName)
}

// Unit test for the following Provider configure scenarios:
//   - Success scenario for provider configure.
//   - clumio_api_base_url is empty in the configure request.
//   - clumio_api_token is empty in the configure request.
func TestProviderConfigure(t *testing.T) {

	ctx := context.Background()
	clumioProvider := New()
	token := "test-token"
	baseUrl := "test-base-url"
	ou := "test-ou"
	apiTokenKey := "clumio_api_token"
	apiBaseUrlKey := "clumio_api_base_url"
	ouContextKey := "clumio_organizational_unit_context"

	mapType := tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			apiBaseUrlKey: tftypes.String,
			apiTokenKey:   tftypes.String,
			ouContextKey:  tftypes.String,
		},
		OptionalAttributes: nil,
	}
	vals := make(map[string]tftypes.Value, 0)
	vals[apiBaseUrlKey] = tftypes.NewValue(tftypes.String, baseUrl)
	vals[apiTokenKey] = tftypes.NewValue(tftypes.String, token)
	vals[ouContextKey] = tftypes.NewValue(tftypes.String, ou)

	// Success scenario for provider configure
	t.Run("Success scenario for provider configure", func(t *testing.T) {

		configResp := &provider.ConfigureResponse{}
		resp := &provider.SchemaResponse{}
		clumioProvider.Schema(context.Background(), provider.SchemaRequest{}, resp)
		clumioProvider.Configure(ctx, provider.ConfigureRequest{
			Config: tfsdk.Config{
				Raw:    tftypes.NewValue(mapType, vals),
				Schema: resp.Schema,
			},
		}, configResp)

		assert.Equal(t, baseUrl, configResp.ResourceData.(*common.ApiClient).ClumioConfig.BaseUrl)
		assert.Equal(t, token, configResp.ResourceData.(*common.ApiClient).ClumioConfig.Token)
		assert.Equal(t, ou,
			configResp.ResourceData.(*common.ApiClient).ClumioConfig.OrganizationalUnitContext)

	})

	// Tests that diagnostics is returned when clumio_api_base_url is empty.
	t.Run("Error when clumio_api_base_url is empty", func(t *testing.T) {
		configResp := &provider.ConfigureResponse{}
		resp := &provider.SchemaResponse{}
		clumioProvider.Schema(context.Background(), provider.SchemaRequest{}, resp)

		vals[apiBaseUrlKey] = tftypes.NewValue(tftypes.String, "")
		clumioProvider.Configure(ctx, provider.ConfigureRequest{
			Config: tfsdk.Config{
				Raw:    tftypes.NewValue(mapType, vals),
				Schema: resp.Schema,
			},
		}, configResp)

		assert.True(t, configResp.Diagnostics.HasError())

		//Reset the base url at the end of test.
		vals[apiBaseUrlKey] = tftypes.NewValue(tftypes.String, "")
	})

	// Tests that diagnostics is returned when clumio_api_token is empty.
	t.Run("Error when clumio_api_token is empty", func(t *testing.T) {
		configResp := &provider.ConfigureResponse{}
		resp := &provider.SchemaResponse{}
		clumioProvider.Schema(context.Background(), provider.SchemaRequest{}, resp)

		vals[apiBaseUrlKey] = tftypes.NewValue(tftypes.String, baseUrl)
		vals[apiTokenKey] = tftypes.NewValue(tftypes.String, "")
		clumioProvider.Configure(ctx, provider.ConfigureRequest{
			Config: tfsdk.Config{
				Raw:    tftypes.NewValue(mapType, vals),
				Schema: resp.Schema,
			},
		}, configResp)

		assert.True(t, configResp.Diagnostics.HasError())
	})

}

// Unit test for the provider Resources function.
func TestProviderResources(t *testing.T) {

	ctx := context.Background()
	clumioProvider := New()

	resp := clumioProvider.Resources(ctx)
	assert.Equal(t, 15, len(resp))
}

// Unit test for the provider DataSources function.
func TestProviderDataSources(t *testing.T) {

	ctx := context.Background()
	clumioProvider := New()

	resp := clumioProvider.DataSources(ctx)
	assert.Equal(t, 10, len(resp))
}
