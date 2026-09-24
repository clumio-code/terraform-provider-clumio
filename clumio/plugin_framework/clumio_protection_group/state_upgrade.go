// Copyright 2026. Clumio, Inc.

package clumio_protection_group

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// UpgradeState migrates the object_filter and prefix_filters nested blocks from sets to lists.
// The Go resource model uses slices for both schema representations, so the migration can preserve
// all prior values while re-encoding them with the current schema.
func (r *clumioProtectionGroupResource) UpgradeState(
	_ context.Context) map[int64]resource.StateUpgrader {
	priorSchema := protectionGroupSchemaV0()

	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema: &priorSchema,
			StateUpgrader: func(
				ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				var priorState clumioProtectionGroupResourceModel
				resp.Diagnostics.Append(req.State.Get(ctx, &priorState)...)
				if resp.Diagnostics.HasError() {
					return
				}

				resp.Diagnostics.Append(resp.State.Set(ctx, &priorState)...)
			},
		},
	}
}
