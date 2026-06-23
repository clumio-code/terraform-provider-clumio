// Copyright 2026. Clumio, Inc.

// This file holds various utility functions used by the clumio_gcs_protection_group Terraform
// resource and data source, including the mapping between the Terraform bucket_rule blocks and the
// SDK GCPBucketRuleModel.

package clumio_gcs_protection_group

import (
	"context"

	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// setToStringSlice converts a Terraform string set into a slice of string pointers for SDK request
// bodies. A null or unknown set yields a nil slice.
func setToStringSlice(ctx context.Context, set types.Set) ([]*string, diag.Diagnostics) {
	var diags diag.Diagnostics
	if set.IsNull() || set.IsUnknown() {
		return nil, diags
	}
	var values []string
	diags = set.ElementsAs(ctx, &values, false)
	if diags.HasError() {
		return nil, diags
	}
	result := make([]*string, 0, len(values))
	for i := range values {
		result = append(result, &values[i])
	}
	return result, diags
}

// stringSliceToSet converts a slice of string pointers from an SDK response into a Terraform string
// set (empty set when the slice is empty). Used for computed attributes.
func stringSliceToSet(values []*string) (types.Set, diag.Diagnostics) {
	elems := make([]attr.Value, 0, len(values))
	for _, v := range values {
		elems = append(elems, types.StringPointerValue(v))
	}
	return types.SetValue(types.StringType, elems)
}

// stringSliceToNullableSet behaves like stringSliceToSet but returns a null set for an empty slice.
// Used for optional input fields so an unset value round-trips as null (no spurious diff).
func stringSliceToNullableSet(values []*string) (types.Set, diag.Diagnostics) {
	if len(values) == 0 {
		return types.SetNull(types.StringType), nil
	}
	return stringSliceToSet(values)
}

// mapSchemaBucketRuleToClumioBucketRule converts the bucket_rule block into the SDK
// GCPBucketRuleModel for create/update.
func mapSchemaBucketRuleToClumioBucketRule(
	ctx context.Context, m *bucketRuleModel) (*models.GCPBucketRuleModel, diag.Diagnostics) {

	var diags diag.Diagnostics
	if m == nil || (m.GcpLabel == nil && m.GcpLocation == nil && m.GcpProjectId == nil) {
		return nil, diags
	}
	rule := &models.GCPBucketRuleModel{}
	if m.GcpLabel != nil {
		op, d := mapSchemaLabelOperatorToClumioLabelOperator(ctx, m.GcpLabel)
		diags.Append(d...)
		rule.GcpLabel = op
	}
	if m.GcpLocation != nil {
		op, d := mapSchemaStringOperatorToClumioStringOperator(ctx, m.GcpLocation)
		diags.Append(d...)
		rule.GcpLocation = op
	}
	if m.GcpProjectId != nil {
		op, d := mapSchemaStringOperatorToClumioStringOperator(ctx, m.GcpProjectId)
		diags.Append(d...)
		rule.GcpProjectId = op
	}
	return rule, diags
}

func mapSchemaStringOperatorToClumioStringOperator(
	ctx context.Context, m *gcpStringOperatorModel) (
	*models.GCPStringOperatorModel, diag.Diagnostics) {

	var diags diag.Diagnostics
	in, d := setToStringSlice(ctx, m.In)
	diags.Append(d...)
	notIn, d := setToStringSlice(ctx, m.NotIn)
	diags.Append(d...)
	return &models.GCPStringOperatorModel{
		Eq:    m.Eq.ValueStringPointer(),
		In:    in,
		NotEq: m.NotEq.ValueStringPointer(),
		NotIn: notIn,
	}, diags
}

func mapSchemaLabelOperatorToClumioLabelOperator(
	ctx context.Context, m *gcpLabelOperatorModel) (*models.GCPLabelOperatorModel, diag.Diagnostics) {

	var diags diag.Diagnostics
	eq, d := mapSchemaLabelMapToClumioLabel(ctx, m.Eq)
	diags.Append(d...)
	contains, d := mapSchemaLabelMapToClumioLabel(ctx, m.Contains)
	diags.Append(d...)
	notEq, d := mapSchemaLabelMapToClumioLabel(ctx, m.NotEq)
	diags.Append(d...)
	notContains, d := mapSchemaLabelMapToClumioLabel(ctx, m.NotContains)
	diags.Append(d...)
	all, d := mapSchemaLabelMapToClumioLabels(ctx, m.All)
	diags.Append(d...)
	notAll, d := mapSchemaLabelMapToClumioLabels(ctx, m.NotAll)
	diags.Append(d...)
	return &models.GCPLabelOperatorModel{
		Eq:          eq,
		Contains:    contains,
		In:          mapSchemaLabelsToClumioLabels(m.In),
		All:         all,
		NotEq:       notEq,
		NotContains: notContains,
		NotIn:       mapSchemaLabelsToClumioLabels(m.NotIn),
		NotAll:      notAll,
	}, diags
}

// mapSchemaLabelMapToClumioLabels converts a map(string) of GCP labels into a list of SDK label
// models. A null/empty map yields nil so the operator stays unset.
func mapSchemaLabelMapToClumioLabels(
	ctx context.Context, m types.Map) ([]*models.GCPBucketRuleLabelModel, diag.Diagnostics) {

	var diags diag.Diagnostics
	if m.IsNull() || m.IsUnknown() {
		return nil, diags
	}
	elems := make(map[string]string, len(m.Elements()))
	diags.Append(m.ElementsAs(ctx, &elems, false)...)
	if diags.HasError() {
		return nil, diags
	}
	if len(elems) == 0 {
		return nil, diags
	}
	out := make([]*models.GCPBucketRuleLabelModel, 0, len(elems))
	for k, v := range elems {
		key, val := k, v
		out = append(out, &models.GCPBucketRuleLabelModel{Key: &key, Value: &val})
	}
	return out, diags
}

// mapSchemaLabelMapToClumioLabel converts a single-entry map(string) into one SDK label model.
func mapSchemaLabelMapToClumioLabel(
	ctx context.Context, m types.Map) (*models.GCPBucketRuleLabelModel, diag.Diagnostics) {

	list, diags := mapSchemaLabelMapToClumioLabels(ctx, m)
	if len(list) == 0 {
		return nil, diags
	}
	return list[0], diags
}

// mapSchemaLabelToClumioLabel converts a single {key, value} block (used by in/not_in) into an SDK
// label model.
func mapSchemaLabelToClumioLabel(m *labelModel) *models.GCPBucketRuleLabelModel {
	if m == nil {
		return nil
	}
	return &models.GCPBucketRuleLabelModel{
		Key:   m.Key.ValueStringPointer(),
		Value: m.Value.ValueStringPointer(),
	}
}

func mapSchemaLabelsToClumioLabels(ms []*labelModel) []*models.GCPBucketRuleLabelModel {
	if len(ms) == 0 {
		return nil
	}
	out := make([]*models.GCPBucketRuleLabelModel, 0, len(ms))
	for _, m := range ms {
		out = append(out, mapSchemaLabelToClumioLabel(m))
	}
	return out
}

// mapClumioBucketRuleToSchemaBucketRule converts the SDK GCPBucketRuleModel back into the
// bucket_rule block for read. Returns nil when the rule (or all of its condition fields) is empty,
// so an unset rule stays null.
func mapClumioBucketRuleToSchemaBucketRule(
	ctx context.Context, rule *models.GCPBucketRuleModel) (*bucketRuleModel, diag.Diagnostics) {

	var diags diag.Diagnostics
	if rule == nil {
		return nil, diags
	}
	m := &bucketRuleModel{}
	op, d := mapClumioLabelOperatorToSchemaLabelOperator(ctx, rule.GcpLabel)
	diags.Append(d...)
	if op != nil {
		m.GcpLabel = op
	}
	if rule.GcpLocation != nil {
		op, d := mapClumioStringOperatorToSchemaStringOperator(rule.GcpLocation)
		diags.Append(d...)
		m.GcpLocation = op
	}
	if rule.GcpProjectId != nil {
		op, d := mapClumioStringOperatorToSchemaStringOperator(rule.GcpProjectId)
		diags.Append(d...)
		m.GcpProjectId = op
	}
	if m.GcpLabel == nil && m.GcpLocation == nil && m.GcpProjectId == nil {
		return nil, diags
	}
	return m, diags
}

func mapClumioStringOperatorToSchemaStringOperator(
	op *models.GCPStringOperatorModel) (*gcpStringOperatorModel, diag.Diagnostics) {

	var diags diag.Diagnostics
	in, d := stringSliceToNullableSet(op.In)
	diags.Append(d...)
	notIn, d := stringSliceToNullableSet(op.NotIn)
	diags.Append(d...)
	m := &gcpStringOperatorModel{
		Eq:    types.StringPointerValue(op.Eq),
		In:    in,
		NotEq: types.StringPointerValue(op.NotEq),
		NotIn: notIn,
	}
	// Treat an all-empty operator as unset.
	if m.Eq.IsNull() && m.NotEq.IsNull() && m.In.IsNull() && m.NotIn.IsNull() {
		return nil, diags
	}
	return m, diags
}

func mapClumioLabelOperatorToSchemaLabelOperator(
	ctx context.Context, op *models.GCPLabelOperatorModel) (*gcpLabelOperatorModel, diag.Diagnostics) {

	var diags diag.Diagnostics
	if op == nil {
		return nil, diags
	}
	eq, d := mapClumioLabelToSchemaMap(ctx, op.Eq)
	diags.Append(d...)
	contains, d := mapClumioLabelToSchemaMap(ctx, op.Contains)
	diags.Append(d...)
	notEq, d := mapClumioLabelToSchemaMap(ctx, op.NotEq)
	diags.Append(d...)
	notContains, d := mapClumioLabelToSchemaMap(ctx, op.NotContains)
	diags.Append(d...)
	all, d := mapClumioLabelsToSchemaMap(ctx, op.All)
	diags.Append(d...)
	notAll, d := mapClumioLabelsToSchemaMap(ctx, op.NotAll)
	diags.Append(d...)
	m := &gcpLabelOperatorModel{
		Eq:          eq,
		Contains:    contains,
		All:         all,
		NotEq:       notEq,
		NotContains: notContains,
		NotAll:      notAll,
		In:          mapClumioLabelsToSchemaLabels(op.In),
		NotIn:       mapClumioLabelsToSchemaLabels(op.NotIn),
	}
	if m.Eq.IsNull() && m.Contains.IsNull() && m.NotEq.IsNull() && m.NotContains.IsNull() &&
		m.All.IsNull() && m.NotAll.IsNull() && len(m.In) == 0 && len(m.NotIn) == 0 {
		return nil, diags
	}
	return m, diags
}

// mapClumioLabelsToSchemaMap converts a list of SDK labels into a map(string). An empty list yields
// a null map so the operator stays unset.
func mapClumioLabelsToSchemaMap(
	ctx context.Context, ls []*models.GCPBucketRuleLabelModel) (types.Map, diag.Diagnostics) {

	if len(ls) == 0 {
		return types.MapNull(types.StringType), nil
	}
	elems := make(map[string]string, len(ls))
	for _, l := range ls {
		if l == nil || l.Key == nil {
			continue
		}
		val := ""
		if l.Value != nil {
			val = *l.Value
		}
		elems[*l.Key] = val
	}
	if len(elems) == 0 {
		return types.MapNull(types.StringType), nil
	}
	return types.MapValueFrom(ctx, types.StringType, elems)
}

// mapClumioLabelToSchemaMap converts a single SDK label into a one-entry map(string).
func mapClumioLabelToSchemaMap(
	ctx context.Context, l *models.GCPBucketRuleLabelModel) (types.Map, diag.Diagnostics) {

	if l == nil {
		return types.MapNull(types.StringType), nil
	}
	return mapClumioLabelsToSchemaMap(ctx, []*models.GCPBucketRuleLabelModel{l})
}

// mapClumioLabelToSchemaLabel converts a single SDK label into a {key, value} block (used by
// in/not_in).
func mapClumioLabelToSchemaLabel(l *models.GCPBucketRuleLabelModel) *labelModel {
	if l == nil {
		return nil
	}
	return &labelModel{
		Key:   types.StringPointerValue(l.Key),
		Value: types.StringPointerValue(l.Value),
	}
}

func mapClumioLabelsToSchemaLabels(ls []*models.GCPBucketRuleLabelModel) []*labelModel {
	if len(ls) == 0 {
		return nil
	}
	out := make([]*labelModel, 0, len(ls))
	for _, l := range ls {
		out = append(out, mapClumioLabelToSchemaLabel(l))
	}
	return out
}

// mapProtectionInfo converts the SDK protection info into the Terraform protection_info list.
func mapProtectionInfo(
	protectionInfo *models.GCPProtectionInfoModel) (types.List, diag.Diagnostics) {

	objType := types.ObjectType{AttrTypes: map[string]attr.Type{schemaPolicyId: types.StringType}}
	if protectionInfo == nil {
		return types.ListValueMust(objType, []attr.Value{}), nil
	}
	obj, diags := types.ObjectValue(objType.AttrTypes, map[string]attr.Value{
		schemaPolicyId: types.StringPointerValue(protectionInfo.PolicyId),
	})
	if diags.HasError() {
		return types.ListNull(objType), diags
	}
	listObj, listDiags := types.ListValue(objType, []attr.Value{obj})
	diags.Append(listDiags...)
	return listObj, diags
}

// labelObjectType is the Terraform object type for a single {key, value} label.
func labelObjectType() types.ObjectType {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		schemaKey:   types.StringType,
		schemaValue: types.StringType,
	}}
}

