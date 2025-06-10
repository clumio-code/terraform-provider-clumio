// Copyright 2025. Clumio, Inc.

// This file contains the unit tests for the functions in utils.go

//go:build unit

package clumio_report_configuration

import (
	"testing"

	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
)

// Unit test for the following Notification mapping cases:
//   - Mapping with an email list.
//   - Mapping with an empty email list.
//   - Mapping with nil object.
func TestMapSchemaNotificationToClumioNotification(t *testing.T) {

	// Tests that the schema notification gets converted to a SDK model notification
	t.Run("with an email list", func(t *testing.T) {
		schemaNotification := []*notificationModel{
			{
				EmailList: []types.String{
					basetypes.NewStringValue("email1"),
					basetypes.NewStringValue("email2"),
				},
			},
		}

		clumioNotification := mapSchemaNotificationToClumioNotification(schemaNotification)
		assert.NotNil(t, clumioNotification)
		assert.Equal(t, 2, len(clumioNotification.EmailList))
		assert.Equal(t, "email1", *clumioNotification.EmailList[0])
		assert.Equal(t, "email2", *clumioNotification.EmailList[1])
	})

	// Tests that the schema notification gets converted to a SDK model notification
	// with an empty email list.
	t.Run("with an empty email list", func(t *testing.T) {
		schemaNotificationEmpty := []*notificationModel{{EmailList: []types.String{}}}

		clumioNotification := mapSchemaNotificationToClumioNotification(schemaNotificationEmpty)
		assert.NotNil(t, clumioNotification)
		assert.Equal(t, 0, len(clumioNotification.EmailList))
	})

	// Tests that the schema notification with nil object
	t.Run("with nil object", func(t *testing.T) {
		clumioNotification := mapSchemaNotificationToClumioNotification(nil)
		assert.Nil(t, clumioNotification)
	})
}

// Unit test for the following Schedule mapping cases:
//   - Mapping with all attributes.
//   - Mapping with only the required attributes.
//   - Mapping with nil object.
func TestMapSchemaScheduleToClumioSchedule(t *testing.T) {
	// Tests that the schema schedule gets converted to a SDK model schedule
	t.Run("with all attributes", func(t *testing.T) {
		schemaSchedule := []*scheduleModel{
			{
				DayOfMonth: types.Int64Value(1),
				DayOfWeek:  types.StringValue("Monday"),
				Frequency:  types.StringValue("weekly"),
				StartTime:  types.StringValue("10:00"),
				Timezone:   types.StringValue("UTC"),
			},
		}

		clumioSchedule := mapSchemaScheduleToClumioSchedule(schemaSchedule)
		assert.NotNil(t, clumioSchedule)
		assert.Equal(t, int64(1), *clumioSchedule.DayOfMonth)
		assert.Equal(t, "Monday", *clumioSchedule.DayOfWeek)
		assert.Equal(t, "weekly", *clumioSchedule.Frequency)
		assert.Equal(t, "10:00", *clumioSchedule.StartTime)
		assert.Equal(t, "UTC", *clumioSchedule.Timezone)
	})

	// Tests that the schema schedule gets converted to a SDK model schedule with only the required attributes.
	t.Run("with only required attributes", func(t *testing.T) {
		schemaScheduleRequired := []*scheduleModel{{StartTime: types.StringValue("08:00")}}

		clumioSchedule := mapSchemaScheduleToClumioSchedule(schemaScheduleRequired)
		assert.NotNil(t, clumioSchedule)
		assert.Equal(t, "08:00", *clumioSchedule.StartTime)
	})

	// Tests that the schema notification with nil object
	t.Run("with nil object", func(t *testing.T) {
		clumioSchedule := mapSchemaScheduleToClumioSchedule(nil)
		assert.Nil(t, clumioSchedule)
	})
}

