// Copyright 2026. Clumio, Inc.

// This file contains the unit tests for the functions in data_source.go

//go:build unit

package clumio_gcs_protection_group

import (
	"context"
	"testing"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	apiutils "github.com/clumio-code/clumio-go-sdk/api_utils"
	sdkconfig "github.com/clumio-code/clumio-go-sdk/config"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Unit tests for the data source readProtectionGroup: success, API error, nil response, and
// not-found.
func TestDatasourceReadGCSProtectionGroup(t *testing.T) {
	ctx := context.Background()
	mockClient := sdkclients.NewMockGcpProtectionGroupClient(t)
	id := "test-pg-id"
	name := "test-pg"
	r := clumioGCSProtectionGroupDataSource{
		name:                "test_gcs_protection_group",
		client:              &common.ApiClient{ClumioConfig: sdkconfig.Config{}},
		sdkProtectionGroups: mockClient,
	}
	model := &clumioGCSProtectionGroupDataSourceModel{Name: types.StringValue(name)}

	t.Run("success", func(t *testing.T) {
		count := int64(1)
		mockClient.EXPECT().ListGcpProtectionGroups(
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Times(1).Return(&models.ListGCPProtectionGroupsResponse{
			CurrentCount: &count,
			Embedded: &models.GCPProtectionGroupListEmbedded{
				Items: []*models.GCPProtectionGroup{{Id: &id, Name: &name}},
			},
		}, nil)
		diags := r.readProtectionGroup(ctx, model)
		assert.Nil(t, diags)
		assert.Equal(t, id, model.Id.ValueString())
	})

	t.Run("api error", func(t *testing.T) {
		mockClient.EXPECT().ListGcpProtectionGroups(
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Times(1).Return(nil, &apiutils.APIError{ResponseCode: 500})
		diags := r.readProtectionGroup(ctx, model)
		assert.NotNil(t, diags)
	})

	t.Run("nil response", func(t *testing.T) {
		mockClient.EXPECT().ListGcpProtectionGroups(
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Times(1).Return(nil, nil)
		diags := r.readProtectionGroup(ctx, model)
		assert.NotNil(t, diags)
	})

	t.Run("not found", func(t *testing.T) {
		count := int64(0)
		mockClient.EXPECT().ListGcpProtectionGroups(
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Times(1).Return(&models.ListGCPProtectionGroupsResponse{CurrentCount: &count}, nil)
		diags := r.readProtectionGroup(ctx, model)
		assert.NotNil(t, diags)
	})
}
