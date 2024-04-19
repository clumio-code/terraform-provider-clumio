// Copyright 2024. Clumio, Inc.

// This file contains the unit tests for the functions in utils.go

//go:build unit

package clumio_protection_group

import (
	"testing"

	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
)

// Unit test for the utility function to convert ClumioObjectFilter to SchemaObjectFilter.
func TestMapClumioObjectFilterToSchemaObjectFilter(t *testing.T) {

	latestVersionOnly := true
	modelObjectFilter := &models.ObjectFilter{
		LatestVersionOnly: &latestVersionOnly,
		PrefixFilters: []*models.PrefixFilter{
			{
				ExcludedSubPrefixes: []*string{&exclPrefix, &exclPrefix2},
			},
			{
				ExcludedSubPrefixes: []*string{&exclPrefix, &exclPrefix2},
			},
		},
		StorageClasses: []*string{&storageClass, &storageClass2},
	}

	t.Run("All object filter attributes populated", func(t *testing.T) {
		schemaObjectFilter := mapClumioObjectFilterToSchemaObjectFilter(modelObjectFilter)

		// Ensure the object filter attributes are correct.
		modelOF := modelObjectFilter
		schemaOF := schemaObjectFilter[0]
		assert.Equal(t, *modelOF.LatestVersionOnly, schemaOF.LatestVersionOnly.ValueBool())
		assert.Equal(t, *modelOF.StorageClasses[0], schemaOF.StorageClasses[0].ValueString())
		assert.Equal(t, *modelOF.StorageClasses[1], schemaOF.StorageClasses[1].ValueString())

		// Ensure the first prefix filter attributes are correct.
		modelPF := modelOF.PrefixFilters[0]
		schemaPF := schemaOF.PrefixFilters[0]
		assert.Equal(t, *modelPF.ExcludedSubPrefixes[0],
			schemaPF.ExcludedSubPrefixes[0].ValueString())
		assert.Equal(t, *modelPF.ExcludedSubPrefixes[1],
			schemaPF.ExcludedSubPrefixes[1].ValueString())

		// Ensure the second prefix filter attributes are correct.
		modelPF = modelOF.PrefixFilters[1]
		schemaPF = schemaOF.PrefixFilters[1]
		assert.Equal(t, *modelPF.ExcludedSubPrefixes[0],
			schemaPF.ExcludedSubPrefixes[0].ValueString())
		assert.Equal(t, *modelPF.ExcludedSubPrefixes[1],
			schemaPF.ExcludedSubPrefixes[1].ValueString())
	})

	t.Run("Test for object filter with empty prefix filters", func(t *testing.T) {
		modelObjectFilter.PrefixFilters = nil
		schemaObjectFilter := mapClumioObjectFilterToSchemaObjectFilter(modelObjectFilter)
		assert.Nil(t, schemaObjectFilter[0].PrefixFilters)
	})
}

// Unit test for the utility function to convert SchemaObjectFilter to ClumioObjectFilter.
func TestMapSchemaObjectFilterToClumioObjectFilter(t *testing.T) {

	schemaOF := []*objectFilterModel{
		{
			LatestVersionOnly: basetypes.NewBoolValue(true),
			PrefixFilters: []*prefixFilterModel{
				{
					ExcludedSubPrefixes: []types.String{
						basetypes.NewStringValue(exclPrefix),
						basetypes.NewStringValue(exclPrefix2),
					},
					Prefix: basetypes.NewStringValue(prefix),
				},
				{
					ExcludedSubPrefixes: []types.String{
						basetypes.NewStringValue(exclPrefix),
						basetypes.NewStringValue(exclPrefix2),
					},
					Prefix: basetypes.NewStringValue(prefix),
				},
			},
			StorageClasses: []types.String{
				basetypes.NewStringValue(storageClass),
				basetypes.NewStringValue(storageClass2),
			},
		},
	}

	t.Run("All object filter attributes populated", func(t *testing.T) {
		modelOF := mapSchemaObjectFilterToClumioObjectFilter(schemaOF)

		// Ensure the object filter attributes are correct.
		assert.Equal(t, schemaOF[0].LatestVersionOnly.ValueBool(), *modelOF.LatestVersionOnly)
		assert.Equal(t, schemaOF[0].StorageClasses[0].ValueString(), *modelOF.StorageClasses[0])
		assert.Equal(t, schemaOF[0].StorageClasses[1].ValueString(), *modelOF.StorageClasses[1])

		// Ensure the first prefix filter attributes are correct.
		modelPF := modelOF.PrefixFilters[0]
		schemaPF := schemaOF[0].PrefixFilters[0]
		assert.Equal(t, schemaPF.Prefix.ValueString(), *modelPF.Prefix)
		assert.Equal(t, schemaPF.ExcludedSubPrefixes[0].ValueString(),
			*modelPF.ExcludedSubPrefixes[0])
		assert.Equal(t, schemaPF.ExcludedSubPrefixes[1].ValueString(),
			*modelPF.ExcludedSubPrefixes[1])

		// Ensure the second prefix filter attributes are correct.
		modelPF = modelOF.PrefixFilters[1]
		schemaPF = schemaOF[0].PrefixFilters[1]
		assert.Equal(t, schemaPF.ExcludedSubPrefixes[0].ValueString(),
			*modelPF.ExcludedSubPrefixes[0])
		assert.Equal(t, schemaPF.ExcludedSubPrefixes[1].ValueString(),
			*modelPF.ExcludedSubPrefixes[1])
	})

	t.Run("Test for object filter with empty prefix filters", func(t *testing.T) {
		schemaOF[0].PrefixFilters = nil
		modelOF := mapSchemaObjectFilterToClumioObjectFilter(schemaOF)
		assert.Equal(t, 0, len(modelOF.PrefixFilters))
	})

	t.Run("Test with empty object filter", func(t *testing.T) {
		modelOF := mapSchemaObjectFilterToClumioObjectFilter(nil)
		assert.Nil(t, modelOF)
	})
}

// Unit test for the utility function to convert the Protection Info from the API to the schema
// protection_info.
func TestMapClumioProtectionInfoToSchemaProtectionInfo(t *testing.T) {

	modelProtectionInfo := &models.ProtectionInfoWithRule{
		InheritingEntityId:   &entityId,
		InheritingEntityType: &entityType,
		PolicyId:             &policyId,
	}

	t.Run("All protection info attributes populated", func(t *testing.T) {
		schemaList, diags := mapClumioProtectionInfoToSchemaProtectionInfo(modelProtectionInfo)
		assert.Nil(t, diags)
		assert.Equal(t, 1, len(schemaList.Elements()))

		schemaProtectionInfoObject := schemaList.Elements()[0].(types.Object)
		schemaProtectionInfo := make(map[string]*string)
		for key, val := range schemaProtectionInfoObject.Attributes() {
			valStr := val.(types.String).ValueString()
			schemaProtectionInfo[key] = &valStr
		}

		assert.Equal(t, *schemaProtectionInfo[schemaInheritingEntityId],
			*modelProtectionInfo.InheritingEntityId)
		assert.Equal(t, *schemaProtectionInfo[schemaInheritingEntityType],
			*modelProtectionInfo.InheritingEntityType)
		assert.Equal(t, *schemaProtectionInfo[schemaPolicyId], *modelProtectionInfo.PolicyId)
	})

	t.Run("Test for empty protection info", func(t *testing.T) {
		schemaList, diags := mapClumioProtectionInfoToSchemaProtectionInfo(nil)
		assert.Nil(t, diags)
		assert.Equal(t, 0, len(schemaList.Elements()))
	})
}