// Unit test for the Parameter mapping cases:
//   - Mapping with all attributes.
//   - Mapping with nil object.
func TestMapSchemaParameterToClumioParameter(t *testing.T) {
	testTimeUnit := []*timeUnitModel{
		{
			Unit:  types.StringValue("days"),
			Value: types.Int32Value(7),
		},
	}
	// Tests that the schema parameter gets converted to a SDK model parameter
	t.Run("with all attributes", func(t *testing.T) {
		schemaParameter := []*parameterModel{
			{
				Controls: []*controlsModel{
					{
						AssetBackupControl: []*assetBackupControl{
							{
								LookBackPeriod:           testTimeUnit,
								MinimumRetentionDuration: testTimeUnit,
								WindowSize:               testTimeUnit,
							},
						},
						AssetProtectionControl: []*assetProtectionControl{
							{
								ShouldIgnoreDeactivatedPolicy: types.BoolValue(true),
							},
						},
						PolicyControl: []*policyControl{
							{
								MinimumRetentionDuration: testTimeUnit,
								MinimumRpoFrequency:      testTimeUnit,
							},
						},
					},
				},
				Filters: []*filtersModel{
					{
						AssetFilter: []*assetFilter{
							{
								GroupsFilter: []*assetGroupFilter{
									{
										ID:     types.StringValue("group-id-123"),
										Region: types.StringValue("us-west-2"),
										Type:   types.StringValue("ec2"),
									},
									{
										ID:     types.StringValue("group-id-456"),
										Region: types.StringValue("us-east-1"),
										Type:   types.StringValue("s3"),
									},
								},
								TagOpMode: types.StringValue("equal"),
								TagsFilter: []*assetTagFilter{
									{
										Key:   types.StringValue("Environment"),
										Value: types.StringValue("Production"),
									},
									{
										Key:   types.StringValue("Department"),
										Value: types.StringValue("Engineering"),
									},
								},
							},
						},
						CommonFilter: []*commonFilter{
							{
								AssetTypes: []types.String{
									types.StringValue("aws_ec2_instance"),
									types.StringValue("microsoft365_drive"),
								},
								DataSources: []types.String{
									types.StringValue("aws"),
									types.StringValue("microsoft365"),
								},
								OrganizationalUnits: []types.String{
									types.StringValue("ou-123456"),
									types.StringValue("ou-789012"),
								},
							},
						},
					},
				},
			},
		}

		clumioParamter := mapSchemaParameterToClumioParameter(schemaParameter)
		assert.NotNil(t, clumioParamter)
		assert.Equal(t, "days", *clumioParamter.Controls.AssetBackup.LookBackPeriod.Unit)
		assert.Equal(t, int32(7), *clumioParamter.Controls.AssetBackup.LookBackPeriod.Value)
		assert.Equal(t, "days", *clumioParamter.Controls.AssetBackup.MinimumRetentionDuration.Unit)
		assert.Equal(t, int32(7),
			*clumioParamter.Controls.AssetBackup.MinimumRetentionDuration.Value)
		assert.Equal(t, "days", *clumioParamter.Controls.AssetBackup.WindowSize.Unit)
		assert.Equal(t, int32(7), *clumioParamter.Controls.AssetBackup.WindowSize.Value)
		assert.True(t, *clumioParamter.Controls.AssetProtection.ShouldIgnoreDeactivatedPolicy)
		assert.Equal(t, "days", *clumioParamter.Controls.Policy.MinimumRetentionDuration.Unit)
		assert.Equal(t, int32(7), *clumioParamter.Controls.Policy.MinimumRetentionDuration.Value)
		assert.Equal(t, "days", *clumioParamter.Controls.Policy.MinimumRpoFrequency.Unit)
		assert.Equal(t, int32(7), *clumioParamter.Controls.Policy.MinimumRpoFrequency.Value)
		assert.Equal(t, 2, len(clumioParamter.Filters.Asset.Groups))
		assert.Equal(t, "group-id-123", *clumioParamter.Filters.Asset.Groups[0].Id)
		assert.Equal(t, "us-west-2", *clumioParamter.Filters.Asset.Groups[0].Region)
		assert.Equal(t, "ec2", *clumioParamter.Filters.Asset.Groups[0].ClumioType)
		assert.Equal(t, "group-id-456", *clumioParamter.Filters.Asset.Groups[1].Id)
		assert.Equal(t, "us-east-1", *clumioParamter.Filters.Asset.Groups[1].Region)
		assert.Equal(t, "s3", *clumioParamter.Filters.Asset.Groups[1].ClumioType)
		assert.Equal(t, 2, len(clumioParamter.Filters.Asset.Tags))
		assert.Equal(t, "Environment", *clumioParamter.Filters.Asset.Tags[0].Key)
		assert.Equal(t, "Production", *clumioParamter.Filters.Asset.Tags[0].Value)
		assert.Equal(t, "Department", *clumioParamter.Filters.Asset.Tags[1].Key)
		assert.Equal(t, "Engineering", *clumioParamter.Filters.Asset.Tags[1].Value)
		assert.Equal(t, "equal", *clumioParamter.Filters.Asset.TagOpMode)
		assert.Equal(t, 2, len(clumioParamter.Filters.Common.AssetTypes))
		assert.Equal(t, "aws_ec2_instance", *clumioParamter.Filters.Common.AssetTypes[0])
		assert.Equal(t, "microsoft365_drive", *clumioParamter.Filters.Common.AssetTypes[1])
		assert.Equal(t, 2, len(clumioParamter.Filters.Common.DataSources))
		assert.Equal(t, "aws", *clumioParamter.Filters.Common.DataSources[0])
		assert.Equal(t, "microsoft365", *clumioParamter.Filters.Common.DataSources[1])
		assert.Equal(t, 2, len(clumioParamter.Filters.Common.OrganizationalUnits))
		assert.Equal(t, "ou-123456", *clumioParamter.Filters.Common.OrganizationalUnits[0])
		assert.Equal(t, "ou-789012", *clumioParamter.Filters.Common.OrganizationalUnits[1])
	})

	// Tests that the schema parameter with nil object returns nil
	t.Run("with nil object", func(t *testing.T) {
		clumioParameter := mapSchemaParameterToClumioParameter(nil)

		assert.Nil(t, clumioParameter)
	})
}

