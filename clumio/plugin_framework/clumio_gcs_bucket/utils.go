// Copyright 2026. Clumio, Inc.

// This file holds various utility functions used by the clumio_gcs_bucket Terraform data
// source.

package clumio_gcs_bucket

import (
	"context"

	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// gcsBucketAttrTypes returns the attribute types of a single GCS bucket object in the gcs_buckets
// attribute. It is used both for the object conversion and for the resulting set type.
func gcsBucketAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		schemaId:                   types.StringType,
		schemaBucketName:           types.StringType,
		schemaProjectId:            types.StringType,
		schemaProjectUuid:          types.StringType,
		schemaLocation:             types.StringType,
		schemaLocationType:         types.StringType,
		schemaLocationUuid:         types.StringType,
		schemaOrganizationalUnitId: types.StringType,
		schemaObjectCount:          types.Int64Type,
		schemaSizeBytes:            types.Int64Type,
		schemaProtectionGroupCount: types.Int64Type,
		schemaIsDeleted:            types.BoolType,
		schemaIsVersioningEnabled:  types.BoolType,
		schemaCreatedTimestamp:     types.StringType,
		schemaLastBackupTimestamp:  types.StringType,
		schemaLabels:               types.SetType{ElemType: labelObjectType()},
	}
}

// labelObjectType is the Terraform object type for a single {key, value} GCP label.
func labelObjectType() types.ObjectType {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		schemaKey:   types.StringType,
		schemaValue: types.StringType,
	}}
}

// gcpLabelsToSet converts the SDK GCP labels into a Terraform set of {key, value} objects.
func gcpLabelsToSet(labels []*models.GcpLabelModel) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics
	objType := labelObjectType()
	elems := make([]attr.Value, 0, len(labels))
	for _, l := range labels {
		obj, d := types.ObjectValue(objType.AttrTypes, map[string]attr.Value{
			schemaKey:   basetypes.NewStringPointerValue(l.Key),
			schemaValue: basetypes.NewStringPointerValue(l.Value),
		})
		diags.Append(d...)
		if diags.HasError() {
			return types.SetNull(objType), diags
		}
		elems = append(elems, obj)
	}
	set, d := types.SetValue(objType, elems)
	diags.Append(d...)
	return set, diags
}

// populateGCSBucketsInDataSourceModel populates the gcs_buckets attribute in the data source model
// from the items returned in the API response.
func populateGCSBucketsInDataSourceModel(ctx context.Context,
	model *clumioGCSBucketDataSourceModel,
	items []*models.GCSBucket) diag.Diagnostics {

	var diags diag.Diagnostics

	attrTypes := gcsBucketAttrTypes()
	objType := types.ObjectType{AttrTypes: attrTypes}

	attrVals := make([]attr.Value, 0, len(items))
	for _, item := range items {
		labelsSet, labelDiags := gcpLabelsToSet(item.Labels)
		diags.Append(labelDiags...)
		if diags.HasError() {
			return diags
		}
		attrValues := map[string]attr.Value{
			schemaId:                   basetypes.NewStringPointerValue(item.Id),
			schemaBucketName:           basetypes.NewStringPointerValue(item.BucketName),
			schemaProjectId:            basetypes.NewStringPointerValue(item.ProjectId),
			schemaProjectUuid:          basetypes.NewStringPointerValue(item.ProjectUuid),
			schemaLocation:             basetypes.NewStringPointerValue(item.Location),
			schemaLocationType:         basetypes.NewStringPointerValue(item.LocationType),
			schemaLocationUuid:         basetypes.NewStringPointerValue(item.LocationUuid),
			schemaOrganizationalUnitId: basetypes.NewStringPointerValue(item.OrganizationalUnitId),
			schemaObjectCount:          basetypes.NewInt64PointerValue(item.ObjectCount),
			schemaSizeBytes:            basetypes.NewInt64PointerValue(item.SizeBytes),
			schemaProtectionGroupCount: basetypes.NewInt64PointerValue(item.ProtectionGroupCount),
			schemaIsDeleted:            basetypes.NewBoolPointerValue(item.IsDeleted),
			schemaIsVersioningEnabled:  basetypes.NewBoolPointerValue(item.IsVersioningEnabled),
			schemaCreatedTimestamp:     basetypes.NewStringPointerValue(item.CreatedTimestamp),
			schemaLastBackupTimestamp:  basetypes.NewStringPointerValue(item.LastBackupTimestamp),
			schemaLabels:               labelsSet,
		}

		obj, conversionDiags := types.ObjectValue(attrTypes, attrValues)
		diags.Append(conversionDiags...)
		if diags.HasError() {
			return diags
		}
		attrVals = append(attrVals, obj)
	}

	setObj, setDiags := types.SetValue(objType, attrVals)
	diags.Append(setDiags...)
	model.GCSBuckets = setObj

	return diags
}
