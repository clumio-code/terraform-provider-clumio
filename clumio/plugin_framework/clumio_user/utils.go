// Copyright 2024. Clumio, Inc.

// This file hold various utility functions used by the clumio_user Terraform resource.

package clumio_user

import (
	"context"

	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// getAccessControlCfgFromHTTPRes parses the accessControlCfg from the http response  and converts
// it into the corresponding schema attribute value. Along with that it also returns the
// assignedRole and OUIds.
func getAccessControlCfgFromHTTPRes(
	ctx context.Context, accessControlCfg []*models.RoleForOrganizationalUnits,
	diag *diag.Diagnostics) (basetypes.SetValue, basetypes.StringValue, basetypes.SetValue) {

	accessControl := make([]roleForOrganizationalUnitModel, len(accessControlCfg))
	organizationalUnitIds := make([]*string, 0)

	var assignedRole string
	for idx, element := range accessControlCfg {
		if element.RoleId != nil {
			assignedRole = *element.RoleId
		} else {
			assignedRole = ""
		}
		organizationalUnitIds = append(organizationalUnitIds, element.OrganizationalUnitIds...)
		ouIds, conversionDiags := types.SetValueFrom(
			ctx, types.StringType, element.OrganizationalUnitIds)
		diag.Append(conversionDiags...)
		accessControl[idx] = roleForOrganizationalUnitModel{
			RoleId:                types.StringValue(assignedRole),
			OrganizationalUnitIds: ouIds,
		}
	}

	ouIdSet, conversionDiags := types.SetValueFrom(ctx, types.StringType, organizationalUnitIds)
	diag.Append(conversionDiags...)

	accessControlList, conversionDiags := types.SetValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			schemaRoleId: types.StringType,
			schemaOrganizationalUnitIds: types.SetType{
				ElemType: types.StringType,
			},
		},
	}, accessControl)
	diag.Append(conversionDiags...)

	return accessControlList, types.StringValue(assignedRole), ouIdSet
}

// makeAccessControlCfgMap creates an access control config map from the schema attribute value for
// the same.
func makeAccessControlCfgMap(
	ctx context.Context, accessControlCfg []attr.Value) map[string][]string {

	accessControlCfgMap := make(map[string][]string)
	for _, element := range accessControlCfg {
		roleForOU := roleForOrganizationalUnitModel{}
		element.(types.Object).As(ctx, &roleForOU, basetypes.ObjectAsOptions{})
		roleId := roleForOU.RoleId.ValueString()
		ouIds := make([]string, 0)
		_ = roleForOU.OrganizationalUnitIds.ElementsAs(ctx, &ouIds, false)
		accessControlCfgMap[roleId] = ouIds
	}
	return accessControlCfgMap
}

// getAccessControlCfgMapDiff generates the difference between the two given maps.
func getAccessControlCfgMapDiff(map1 map[string][]string,
	map2 map[string][]string) []*models.RoleForOrganizationalUnits {

	mapDiff := make([]*models.RoleForOrganizationalUnits, 0)
	for key := range map1 {
		roleId := key
		if _, ok := map2[roleId]; !ok {
			mapDiff = append(mapDiff, &models.RoleForOrganizationalUnits{
				RoleId:                &roleId,
				OrganizationalUnitIds: common.GetStringPtrSliceFromStringSlice(map1[roleId]),
			})
			continue
		}
		diff := common.SliceDifferenceString(map1[roleId], map2[roleId])
		if len(diff) > 0 {
			mapDiff = append(mapDiff, &models.RoleForOrganizationalUnits{
				RoleId:                &roleId,
				OrganizationalUnitIds: common.GetStringPtrSliceFromStringSlice(diff),
			})
		}
	}
	return mapDiff
}

// getAccessControlCfgUpdates compares the old and new configs and returns the access control config
// maps to be added and removed.
func getAccessControlCfgUpdates(ctx context.Context, oldCfg, newCfg []attr.Value) (
	[]*models.RoleForOrganizationalUnits, []*models.RoleForOrganizationalUnits) {

	oldCfgMap := makeAccessControlCfgMap(ctx, oldCfg)
	newCfgMap := makeAccessControlCfgMap(ctx, newCfg)

	add := getAccessControlCfgMapDiff(newCfgMap, oldCfgMap)
	remove := getAccessControlCfgMapDiff(oldCfgMap, newCfgMap)

	return add, remove
}
