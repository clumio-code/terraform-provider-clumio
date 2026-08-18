// Copyright 2026. Clumio, Inc.
package common

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = TimestampValidator{}

// TimestampValidator validates that a string attribute, if set, is a valid timestamp. Clumio
// timestamps are RFC-3339 formatted, so the value is parsed with the stdlib RFC-3339 layout to
// catch malformed timestamps at plan time instead of at the API call.
type TimestampValidator struct{}

// Description returns the plain-text description of the validation performed.
func (v TimestampValidator) Description(_ context.Context) string {
	return "value must be a valid RFC-3339 timestamp (e.g. 2026-07-27T00:00:00Z)"
}

// MarkdownDescription returns the Markdown description of the validation performed.
func (v TimestampValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString runs the validation logic on the provided request and adds diagnostics as
// required to the response. Null and unknown values are skipped.
func (v TimestampValidator) ValidateString(
	ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {

	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	if _, err := time.Parse(time.RFC3339, req.ConfigValue.ValueString()); err != nil {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid timestamp",
			fmt.Sprintf("Attribute %s: %s.", req.Path, v.Description(ctx)))
	}
}

// IsTimestamp returns a validator which ensures that a string attribute is a valid timestamp.
func IsTimestamp() validator.String {
	return TimestampValidator{}
}