// Unit test for the Controls nil check.
func TestMapSchemaControlsToClumioControls(t *testing.T) {
	// Tests that the schema controls with nil object returns nil
	t.Run("with nil object", func(t *testing.T) {
		clumioControls := mapSchemaControlsToClumioControls(nil)
		assert.Nil(t, clumioControls)
	})

	// Tests that the schema controls with empty object returns empty object
	t.Run("with empty object", func(t *testing.T) {
		schemaControls := []*controlsModel{
			{
				AssetBackupControl:     nil,
				AssetProtectionControl: nil,
				PolicyControl:          nil,
			},
		}

		clumioControls := mapSchemaControlsToClumioControls(schemaControls)
		assert.NotNil(t, clumioControls)
		assert.Nil(t, clumioControls.AssetBackup)
		assert.Nil(t, clumioControls.AssetProtection)
		assert.Nil(t, clumioControls.Policy)
	})
}

// Unit test for the Filters nil check.
func TestMapSchemaFiltersToClumioFilters(t *testing.T) {
	// Tests that the schema filters with nil object returns nil
	t.Run("with nil object", func(t *testing.T) {
		clumioFilters := mapSchemaFiltersToClumioFilters(nil)
		assert.Nil(t, clumioFilters)
	})

	// Tests that the schema filters with empty object returns empty object
	t.Run("with empty object", func(t *testing.T) {
		schemaFilters := []*filtersModel{
			{
				AssetFilter:  nil,
				CommonFilter: nil,
			},
		}

		clumioFilters := mapSchemaFiltersToClumioFilters(schemaFilters)
		assert.NotNil(t, clumioFilters)
		assert.Nil(t, clumioFilters.Asset)
		assert.Nil(t, clumioFilters.Common)
	})
}

