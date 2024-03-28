// Code generated by mockery. DO NOT EDIT.

package sdkclients

import (
	apiutils "github.com/clumio-code/clumio-go-sdk/api_utils"
	mock "github.com/stretchr/testify/mock"

	models "github.com/clumio-code/clumio-go-sdk/models"
)

// MockPolicyAssignmentClient is an autogenerated mock type for the PolicyAssignmentClient type
type MockPolicyAssignmentClient struct {
	mock.Mock
}

type MockPolicyAssignmentClient_Expecter struct {
	mock *mock.Mock
}

func (_m *MockPolicyAssignmentClient) EXPECT() *MockPolicyAssignmentClient_Expecter {
	return &MockPolicyAssignmentClient_Expecter{mock: &_m.Mock}
}

// SetPolicyAssignments provides a mock function with given fields: body
func (_m *MockPolicyAssignmentClient) SetPolicyAssignments(body *models.SetPolicyAssignmentsV1Request) (*models.SetAssignmentsResponse, *apiutils.APIError) {
	ret := _m.Called(body)

	if len(ret) == 0 {
		panic("no return value specified for SetPolicyAssignments")
	}

	var r0 *models.SetAssignmentsResponse
	var r1 *apiutils.APIError
	if rf, ok := ret.Get(0).(func(*models.SetPolicyAssignmentsV1Request) (*models.SetAssignmentsResponse, *apiutils.APIError)); ok {
		return rf(body)
	}
	if rf, ok := ret.Get(0).(func(*models.SetPolicyAssignmentsV1Request) *models.SetAssignmentsResponse); ok {
		r0 = rf(body)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.SetAssignmentsResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.SetPolicyAssignmentsV1Request) *apiutils.APIError); ok {
		r1 = rf(body)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*apiutils.APIError)
		}
	}

	return r0, r1
}

// MockPolicyAssignmentClient_SetPolicyAssignments_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetPolicyAssignments'
type MockPolicyAssignmentClient_SetPolicyAssignments_Call struct {
	*mock.Call
}

// SetPolicyAssignments is a helper method to define mock.On call
//   - body *models.SetPolicyAssignmentsV1Request
func (_e *MockPolicyAssignmentClient_Expecter) SetPolicyAssignments(body interface{}) *MockPolicyAssignmentClient_SetPolicyAssignments_Call {
	return &MockPolicyAssignmentClient_SetPolicyAssignments_Call{Call: _e.mock.On("SetPolicyAssignments", body)}
}

func (_c *MockPolicyAssignmentClient_SetPolicyAssignments_Call) Run(run func(body *models.SetPolicyAssignmentsV1Request)) *MockPolicyAssignmentClient_SetPolicyAssignments_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*models.SetPolicyAssignmentsV1Request))
	})
	return _c
}

func (_c *MockPolicyAssignmentClient_SetPolicyAssignments_Call) Return(_a0 *models.SetAssignmentsResponse, _a1 *apiutils.APIError) *MockPolicyAssignmentClient_SetPolicyAssignments_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockPolicyAssignmentClient_SetPolicyAssignments_Call) RunAndReturn(run func(*models.SetPolicyAssignmentsV1Request) (*models.SetAssignmentsResponse, *apiutils.APIError)) *MockPolicyAssignmentClient_SetPolicyAssignments_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockPolicyAssignmentClient creates a new instance of MockPolicyAssignmentClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockPolicyAssignmentClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockPolicyAssignmentClient {
	mock := &MockPolicyAssignmentClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
