// Copyright 2026. Clumio, Inc.

//go:build unit

package clumio_protection_group

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpgradeStateV0(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	protectionInfo, diags := mapClumioProtectionInfoToSchemaProtectionInfo(nil)
	require.False(t, diags.HasError())
	excludedSubPrefixes := make([]types.String, 656)
	for i := range excludedSubPrefixes {
		excludedSubPrefixes[i] = types.StringValue(fmt.Sprintf("excluded-prefix-%04d/", i))
	}

	priorStateData := clumioProtectionGroupResourceModel{
		ID:               types.StringValue(id),
		Name:             types.StringValue(name),
		Description:      types.StringValue(description),
		BucketRule:       types.StringValue(bucketRule),
		ProtectionStatus: types.StringValue(protectionStatus),
		ProtectionInfo:   protectionInfo,
		ObjectFilter: []*objectFilterModel{
			{
				LatestVersionOnly:             types.BoolValue(true),
				EarliestLastModifiedTimestamp: types.StringValue(earliestLastModifiedTimestamp),
				StorageClasses: []types.String{
					types.StringValue(storageClass),
				},
				PrefixFilters: []*prefixFilterModel{
					{
						Prefix:              types.StringValue(prefix),
						ExcludedSubPrefixes: excludedSubPrefixes,
					},
				},
			},
		},
	}

	priorSchema := protectionGroupSchemaV0()
	priorState := tfsdk.State{Schema: priorSchema}
	diags = priorState.Set(ctx, &priorStateData)
	require.False(t, diags.HasError(), diags)

	currentSchema := protectionGroupSchemaV1()
	resp := resource.UpgradeStateResponse{
		State: tfsdk.State{Schema: currentSchema},
	}
	upgrader := (&clumioProtectionGroupResource{}).UpgradeState(ctx)[0]
	upgrader.StateUpgrader(ctx, resource.UpgradeStateRequest{State: &priorState}, &resp)
	require.False(t, resp.Diagnostics.HasError(), resp.Diagnostics)

	var upgradedStateData clumioProtectionGroupResourceModel
	diags = resp.State.Get(ctx, &upgradedStateData)
	require.False(t, diags.HasError(), diags)
	assert.Equal(t, priorStateData, upgradedStateData)
	assert.Len(t, upgradedStateData.ObjectFilter[0].PrefixFilters[0].ExcludedSubPrefixes, 656)
}