// Unit test for the TimeUnit nil check.
func TestMapSchemaTimeUnitToClumioTimeUnit(t *testing.T) {
	// Tests that the schema time unit with nil object returns nil
	t.Run("with nil object", func(t *testing.T) {
		clumioTimeUnit := mapSchemaTimeUnitToClumioTimeUnit(nil)
		assert.Nil(t, clumioTimeUnit)
	})
}

// Unit test for the following Notification mapping cases:
//   - Mapping with an email list.
//   - Mapping with an empty email list.
//   - Mapping with nil object.
func TestMapClumioNotificationToSchemaNotification(t *testing.T) {
	// Tests that the SDK model notification gets converted to a schema notification
	t.Run("with an email list", func(t *testing.T) {
		email1 := "email1"
		email2 := "email2"
		clumioNotification := &models.NotificationSetting{
			EmailList: []*string{&email1, &email2},
		}

		schemaNotification := mapClumioNotificationToSchemaNotification(clumioNotification)
		assert.NotNil(t, schemaNotification)
		assert.Equal(t, 2, len(schemaNotification[0].EmailList))
		assert.Equal(t, email1, schemaNotification[0].EmailList[0].ValueString())
		assert.Equal(t, email2, schemaNotification[0].EmailList[1].ValueString())
	})

	// Tests that the SDK model notification gets converted with an empty email list
	t.Run("with an empty email list", func(t *testing.T) {
		clumioNotificationEmpty := &models.NotificationSetting{EmailList: nil}

		schemaNotification := mapClumioNotificationToSchemaNotification(clumioNotificationEmpty)
		assert.NotNil(t, schemaNotification)
		assert.Equal(t, 0, len(schemaNotification[0].EmailList))
	})

	// Tests that the SDK model notification with nil object
	t.Run("with nil object", func(t *testing.T) {
		schemaNotification := mapClumioNotificationToSchemaNotification(nil)
		assert.Nil(t, schemaNotification)
	})
}

// Unit test for the following Schedule mapping cases:
//   - Mapping with all attributes.
//   - Mapping with only the required attributes.
//   - Mapping with nil object.
func TestMapClumioScheduleToSchemaSchedule(t *testing.T) {
	dayOfMonth := int64(1)
	dayOfWeek := "monday"
	frequency := "weekly"
	startTime := "10:00"
	timezone := "UTC"

	// Tests that the SDK model schedule gets converted to a schema schedule
	t.Run("with all attributes", func(t *testing.T) {
		clumioSchedule := &models.ScheduleSetting{
			DayOfMonth: &dayOfMonth,
			DayOfWeek:  &dayOfWeek,
			Frequency:  &frequency,
			StartTime:  &startTime,
			Timezone:   &timezone,
		}

		schemaSchedule := mapClumioScheduleToSchemaSchedule(clumioSchedule)
		assert.NotNil(t, schemaSchedule)
		assert.Equal(t, dayOfMonth, schemaSchedule[0].DayOfMonth.ValueInt64())
		assert.Equal(t, dayOfWeek, schemaSchedule[0].DayOfWeek.ValueString())
		assert.Equal(t, frequency, schemaSchedule[0].Frequency.ValueString())
		assert.Equal(t, startTime, schemaSchedule[0].StartTime.ValueString())
		assert.Equal(t, timezone, schemaSchedule[0].Timezone.ValueString())
	})

	// Tests that the SDK model schedule gets converted to a schema schedule with only the required attributes.
	t.Run("with only required attributes", func(t *testing.T) {
		clumioScheduleRequired := &models.ScheduleSetting{StartTime: &startTime}

		schemaSchedule := mapClumioScheduleToSchemaSchedule(clumioScheduleRequired)
		assert.NotNil(t, schemaSchedule)
		assert.Equal(t, startTime, schemaSchedule[0].StartTime.ValueString())
	})

	// Tests that the SDK model schedule with nil object
	t.Run("with nil object", func(t *testing.T) {
		schemaSchedule := mapClumioScheduleToSchemaSchedule(nil)
		assert.Nil(t, schemaSchedule)
	})
}

