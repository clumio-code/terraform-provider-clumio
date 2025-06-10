// Copyright 2025. Clumio, Inc.

// This file contains the unit test for the Schema function in schema.go.

//go:build unit

package clumio_report_configuration

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
)

// TestSchema checks the schema returned for a given resource.
func TestSchema(t *testing.T) {

	res := &clumioReportConfigurationResource{}
	resp := &resource.SchemaResponse{}
	res.Schema(context.Background(), resource.SchemaRequest{}, resp)
	assert.NotNil(t, resp.Schema)

	// Ensure that all attributes have a description set.
	for _, attr := range resp.Schema.Attributes {
		assert.NotEmpty(t, attr.GetDescription())
	}
}

// TestIfTagExists checks validation for the tag attribute.
func TestIfTagExists(t *testing.T) {
	ctx := context.Background()
	testValidator := ifTagExist()
	testRoot := path.Root("assets")
	testPath := testRoot.AtName("tags")
	assert.Equal(t, testPath.ParentPath(), testRoot)

	res := &clumioReportConfigurationResource{}
	schemaResp := &resource.SchemaResponse{}
	res.Schema(ctx, resource.SchemaRequest{}, schemaResp)

	t.Run("ConfigValue is unknown", func(t *testing.T) {
		req := validator.StringRequest{
			ConfigValue: types.StringUnknown(),
		}
		resp := &validator.StringResponse{}
		testValidator.ValidateString(ctx, req, resp)
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("ConfigValue is null", func(t *testing.T) {
		req := validator.StringRequest{
			ConfigValue: types.StringNull(),
		}
		resp := &validator.StringResponse{}
		testValidator.ValidateString(ctx, req, resp)
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("GetAttribute returns diagnostics error", func(t *testing.T) {
		req := validator.StringRequest{
			ConfigValue: types.StringValue("any"),
			Config: tfsdk.Config{
				Raw:    tftypes.NewValue(tftypes.String, "any"),
				Schema: schemaResp.Schema,
			},
			Path: testPath,
		}
		resp := &validator.StringResponse{}
		testValidator.ValidateString(ctx, req, resp)
		assert.True(t, resp.Diagnostics.HasError())
	})
}
