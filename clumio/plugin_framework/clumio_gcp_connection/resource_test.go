// Copyright (c) 2025 Clumio, a Commvault Company All Rights Reserved

//go:build unit

package clumio_gcp_connection

import (
	"context"
	"testing"

	apiutils "github.com/clumio-code/clumio-go-sdk/api_utils"
	sdkconfig "github.com/clumio-code/clumio-go-sdk/config"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	apiError = &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte("Test Error"),
	}
	apiErrorStatusNotFound = &apiutils.APIError{
		ResponseCode: 404,
		Reason:       "Not Found",
		Response:     []byte("Test Error"),
	}
	controlPlaneId   = "controlPlaneId"
	controlPlaneRole = "controlPlaneRole"
	token            = "token"
)

// Unit test for the following cases:
//   - Create GCP connection success scenario.
//   - SDK API for create GCP connection returns an error.
func TestCreateGcpConnection(t *testing.T) {
	ctx := context.Background()
	mockSdkConnection := sdkclients.NewMockGcpConnectionClient(t)

	r := &clumioGCPConnectionResource{
		name: "test_clumio_post_process_gcp_connection",
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkConnections: mockSdkConnection,
	}

	model := &clumioGCPConnectionResourceModel{
		ProjectID:   types.StringValue("1234"),
		Description: types.StringValue("Description"),
	}

	req := &models.CreateGcpConnectionV1Request{
		Description: model.Description.ValueStringPointer(),
		ProjectId:   model.ProjectID.ValueStringPointer(),
	}

	resp := &models.CreateGCPConnectionResponse{
		Links:                 nil,
		ControlPlaneId:        &controlPlaneId,
		ControlPlaneRole:      &controlPlaneRole,
		Configuration:         nil,
		ConnectionStatus:      nil,
		ConnectionType:        nil,
		CreatedTimestamp:      nil,
		Description:           nil,
		OrganizationalUnitId:  nil,
		ProjectId:             nil,
		ProjectNumber:         nil,
		TemplatePermissionSet: nil,
		Token:                 &token,
		UpdatedTimestamp:      nil,
	}

	t.Run("Create GCP connection success scenario", func(t *testing.T) {
		mockSdkConnection.EXPECT().CreateGcpConnection(req).Times(1).
			Return(resp, nil)

		diags := r.createGcpConnection(ctx, model)

		assert.Equal(t, model.ClumioControlPlaneId.ValueString(), controlPlaneId)
		assert.Equal(t, model.ClumioControlPlaneRole.ValueString(), controlPlaneRole)
		assert.Equal(t, model.Token.ValueString(), token)

		assert.False(t, diags.HasError())

	})

	t.Run("SDK API for create GCP connection returns an error", func(t *testing.T) {
		mockSdkConnection.EXPECT().CreateGcpConnection(mock.Anything).Times(1).
			Return(nil, apiError)

		diags := r.createGcpConnection(ctx, model)
		assert.NotNil(t, diags)
	})
}

// Unit test for the following cases:
//   - Update GCP connection success scenario.
//   - SDK API for update GCP connection returns an error.
func TestUpdateGcpConnection(t *testing.T) {
	ctx := context.Background()
	mockSdkConnection := sdkclients.NewMockGcpConnectionClient(t)

	r := &clumioGCPConnectionResource{
		name: "test_clumio_post_process_gcp_connection",
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkConnections: mockSdkConnection,
	}

	model := &clumioGCPConnectionResourceModel{
		ProjectID:   types.StringValue("1234"),
		Description: types.StringValue("Description"),
	}

	req := &models.UpdateGcpConnectionV1Request{
		Description: model.Description.ValueStringPointer(),
	}

	t.Run("Update GCP connection success scenario", func(t *testing.T) {
		mockSdkConnection.EXPECT().UpdateGcpConnection(model.ProjectID.ValueString(), req).Times(1).
			Return(&models.UpdateGCPConnectionResponse{}, nil)

		diags := r.updateGcpConnection(ctx, model)

		assert.False(t, diags.HasError())

	})

	t.Run("SDK API for update GCP connection returns an error", func(t *testing.T) {
		mockSdkConnection.EXPECT().UpdateGcpConnection(mock.Anything, mock.Anything).Times(1).
			Return(nil, apiError)

		diags := r.updateGcpConnection(ctx, model)
		assert.NotNil(t, diags)
	})
}

