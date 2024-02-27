// Copyright 2024. Clumio, Inc.

// This file hold various utility functions used by the clumio_protection_group Terraform resource.

package clumio_protection_group

import (
	"context"
	"errors"
	sdkProtectionGroups "github.com/clumio-code/clumio-go-sdk/controllers/protection_groups"
	"strings"
	"time"

	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// mapSchemaObjectFilterToClumioObjectFilter converts the schema object_filter to the model
// ObjectFilter
func mapSchemaObjectFilterToClumioObjectFilter(
	objectFilterSlice []*objectFilterModel) *models.ObjectFilter {
	if objectFilterSlice == nil || len(objectFilterSlice) == 0 {
		return nil
	}
	objectFilter := objectFilterSlice[0]
	latestVersionOnly := objectFilter.LatestVersionOnly.ValueBool()

	// Loop over StorageClasses field and map each item inside array to model
	storageClasses := make([]*string, 0)
	if objectFilter.StorageClasses != nil {
		for _, storageClass := range objectFilter.StorageClasses {
			storageClassStr := storageClass.ValueString()
			storageClasses = append(storageClasses, &storageClassStr)
		}
	}
	modelPrefixFilters := make([]*models.PrefixFilter, 0)
	// Loop over PrefixFilters field and map each item inside array to model
	if objectFilter.PrefixFilters != nil {
		for _, prefixFilter := range objectFilter.PrefixFilters {
			excludedSubPrefixesList := make([]*string, 0)
			for _, excludedSubPrefix := range prefixFilter.ExcludedSubPrefixes {
				excludedSubPrefixStr := excludedSubPrefix.ValueString()
				excludedSubPrefixesList = append(
					excludedSubPrefixesList, &excludedSubPrefixStr)
			}
			prefix := prefixFilter.Prefix.ValueString()
			modelPrefixFilter := &models.PrefixFilter{
				ExcludedSubPrefixes: excludedSubPrefixesList,
				Prefix:              &prefix,
			}
			modelPrefixFilters = append(modelPrefixFilters, modelPrefixFilter)
		}
	}
	return &models.ObjectFilter{
		LatestVersionOnly: &latestVersionOnly,
		PrefixFilters:     modelPrefixFilters,
		StorageClasses:    storageClasses,
	}
}

// pollForProtectionGroup polls till the protection group becomes available after create  or update
// protection group as they are asynchronous operations.
func pollForProtectionGroup(
	ctx context.Context, id string, protectionGroup sdkProtectionGroups.ProtectionGroupsV1Client) error {

	interval := time.Duration(intervalInSec) * time.Second
	ticker := time.NewTicker(interval)
	timeout := time.After(time.Duration(timeoutInSec) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return errors.New("context done")
		case <-ticker.C:
			_, err := protectionGroup.ReadProtectionGroup(id)
			if err != nil {
				errResponseStr := string(err.Response)
				if !strings.Contains(errResponseStr,
					"A resource with the requested ID could not be found.") {
					return errors.New(
						"error reading protection-group which was created")
				}
				continue
			}
			return nil
		case <-timeout:
			return errors.New("polling timeout")
		}
	}
}

// mapClumioObjectFilterToSchemaObjectFilter converts the Object Filter from the API to the schema
// object_filter
func mapClumioObjectFilterToSchemaObjectFilter(
	modelObjectFilter *models.ObjectFilter) []*objectFilterModel {
	schemaObjFilter := &objectFilterModel{}
	if modelObjectFilter.LatestVersionOnly != nil {
		schemaObjFilter.LatestVersionOnly = types.BoolValue(
			*modelObjectFilter.LatestVersionOnly)
	}
	// Loop over PrefixFilters field and map each item inside array to schema
	if modelObjectFilter.PrefixFilters != nil {
		prefixFilters := make([]*prefixFilterModel, 0)
		for _, modelPrefixFilter := range modelObjectFilter.PrefixFilters {
			prefixFilter := &prefixFilterModel{}
			prefixFilter.Prefix = types.StringPointerValue(modelPrefixFilter.Prefix)
			if modelPrefixFilter.ExcludedSubPrefixes != nil {
				excludedSubPrefixes := make([]types.String, 0)
				for _, excludeSubPrefix := range modelPrefixFilter.ExcludedSubPrefixes {
					excludeSubPrefixStr := types.StringPointerValue(excludeSubPrefix)
					excludedSubPrefixes = append(excludedSubPrefixes, excludeSubPrefixStr)
				}
				prefixFilter.ExcludedSubPrefixes = excludedSubPrefixes
			}
			prefixFilters = append(prefixFilters, prefixFilter)
		}
		schemaObjFilter.PrefixFilters = prefixFilters
	}
	// Loop over StorageClasses field and map each item inside array to schema
	storageClasses := make([]types.String, 0)
	for _, storageClass := range modelObjectFilter.StorageClasses {
		storageClassStrType := types.StringPointerValue(storageClass)
		storageClasses = append(storageClasses, storageClassStrType)
	}
	schemaObjFilter.StorageClasses = storageClasses
	return []*objectFilterModel{schemaObjFilter}
}

// mapClumioProtectionInfoToSchemaProtectionInfo converts the Protection Info from the API to the
// schema protection_info.
func mapClumioProtectionInfoToSchemaProtectionInfo(
	modelProtectionInfo *models.ProtectionInfoWithRule) (types.List, diag.Diagnostics) {

	objtype := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			schemaPolicyId:             types.StringType,
			schemaInheritingEntityId:   types.StringType,
			schemaInheritingEntityType: types.StringType,
		},
	}
	if modelProtectionInfo == nil {
		return types.ListValue(objtype, []attr.Value{})
	}
	attrTypes := make(map[string]attr.Type)
	attrTypes[schemaPolicyId] = types.StringType
	attrTypes[schemaInheritingEntityType] = types.StringType
	attrTypes[schemaInheritingEntityId] = types.StringType

	attrValues := make(map[string]attr.Value)
	attrValues[schemaPolicyId] = types.StringValue("")
	attrValues[schemaInheritingEntityType] = types.StringValue("")
	attrValues[schemaInheritingEntityId] = types.StringValue("")
	if modelProtectionInfo != nil {
		if modelProtectionInfo.PolicyId != nil {
			attrValues[schemaPolicyId] = types.StringPointerValue(modelProtectionInfo.PolicyId)
		}
		if modelProtectionInfo.InheritingEntityType != nil {
			attrValues[schemaInheritingEntityType] =
				types.StringPointerValue(modelProtectionInfo.InheritingEntityType)
		}
		if modelProtectionInfo.InheritingEntityId != nil {
			attrValues[schemaInheritingEntityId] =
				types.StringPointerValue(modelProtectionInfo.InheritingEntityId)
		}
	}
	obj, diags := types.ObjectValue(attrTypes, attrValues)

	listobj, listdiag := types.ListValue(objtype, []attr.Value{obj})
	listdiag.Append(diags...)
	return listobj, listdiag
}

// clearOUContext resets the OrganizationalUnitContext in the client.
func (r *clumioProtectionGroupResource) clearOUContext() {
	r.client.ClumioConfig.OrganizationalUnitContext = ""
}
