// Copyright 2024. Clumio, Inc.
package common

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.Set = WrapperSetValidator{}

// WrapperSetValidator acts as a wrapper and validates the input against the provider validator.
// If there's an error it returns the provided custom error message.
type WrapperSetValidator struct {
	validator.Set
}

// ValidateSet runs validation logic on the provided request and adds diagnostics as required to the response
func (v WrapperSetValidator) ValidateSet(ctx context.Context, req validator.SetRequest, res *validator.SetResponse) {
	// Run the provided validator after simplifying the attribute path
	fieldName := GetFieldNameFromNestedBlockPath(req)
	modifiedReq := validator.SetRequest{
		Path:           path.Root(fieldName),
		PathExpression: req.PathExpression,
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
	}
	v.Set.ValidateSet(ctx, modifiedReq, res)
}

// Function wrapper around the validator
func WrapSetValidator(validator validator.Set) validator.Set {
	return WrapperSetValidator{validator}
}
