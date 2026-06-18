// Copyright 2026. Clumio, Inc.

// This file contains the unit tests for the bucket_rule mapping helpers in utils.go.

//go:build unit

package clumio_gcs_protection_group

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

// TestBucketRuleLabelRoundTrip verifies the gcp_label operator mapping in both directions:
//   - eq/contains/not_eq/not_contains and all/not_all are maps of key => value.
//   - in/not_in are repeated {key, value} blocks that may carry the same key with multiple values.
func TestBucketRuleLabelRoundTrip(t *testing.T) {
	ctx := context.Background()

	eqMap, d := types.MapValueFrom(ctx, types.StringType, map[string]string{"dongjoo": "test"})
	assert.False(t, d.HasError())
	allMap, d := types.MapValueFrom(ctx, types.StringType,
		map[string]string{"env": "prod", "team": "x"})
	assert.False(t, d.HasError())

	m := &bucketRuleModel{
		GcpLabel: &gcpLabelOperatorModel{
			Eq:          eqMap,
			Contains:    types.MapNull(types.StringType),
			All:         allMap,
			NotEq:       types.MapNull(types.StringType),
			NotContains: types.MapNull(types.StringType),
			NotAll:      types.MapNull(types.StringType),
			In: []*labelModel{
				{Key: types.StringValue("env"), Value: types.StringValue("prod")},
				{Key: types.StringValue("env"), Value: types.StringValue("staging")},
			},
		},
	}

	// Terraform -> SDK.
	sdk, diags := mapSchemaBucketRuleToClumioBucketRule(ctx, m)
	assert.False(t, diags.HasError())
	assert.NotNil(t, sdk.GcpLabel)
	assert.Equal(t, "dongjoo", *sdk.GcpLabel.Eq.Key)
	assert.Equal(t, "test", *sdk.GcpLabel.Eq.Value)
	assert.Len(t, sdk.GcpLabel.All, 2)
	// in keeps both values for the same key (a map could not).
	assert.Len(t, sdk.GcpLabel.In, 2)

	// SDK -> Terraform.
	back, diags := mapClumioBucketRuleToSchemaBucketRule(ctx, sdk)
	assert.False(t, diags.HasError())
	assert.NotNil(t, back.GcpLabel)

	var eqBack map[string]string
	back.GcpLabel.Eq.ElementsAs(ctx, &eqBack, false)
	assert.Equal(t, map[string]string{"dongjoo": "test"}, eqBack)

	var allBack map[string]string
	back.GcpLabel.All.ElementsAs(ctx, &allBack, false)
	assert.Equal(t, map[string]string{"env": "prod", "team": "x"}, allBack)

	assert.Len(t, back.GcpLabel.In, 2)
	// Unset operators stay null.
	assert.True(t, back.GcpLabel.Contains.IsNull())
	assert.True(t, back.GcpLabel.NotAll.IsNull())
}

// TestBucketRuleToSDKEmpty verifies that an absent or empty bucket_rule (no condition fields set)
// maps to a nil SDK rule, so the caller can clear it via clear_bucket_rule rather than sending an
// empty rule object.
func TestBucketRuleToSDKEmpty(t *testing.T) {
	ctx := context.Background()

	rule, diags := mapSchemaBucketRuleToClumioBucketRule(ctx, nil)
	assert.False(t, diags.HasError())
	assert.Nil(t, rule)

	rule, diags = mapSchemaBucketRuleToClumioBucketRule(ctx, &bucketRuleModel{})
	assert.False(t, diags.HasError())
	assert.Nil(t, rule)
}

// TestBucketRuleNotEmptyValidator verifies the plan-time validator: an absent rule is allowed, an
// empty rule (no condition fields) is rejected, and a rule with one condition is allowed.
func TestBucketRuleNotEmptyValidator(t *testing.T) {
	ctx := context.Background()
	v := bucketRuleNotEmptyValidator{}
	childTypes := map[string]attr.Type{}
	brTypes := map[string]attr.Type{
		schemaGcpLabel:     types.ObjectType{AttrTypes: childTypes},
		schemaGcpLocation:  types.ObjectType{AttrTypes: childTypes},
		schemaGcpProjectId: types.ObjectType{AttrTypes: childTypes},
	}

	// Absent bucket_rule -> allowed (optional block).
	resp := &validator.ObjectResponse{}
	v.ValidateObject(ctx, validator.ObjectRequest{
		Path: path.Root("bucket_rule"), ConfigValue: types.ObjectNull(brTypes),
	}, resp)
	assert.False(t, resp.Diagnostics.HasError())

	// Empty bucket_rule (all conditions null) -> error.
	emptyVal := types.ObjectValueMust(brTypes, map[string]attr.Value{
		schemaGcpLabel:     types.ObjectNull(childTypes),
		schemaGcpLocation:  types.ObjectNull(childTypes),
		schemaGcpProjectId: types.ObjectNull(childTypes),
	})
	resp = &validator.ObjectResponse{}
	v.ValidateObject(ctx, validator.ObjectRequest{
		Path: path.Root("bucket_rule"), ConfigValue: emptyVal,
	}, resp)
	assert.True(t, resp.Diagnostics.HasError())

	// One condition set -> allowed.
	oneSet := types.ObjectValueMust(brTypes, map[string]attr.Value{
		schemaGcpLabel:     types.ObjectValueMust(childTypes, map[string]attr.Value{}),
		schemaGcpLocation:  types.ObjectNull(childTypes),
		schemaGcpProjectId: types.ObjectNull(childTypes),
	})
	resp = &validator.ObjectResponse{}
	v.ValidateObject(ctx, validator.ObjectRequest{
		Path: path.Root("bucket_rule"), ConfigValue: oneSet,
	}, resp)
	assert.False(t, resp.Diagnostics.HasError())
}
