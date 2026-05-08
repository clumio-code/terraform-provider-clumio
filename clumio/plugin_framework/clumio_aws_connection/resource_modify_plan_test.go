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
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
)

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

	objectType := tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			schemaId:                   tftypes.String,
			schemaAccountNativeId:      tftypes.String,
			schemaAwsRegion:            tftypes.String,
			schemaDescription:          tftypes.String,
			schemaOrganizationalUnitId: tftypes.String,
			schemaConnectionStatus:     tftypes.String,
			schemaToken:                tftypes.String,
			schemaNamespace:            tftypes.String,
			schemaClumioAwsAccountId:   tftypes.String,
			schemaClumioAwsRegion:      tftypes.String,
			schemaExternalId:           tftypes.String,
			schemaDataPlaneAccountId:   tftypes.String,
		},
	}

	values := map[string]tftypes.Value{
		schemaId:                   tftypes.NewValue(tftypes.String, nil),
		schemaAccountNativeId:      tftypes.NewValue(tftypes.String, accountId),
		schemaAwsRegion:            tftypes.NewValue(tftypes.String, region),
		schemaDescription:          tftypes.NewValue(tftypes.String, description),
		schemaOrganizationalUnitId: tftypes.NewValue(tftypes.String, nil),
		schemaConnectionStatus:     tftypes.NewValue(tftypes.String, nil),
		schemaToken:                tftypes.NewValue(tftypes.String, nil),
		schemaNamespace:            tftypes.NewValue(tftypes.String, nil),
		schemaClumioAwsAccountId:   tftypes.NewValue(tftypes.String, nil),
		schemaClumioAwsRegion:      tftypes.NewValue(tftypes.String, nil),
		schemaExternalId:           tftypes.NewValue(tftypes.String, nil),
		schemaDataPlaneAccountId:   tftypes.NewValue(tftypes.String, nil),
	}

	plan := tfsdk.Plan{
		Raw:    tftypes.NewValue(objectType, values),
		Schema: schemaResp.Schema,
	}
	resp := &resource.ModifyPlanResponse{
		Plan: plan,
	}

	res.ModifyPlan(ctx, resource.ModifyPlanRequest{
		Plan: plan,
	}, resp)

	var got clumioAWSConnectionResourceModel
	diags := resp.Plan.Get(ctx, &got)
	assert.False(t, diags.HasError())
	assert.Equal(t, basetypes.NewStringValue(ouId), got.OrganizationalUnitID)
}
