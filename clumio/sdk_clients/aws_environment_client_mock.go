// Code generated by mockery. DO NOT EDIT.

package sdkclients

import (
	apiutils "github.com/clumio-code/clumio-go-sdk/api_utils"
	mock "github.com/stretchr/testify/mock"

	models "github.com/clumio-code/clumio-go-sdk/models"
)

// MockAWSEnvironmentClient is an autogenerated mock type for the AWSEnvironmentClient type
type MockAWSEnvironmentClient struct {
	mock.Mock
}

type MockAWSEnvironmentClient_Expecter struct {
	mock *mock.Mock
}

func (_m *MockAWSEnvironmentClient) EXPECT() *MockAWSEnvironmentClient_Expecter {
	return &MockAWSEnvironmentClient_Expecter{mock: &_m.Mock}
}

// ListAwsEnvironments provides a mock function with given fields: limit, start, filter, embed
func (_m *MockAWSEnvironmentClient) ListAwsEnvironments(limit *int64, start *string, filter *string, embed *string) (*models.ListAWSEnvironmentsResponse, *apiutils.APIError) {
	ret := _m.Called(limit, start, filter, embed)

	if len(ret) == 0 {
		panic("no return value specified for ListAwsEnvironments")
	}

	var r0 *models.ListAWSEnvironmentsResponse
	var r1 *apiutils.APIError
	if rf, ok := ret.Get(0).(func(*int64, *string, *string, *string) (*models.ListAWSEnvironmentsResponse, *apiutils.APIError)); ok {
		return rf(limit, start, filter, embed)
	}
	if rf, ok := ret.Get(0).(func(*int64, *string, *string, *string) *models.ListAWSEnvironmentsResponse); ok {
		r0 = rf(limit, start, filter, embed)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.ListAWSEnvironmentsResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(*int64, *string, *string, *string) *apiutils.APIError); ok {
		r1 = rf(limit, start, filter, embed)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*apiutils.APIError)
		}
	}

	return r0, r1
}

// MockAWSEnvironmentClient_ListAwsEnvironments_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListAwsEnvironments'
type MockAWSEnvironmentClient_ListAwsEnvironments_Call struct {
	*mock.Call
}

// ListAwsEnvironments is a helper method to define mock.On call
//   - limit *int64
//   - start *string
//   - filter *string
//   - embed *string
func (_e *MockAWSEnvironmentClient_Expecter) ListAwsEnvironments(limit interface{}, start interface{}, filter interface{}, embed interface{}) *MockAWSEnvironmentClient_ListAwsEnvironments_Call {
	return &MockAWSEnvironmentClient_ListAwsEnvironments_Call{Call: _e.mock.On("ListAwsEnvironments", limit, start, filter, embed)}
}

func (_c *MockAWSEnvironmentClient_ListAwsEnvironments_Call) Run(run func(limit *int64, start *string, filter *string, embed *string)) *MockAWSEnvironmentClient_ListAwsEnvironments_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*int64), args[1].(*string), args[2].(*string), args[3].(*string))
	})
	return _c
}

func (_c *MockAWSEnvironmentClient_ListAwsEnvironments_Call) Return(_a0 *models.ListAWSEnvironmentsResponse, _a1 *apiutils.APIError) *MockAWSEnvironmentClient_ListAwsEnvironments_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAWSEnvironmentClient_ListAwsEnvironments_Call) RunAndReturn(run func(*int64, *string, *string, *string) (*models.ListAWSEnvironmentsResponse, *apiutils.APIError)) *MockAWSEnvironmentClient_ListAwsEnvironments_Call {
	_c.Call.Return(run)
	return _c
}

// ReadAwsEnvironment provides a mock function with given fields: environmentId, embed
func (_m *MockAWSEnvironmentClient) ReadAwsEnvironment(environmentId string, embed *string) (*models.ReadAWSEnvironmentResponse, *apiutils.APIError) {
	ret := _m.Called(environmentId, embed)

	if len(ret) == 0 {
		panic("no return value specified for ReadAwsEnvironment")
	}

	var r0 *models.ReadAWSEnvironmentResponse
	var r1 *apiutils.APIError
	if rf, ok := ret.Get(0).(func(string, *string) (*models.ReadAWSEnvironmentResponse, *apiutils.APIError)); ok {
		return rf(environmentId, embed)
	}
	if rf, ok := ret.Get(0).(func(string, *string) *models.ReadAWSEnvironmentResponse); ok {
		r0 = rf(environmentId, embed)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.ReadAWSEnvironmentResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(string, *string) *apiutils.APIError); ok {
		r1 = rf(environmentId, embed)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*apiutils.APIError)
		}
	}

	return r0, r1
}

// MockAWSEnvironmentClient_ReadAwsEnvironment_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ReadAwsEnvironment'
type MockAWSEnvironmentClient_ReadAwsEnvironment_Call struct {
	*mock.Call
}

// ReadAwsEnvironment is a helper method to define mock.On call
//   - environmentId string
//   - embed *string
func (_e *MockAWSEnvironmentClient_Expecter) ReadAwsEnvironment(environmentId interface{}, embed interface{}) *MockAWSEnvironmentClient_ReadAwsEnvironment_Call {
	return &MockAWSEnvironmentClient_ReadAwsEnvironment_Call{Call: _e.mock.On("ReadAwsEnvironment", environmentId, embed)}
}

func (_c *MockAWSEnvironmentClient_ReadAwsEnvironment_Call) Run(run func(environmentId string, embed *string)) *MockAWSEnvironmentClient_ReadAwsEnvironment_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(*string))
	})
	return _c
}

func (_c *MockAWSEnvironmentClient_ReadAwsEnvironment_Call) Return(_a0 *models.ReadAWSEnvironmentResponse, _a1 *apiutils.APIError) *MockAWSEnvironmentClient_ReadAwsEnvironment_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAWSEnvironmentClient_ReadAwsEnvironment_Call) RunAndReturn(run func(string, *string) (*models.ReadAWSEnvironmentResponse, *apiutils.APIError)) *MockAWSEnvironmentClient_ReadAwsEnvironment_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockAWSEnvironmentClient creates a new instance of MockAWSEnvironmentClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockAWSEnvironmentClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockAWSEnvironmentClient {
	mock := &MockAWSEnvironmentClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
