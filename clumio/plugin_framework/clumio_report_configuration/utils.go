// Copyright 2025. Clumio, Inc.

// This file hold various utility functions used by the clumio_report_configuration Terraform resource.

package clumio_report_configuration

import (
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// mapSchemaNotificationToClumioNotification maps the schema notification to the Clumio API request
// notification.
func mapSchemaNotificationToClumioNotification(
	notificationSlice []*notificationModel) *models.NotificationSetting {

	if len(notificationSlice) == 0 {
		return nil
	}
	notification := notificationSlice[0]
	emailList := make([]*string, 0)
	for _, email := range notification.EmailList {
		emailList = append(emailList, email.ValueStringPointer())
	}

	return &models.NotificationSetting{
		EmailList: emailList,
	}
}

// mapSchemaScheduleToClumioSchedule maps the schema schedule to the Clumio API request schedule.
func mapSchemaScheduleToClumioSchedule(scheduleSlice []*scheduleModel) *models.ScheduleSetting {

	if len(scheduleSlice) == 0 {
		return nil
	}
	schedule := scheduleSlice[0]

	return &models.ScheduleSetting{
		DayOfMonth: schedule.DayOfMonth.ValueInt64Pointer(),
		DayOfWeek:  schedule.DayOfWeek.ValueStringPointer(),
		Frequency:  schedule.Frequency.ValueStringPointer(),
		StartTime:  schedule.StartTime.ValueStringPointer(),
		Timezone:   schedule.Timezone.ValueStringPointer(),
	}
}

// mapSchemaParameterToClumioParameter maps the schema parameter to the Clumio API request
// parameter.
func mapSchemaParameterToClumioParameter(parameterSlice []*parameterModel) *models.Parameter {

	if len(parameterSlice) == 0 {
		return nil
	}
	parameter := parameterSlice[0]

	return &models.Parameter{
		Controls: mapSchemaControlsToClumioControls(parameter.Controls),
		Filters:  mapSchemaFiltersToClumioFilters(parameter.Filters),
	}
}

// mapSchemaControlsToClumioControls maps the schema controls to the Clumio API request controls.
func mapSchemaControlsToClumioControls(controlsSlice []*controlsModel) *models.ComplianceControls {

	if len(controlsSlice) == 0 {
		return nil
	}
	controls := controlsSlice[0]

	var assetBackupControl *models.AssetBackupControl
	if controls.AssetBackupControl != nil {
		schemaAssetBackup := controls.AssetBackupControl[0]
		assetBackupControl = &models.AssetBackupControl{
			LookBackPeriod: mapSchemaTimeUnitToClumioTimeUnit(schemaAssetBackup.LookBackPeriod),
			MinimumRetentionDuration: mapSchemaTimeUnitToClumioTimeUnit(
				schemaAssetBackup.MinimumRetentionDuration),
			WindowSize: mapSchemaTimeUnitToClumioTimeUnit(schemaAssetBackup.WindowSize),
		}
	}

	var assetProtectionControl *models.AssetProtectionControl
	if controls.AssetProtectionControl != nil {
		assetProtectionControl = &models.AssetProtectionControl{
			ShouldIgnoreDeactivatedPolicy: controls.AssetProtectionControl[0].
				ShouldIgnoreDeactivatedPolicy.ValueBoolPointer(),
		}
	}

	var policyControl *models.PolicyControl
	if controls.PolicyControl != nil {
		schemaPolicyControl := controls.PolicyControl[0]
		policyControl = &models.PolicyControl{
			MinimumRetentionDuration: mapSchemaTimeUnitToClumioTimeUnit(
				schemaPolicyControl.MinimumRetentionDuration),
			MinimumRpoFrequency: mapSchemaTimeUnitToClumioTimeUnit(
				schemaPolicyControl.MinimumRpoFrequency),
		}
	}

	return &models.ComplianceControls{
		AssetBackup:     assetBackupControl,
		AssetProtection: assetProtectionControl,
		Policy:          policyControl,
	}
}

// mapSchemaFiltersToClumioFilters maps the schema filters to the Clumio API request filters.
func mapSchemaFiltersToClumioFilters(filtersSlice []*filtersModel) *models.ComplianceFilters {

	if len(filtersSlice) == 0 {
		return nil
	}
	filters := filtersSlice[0]

	var assetFilter *models.AssetFilter
	if filters.AssetFilter != nil {
		schemaAssetFilter := filters.AssetFilter[0]
		groupFilterList := make([]*models.AssetGroupFilter, 0)
		for _, groupFilter := range schemaAssetFilter.GroupsFilter {
			groupFilterList = append(groupFilterList, &models.AssetGroupFilter{
				Id:         groupFilter.ID.ValueStringPointer(),
				Region:     groupFilter.Region.ValueStringPointer(),
				ClumioType: groupFilter.Type.ValueStringPointer(),
			})
		}
		tagFilterList := make([]*models.Tag, 0)
		for _, tagFilter := range schemaAssetFilter.TagsFilter {
			tagFilterList = append(tagFilterList, &models.Tag{
				Key:   tagFilter.Key.ValueStringPointer(),
				Value: tagFilter.Value.ValueStringPointer(),
			})
		}
		assetFilter = &models.AssetFilter{
			Groups:    groupFilterList,
			TagOpMode: schemaAssetFilter.TagOpMode.ValueStringPointer(),
			Tags:      tagFilterList,
		}
	}

	var commonFilter *models.CommonFilter
	if filters.CommonFilter != nil {
		schemaCommonFilter := filters.CommonFilter[0]
		assetTypes := make([]*string, 0)
		for _, assetType := range schemaCommonFilter.AssetTypes {
			assetTypes = append(assetTypes, assetType.ValueStringPointer())
		}
		dataSources := make([]*string, 0)
		for _, dataSource := range schemaCommonFilter.DataSources {
			dataSources = append(dataSources, dataSource.ValueStringPointer())
		}
		organizationalUnits := make([]*string, 0)
		for _, ou := range schemaCommonFilter.OrganizationalUnits {
			organizationalUnits = append(organizationalUnits, ou.ValueStringPointer())
		}
		commonFilter = &models.CommonFilter{
			AssetTypes:          assetTypes,
			DataSources:         dataSources,
			OrganizationalUnits: organizationalUnits,
		}
	}

	return &models.ComplianceFilters{
		Asset:  assetFilter,
		Common: commonFilter,
	}
}

// mapSchemaTimeUnitToClumioTimeUnit converts the schema time unit to the Clumio API request time
// unit.
func mapSchemaTimeUnitToClumioTimeUnit(timeUnitSlice []*timeUnitModel) *models.TimeUnitParam {

	if len(timeUnitSlice) == 0 {
		return nil
	}
	timeUnit := timeUnitSlice[0]

	return &models.TimeUnitParam{
		Unit:  timeUnit.Unit.ValueStringPointer(),
		Value: timeUnit.Value.ValueInt32Pointer(),
	}
}

// mapClumioNotificationToSchemaNotification converts the Clumio API notification to the schema
// notification.
func mapClumioNotificationToSchemaNotification(
	notification *models.NotificationSetting) []*notificationModel {

	if notification == nil {
		return nil
	}

	schemaNotification := &notificationModel{}
	emailList := make([]types.String, 0)
	for _, email := range notification.EmailList {
		emailList = append(emailList, types.StringPointerValue(email))
	}
	schemaNotification.EmailList = emailList

	return []*notificationModel{schemaNotification}
}

// mapClumioParameterToSchemaParameter converts the Clumio API parameter to the schema parameter.
func mapClumioParameterToSchemaParameter(parameter *models.Parameter) []*parameterModel {

	if parameter == nil {
		return nil
	}

	schemaParameter := &parameterModel{
		Controls: mapClumioControlsToSchemaControls(parameter.Controls),
		Filters:  mapClumioFiltersToSchemaFilters(parameter.Filters),
	}

	return []*parameterModel{schemaParameter}
}

// mapClumioControlsToSchemaControls converts the Clumio API controls to the schema controls.
func mapClumioControlsToSchemaControls(controls *models.ComplianceControls) []*controlsModel {

	if controls == nil {
		return nil
	}

	schemaControls := &controlsModel{}
	if controls.AssetBackup != nil {
		schemaControls.AssetBackupControl = []*assetBackupControl{
			{
				LookBackPeriod: mapClumioTimeUnitToSchemaTimeUnit(
					controls.AssetBackup.LookBackPeriod),
				MinimumRetentionDuration: mapClumioTimeUnitToSchemaTimeUnit(
					controls.AssetBackup.MinimumRetentionDuration),
				WindowSize: mapClumioTimeUnitToSchemaTimeUnit(controls.AssetBackup.WindowSize),
			},
		}
	}
	if controls.AssetProtection != nil {
		schemaControls.AssetProtectionControl = []*assetProtectionControl{
			{
				ShouldIgnoreDeactivatedPolicy: types.BoolValue(
					*controls.AssetProtection.ShouldIgnoreDeactivatedPolicy),
			},
		}
	}
	if controls.Policy != nil {
		schemaControls.PolicyControl = []*policyControl{
			{
				MinimumRetentionDuration: mapClumioTimeUnitToSchemaTimeUnit(
					controls.Policy.MinimumRetentionDuration),
				MinimumRpoFrequency: mapClumioTimeUnitToSchemaTimeUnit(
					controls.Policy.MinimumRpoFrequency),
			},
		}
	}

	return []*controlsModel{schemaControls}
}

// mapClumioFiltersToSchemaFilters converts the Clumio API filters to the schema filters.
func mapClumioFiltersToSchemaFilters(filters *models.ComplianceFilters) []*filtersModel {

	if filters == nil {
		return nil
	}

	schemaFilters := &filtersModel{
		AssetFilter:  mapClumioAssetFilterToSchemaAssetFilter(filters.Asset),
		CommonFilter: mapClumioCommonFilterToSchemaCommonFilter(filters.Common),
	}

	if schemaFilters.AssetFilter != nil || schemaFilters.CommonFilter != nil {
		return []*filtersModel{schemaFilters}
	} else {
		return nil
	}
}

// mapClumioAssetFilterToSchemaAssetFilter converts the Clumio API asset filter to the schema asset
// filter.
func mapClumioAssetFilterToSchemaAssetFilter(asset *models.AssetFilter) []*assetFilter {

	if asset != nil {
		schemaAssetFilter := &assetFilter{}
		schemaGroupFilters := make([]*assetGroupFilter, 0)
		for _, groupFilter := range asset.Groups {
			schemaGroupFilters = append(schemaGroupFilters, &assetGroupFilter{
				ID:     types.StringPointerValue(groupFilter.Id),
				Region: types.StringPointerValue(groupFilter.Region),
				Type:   types.StringPointerValue(groupFilter.ClumioType),
			})
		}
		schemaAssetFilter.GroupsFilter = schemaGroupFilters
		schemaTagFilters := make([]*assetTagFilter, 0)
		for _, tagFilter := range asset.Tags {
			schemaTagFilters = append(schemaTagFilters, &assetTagFilter{
				Key:   types.StringPointerValue(tagFilter.Key),
				Value: types.StringPointerValue(tagFilter.Value),
			})
		}
		schemaAssetFilter.TagsFilter = schemaTagFilters
		if asset.TagOpMode != nil {
			schemaAssetFilter.TagOpMode = types.StringValue(*asset.TagOpMode)
		}
		// Only add the asset filter if at least one field is set.
		// This avoids creating an empty asset filter in the schema.
		if len(schemaGroupFilters) > 0 || len(schemaTagFilters) > 0 ||
			schemaAssetFilter.TagOpMode.ValueString() != "" {
			return []*assetFilter{schemaAssetFilter}
		}
	}

	return nil
}

func mapClumioCommonFilterToSchemaCommonFilter(common *models.CommonFilter) []*commonFilter {

	if common != nil {
		schemaCommonFilter := &commonFilter{}
		if common.AssetTypes != nil {
			schemaAssetTypes := make([]types.String, 0)
			for _, assetType := range common.AssetTypes {
				schemaAssetTypes = append(schemaAssetTypes, types.StringPointerValue(assetType))
			}
			schemaCommonFilter.AssetTypes = schemaAssetTypes
		}
		if common.DataSources != nil {
			schemaDataSources := make([]types.String, 0)
			for _, dataSource := range common.DataSources {
				schemaDataSources = append(schemaDataSources, types.StringPointerValue(dataSource))
			}
			schemaCommonFilter.DataSources = schemaDataSources
		}
		if common.OrganizationalUnits != nil {
			schemaOrganizationalUnits := make([]types.String, 0)
			for _, ou := range common.OrganizationalUnits {
				schemaOrganizationalUnits = append(schemaOrganizationalUnits,
					types.StringPointerValue(ou))
			}
			schemaCommonFilter.OrganizationalUnits = schemaOrganizationalUnits
		}
		// Only add the common filter if at least one field is set.
		// This avoids creating an empty common filter in the schema.
		if schemaCommonFilter.AssetTypes != nil || schemaCommonFilter.DataSources != nil ||
			schemaCommonFilter.OrganizationalUnits != nil {
			return []*commonFilter{schemaCommonFilter}
		}
	}

	return nil
}

// mapClumioScheduleToSchemaSchedule converts the Clumio API schedule to the schema schedule.
func mapClumioScheduleToSchemaSchedule(schedule *models.ScheduleSetting) []*scheduleModel {

	if schedule == nil {
		return nil
	}

	schemaSchedule := &scheduleModel{}
	if schedule.DayOfMonth != nil {
		schemaSchedule.DayOfMonth = types.Int64Value(*schedule.DayOfMonth)
	}
	if schedule.DayOfWeek != nil {
		schemaSchedule.DayOfWeek = types.StringValue(*schedule.DayOfWeek)
	}
	if schedule.Frequency != nil {
		schemaSchedule.Frequency = types.StringValue(*schedule.Frequency)
	}
	if schedule.StartTime != nil {
		schemaSchedule.StartTime = types.StringValue(*schedule.StartTime)
	}
	if schedule.Timezone != nil {
		schemaSchedule.Timezone = types.StringValue(*schedule.Timezone)
	}

	return []*scheduleModel{schemaSchedule}
}

// mapClumioTimeUnitToSchemaTimeUnit converts the Clumio API time unit to the schema time unit.
func mapClumioTimeUnitToSchemaTimeUnit(
	timeUnit *models.TimeUnitParam) []*timeUnitModel {

	if timeUnit == nil {
		return nil
	}

	schemaTimeUnit := &timeUnitModel{
		Unit:  types.StringValue(*timeUnit.Unit),
		Value: types.Int32Value(*timeUnit.Value),
	}

	return []*timeUnitModel{schemaTimeUnit}
}
