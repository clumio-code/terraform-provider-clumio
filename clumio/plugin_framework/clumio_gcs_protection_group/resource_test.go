// Copyright 2026. Clumio, Inc.

// This file contains the unit tests for the functions in resource.go

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

// newTestModel returns a resource model populated with the minimum required inputs, including a
// structured bucket_rule matching on gcp_project_id.
func newTestModel(id, name string) *clumioGCSProtectionGroupResourceModel {
	return &clumioGCSProtectionGroupResourceModel{
		ID:   types.StringValue(id),
		Name: types.StringValue(name),
		BucketRule: &bucketRuleModel{
			GcpProjectId: &gcpStringOperatorModel{
				Eq:    types.StringValue("my-project"),
				In:    types.SetNull(types.StringType),
				NotEq: types.StringNull(),
				NotIn: types.SetNull(types.StringType),
			},
		},
		IncludePrefixes:   types.SetNull(types.StringType),
		ExcludePrefixes:   types.SetNull(types.StringType),
		LatestVersionOnly: types.BoolValue(true),
	}
}

// Unit tests for createProtectionGroup: success (with read-after-write), API error, and nil
// response.
func TestCreateGCSProtectionGroup(t *testing.T) {
	ctx := context.Background()
	mockClient := sdkclients.NewMockGcpProtectionGroupClient(t)
	id := "test-pg-id"
	name := "test-pg"
	r := clumioGCSProtectionGroupResource{
		name:                "test_gcs_protection_group",
		client:              &common.ApiClient{ClumioConfig: sdkconfig.Config{}},
		sdkProtectionGroups: mockClient,
	}
	apiError := &apiutils.APIError{ResponseCode: 500, Reason: "test", Response: []byte("err")}

	t.Run("success", func(t *testing.T) {
		plan := newTestModel("", name)
		mockClient.EXPECT().CreateGcpProtectionGroup(mock.Anything).Times(1).Return(
			&models.CreateGCPProtectionGroupResponse{Id: &id, Name: &name}, nil)
		// Create reads back the protection group to populate computed fields.
		mockClient.EXPECT().ReadGcpProtectionGroup(id, mock.Anything).Times(1).Return(
			&models.ReadGCPProtectionGroupResponse{Id: &id, Name: &name}, nil)
		diags := r.createProtectionGroup(ctx, plan)
		assert.Nil(t, diags)
		assert.Equal(t, id, plan.ID.ValueString())
	})

	t.Run("api error", func(t *testing.T) {
		plan := newTestModel("", name)
		mockClient.EXPECT().CreateGcpProtectionGroup(mock.Anything).Times(1).Return(nil, apiError)
		diags := r.createProtectionGroup(ctx, plan)
		assert.NotNil(t, diags)
	})

	t.Run("nil response", func(t *testing.T) {
		plan := newTestModel("", name)
		mockClient.EXPECT().CreateGcpProtectionGroup(mock.Anything).Times(1).Return(nil, nil)
		diags := r.createProtectionGroup(ctx, plan)
		assert.NotNil(t, diags)
	})
}

// Unit tests for readProtectionGroup: success, not found (404), externally deleted, and API error.
func TestReadGCSProtectionGroup(t *testing.T) {
	ctx := context.Background()
	mockClient := sdkclients.NewMockGcpProtectionGroupClient(t)
	id := "test-pg-id"
	name := "test-pg"
	r := clumioGCSProtectionGroupResource{
		name:                "test_gcs_protection_group",
		client:              &common.ApiClient{ClumioConfig: sdkconfig.Config{}},
		sdkProtectionGroups: mockClient,
	}

	t.Run("success", func(t *testing.T) {
		state := newTestModel(id, name)
		mockClient.EXPECT().ReadGcpProtectionGroup(id, mock.Anything).Times(1).Return(
			&models.ReadGCPProtectionGroupResponse{Id: &id, Name: &name}, nil)
		remove, diags := r.readProtectionGroup(ctx, state)
		assert.Nil(t, diags)
		assert.False(t, remove)
	})

	t.Run("not found removes from state", func(t *testing.T) {
		state := newTestModel(id, name)
		mockClient.EXPECT().ReadGcpProtectionGroup(id, mock.Anything).Times(1).Return(
			nil, &apiutils.APIError{ResponseCode: 404})
		remove, diags := r.readProtectionGroup(ctx, state)
		assert.Nil(t, diags)
		assert.True(t, remove)
	})

	t.Run("externally deleted removes from state", func(t *testing.T) {
		state := newTestModel(id, name)
		isDeleted := true
		mockClient.EXPECT().ReadGcpProtectionGroup(id, mock.Anything).Times(1).Return(
			&models.ReadGCPProtectionGroupResponse{Id: &id, IsDeleted: &isDeleted}, nil)
		remove, diags := r.readProtectionGroup(ctx, state)
		assert.Nil(t, diags)
		assert.True(t, remove)
	})

	t.Run("api error", func(t *testing.T) {
		state := newTestModel(id, name)
		mockClient.EXPECT().ReadGcpProtectionGroup(id, mock.Anything).Times(1).Return(
			nil, &apiutils.APIError{ResponseCode: 500})
		remove, diags := r.readProtectionGroup(ctx, state)
		assert.NotNil(t, diags)
		assert.False(t, remove)
	})
}

