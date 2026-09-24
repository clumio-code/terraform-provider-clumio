// Copyright 2024. Clumio, Inc.

// This file contains the unit test for the Schema function in schema.go.

//go:build unit

package clumio_protection_group

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSchema checks the schema returned for a given resource.
func TestSchema(t *testing.T) {

	res := &clumioProtectionGroupResource{}
	resp := &resource.SchemaResponse{}
	res.Schema(context.Background(), resource.SchemaRequest{}, resp)
	assert.NotNil(t, resp.Schema)
	assert.Equal(t, int64(1), resp.Schema.Version)

	objectFilterBlock, ok := resp.Schema.Blocks[schemaObjectFilter].(schema.ListNestedBlock)
	require.True(t, ok, "object_filter must use list nesting to avoid value-based set paths")
	_, ok = objectFilterBlock.NestedObject.Blocks[schemaPrefixFilters].(schema.ListNestedBlock)
	require.True(t, ok, "prefix_filters must use list nesting to avoid value-based set paths")

	// Ensure that all attributes have a description set.
	for _, attr := range resp.Schema.Attributes {
		assert.NotEmpty(t, attr.GetDescription())
	}
}