// mapLabels converts the SDK labels into the computed labels set.
func mapLabels(labels []*models.GcpLabelModel) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics
	objType := labelObjectType()
	elems := make([]attr.Value, 0, len(labels))
	for _, l := range labels {
		obj, d := types.ObjectValue(objType.AttrTypes, map[string]attr.Value{
			schemaKey:   types.StringPointerValue(l.Key),
			schemaValue: types.StringPointerValue(l.Value),
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

// backupStatusStatsObjectType is the Terraform object type for backup_status_stats.
func backupStatusStatsObjectType() types.ObjectType {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		schemaFailureCount:        types.Int64Type,
		schemaNoBackupCount:       types.Int64Type,
		schemaPartialSuccessCount: types.Int64Type,
		schemaSuccessCount:        types.Int64Type,
	}}
}

// mapBackupStatusStats converts the SDK backup status stats into the computed object.
func mapBackupStatusStats(stats *models.BackupStatusStats) (types.Object, diag.Diagnostics) {
	objType := backupStatusStatsObjectType()
	if stats == nil {
		return types.ObjectNull(objType.AttrTypes), nil
	}
	return types.ObjectValue(objType.AttrTypes, map[string]attr.Value{
		schemaFailureCount:        types.Int64PointerValue(stats.FailureCount),
		schemaNoBackupCount:       types.Int64PointerValue(stats.NoBackupCount),
		schemaPartialSuccessCount: types.Int64PointerValue(stats.PartialSuccessCount),
		schemaSuccessCount:        types.Int64PointerValue(stats.SuccessCount),
	})
}