// Unit tests for updateProtectionGroup: success (with read-after-write) and API error.
func TestUpdateGCSProtectionGroup(t *testing.T) {
	ctx := context.Background()
	mockClient := sdkclients.NewMockGcpProtectionGroupClient(t)
	id := "test-pg-id"
	name := "test-pg"
	r := clumioGCSProtectionGroupResource{
		name:                "test_gcs_protection_group",
		client:              &common.ApiClient{ClumioConfig: sdkconfig.Config{}},
		sdkProtectionGroups: mockClient,
	}

	t.Run("success", func(t *testing.T) {
		plan := newTestModel(id, name)
		mockClient.EXPECT().UpdateGcpProtectionGroup(id, mock.Anything).Times(1).Return(
			&models.UpdateGCPProtectionGroupResponse{Id: &id, Name: &name}, nil)
		mockClient.EXPECT().ReadGcpProtectionGroup(id, mock.Anything).Times(1).Return(
			&models.ReadGCPProtectionGroupResponse{Id: &id, Name: &name}, nil)
		diags := r.updateProtectionGroup(ctx, plan)
		assert.Nil(t, diags)
	})

	t.Run("api error", func(t *testing.T) {
		plan := newTestModel(id, name)
		mockClient.EXPECT().UpdateGcpProtectionGroup(id, mock.Anything).Times(1).Return(
			nil, &apiutils.APIError{ResponseCode: 500})
		diags := r.updateProtectionGroup(ctx, plan)
		assert.NotNil(t, diags)
	})
}

// Unit tests for deleteProtectionGroup: success, already deleted (404), and API error.
func TestDeleteGCSProtectionGroup(t *testing.T) {
	ctx := context.Background()
	mockClient := sdkclients.NewMockGcpProtectionGroupClient(t)
	id := "test-pg-id"
	r := clumioGCSProtectionGroupResource{
		name:                "test_gcs_protection_group",
		client:              &common.ApiClient{ClumioConfig: sdkconfig.Config{}},
		sdkProtectionGroups: mockClient,
	}

	t.Run("success", func(t *testing.T) {
		state := newTestModel(id, "test-pg")
		mockClient.EXPECT().DeleteGcpProtectionGroup(id).Times(1).Return(nil, nil)
		diags := r.deleteProtectionGroup(ctx, state)
		assert.Nil(t, diags)
	})

	t.Run("already deleted", func(t *testing.T) {
		state := newTestModel(id, "test-pg")
		mockClient.EXPECT().DeleteGcpProtectionGroup(id).Times(1).Return(
			nil, &apiutils.APIError{ResponseCode: 404})
		diags := r.deleteProtectionGroup(ctx, state)
		assert.Nil(t, diags)
	})

	t.Run("api error", func(t *testing.T) {
		state := newTestModel(id, "test-pg")
		mockClient.EXPECT().DeleteGcpProtectionGroup(id).Times(1).Return(
			nil, &apiutils.APIError{ResponseCode: 500})
		diags := r.deleteProtectionGroup(ctx, state)
		assert.NotNil(t, diags)
	})
}