// Unit test for the Parameter mapping cases:
//   - Mapping with all attributes.
//   - Mapping with nil object.
func TestMapClumioParameterToSchemaParameter(t *testing.T) {
	unit := "days"
	value := int32(7)
	testTimeUnit := &models.TimeUnitParam{
		Unit:  &unit,
		Value: &value,
	}
	shouldIgnoreDeactivatedPolicy := true
	// Common filter attributes
	assetType1 := "aws_ec2_instance"
	assetType2 := "microsoft365_drive"
	dataSource1 := "aws"
	dataSource2 := "microsoft365"
	organizationalUnit1 := "ou-123456"
	organizationalUnit2 := "ou-789012"
	// Asset group filter attributes
	assetGroupFilterId1 := "group-id-123"
	assetGroupFilterRegion1 := "us-west-2"
	assetGroupFilterType1 := "ec2"
	assetGroupFilterId2 := "group-id-456"
	assetGroupFilterRegion2 := "us-east-1"
	assetGroupFilterType2 := "s3"
	// Asset tag filter attributes
	assetTagFilterKey1 := "Environment"
	assetTagFilterValue1 := "Production"
	assetTagFilterKey2 := "Department"
	assetTagFilterValue2 := "Engineering"
	tagOpMode := "equal"

	// Tests that the SDK model parameter gets converted to a schema parameter
	t.Run("with all attributes", func(t *testing.T) {
		clumioParameter := &models.Parameter{
			Controls: &models.ComplianceControls{
				AssetBackup: &models.AssetBackupControl{
					LookBackPeriod:           testTimeUnit,
					MinimumRetentionDuration: testTimeUnit,
					WindowSize:               testTimeUnit,
				},
				AssetProtection: &models.AssetProtectionControl{
					ShouldIgnoreDeactivatedPolicy: &shouldIgnoreDeactivatedPolicy,
				},
				Policy: &models.PolicyControl{
					MinimumRetentionDuration: testTimeUnit,
					MinimumRpoFrequency:      testTimeUnit,
				},
			},
			Filters: &models.ComplianceFilters{
				Common: &models.CommonFilter{
					AssetTypes:          []*string{&assetType1, &assetType2},
					DataSources:         []*string{&dataSource1, &dataSource2},
					OrganizationalUnits: []*string{&organizationalUnit1, &organizationalUnit2},
				},
				Asset: &models.AssetFilter{
					Groups: []*models.AssetGroupFilter{
						{Id: &assetGroupFilterId1, Region: &assetGroupFilterRegion1, ClumioType: &assetGroupFilterType1},
						{Id: &assetGroupFilterId2, Region: &assetGroupFilterRegion2, ClumioType: &assetGroupFilterType2},
					},
					TagOpMode: &tagOpMode,
					Tags: []*models.Tag{
						{Key: &assetTagFilterKey1, Value: &assetTagFilterValue1},
						{Key: &assetTagFilterKey2, Value: &assetTagFilterValue2},
					},
				},
			},
		}

		schemaParameter := mapClumioParameterToSchemaParameter(clumioParameter)

		assert.NotNil(t, schemaParameter)
		assert.Equal(t, unit, schemaParameter[0].Controls[0].AssetBackupControl[0].
			LookBackPeriod[0].Unit.ValueString())
		assert.Equal(t, value, schemaParameter[0].Controls[0].AssetBackupControl[0].
			LookBackPeriod[0].Value.ValueInt32())
		assert.Equal(t, unit, schemaParameter[0].Controls[0].AssetBackupControl[0].
			MinimumRetentionDuration[0].Unit.ValueString())
		assert.Equal(t, value, schemaParameter[0].Controls[0].AssetBackupControl[0].
			MinimumRetentionDuration[0].Value.ValueInt32())
		assert.Equal(t, unit, schemaParameter[0].Controls[0].AssetBackupControl[0].WindowSize[0].
			Unit.ValueString())
		assert.Equal(t, value, schemaParameter[0].Controls[0].AssetBackupControl[0].WindowSize[0].
			Value.ValueInt32())

		assert.True(t, schemaParameter[0].Controls[0].AssetProtectionControl[0].
			ShouldIgnoreDeactivatedPolicy.ValueBool())

		assert.Equal(t, unit, schemaParameter[0].Controls[0].PolicyControl[0].
			MinimumRetentionDuration[0].Unit.ValueString())
		assert.Equal(t, value, schemaParameter[0].Controls[0].PolicyControl[0].
			MinimumRetentionDuration[0].Value.ValueInt32())
		assert.Equal(t, unit, schemaParameter[0].Controls[0].PolicyControl[0].
			MinimumRpoFrequency[0].Unit.ValueString())
		assert.Equal(t, value, schemaParameter[0].Controls[0].PolicyControl[0].
			MinimumRpoFrequency[0].Value.ValueInt32())

		assert.Equal(t, 2, len(schemaParameter[0].Filters[0].AssetFilter[0].GroupsFilter))
		assert.Equal(t, assetGroupFilterId1, schemaParameter[0].Filters[0].AssetFilter[0].
			GroupsFilter[0].ID.ValueString())
		assert.Equal(t, assetGroupFilterRegion1, schemaParameter[0].Filters[0].AssetFilter[0].
			GroupsFilter[0].Region.ValueString())
		assert.Equal(t, assetGroupFilterType1, schemaParameter[0].Filters[0].AssetFilter[0].
			GroupsFilter[0].Type.ValueString())
		assert.Equal(t, assetGroupFilterId2, schemaParameter[0].Filters[0].AssetFilter[0].
			GroupsFilter[1].ID.ValueString())
		assert.Equal(t, assetGroupFilterRegion2, schemaParameter[0].Filters[0].AssetFilter[0].
			GroupsFilter[1].Region.ValueString())
		assert.Equal(t, assetGroupFilterType2, schemaParameter[0].Filters[0].AssetFilter[0].
			GroupsFilter[1].Type.ValueString())

		assert.Equal(t, 2, len(schemaParameter[0].Filters[0].AssetFilter[0].TagsFilter))
		assert.Equal(t, assetTagFilterKey1, schemaParameter[0].Filters[0].AssetFilter[0].
			TagsFilter[0].Key.ValueString())
		assert.Equal(t, assetTagFilterValue1, schemaParameter[0].Filters[0].AssetFilter[0].
			TagsFilter[0].Value.ValueString())
		assert.Equal(t, assetTagFilterKey2, schemaParameter[0].Filters[0].AssetFilter[0].
			TagsFilter[1].Key.ValueString())
		assert.Equal(t, assetTagFilterValue2, schemaParameter[0].Filters[0].AssetFilter[0].
			TagsFilter[1].Value.ValueString())

		assert.Equal(t, tagOpMode, schemaParameter[0].Filters[0].AssetFilter[0].
			TagOpMode.ValueString())

		assert.Equal(t, 2, len(schemaParameter[0].Filters[0].CommonFilter[0].AssetTypes))
		assert.Equal(t, assetType1, schemaParameter[0].Filters[0].CommonFilter[0].
			AssetTypes[0].ValueString())
		assert.Equal(t, assetType2, schemaParameter[0].Filters[0].CommonFilter[0].
			AssetTypes[1].ValueString())

		assert.Equal(t, 2, len(schemaParameter[0].Filters[0].CommonFilter[0].DataSources))
		assert.Equal(t, dataSource1, schemaParameter[0].Filters[0].CommonFilter[0].
			DataSources[0].ValueString())
		assert.Equal(t, dataSource2, schemaParameter[0].Filters[0].CommonFilter[0].
			DataSources[1].ValueString())

		assert.Equal(t, 2, len(schemaParameter[0].Filters[0].CommonFilter[0].OrganizationalUnits))
		assert.Equal(t, organizationalUnit1, schemaParameter[0].Filters[0].CommonFilter[0].
			OrganizationalUnits[0].ValueString())
		assert.Equal(t, organizationalUnit2, schemaParameter[0].Filters[0].CommonFilter[0].
			OrganizationalUnits[1].ValueString())
	})

	// Tests that the SDK model parameter with nil object returns nil
	t.Run("with nil object", func(t *testing.T) {
		schemaParameter := mapClumioParameterToSchemaParameter(nil)
		assert.Nil(t, schemaParameter)
	})
}