// Unit test for the following cases:
//   - Delete GCP connection success scenario.
//   - SDK API for delete GCP connection returns an error.
func TestDeleteGcpConnection(t *testing.T) {
	ctx := context.Background()
	mockSdkConnection := sdkclients.NewMockGcpConnectionClient(t)

	r := &clumioGCPConnectionResource{
		name: "test_clumio_post_process_gcp_connection",
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkConnections: mockSdkConnection,
	}

	model := &clumioGCPConnectionResourceModel{
		ProjectID: types.StringValue("1234"),
	}

	t.Run("Delete GCP connection success scenario", func(t *testing.T) {
		mockSdkConnection.EXPECT().DeleteGcpConnection(model.ProjectID.ValueString()).Times(1).
			Return(nil, nil)

		diags := r.deleteGcpConnection(ctx, model)

		assert.False(t, diags.HasError())

	})

	t.Run("SDK API for delete GCP connection returns an error", func(t *testing.T) {
		mockSdkConnection.EXPECT().DeleteGcpConnection(mock.Anything).Times(1).
			Return(nil, apiError)

		diags := r.deleteGcpConnection(ctx, model)
		assert.NotNil(t, diags)
	})
}

// Unit test for the following cases:
//   - Read GCP connection success scenario.
//   - SDK API for read GCP connection returns an error
//   - SDK API not found error return remove bool as true
func TestReadGcpConnection(t *testing.T) {
	ctx := context.Background()
	mockSdkConnection := sdkclients.NewMockGcpConnectionClient(t)

	r := &clumioGCPConnectionResource{
		name: "test_clumio_post_process_gcp_connection",
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkConnections: mockSdkConnection,
	}

	model := &clumioGCPConnectionResourceModel{
		ProjectID:              types.StringValue("1234"),
		ClumioControlPlaneId:   types.StringValue(""),
		ClumioControlPlaneRole: types.StringValue(""),
	}

	res := &models.ReadGCPConnectionResponse{
		Token:            types.StringValue("token").ValueStringPointer(),
		ControlPlaneId:   types.StringValue("controlPlaneId").ValueStringPointer(),
		ControlPlaneRole: types.StringValue("controlPlaneRole").ValueStringPointer(),
	}

	t.Run("Read GCP connection success scenario", func(t *testing.T) {
		mockSdkConnection.EXPECT().ReadGcpConnection(model.ProjectID.ValueString()).Times(1).
			Return(res, nil)

		remove, diags := r.readGcpConnection(ctx, model)
		assert.False(t, diags.HasError())
		assert.False(t, remove)
		assert.Equal(t, model.Token.ValueString(), *res.Token)
		assert.Equal(t, model.ClumioControlPlaneId.ValueString(), *res.ControlPlaneId)
		assert.Equal(t, model.ClumioControlPlaneRole.ValueString(), *res.ControlPlaneRole)
	})

	t.Run("SDK API for read GCP connection returns an error", func(t *testing.T) {
		mockSdkConnection.EXPECT().ReadGcpConnection(mock.Anything).Times(1).
			Return(nil, apiError)

		remove, diags := r.readGcpConnection(ctx, model)
		assert.NotNil(t, diags)
		assert.False(t, remove)
	})

	t.Run("SDK API not found error return remove bool as true", func(t *testing.T) {
		mockSdkConnection.EXPECT().ReadGcpConnection(mock.Anything).Times(1).
			Return(nil, apiErrorStatusNotFound)

		remove, diags := r.readGcpConnection(ctx, model)
		assert.False(t, diags.HasError())
		assert.True(t, remove)
	})
}
