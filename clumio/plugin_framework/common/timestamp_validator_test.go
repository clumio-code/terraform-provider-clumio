// Copyright 2026. Clumio, Inc.

// This file contains the unit tests for the timestamp validator.

//go:build unit

package common

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

// Unit test for the following cases:
//   - Valid RFC-3339 timestamp passes validation.
//   - Malformed timestamp fails validation.
//   - Null and unknown values are skipped.
func TestIsTimestamp(t *testing.T) {

	ctx := context.Background()

	tests := []struct {
		name      string
		value     types.String
		wantDiags bool
	}{
		{name: "valid timestamp", value: types.StringValue("2026-07-27T00:00:00Z")},
		{name: "valid timestamp with offset", value: types.StringValue("2026-07-27T09:00:00+09:00")},
		{name: "malformed timestamp", value: types.StringValue("2026-07-27"), wantDiags: true},
		{name: "not a timestamp", value: types.StringValue("latest"), wantDiags: true},
		{name: "null is skipped", value: types.StringNull()},
		{name: "unknown is skipped", value: types.StringUnknown()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &validator.StringResponse{}
			IsTimestamp().ValidateString(ctx, validator.StringRequest{
				Path:        path.Root("timestamp"),
				ConfigValue: tt.value,
			}, resp)
			assert.Equal(t, tt.wantDiags, resp.Diagnostics.HasError())
		})
	}
}