// Unit test for the Controls nil check.
func TestMapClumioControlsToSchemaControls(t *testing.T) {
	// Tests that the SDK model controls with nil object returns nil
	t.Run("with nil object", func(t *testing.T) {
		schemaControls := mapClumioControlsToSchemaControls(nil)
		assert.Nil(t, schemaControls)
	})

	// Tests that the SDK model controls with empty object returns empty object
	t.Run("with empty object", func(t *testing.T) {
		clumioControls := &models.ComplianceControls{
			AssetBackup:     nil,
			AssetProtection: nil,
			Policy:          nil,
		}

		schemaControls := mapClumioControlsToSchemaControls(clumioControls)
		assert.NotNil(t, schemaControls)
		assert.Nil(t, schemaControls[0].AssetBackupControl)
		assert.Nil(t, schemaControls[0].AssetProtectionControl)
		assert.Nil(t, schemaControls[0].PolicyControl)
	})
}

// Unit test for the Filters nil check.
func TestMapClumioFiltersToSchemaFilters(t *testing.T) {
	// Tests that the SDK model filters with nil object returns nil
	t.Run("with nil object", func(t *testing.T) {
		schemaFilters := mapClumioFiltersToSchemaFilters(nil)
		assert.Nil(t, schemaFilters)
	})
}

