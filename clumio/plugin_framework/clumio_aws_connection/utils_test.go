// Copyright 2024. Clumio, Inc.

// This file contains the unit tests for the functions in utils.go

//go:build unit

package clumio_aws_connection

import (
	"context"
	"fmt"
	"testing"
	"time"

	apiutils "github.com/clumio-code/clumio-go-sdk/api_utils"
	sdkconfig "github.com/clumio-code/clumio-go-sdk/config"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Unit test for the following cases:
//   - Success scenario for getting the AWS environment for the connection.
//   - SDK API for listing AWS environments returns an error.
//   - SDK API for listing AWS environments returns no AWS environment.
//   - SDK API for listing AWS environments returns more than one AWS environment.
func TestGetEnvironmentForConnection(t *testing.T) {

	mockAwsEnvClient := sdkclients.NewMockAWSEnvironmentClient(t)
	ctx := context.Background()
	cr := &clumioAWSConnectionResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkEnvironments: mockAwsEnvClient,
	}

	// Populate the Clumio AWS connection resource model to be used as input for
	// getEnvironmentForConnection().
	state := &clumioAWSConnectionResourceModel{
		ID:              basetypes.NewStringValue(id),
		AccountNativeID: basetypes.NewStringValue(accountId),
		AWSRegion:       basetypes.NewStringValue(region),
		Description:     basetypes.NewStringValue(description),
	}

	// Tests that the get AWS environment is successful.
	t.Run("success scenario for get AWS environment", func(t *testing.T) {

		listEnvsResponse := &models.ListAWSEnvironmentsResponse{
			Embedded: &models.AWSEnvironmentListEmbedded{
				Items: []*models.AWSEnvironment{
					{
						Id: &envId,
					},
				},
			},
		}

		// Setup Expectations
		mockAwsEnvClient.EXPECT().ListAwsEnvironments(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Times(1).Return(listEnvsResponse, nil)

		env, err := getEnvironmentForConnection(ctx, cr, state)
		assert.Nil(t, err)
		assert.Equal(t, envId, *env.Id)
	})

	// Tests that the get AWS environment fails due to ListAwsEnvironments API returning an error.
	t.Run("list aws environments returns an error", func(t *testing.T) {

		apiError := &apiutils.APIError{
			ResponseCode: 500,
			Reason:       "test",
			Response:     []byte(testError),
		}

		// Setup Expectations
		mockAwsEnvClient.EXPECT().ListAwsEnvironments(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Times(1).Return(nil, apiError)

		env, err := getEnvironmentForConnection(ctx, cr, state)
		assert.NotNil(t, err)
		assert.Nil(t, env)
	})

	// Tests that the get AWS environment fails due to ListAwsEnvironments API returning no
	// environment.
	t.Run("list aws environments returns no environment", func(t *testing.T) {

		// Setup Expectations
		mockAwsEnvClient.EXPECT().ListAwsEnvironments(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Times(1).Return(&models.ListAWSEnvironmentsResponse{}, nil)

		env, err := getEnvironmentForConnection(ctx, cr, state)
		assert.NotNil(t, err)
		assert.Nil(t, env)
	})

	t.Run("falls back to global OU context when current context cannot see environment", func(t *testing.T) {
		previousFactory := newAWSEnvironmentClient
		previousClient := cr.client
		defer func() {
			newAWSEnvironmentClient = previousFactory
			cr.client = previousClient
		}()

		mockGlobalAwsEnvClient := sdkclients.NewMockAWSEnvironmentClient(t)
		newAWSEnvironmentClient = func(config sdkconfig.Config) sdkclients.AWSEnvironmentClient {
			assert.Equal(t, "", config.OrganizationalUnitContext)
			return mockGlobalAwsEnvClient
		}

		cr.client = &common.ApiClient{
			ClumioConfig: sdkconfig.Config{
				OrganizationalUnitContext: ouId,
			},
		}

		mockAwsEnvClient.EXPECT().ListAwsEnvironments(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Times(1).Return(&models.ListAWSEnvironmentsResponse{}, nil)
		mockGlobalAwsEnvClient.EXPECT().ListAwsEnvironments(
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Times(1).
			Return(&models.ListAWSEnvironmentsResponse{
				Embedded: &models.AWSEnvironmentListEmbedded{
					Items: []*models.AWSEnvironment{{Id: &envId}},
				},
			}, nil)

		env, err := getEnvironmentForConnection(ctx, cr, state)
		assert.Nil(t, err)
		assert.Equal(t, envId, *env.Id)
	})

	// Tests that the get AWS environment fails due to ListAwsEnvironments API returning more than
	// one environment.
	t.Run("list aws environments returns more than one environment", func(t *testing.T) {

		currentCount := int64(2)
		listEnvsResponse := &models.ListAWSEnvironmentsResponse{
			Embedded: &models.AWSEnvironmentListEmbedded{
				Items: []*models.AWSEnvironment{
					{
						Id: &envId,
					},
					{
						Id: &envId,
					},
				},
			},
			CurrentCount: &currentCount,
		}

		// Setup Expectations
		mockAwsEnvClient.EXPECT().ListAwsEnvironments(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Times(1).Return(listEnvsResponse, nil)

		env, err := getEnvironmentForConnection(ctx, cr, state)
		assert.NotNil(t, err)
		assert.Nil(t, env)
	})

	t.Run("list aws environments returns environment without id", func(t *testing.T) {
		listEnvsResponse := &models.ListAWSEnvironmentsResponse{
			Embedded: &models.AWSEnvironmentListEmbedded{
				Items: []*models.AWSEnvironment{{}},
			},
		}

		mockAwsEnvClient.EXPECT().ListAwsEnvironments(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Times(1).Return(listEnvsResponse, nil)

		env, err := getEnvironmentForConnection(ctx, cr, state)
		assert.NotNil(t, err)
		assert.Nil(t, env)
	})

}

func TestSetOrganizationalUnitID(t *testing.T) {
	t.Run("uses response OU when present", func(t *testing.T) {
		state := &clumioAWSConnectionResourceModel{}
		setOrganizationalUnitID(state, &ouId)
		assert.Equal(t, ouId, state.OrganizationalUnitID.ValueString())
	})

	t.Run("normalizes empty OU to global", func(t *testing.T) {
		state := &clumioAWSConnectionResourceModel{}
		empty := ""
		setOrganizationalUnitID(state, &empty)
		assert.Equal(t, defaultOrgUnitId, state.OrganizationalUnitID.ValueString())
	})

	t.Run("normalizes nil OU to global", func(t *testing.T) {
		state := &clumioAWSConnectionResourceModel{}
		setOrganizationalUnitID(state, nil)
		assert.Equal(t, defaultOrgUnitId, state.OrganizationalUnitID.ValueString())
	})
}

func TestGetOrgUnitForConnection(t *testing.T) {
	mockOrgUnitsClient := sdkclients.NewMockOrganizationalUnitClient(t)

	t.Run("success scenario for get OU", func(t *testing.T) {
		mockOrgUnitsClient.EXPECT().ReadOrganizationalUnit(ouId, mock.Anything).Times(1).Return(
			&models.ReadOrganizationalUnitResponse{Id: &ouId}, nil)

		orgUnit, err := getOrgUnitForConnection(mockOrgUnitsClient, ouId)
		assert.Nil(t, err)
		assert.Equal(t, ouId, *orgUnit.Id)
	})

	t.Run("read OU returns error", func(t *testing.T) {
		apiError := &apiutils.APIError{
			ResponseCode: 500,
			Reason:       "test",
			Response:     []byte(testError),
		}

		mockOrgUnitsClient.EXPECT().ReadOrganizationalUnit(parentOuId, mock.Anything).Times(1).Return(
			nil, apiError)

		orgUnit, err := getOrgUnitForConnection(mockOrgUnitsClient, parentOuId)
		assert.NotNil(t, err)
		assert.Nil(t, orgUnit)
	})
}

func TestUpdateOrgUnitForConnection(t *testing.T) {
	ctx := context.Background()

	t.Run("moves connection from global OU to child OU by add", func(t *testing.T) {
		mockAwsEnvClient := sdkclients.NewMockAWSEnvironmentClient(t)
		mockOrgUnitsClient := sdkclients.NewMockOrganizationalUnitClient(t)
		mockTaskClient := sdkclients.NewMockTaskClient(t)

		cr := &clumioAWSConnectionResource{
			client:          &common.ApiClient{ClumioConfig: sdkconfig.Config{}},
			sdkEnvironments: mockAwsEnvClient,
			sdkOrgUnits:     mockOrgUnitsClient,
			sdkTasks:        mockTaskClient,
			pollTimeout:     time.Second,
			pollInterval:    time.Millisecond,
		}

		plan := &clumioAWSConnectionResourceModel{
			OrganizationalUnitID: basetypes.NewStringValue(parentOuId),
		}
		state := &clumioAWSConnectionResourceModel{
			AccountNativeID:      basetypes.NewStringValue(accountId),
			AWSRegion:            basetypes.NewStringValue(region),
			OrganizationalUnitID: basetypes.NewStringValue(defaultOrgUnitId),
		}

		listEnvsResponse := &models.ListAWSEnvironmentsResponse{
			Embedded: &models.AWSEnvironmentListEmbedded{
				Items: []*models.AWSEnvironment{{Id: &envId}},
			},
		}
		patchResponse := &models.PatchOrganizationalUnitResponseWrapper{
			StatusCode: 202,
			Http202: &models.PatchOrganizationalUnitResponse{
				TaskId: &taskId,
			},
		}
		taskResponse := &models.ReadTaskResponse{Status: &status}

		mockAwsEnvClient.EXPECT().ListAwsEnvironments(
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Times(1).
			Return(listEnvsResponse, nil)
		mockOrgUnitsClient.EXPECT().ReadOrganizationalUnit(parentOuId, mock.Anything).Times(1).Return(
			&models.ReadOrganizationalUnitResponse{Id: &parentOuId}, nil)
		mockOrgUnitsClient.EXPECT().PatchOrganizationalUnit(parentOuId, mock.Anything, mock.Anything).
			Run(func(id string, _ *string, body *models.PatchOrganizationalUnitV2Request) {
				assert.Len(t, body.Entities.Add, 1)
				assert.Nil(t, body.Entities.Remove)
				assert.Equal(t, envId, *body.Entities.Add[0].PrimaryEntity.Id)
				assert.Equal(t, awsEnvironment, *body.Entities.Add[0].PrimaryEntity.ClumioType)
			}).Return(patchResponse, nil).Once()
		mockTaskClient.EXPECT().ReadTask(taskId).Times(1).Return(taskResponse, nil)

		err := updateOrgUnitForConnection(ctx, cr, plan, state)
		assert.Nil(t, err)
	})

	t.Run("moves connection between sibling OUs by remove then add", func(t *testing.T) {
		mockAwsEnvClient := sdkclients.NewMockAWSEnvironmentClient(t)
		mockOrgUnitsClient := sdkclients.NewMockOrganizationalUnitClient(t)
		mockTaskClient := sdkclients.NewMockTaskClient(t)

		cr := &clumioAWSConnectionResource{
			client:          &common.ApiClient{ClumioConfig: sdkconfig.Config{}},
			sdkEnvironments: mockAwsEnvClient,
			sdkOrgUnits:     mockOrgUnitsClient,
			sdkTasks:        mockTaskClient,
			pollTimeout:     time.Second,
			pollInterval:    time.Millisecond,
		}

		plan := &clumioAWSConnectionResourceModel{
			OrganizationalUnitID: basetypes.NewStringValue(parentOuId),
		}
		state := &clumioAWSConnectionResourceModel{
			AccountNativeID:      basetypes.NewStringValue(accountId),
			AWSRegion:            basetypes.NewStringValue(region),
			OrganizationalUnitID: basetypes.NewStringValue(ouId),
		}

		listEnvsResponse := &models.ListAWSEnvironmentsResponse{
			Embedded: &models.AWSEnvironmentListEmbedded{
				Items: []*models.AWSEnvironment{{Id: &envId}},
			},
		}
		patchResponse := &models.PatchOrganizationalUnitResponseWrapper{
			StatusCode: 202,
			Http202: &models.PatchOrganizationalUnitResponse{
				TaskId: &taskId,
			},
		}
		taskResponse := &models.ReadTaskResponse{Status: &status}

		mockAwsEnvClient.EXPECT().ListAwsEnvironments(
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Times(1).
			Return(listEnvsResponse, nil)
		mockOrgUnitsClient.EXPECT().ReadOrganizationalUnit(ouId, mock.Anything).Times(1).Return(
			&models.ReadOrganizationalUnitResponse{Id: &ouId}, nil)
		mockOrgUnitsClient.EXPECT().ReadOrganizationalUnit(parentOuId, mock.Anything).Times(1).Return(
			&models.ReadOrganizationalUnitResponse{Id: &parentOuId}, nil)
		var patched []string
		mockOrgUnitsClient.EXPECT().PatchOrganizationalUnit(ouId, mock.Anything, mock.Anything).
			Run(func(id string, _ *string, body *models.PatchOrganizationalUnitV2Request) {
				patched = append(patched, id)
				assert.Len(t, body.Entities.Remove, 1)
				assert.Nil(t, body.Entities.Add)
				assert.Equal(t, envId, *body.Entities.Remove[0].PrimaryEntity.Id)
			}).Return(patchResponse, nil).Once()
		mockOrgUnitsClient.EXPECT().PatchOrganizationalUnit(parentOuId, mock.Anything, mock.Anything).
			Run(func(id string, _ *string, body *models.PatchOrganizationalUnitV2Request) {
				patched = append(patched, id)
				assert.Len(t, body.Entities.Add, 1)
				assert.Nil(t, body.Entities.Remove)
				assert.Equal(t, envId, *body.Entities.Add[0].PrimaryEntity.Id)
			}).Return(patchResponse, nil).Once()
		mockTaskClient.EXPECT().ReadTask(taskId).Times(2).Return(taskResponse, nil)

		err := updateOrgUnitForConnection(ctx, cr, plan, state)
		assert.Nil(t, err)
		assert.Equal(t, []string{ouId, parentOuId}, patched)
	})

	t.Run("moves connection between cousin OUs one level per patch", func(t *testing.T) {
		mockAwsEnvClient := sdkclients.NewMockAWSEnvironmentClient(t)
		mockOrgUnitsClient := sdkclients.NewMockOrganizationalUnitClient(t)
		mockTaskClient := sdkclients.NewMockTaskClient(t)

		cr := &clumioAWSConnectionResource{
			client:          &common.ApiClient{ClumioConfig: sdkconfig.Config{}},
			sdkEnvironments: mockAwsEnvClient,
			sdkOrgUnits:     mockOrgUnitsClient,
			sdkTasks:        mockTaskClient,
			pollTimeout:     time.Second,
			pollInterval:    time.Millisecond,
		}

		// root -> a -> b holds the connection; root -> c is the target.
		rootId, aId, bId, cId, globalId := "root", "a", "b", "c", defaultOrgUnitId
		plan := &clumioAWSConnectionResourceModel{
			OrganizationalUnitID: basetypes.NewStringValue(cId),
		}
		state := &clumioAWSConnectionResourceModel{
			AccountNativeID:      basetypes.NewStringValue(accountId),
			AWSRegion:            basetypes.NewStringValue(region),
			OrganizationalUnitID: basetypes.NewStringValue(bId),
		}

		listEnvsResponse := &models.ListAWSEnvironmentsResponse{
			Embedded: &models.AWSEnvironmentListEmbedded{
				Items: []*models.AWSEnvironment{{Id: &envId}},
			},
		}
		patchResponse := &models.PatchOrganizationalUnitResponseWrapper{
			StatusCode: 202,
			Http202: &models.PatchOrganizationalUnitResponse{
				TaskId: &taskId,
			},
		}
		taskResponse := &models.ReadTaskResponse{Status: &status}

		mockAwsEnvClient.EXPECT().ListAwsEnvironments(
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Times(1).
			Return(listEnvsResponse, nil)
		mockOrgUnitsClient.EXPECT().ReadOrganizationalUnit(bId, mock.Anything).Times(1).Return(
			&models.ReadOrganizationalUnitResponse{Id: &bId, ParentId: &aId}, nil)
		mockOrgUnitsClient.EXPECT().ReadOrganizationalUnit(aId, mock.Anything).Times(1).Return(
			&models.ReadOrganizationalUnitResponse{Id: &aId, ParentId: &rootId}, nil)
		mockOrgUnitsClient.EXPECT().ReadOrganizationalUnit(cId, mock.Anything).Times(1).Return(
			&models.ReadOrganizationalUnitResponse{Id: &cId, ParentId: &rootId}, nil)
		mockOrgUnitsClient.EXPECT().ReadOrganizationalUnit(rootId, mock.Anything).Times(1).Return(
			&models.ReadOrganizationalUnitResponse{
				Id: &rootId, ParentId: &globalId,
				DescendantIds: []*string{&aId, &bId, &cId}}, nil)
		var patched []string
		mockOrgUnitsClient.EXPECT().PatchOrganizationalUnit(mock.Anything, mock.Anything, mock.Anything).
			Run(func(id string, _ *string, body *models.PatchOrganizationalUnitV2Request) {
				op := "remove"
				if body.Entities.Add != nil {
					op = "add"
				}
				patched = append(patched, op+":"+id)
			}).Return(patchResponse, nil).Times(3)
		mockTaskClient.EXPECT().ReadTask(taskId).Times(3).Return(taskResponse, nil)

		err := updateOrgUnitForConnection(ctx, cr, plan, state)
		assert.Nil(t, err)
		assert.Equal(t, []string{"remove:b", "remove:a", "add:c"}, patched)
	})

	t.Run("moves connection from global OU down several levels in order", func(t *testing.T) {
		mockAwsEnvClient := sdkclients.NewMockAWSEnvironmentClient(t)
		mockOrgUnitsClient := sdkclients.NewMockOrganizationalUnitClient(t)
		mockTaskClient := sdkclients.NewMockTaskClient(t)

		cr := &clumioAWSConnectionResource{
			client:          &common.ApiClient{ClumioConfig: sdkconfig.Config{}},
			sdkEnvironments: mockAwsEnvClient,
			sdkOrgUnits:     mockOrgUnitsClient,
			sdkTasks:        mockTaskClient,
			pollTimeout:     time.Second,
			pollInterval:    time.Millisecond,
		}

		// root -> a -> b is the target; root reports an empty parent, which means Global.
		rootId, aId, bId, emptyId := "root", "a", "b", ""
		plan := &clumioAWSConnectionResourceModel{
			OrganizationalUnitID: basetypes.NewStringValue(bId),
		}
		state := &clumioAWSConnectionResourceModel{
			AccountNativeID:      basetypes.NewStringValue(accountId),
			AWSRegion:            basetypes.NewStringValue(region),
			OrganizationalUnitID: basetypes.NewStringValue(defaultOrgUnitId),
		}

		listEnvsResponse := &models.ListAWSEnvironmentsResponse{
			Embedded: &models.AWSEnvironmentListEmbedded{
				Items: []*models.AWSEnvironment{{Id: &envId}},
			},
		}
		patchResponse := &models.PatchOrganizationalUnitResponseWrapper{
			StatusCode: 202,
			Http202: &models.PatchOrganizationalUnitResponse{
				TaskId: &taskId,
			},
		}
		taskResponse := &models.ReadTaskResponse{Status: &status}

		mockAwsEnvClient.EXPECT().ListAwsEnvironments(
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Times(1).
			Return(listEnvsResponse, nil)
		mockOrgUnitsClient.EXPECT().ReadOrganizationalUnit(bId, mock.Anything).Times(1).Return(
			&models.ReadOrganizationalUnitResponse{Id: &bId, ParentId: &aId}, nil)
		mockOrgUnitsClient.EXPECT().ReadOrganizationalUnit(aId, mock.Anything).Times(1).Return(
			&models.ReadOrganizationalUnitResponse{Id: &aId, ParentId: &rootId}, nil)
		mockOrgUnitsClient.EXPECT().ReadOrganizationalUnit(rootId, mock.Anything).Times(1).Return(
			&models.ReadOrganizationalUnitResponse{Id: &rootId, ParentId: &emptyId}, nil)
		var patched []string
		mockOrgUnitsClient.EXPECT().PatchOrganizationalUnit(mock.Anything, mock.Anything, mock.Anything).
			Run(func(id string, _ *string, body *models.PatchOrganizationalUnitV2Request) {
				op := "remove"
				if body.Entities.Add != nil {
					op = "add"
				}
				patched = append(patched, op+":"+id)
			}).Return(patchResponse, nil).Times(3)
		mockTaskClient.EXPECT().ReadTask(taskId).Times(3).Return(taskResponse, nil)

		err := updateOrgUnitForConnection(ctx, cr, plan, state)
		assert.Nil(t, err)
		assert.Equal(t, []string{"add:root", "add:a", "add:b"}, patched)
	})

	// root -> a -> b holds the connection; root -> c is the target. The walk is remove:b,
	// remove:a, add:c, and a failed task on any hop must stop the walk and name that hop.
	for failedHop, wantErr := range []string{
		"unable to remove the connection from organizational unit b",
		"unable to remove the connection from organizational unit a",
		"unable to add the connection to organizational unit c",
	} {
		t.Run(fmt.Sprintf("stops the walk when hop %d fails", failedHop+1), func(t *testing.T) {
			mockAwsEnvClient := sdkclients.NewMockAWSEnvironmentClient(t)
			mockOrgUnitsClient := sdkclients.NewMockOrganizationalUnitClient(t)
			mockTaskClient := sdkclients.NewMockTaskClient(t)

			cr := &clumioAWSConnectionResource{
				client:          &common.ApiClient{ClumioConfig: sdkconfig.Config{}},
				sdkEnvironments: mockAwsEnvClient,
				sdkOrgUnits:     mockOrgUnitsClient,
				sdkTasks:        mockTaskClient,
				pollTimeout:     time.Second,
				pollInterval:    time.Millisecond,
			}

			rootId, aId, bId, cId, globalId := "root", "a", "b", "c", defaultOrgUnitId
			plan := &clumioAWSConnectionResourceModel{
				OrganizationalUnitID: basetypes.NewStringValue(cId),
			}
			state := &clumioAWSConnectionResourceModel{
				AccountNativeID:      basetypes.NewStringValue(accountId),
				AWSRegion:            basetypes.NewStringValue(region),
				OrganizationalUnitID: basetypes.NewStringValue(bId),
			}

			listEnvsResponse := &models.ListAWSEnvironmentsResponse{
				Embedded: &models.AWSEnvironmentListEmbedded{
					Items: []*models.AWSEnvironment{{Id: &envId}},
				},
			}
			failedStatus := common.TaskFailed

			mockAwsEnvClient.EXPECT().ListAwsEnvironments(
				mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Times(1).
				Return(listEnvsResponse, nil)
			mockOrgUnitsClient.EXPECT().ReadOrganizationalUnit(bId, mock.Anything).Times(1).Return(
				&models.ReadOrganizationalUnitResponse{Id: &bId, ParentId: &aId}, nil)
			mockOrgUnitsClient.EXPECT().ReadOrganizationalUnit(aId, mock.Anything).Times(1).Return(
				&models.ReadOrganizationalUnitResponse{Id: &aId, ParentId: &rootId}, nil)
			mockOrgUnitsClient.EXPECT().ReadOrganizationalUnit(cId, mock.Anything).Times(1).Return(
				&models.ReadOrganizationalUnitResponse{Id: &cId, ParentId: &rootId}, nil)
			mockOrgUnitsClient.EXPECT().ReadOrganizationalUnit(rootId, mock.Anything).Times(1).Return(
				&models.ReadOrganizationalUnitResponse{
					Id: &rootId, ParentId: &globalId,
					DescendantIds: []*string{&aId, &bId, &cId}}, nil)
			var patched []string
			mockOrgUnitsClient.EXPECT().PatchOrganizationalUnit(mock.Anything, mock.Anything, mock.Anything).
				RunAndReturn(func(id string, _ *string, _ *models.PatchOrganizationalUnitV2Request) (
					*models.PatchOrganizationalUnitResponseWrapper, *apiutils.APIError) {
					patched = append(patched, id)
					hopTaskId := fmt.Sprint(len(patched))
					return &models.PatchOrganizationalUnitResponseWrapper{
						StatusCode: 202,
						Http202:    &models.PatchOrganizationalUnitResponse{TaskId: &hopTaskId},
					}, nil
				}).Times(failedHop + 1)
			mockTaskClient.EXPECT().ReadTask(mock.Anything).
				RunAndReturn(func(id string) (*models.ReadTaskResponse, *apiutils.APIError) {
					if id == fmt.Sprint(failedHop+1) {
						return &models.ReadTaskResponse{Status: &failedStatus}, nil
					}
					return &models.ReadTaskResponse{Status: &status}, nil
				}).Times(failedHop + 1)

			err := updateOrgUnitForConnection(ctx, cr, plan, state)
			assert.ErrorContains(t, err, wantErr)
			assert.Equal(t, []string{"b", "a", "c"}[:failedHop+1], patched)
		})
	}

	t.Run("fails when the OU chain loops", func(t *testing.T) {
		mockOrgUnitsClient := sdkclients.NewMockOrganizationalUnitClient(t)
		mockOrgUnitsClient.EXPECT().ReadOrganizationalUnit(ouId, mock.Anything).Times(1).Return(
			&models.ReadOrganizationalUnitResponse{Id: &ouId, ParentId: &parentOuId}, nil)
		mockOrgUnitsClient.EXPECT().ReadOrganizationalUnit(parentOuId, mock.Anything).Times(1).Return(
			&models.ReadOrganizationalUnitResponse{Id: &parentOuId, ParentId: &ouId}, nil)

		_, _, err := getOrgUnitsUpToCommonAncestor(mockOrgUnitsClient, ouId, "target")
		assert.ErrorContains(t, err, "its own ancestor")
	})

	// The API key can read only OU x and the OUs below it; t, the parent of x, is unreadable. The
	// strict mock fails the test on any read of t.
	t.Run("moves connection inside a scoped key's OU without reading above it", func(t *testing.T) {
		cr, mockOrgUnitsClient, mockTaskClient := newOrgUnitMoveTestResource(t)
		tId, xId, pId, qId := "t", "x", "p", "q"
		plan := &clumioAWSConnectionResourceModel{
			OrganizationalUnitID: basetypes.NewStringValue(qId),
		}
		state := &clumioAWSConnectionResourceModel{
			AccountNativeID:      basetypes.NewStringValue(accountId),
			AWSRegion:            basetypes.NewStringValue(region),
			OrganizationalUnitID: basetypes.NewStringValue(pId),
		}

		mockOrgUnitsClient.EXPECT().ReadOrganizationalUnit(pId, mock.Anything).Times(1).Return(
			&models.ReadOrganizationalUnitResponse{Id: &pId, ParentId: &xId}, nil)
		mockOrgUnitsClient.EXPECT().ReadOrganizationalUnit(xId, mock.Anything).Times(1).Return(
			&models.ReadOrganizationalUnitResponse{
				Id: &xId, ParentId: &tId, DescendantIds: []*string{&pId, &qId}}, nil)
		mockOrgUnitsClient.EXPECT().ReadOrganizationalUnit(qId, mock.Anything).Times(1).Return(
			&models.ReadOrganizationalUnitResponse{Id: &qId, ParentId: &xId}, nil)
		var patched []string
		mockOrgUnitsClient.EXPECT().PatchOrganizationalUnit(mock.Anything, mock.Anything, mock.Anything).
			Run(func(id string, _ *string, body *models.PatchOrganizationalUnitV2Request) {
				op := "remove"
				if body.Entities.Add != nil {
					op = "add"
				}
				patched = append(patched, op+":"+id)
			}).Return(&models.PatchOrganizationalUnitResponseWrapper{
			StatusCode: 202, Http202: &models.PatchOrganizationalUnitResponse{TaskId: &taskId},
		}, nil).Times(2)
		mockTaskClient.EXPECT().ReadTask(taskId).Times(2).Return(
			&models.ReadTaskResponse{Status: &status}, nil)

		err := updateOrgUnitForConnection(ctx, cr, plan, state)
		assert.Nil(t, err)
		assert.Equal(t, []string{"remove:p", "add:q"}, patched)
	})

	// The strict mock has no PATCH expectation, so the test fails if any PATCH is sent.
	t.Run("sends no patch when the target OU cannot be read", func(t *testing.T) {
		cr, mockOrgUnitsClient, _ := newOrgUnitMoveTestResource(t)
		plan := &clumioAWSConnectionResourceModel{
			OrganizationalUnitID: basetypes.NewStringValue(parentOuId),
		}
		state := &clumioAWSConnectionResourceModel{
			AccountNativeID:      basetypes.NewStringValue(accountId),
			AWSRegion:            basetypes.NewStringValue(region),
			OrganizationalUnitID: basetypes.NewStringValue(ouId),
		}

		mockOrgUnitsClient.EXPECT().ReadOrganizationalUnit(ouId, mock.Anything).Times(1).Return(
			&models.ReadOrganizationalUnitResponse{Id: &ouId}, nil)
		mockOrgUnitsClient.EXPECT().ReadOrganizationalUnit(parentOuId, mock.Anything).Times(1).Return(
			nil, &apiutils.APIError{ResponseCode: 500, Reason: "test", Response: []byte(testError)})

		err := updateOrgUnitForConnection(ctx, cr, plan, state)
		assert.ErrorContains(t, err, "unable to retrieve Organizational Unit")
	})

	// a lists q as a descendant, but q's parent chain reaches Global without passing a. The walk
	// must fail instead of looping, and send no PATCH (the strict mock has no PATCH expectation).
	t.Run("fails when the target is not below the common ancestor", func(t *testing.T) {
		cr, mockOrgUnitsClient, _ := newOrgUnitMoveTestResource(t)
		aId, pId, qId := "a", "p", "q"
		plan := &clumioAWSConnectionResourceModel{
			OrganizationalUnitID: basetypes.NewStringValue(qId),
		}
		state := &clumioAWSConnectionResourceModel{
			AccountNativeID:      basetypes.NewStringValue(accountId),
			AWSRegion:            basetypes.NewStringValue(region),
			OrganizationalUnitID: basetypes.NewStringValue(pId),
		}

		mockOrgUnitsClient.EXPECT().ReadOrganizationalUnit(pId, mock.Anything).Times(1).Return(
			&models.ReadOrganizationalUnitResponse{Id: &pId, ParentId: &aId}, nil)
		mockOrgUnitsClient.EXPECT().ReadOrganizationalUnit(aId, mock.Anything).Times(1).Return(
			&models.ReadOrganizationalUnitResponse{Id: &aId, DescendantIds: []*string{nil, &pId, &qId}},
			nil)
		mockOrgUnitsClient.EXPECT().ReadOrganizationalUnit(qId, mock.Anything).Times(1).Return(
			&models.ReadOrganizationalUnitResponse{Id: &qId}, nil)

		err := updateOrgUnitForConnection(ctx, cr, plan, state)
		assert.ErrorContains(t, err, "organizational unit q is not below organizational unit a")
	})

	t.Run("uses global OU client for OU operations when provider context is set", func(t *testing.T) {
		mockAwsEnvClient := sdkclients.NewMockAWSEnvironmentClient(t)
		mockOrgUnitsClient := sdkclients.NewMockOrganizationalUnitClient(t)
		mockGlobalOrgUnitsClient := sdkclients.NewMockOrganizationalUnitClient(t)
		mockTaskClient := sdkclients.NewMockTaskClient(t)
		mockGlobalTaskClient := sdkclients.NewMockTaskClient(t)
		previousFactory := newOrganizationalUnitClient
		previousTaskFactory := newTaskClient
		defer func() {
			newOrganizationalUnitClient = previousFactory
			newTaskClient = previousTaskFactory
		}()

		newOrganizationalUnitClient = func(config sdkconfig.Config) sdkclients.OrganizationalUnitClient {
			assert.Equal(t, "", config.OrganizationalUnitContext)
			return mockGlobalOrgUnitsClient
		}
		newTaskClient = func(config sdkconfig.Config) sdkclients.TaskClient {
			assert.Equal(t, "", config.OrganizationalUnitContext)
			return mockGlobalTaskClient
		}

		cr := &clumioAWSConnectionResource{
			client: &common.ApiClient{ClumioConfig: sdkconfig.Config{
				OrganizationalUnitContext: ouId,
			}},
			sdkEnvironments: mockAwsEnvClient,
			sdkOrgUnits:     mockOrgUnitsClient,
			sdkTasks:        mockTaskClient,
			pollTimeout:     time.Second,
			pollInterval:    time.Millisecond,
		}

		plan := &clumioAWSConnectionResourceModel{
			OrganizationalUnitID: basetypes.NewStringValue(parentOuId),
		}
		state := &clumioAWSConnectionResourceModel{
			AccountNativeID:      basetypes.NewStringValue(accountId),
			AWSRegion:            basetypes.NewStringValue(region),
			OrganizationalUnitID: basetypes.NewStringValue(ouId),
		}

		listEnvsResponse := &models.ListAWSEnvironmentsResponse{
			Embedded: &models.AWSEnvironmentListEmbedded{
				Items: []*models.AWSEnvironment{{Id: &envId}},
			},
		}
		patchResponse := &models.PatchOrganizationalUnitResponseWrapper{
			StatusCode: 202,
			Http202: &models.PatchOrganizationalUnitResponse{
				TaskId: &taskId,
			},
		}
		taskResponse := &models.ReadTaskResponse{Status: &status}

		mockAwsEnvClient.EXPECT().ListAwsEnvironments(
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Times(1).
			Return(listEnvsResponse, nil)
		mockGlobalOrgUnitsClient.EXPECT().ReadOrganizationalUnit(ouId, mock.Anything).Times(1).Return(
			&models.ReadOrganizationalUnitResponse{Id: &ouId}, nil)
		mockGlobalOrgUnitsClient.EXPECT().ReadOrganizationalUnit(parentOuId, mock.Anything).Times(1).Return(
			&models.ReadOrganizationalUnitResponse{Id: &parentOuId}, nil)
		mockGlobalOrgUnitsClient.EXPECT().PatchOrganizationalUnit(ouId, mock.Anything, mock.Anything).
			Run(func(id string, _ *string, body *models.PatchOrganizationalUnitV2Request) {
				assert.Len(t, body.Entities.Remove, 1)
			}).Return(patchResponse, nil).Once()
		mockGlobalOrgUnitsClient.EXPECT().PatchOrganizationalUnit(parentOuId, mock.Anything, mock.Anything).
			Run(func(id string, _ *string, body *models.PatchOrganizationalUnitV2Request) {
				assert.Len(t, body.Entities.Add, 1)
			}).Return(patchResponse, nil).Once()
		mockGlobalTaskClient.EXPECT().ReadTask(taskId).Times(2).Return(taskResponse, nil)

		err := updateOrgUnitForConnection(ctx, cr, plan, state)
		assert.Nil(t, err)
	})

	t.Run("moves connection to parent OU by remove", func(t *testing.T) {
		mockAwsEnvClient := sdkclients.NewMockAWSEnvironmentClient(t)
		mockOrgUnitsClient := sdkclients.NewMockOrganizationalUnitClient(t)
		mockTaskClient := sdkclients.NewMockTaskClient(t)

		cr := &clumioAWSConnectionResource{
			client:          &common.ApiClient{ClumioConfig: sdkconfig.Config{}},
			sdkEnvironments: mockAwsEnvClient,
			sdkOrgUnits:     mockOrgUnitsClient,
			sdkTasks:        mockTaskClient,
			pollTimeout:     time.Second,
			pollInterval:    time.Millisecond,
		}

		plan := &clumioAWSConnectionResourceModel{
			OrganizationalUnitID: basetypes.NewStringValue(parentOuId),
		}
		state := &clumioAWSConnectionResourceModel{
			AccountNativeID:      basetypes.NewStringValue(accountId),
			AWSRegion:            basetypes.NewStringValue(region),
			OrganizationalUnitID: basetypes.NewStringValue(ouId),
		}

		listEnvsResponse := &models.ListAWSEnvironmentsResponse{
			Embedded: &models.AWSEnvironmentListEmbedded{
				Items: []*models.AWSEnvironment{{Id: &envId}},
			},
		}
		patchResponse := &models.PatchOrganizationalUnitResponseWrapper{
			StatusCode: 202,
			Http202: &models.PatchOrganizationalUnitResponse{
				TaskId: &taskId,
			},
		}
		taskResponse := &models.ReadTaskResponse{Status: &status}

		mockAwsEnvClient.EXPECT().ListAwsEnvironments(
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Times(1).
			Return(listEnvsResponse, nil)
		mockOrgUnitsClient.EXPECT().ReadOrganizationalUnit(ouId, mock.Anything).Times(1).Return(
			&models.ReadOrganizationalUnitResponse{Id: &ouId, ParentId: &parentOuId}, nil)
		mockOrgUnitsClient.EXPECT().PatchOrganizationalUnit(ouId, mock.Anything, mock.Anything).
			Run(func(id string, _ *string, body *models.PatchOrganizationalUnitV2Request) {
				assert.Len(t, body.Entities.Remove, 1)
				assert.Nil(t, body.Entities.Add)
				assert.Equal(t, envId, *body.Entities.Remove[0].PrimaryEntity.Id)
			}).Return(patchResponse, nil).Once()
		mockTaskClient.EXPECT().ReadTask(taskId).Times(1).Return(taskResponse, nil)

		err := updateOrgUnitForConnection(ctx, cr, plan, state)
		assert.Nil(t, err)
	})

	t.Run("moves connection to global OU by remove from current OU", func(t *testing.T) {
		mockAwsEnvClient := sdkclients.NewMockAWSEnvironmentClient(t)
		mockOrgUnitsClient := sdkclients.NewMockOrganizationalUnitClient(t)
		mockTaskClient := sdkclients.NewMockTaskClient(t)

		cr := &clumioAWSConnectionResource{
			client:          &common.ApiClient{ClumioConfig: sdkconfig.Config{}},
			sdkEnvironments: mockAwsEnvClient,
			sdkOrgUnits:     mockOrgUnitsClient,
			sdkTasks:        mockTaskClient,
			pollTimeout:     time.Second,
			pollInterval:    time.Millisecond,
		}

		plan := &clumioAWSConnectionResourceModel{
			OrganizationalUnitID: basetypes.NewStringValue(defaultOrgUnitId),
		}
		state := &clumioAWSConnectionResourceModel{
			AccountNativeID:      basetypes.NewStringValue(accountId),
			AWSRegion:            basetypes.NewStringValue(region),
			OrganizationalUnitID: basetypes.NewStringValue(ouId),
		}

		listEnvsResponse := &models.ListAWSEnvironmentsResponse{
			Embedded: &models.AWSEnvironmentListEmbedded{
				Items: []*models.AWSEnvironment{{Id: &envId}},
			},
		}
		patchResponse := &models.PatchOrganizationalUnitResponseWrapper{
			StatusCode: 202,
			Http202: &models.PatchOrganizationalUnitResponse{
				TaskId: &taskId,
			},
		}
		taskResponse := &models.ReadTaskResponse{Status: &status}

		mockAwsEnvClient.EXPECT().ListAwsEnvironments(
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Times(1).
			Return(listEnvsResponse, nil)
		mockOrgUnitsClient.EXPECT().ReadOrganizationalUnit(ouId, mock.Anything).Times(1).Return(
			&models.ReadOrganizationalUnitResponse{Id: &ouId}, nil)
		mockOrgUnitsClient.EXPECT().PatchOrganizationalUnit(ouId, mock.Anything, mock.Anything).
			Run(func(id string, _ *string, body *models.PatchOrganizationalUnitV2Request) {
				assert.Len(t, body.Entities.Remove, 1)
				assert.Nil(t, body.Entities.Add)
				assert.Equal(t, envId, *body.Entities.Remove[0].PrimaryEntity.Id)
			}).Return(patchResponse, nil).Once()
		mockTaskClient.EXPECT().ReadTask(taskId).Times(1).Return(taskResponse, nil)

		err := updateOrgUnitForConnection(ctx, cr, plan, state)
		assert.Nil(t, err)
	})
}

// Unit test for setExternalId that checks and sets the ExternalID.
func TestSetExternalId(t *testing.T) {

	// Populate the Clumio AWS connection resource model to be used as input for
	// setExternalId().
	state := &clumioAWSConnectionResourceModel{
		ExternalID: basetypes.NewStringValue(externalId),
	}

	setExternalId(state, nil, &token)
	assert.Equal(t, fmt.Sprintf(externalIDFmt, token), state.ExternalID.ValueString())
}

// Unit test for setDataPlaneAccountId that checks and sets the DataPlaneAccountID.
func TestSetDataPlaneAccountId(t *testing.T) {

	// Populate the Clumio AWS connection resource model to be used as input for
	// setExternalId().
	state := &clumioAWSConnectionResourceModel{
		DataPlaneAccountID: basetypes.NewStringValue(dataplaneAccountId),
	}

	setDataPlaneAccountId(state, nil)
	assert.Equal(t, defaultDataPlaneAccountId, state.DataPlaneAccountID.ValueString())
}

// newOrgUnitMoveTestResource returns a resource whose environment list finds envId, plus its OU and
// task mocks.
func newOrgUnitMoveTestResource(t *testing.T) (*clumioAWSConnectionResource,
	*sdkclients.MockOrganizationalUnitClient, *sdkclients.MockTaskClient) {

	mockAwsEnvClient := sdkclients.NewMockAWSEnvironmentClient(t)
	mockOrgUnitsClient := sdkclients.NewMockOrganizationalUnitClient(t)
	mockTaskClient := sdkclients.NewMockTaskClient(t)
	mockAwsEnvClient.EXPECT().ListAwsEnvironments(
		mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Times(1).
		Return(&models.ListAWSEnvironmentsResponse{
			Embedded: &models.AWSEnvironmentListEmbedded{
				Items: []*models.AWSEnvironment{{Id: &envId}},
			},
		}, nil)
	return &clumioAWSConnectionResource{
		client:          &common.ApiClient{ClumioConfig: sdkconfig.Config{}},
		sdkEnvironments: mockAwsEnvClient,
		sdkOrgUnits:     mockOrgUnitsClient,
		sdkTasks:        mockTaskClient,
		pollTimeout:     time.Second,
		pollInterval:    time.Millisecond,
	}, mockOrgUnitsClient, mockTaskClient
}
