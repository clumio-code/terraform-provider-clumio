// Copyright (c) 2025 Clumio, a Commvault Company All Rights Reserved

//go:build unit

package clumio_post_process_gcp_connection

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	apiutils "github.com/clumio-code/clumio-go-sdk/api_utils"
	sdkconfig "github.com/clumio-code/clumio-go-sdk/config"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"
)

var (
	apiError = &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte("Test Error"),
	}
	PropertyKey   = "key"
	PropertyValue = "value"
)

func setupTestModel(t *testing.T) *clumioPostProcessGCPConnectionResourceModel {
	props, diags := basetypes.NewMapValueFrom(context.Background(), types.StringType, map[string]attr.Value{
		PropertyKey: basetypes.NewStringValue(PropertyValue),
	})
	assert.Nil(t, diags)

	return &clumioPostProcessGCPConnectionResourceModel{
		ProjectID:           basetypes.NewStringValue("ProjectId"),
		ProjectName:         basetypes.NewStringValue("ProjectName"),
		ProjectNumber:       basetypes.NewStringValue("ProjectNumber"),
		Token:               basetypes.NewStringValue("Token"),
		ServiceAccountEmail: basetypes.NewStringValue("ServiceAccountEmail"),
		WifPoolId:           basetypes.NewStringValue("WifPoolId"),
		WifProviderId:       basetypes.NewStringValue("WifProviderId"),
		ConfigVersion:       basetypes.NewStringValue("1.1"),
		ProtectGcsVersion:   basetypes.NewStringValue("1.1"),
		Properties:          props,
	}
}

// Unit test for the following cases:
//   - Post-process GCP connection success scenario.
//   - Get template configuration returns a version error.
//   - Get template configuration returns a marshal error.
//   - SDK API for post-process GCP connection returns an error.
func TestCreateUpdatePostProcessGcpConnection(t *testing.T) {
	ctx := context.Background()
	mockSdkConnection := sdkclients.NewMockGcpConnectionClient(t)
	r := &clumioPostProcessGCPConnectionResource{
		name: "test_clumio_post_process_gcp_connection",
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkConnections: mockSdkConnection,
	}

	model := setupTestModel(t)
	templateConfig, err := GetTemplateConfiguration(model)
	assert.Nil(t, err)

	configBytes, err := json.Marshal(templateConfig)
	assert.Nil(t, err)
	configuration := string(configBytes)

	req := &models.PostProcessGcpConnectionV1Request{
		Configuration: &configuration,
		ProjectId:     model.ProjectID.ValueStringPointer(),
		ProjectName:   model.ProjectName.ValueStringPointer(),
		ProjectNumber: model.ProjectNumber.ValueStringPointer(),
		RequestType:   types.StringValue(createRequestType).ValueStringPointer(),
		ResourceProperties: map[string]*string{
			PropertyKey: &PropertyValue,
		},
		ServiceAccountEmail: model.ServiceAccountEmail.ValueStringPointer(),
		Token:               model.Token.ValueStringPointer(),
		WifPoolId:           model.WifPoolId.ValueStringPointer(),
		WifProviderId:       model.WifProviderId.ValueStringPointer(),
	}

	t.Run("Success scenario for post-process create", func(t *testing.T) {
		//Setup expectations.
		mockSdkConnection.EXPECT().PostProcessGcpConnection(req).Times(1).
			Return(nil, nil)

		diags := r.createUpdatePostProcessGcpConnection(ctx, model, createRequestType)
		assert.Nil(t, diags)
	})

	t.Run("Success scenario for post-process update", func(t *testing.T) {
		//Setup expectations.
		req.RequestType = types.StringValue(updateRequestType).ValueStringPointer()
		mockSdkConnection.EXPECT().PostProcessGcpConnection(req).Times(1).
			Return(nil, nil)

		diags := r.createUpdatePostProcessGcpConnection(ctx, model, updateRequestType)
		assert.Nil(t, diags)
	})

	t.Run("GetTemplateConfiguration returns a version error", func(t *testing.T) {
		modelWithInvalidVersion := &clumioPostProcessGCPConnectionResourceModel{
			ConfigVersion: basetypes.NewStringValue("1.2.3"),
		}

		diags := r.createUpdatePostProcessGcpConnection(ctx, modelWithInvalidVersion, updateRequestType)
		assert.NotNil(t, diags)
		assert.Equal(t, "Unable to create template configurations from versions: invalid version 1.2.3", diags.Errors()[0].Detail())
		assert.Equal(t, "Error in invoking Post-process Clumio GCP Connection.", diags.Errors()[0].Summary())
	})

	t.Run("GetTemplateConfiguration returns a marshal error", func(t *testing.T) {
		// Arrange: force json marshal to fail
		old := jsonMarshal
		jsonMarshal = func(v any) ([]byte, error) {
			return nil, fmt.Errorf("forced marshal error")
		}
		defer func() { jsonMarshal = old }()

		model := &clumioPostProcessGCPConnectionResourceModel{
			ConfigVersion: basetypes.NewStringValue("1.2"),
		}

		diags := r.createUpdatePostProcessGcpConnection(ctx, model, updateRequestType)
		assert.NotNil(t, diags)
		assert.Equal(t, "forced marshal error", diags.Errors()[0].Detail())
		assert.Equal(t, "Unable to marshal template configuration", diags.Errors()[0].Summary())
	})

	t.Run("PostProcessGcpConnection returns an error", func(t *testing.T) {
		// Setup expectations
		mockSdkConnection.EXPECT().PostProcessGcpConnection(mock.Anything).Times(1).
			Return(nil, apiError)

		diags := r.createUpdatePostProcessGcpConnection(ctx, model, updateRequestType)
		assert.NotNil(t, diags)
	})

}

// Unit test for the following cases:
//   - delete GCP connection success scenario.
//   - SDK API for post-process GCP connection returns an error.
func TestDeletePostProcessGcpConnection(t *testing.T) {
	ctx := context.Background()
	mockSdkConnection := sdkclients.NewMockGcpConnectionClient(t)
	r := &clumioPostProcessGCPConnectionResource{
		name: "test_clumio_post_process_gcp_connection",
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkConnections: mockSdkConnection,
	}

	model := setupTestModel(t)
	req := &models.PostProcessGcpConnectionV1Request{
		ProjectId:     model.ProjectID.ValueStringPointer(),
		ProjectName:   model.ProjectName.ValueStringPointer(),
		ProjectNumber: model.ProjectNumber.ValueStringPointer(),
		RequestType:   types.StringValue(deleteRequestType).ValueStringPointer(),
		Token:         model.Token.ValueStringPointer(),
	}

	t.Run("Success scenario for post-process delete", func(t *testing.T) {
		//Setup expectations.
		mockSdkConnection.EXPECT().PostProcessGcpConnection(req).Times(1).
			Return(nil, nil)

		diags := r.deletePostProcessGcpConnection(ctx, model)
		assert.Nil(t, diags)
	})

	t.Run("PostProcessGcpConnection returns an error", func(t *testing.T) {
		// Setup expectations
		mockSdkConnection.EXPECT().PostProcessGcpConnection(mock.Anything).Times(1).
			Return(nil, apiError)

		diags := r.deletePostProcessGcpConnection(ctx, model)
		assert.NotNil(t, diags)
	})

}
