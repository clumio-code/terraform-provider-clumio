// Copyright 2026. Clumio, Inc.

//go:build unit

package clumio_aws_connection

import (
	"context"
	"testing"

	sdkconfig "github.com/clumio-code/clumio-go-sdk/config"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
)

// TestModifyPlanSetsOrganizationalUnitID verifies that ModifyPlan pins organizational_unit_id to
// the provider's OU context when one is set.
func TestModifyPlanSetsOrganizationalUnitID(t *testing.T) {
	ctx := context.Background()
	res := NewClumioAWSConnectionResource().(*clumioAWSConnectionResource)
	res.Configure(ctx, resource.ConfigureRequest{
		ProviderData: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{
				OrganizationalUnitContext: ouId,
			},
		},
	}, &resource.ConfigureResponse{})

	schemaResp := &resource.SchemaResponse{}
	res.Schema(ctx, resource.SchemaRequest{}, schemaResp)

	// organizational_unit_id is left unset; ModifyPlan should populate it from the context.
	model := clumioAWSConnectionResourceModel{
		AccountNativeID: basetypes.NewStringValue(accountId),
		AWSRegion:       basetypes.NewStringValue(region),
		Description:     basetypes.NewStringValue(description),
	}
	plan := tfsdk.Plan{Schema: schemaResp.Schema}
	assert.False(t, plan.Set(ctx, &model).HasError())
	resp := &resource.ModifyPlanResponse{Plan: plan}

	res.ModifyPlan(ctx, resource.ModifyPlanRequest{Plan: plan}, resp)

	var got clumioAWSConnectionResourceModel
	diags := resp.Plan.Get(ctx, &got)
	assert.False(t, diags.HasError())
	assert.Equal(t, basetypes.NewStringValue(ouId), got.OrganizationalUnitID)
}

// TestModifyPlanEmptyContextWarnsOnOrganizationalUnitMove verifies that a context-driven OU move
// is applied but surfaced as a plan-time warning rather than done silently.
func TestModifyPlanEmptyContextWarnsOnOrganizationalUnitMove(t *testing.T) {
	ctx := context.Background()
	res := NewClumioAWSConnectionResource().(*clumioAWSConnectionResource)
	// Provider configured with an EMPTY OU context.
	res.Configure(ctx, resource.ConfigureRequest{
		ProviderData: &common.ApiClient{ClumioConfig: sdkconfig.Config{}},
	}, &resource.ConfigureResponse{})

	schemaResp := &resource.SchemaResponse{}
	res.Schema(ctx, resource.SchemaRequest{}, schemaResp)

	// The connection already lives in tenant OU `ouId`; both prior state and the carried-forward
	// plan hold that value.
	model := clumioAWSConnectionResourceModel{
		ID:                   basetypes.NewStringValue("connection-id"),
		AccountNativeID:      basetypes.NewStringValue(accountId),
		AWSRegion:            basetypes.NewStringValue(region),
		Description:          basetypes.NewStringValue(description),
		OrganizationalUnitID: basetypes.NewStringValue(ouId),
	}
	plan := tfsdk.Plan{Schema: schemaResp.Schema}
	assert.False(t, plan.Set(ctx, &model).HasError())
	state := tfsdk.State{Schema: schemaResp.Schema}
	assert.False(t, state.Set(ctx, &model).HasError())
	resp := &resource.ModifyPlanResponse{Plan: plan}

	res.ModifyPlan(ctx, resource.ModifyPlanRequest{Plan: plan, State: state}, resp)

	var got clumioAWSConnectionResourceModel
	diags := resp.Plan.Get(ctx, &got)
	assert.False(t, diags.HasError())
	assert.Equal(t, basetypes.NewStringValue(defaultOrgUnitId), got.OrganizationalUnitID)
	assert.Len(t, resp.Diagnostics.Warnings(), 1)
}