// Unit test for the AssetFilter nil check.
func TestMapClumioAssetFilterToSchemaAssetFilter(t *testing.T) {
	// Tests that the SDK model asset filter with nil object returns nil
	t.Run("with nil object", func(t *testing.T) {
		schemaAssetFilter := mapClumioAssetFilterToSchemaAssetFilter(nil)
		assert.Nil(t, schemaAssetFilter)
	})
	// Tests that the SDK model asset filter with empty object returns empty object
	t.Run("with empty object", func(t *testing.T) {
		clumioAssetFilter := &models.AssetFilter{
			Groups:    nil,
			TagOpMode: nil,
			Tags:      nil,
		}

		schemaAssetFilter := mapClumioAssetFilterToSchemaAssetFilter(clumioAssetFilter)
		assert.Nil(t, schemaAssetFilter)
	})
}

// Unit test for the CommonFilter nil check.
func TestMapClumioCommonFilterToSchemaCommonFilter(t *testing.T) {
	// Tests that the SDK model common filter with nil object returns nil
	t.Run("with nil object", func(t *testing.T) {
		schemaCommonFilter := mapClumioCommonFilterToSchemaCommonFilter(nil)
		assert.Nil(t, schemaCommonFilter)
	})

	// Tests that the SDK model common filter with empty object returns empty object
	t.Run("with empty object", func(t *testing.T) {
		clumioCommonFilter := &models.CommonFilter{
			AssetTypes:          nil,
			DataSources:         nil,
			OrganizationalUnits: nil,
		}

		schemaCommonFilter := mapClumioCommonFilterToSchemaCommonFilter(clumioCommonFilter)
		assert.Nil(t, schemaCommonFilter)
	})
}

// Unit test for the TimeUnit nil check.
func TestMapClumioTimeUnitToSchemaTimeUnit(t *testing.T) {
	// Tests that the SDK model time unit with nil object returns nil
	t.Run("with nil object", func(t *testing.T) {
		schemaTimeUnit := mapClumioTimeUnitToSchemaTimeUnit(nil)
		assert.Nil(t, schemaTimeUnit)
	})
}
