// Copyright 2024. Clumio, Inc.

// This file contains the unit tests for the functions in utils.go

//go:build unit

package clumio_user

import (
	"context"
	"testing"

	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
)

// Unit test for converting the SDK RoleForOrganizationalUnits model into the clumio_user resource
// roleForOrganizationalUnitModel.
func TestGetAccessControlCfgFromHTTPRes(t *testing.T) {

	ctx := context.Background()
	accessControlCfg := []*models.RoleForOrganizationalUnits{
		{
			RoleId:                &roleId,
			OrganizationalUnitIds: []*string{&ou},
		},
	}

	var diags diag.Diagnostics
	var roleForOUs []*roleForOrganizationalUnitModel
	accCtrlCfg := getAccessControlCfgFromHTTPRes(ctx, accessControlCfg, &diags)
	diags = accCtrlCfg.ElementsAs(ctx, &roleForOUs, true)
	assert.Nil(t, diags)
	assert.Equal(t, roleId, roleForOUs[0].RoleId.ValueString())
	var ous []*string
	diags = roleForOUs[0].OrganizationalUnitIds.ElementsAs(ctx, &ous, true)
	assert.Nil(t, diags)
	assert.Equal(t, &ou, ous[0])
}

// Unit test that checks if the resource old and roleForOrganizationalUnitModel gets converted to
// corresponding SDK RoleForOrganizationalUnits lists corresponding to additions and deletions.
func TestGetAccessControlCfgUpdates(t *testing.T) {

	ctx := context.Background()
	ouIdsList := []*string{&ou}
	ouIds, conversionDiags := types.SetValueFrom(
		ctx, types.StringType, ouIdsList)
	assert.Nil(t, conversionDiags)

	accessControlModel := []*roleForOrganizationalUnitModel{
		{
			RoleId:                basetypes.NewStringValue(roleId),
			OrganizationalUnitIds: ouIds,
		},
	}
	firstAccessControlList, conversionDiags := types.SetValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			schemaRoleId: types.StringType,
			schemaOrganizationalUnitIds: types.SetType{
				ElemType: types.StringType,
			},
		},
	}, accessControlModel)
	assert.Nil(t, conversionDiags)

	ouIdsList = []*string{&ouUpdated}
	ouIds, conversionDiags = types.SetValueFrom(
		ctx, types.StringType, ouIdsList)
	assert.Nil(t, conversionDiags)

	accessControlModel = []*roleForOrganizationalUnitModel{
		{
			RoleId:                basetypes.NewStringValue(roleId2),
			OrganizationalUnitIds: ouIds,
		},
	}
	secondAccessControlList, conversionDiags := types.SetValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			schemaRoleId: types.StringType,
			schemaOrganizationalUnitIds: types.SetType{
				ElemType: types.StringType,
			},
		},
	}, accessControlModel)
	assert.Nil(t, conversionDiags)

	added, removed := getAccessControlCfgUpdates(
		ctx, firstAccessControlList.Elements(), secondAccessControlList.Elements())

	assert.Equal(t, roleId2, *added[0].RoleId)
	assert.Equal(t, ouUpdated, *added[0].OrganizationalUnitIds[0])

	assert.Equal(t, roleId, *removed[0].RoleId)
	assert.Equal(t, ou, *removed[0].OrganizationalUnitIds[0])

}
