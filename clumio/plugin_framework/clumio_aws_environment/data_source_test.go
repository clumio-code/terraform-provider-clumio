// Copyright 2026. Clumio, Inc.

// This file contains the unit tests for the functions in data_source.go

//go:build unit

package clumio_aws_environment

import (
	"context"
	"testing"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	apiutils "github.com/clumio-code/clumio-go-sdk/api_utils"
	sdkconfig "github.com/clumio-code/clumio-go-sdk/config"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Unit test for the following cases:
//   - Read AWS environment success scenario populates the environment ID.
//   - The environment lookup returns an error.
//
// The detailed lookup branches (nil response, empty items, more than one match, missing ID) are
// covered by common.TestLookupAWSEnvironment.
func TestDatasourceReadAWSEnvironment(t *testing.T) {

	ctx := context.Background()
	envClient := sdkclients.NewMockAWSEnvironmentClient(t)
	envId := "test-env-id"

	rds := clumioAWSEnvironmentDataSource{
		name:      "test_aws_environment",
		client:    &common.ApiClient{ClumioConfig: sdkconfig.Config{}},
		envClient: envClient,
	}
	rdsm := &clumioAWSEnvironmentDataSourceModel{
		AccountNativeID: basetypes.NewStringValue("test-account-native-id"),
		AWSRegion:       basetypes.NewStringValue("test-region"),
	}

	// Tests the success scenario for AWS environment read. It should not return Diagnostics and
	// the ID in the model should be populated from the API response.
	t.Run("Basic success scenario for read AWS environment", func(t *testing.T) {

		readResponse := &models.ListAWSEnvironmentsResponse{
			Embedded: &models.AWSEnvironmentListEmbedded{
				Items: []*models.AWSEnvironment{{Id: &envId}},
			},
		}
		envClient.EXPECT().ListAwsEnvironments(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Times(1).Return(readResponse, nil)

		diags := rds.readAWSEnvironment(ctx, rdsm)
		assert.Nil(t, diags)
		assert.Equal(t, envId, rdsm.ID.ValueString())
	})

	// Tests that Diagnostics is returned in case the environment lookup fails.
	t.Run("read AWS environment returns an error", func(t *testing.T) {

		envClient.EXPECT().ListAwsEnvironments(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Times(1).Return(nil, &apiutils.APIError{ResponseCode: 500})

		diags := rds.readAWSEnvironment(ctx, rdsm)
		assert.NotNil(t, diags)
	})
}

// Unit test for the datasource Metadata, Configure and Read functions.
func TestDatasourceMetadataConfigureRead(t *testing.T) {

	ctx := context.Background()
	envId := "test-env-id"

	// Tests that the datasource type name is set as part of Metadata().
	t.Run("Metadata sets the datasource type name", func(t *testing.T) {

		ds := &clumioAWSEnvironmentDataSource{}
		resp := &datasource.MetadataResponse{}
		ds.Metadata(ctx, datasource.MetadataRequest{ProviderTypeName: "clumio"}, resp)
		assert.Equal(t, "clumio_aws_environment", resp.TypeName)
	})

	// Tests that Configure() returns early when no provider data is given and sets up the SDK
	// client when it is.
	t.Run("Configure sets up the SDK client", func(t *testing.T) {

		ds := &clumioAWSEnvironmentDataSource{}
		ds.Configure(ctx, datasource.ConfigureRequest{}, &datasource.ConfigureResponse{})
		assert.Nil(t, ds.client)

		ds.Configure(ctx, datasource.ConfigureRequest{
			ProviderData: &common.ApiClient{ClumioConfig: sdkconfig.Config{}},
		}, &datasource.ConfigureResponse{})
		assert.NotNil(t, ds.client)
		assert.NotNil(t, ds.envClient)
	})

	// Tests the success scenario for the datasource Read(). It should not return Diagnostics and
	// should set the environment ID in the Terraform state.
	t.Run("Read sets the state from the API response", func(t *testing.T) {

		envClient := sdkclients.NewMockAWSEnvironmentClient(t)
		ds := &clumioAWSEnvironmentDataSource{name: "clumio_aws_environment", envClient: envClient}

		configType := tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				schemaId:              tftypes.String,
				schemaAccountNativeId: tftypes.String,
				schemaAwsRegion:       tftypes.String,
			},
		}
		configVals := map[string]tftypes.Value{
			schemaId:              tftypes.NewValue(tftypes.String, nil),
			schemaAccountNativeId: tftypes.NewValue(tftypes.String, "test-account-native-id"),
			schemaAwsRegion:       tftypes.NewValue(tftypes.String, "test-region"),
		}
		schemaResp := &datasource.SchemaResponse{}
		ds.Schema(ctx, datasource.SchemaRequest{}, schemaResp)

		envClient.EXPECT().ListAwsEnvironments(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Times(1).Return(&models.ListAWSEnvironmentsResponse{
			Embedded: &models.AWSEnvironmentListEmbedded{
				Items: []*models.AWSEnvironment{{Id: &envId}},
			},
		}, nil)

		resp := &datasource.ReadResponse{
			State: tfsdk.State{Raw: tftypes.NewValue(configType, nil), Schema: schemaResp.Schema},
		}
		ds.Read(ctx, datasource.ReadRequest{
			Config: tfsdk.Config{
				Raw:    tftypes.NewValue(configType, configVals),
				Schema: schemaResp.Schema,
			},
		}, resp)
		assert.False(t, resp.Diagnostics.HasError())
		assert.False(t, resp.State.Raw.IsNull())
	})
}
