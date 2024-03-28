// Code generated by mockery. DO NOT EDIT.

package sdkclients

import (
	apiutils "github.com/clumio-code/clumio-go-sdk/api_utils"
	mock "github.com/stretchr/testify/mock"

	models "github.com/clumio-code/clumio-go-sdk/models"
)

// MockPolicyRuleClient is an autogenerated mock type for the PolicyRuleClient type
type MockPolicyRuleClient struct {
	mock.Mock
}

type MockPolicyRuleClient_Expecter struct {
	mock *mock.Mock
}

func (_m *MockPolicyRuleClient) EXPECT() *MockPolicyRuleClient_Expecter {
	return &MockPolicyRuleClient_Expecter{mock: &_m.Mock}
}

// CreatePolicyRule provides a mock function with given fields: body
func (_m *MockPolicyRuleClient) CreatePolicyRule(body *models.CreatePolicyRuleV1Request) (*models.CreateRuleResponse, *apiutils.APIError) {
	ret := _m.Called(body)

	if len(ret) == 0 {
		panic("no return value specified for CreatePolicyRule")
	}

	var r0 *models.CreateRuleResponse
	var r1 *apiutils.APIError
	if rf, ok := ret.Get(0).(func(*models.CreatePolicyRuleV1Request) (*models.CreateRuleResponse, *apiutils.APIError)); ok {
		return rf(body)
	}
	if rf, ok := ret.Get(0).(func(*models.CreatePolicyRuleV1Request) *models.CreateRuleResponse); ok {
		r0 = rf(body)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.CreateRuleResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.CreatePolicyRuleV1Request) *apiutils.APIError); ok {
		r1 = rf(body)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*apiutils.APIError)
		}
	}

	return r0, r1
}

// MockPolicyRuleClient_CreatePolicyRule_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreatePolicyRule'
type MockPolicyRuleClient_CreatePolicyRule_Call struct {
	*mock.Call
}

// CreatePolicyRule is a helper method to define mock.On call
//   - body *models.CreatePolicyRuleV1Request
func (_e *MockPolicyRuleClient_Expecter) CreatePolicyRule(body interface{}) *MockPolicyRuleClient_CreatePolicyRule_Call {
	return &MockPolicyRuleClient_CreatePolicyRule_Call{Call: _e.mock.On("CreatePolicyRule", body)}
}

func (_c *MockPolicyRuleClient_CreatePolicyRule_Call) Run(run func(body *models.CreatePolicyRuleV1Request)) *MockPolicyRuleClient_CreatePolicyRule_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*models.CreatePolicyRuleV1Request))
	})
	return _c
}

func (_c *MockPolicyRuleClient_CreatePolicyRule_Call) Return(_a0 *models.CreateRuleResponse, _a1 *apiutils.APIError) *MockPolicyRuleClient_CreatePolicyRule_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockPolicyRuleClient_CreatePolicyRule_Call) RunAndReturn(run func(*models.CreatePolicyRuleV1Request) (*models.CreateRuleResponse, *apiutils.APIError)) *MockPolicyRuleClient_CreatePolicyRule_Call {
	_c.Call.Return(run)
	return _c
}

// DeletePolicyRule provides a mock function with given fields: ruleId
func (_m *MockPolicyRuleClient) DeletePolicyRule(ruleId string) (*models.DeleteRuleResponse, *apiutils.APIError) {
	ret := _m.Called(ruleId)

	if len(ret) == 0 {
		panic("no return value specified for DeletePolicyRule")
	}

	var r0 *models.DeleteRuleResponse
	var r1 *apiutils.APIError
	if rf, ok := ret.Get(0).(func(string) (*models.DeleteRuleResponse, *apiutils.APIError)); ok {
		return rf(ruleId)
	}
	if rf, ok := ret.Get(0).(func(string) *models.DeleteRuleResponse); ok {
		r0 = rf(ruleId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.DeleteRuleResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(string) *apiutils.APIError); ok {
		r1 = rf(ruleId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*apiutils.APIError)
		}
	}

	return r0, r1
}

// MockPolicyRuleClient_DeletePolicyRule_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeletePolicyRule'
type MockPolicyRuleClient_DeletePolicyRule_Call struct {
	*mock.Call
}

// DeletePolicyRule is a helper method to define mock.On call
//   - ruleId string
func (_e *MockPolicyRuleClient_Expecter) DeletePolicyRule(ruleId interface{}) *MockPolicyRuleClient_DeletePolicyRule_Call {
	return &MockPolicyRuleClient_DeletePolicyRule_Call{Call: _e.mock.On("DeletePolicyRule", ruleId)}
}

func (_c *MockPolicyRuleClient_DeletePolicyRule_Call) Run(run func(ruleId string)) *MockPolicyRuleClient_DeletePolicyRule_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockPolicyRuleClient_DeletePolicyRule_Call) Return(_a0 *models.DeleteRuleResponse, _a1 *apiutils.APIError) *MockPolicyRuleClient_DeletePolicyRule_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockPolicyRuleClient_DeletePolicyRule_Call) RunAndReturn(run func(string) (*models.DeleteRuleResponse, *apiutils.APIError)) *MockPolicyRuleClient_DeletePolicyRule_Call {
	_c.Call.Return(run)
	return _c
}

// ListPolicyRules provides a mock function with given fields: limit, start, organizationalUnitId, sort, filter
func (_m *MockPolicyRuleClient) ListPolicyRules(limit *int64, start *string, organizationalUnitId *string, sort *string, filter *string) (*models.ListRulesResponse, *apiutils.APIError) {
	ret := _m.Called(limit, start, organizationalUnitId, sort, filter)

	if len(ret) == 0 {
		panic("no return value specified for ListPolicyRules")
	}

	var r0 *models.ListRulesResponse
	var r1 *apiutils.APIError
	if rf, ok := ret.Get(0).(func(*int64, *string, *string, *string, *string) (*models.ListRulesResponse, *apiutils.APIError)); ok {
		return rf(limit, start, organizationalUnitId, sort, filter)
	}
	if rf, ok := ret.Get(0).(func(*int64, *string, *string, *string, *string) *models.ListRulesResponse); ok {
		r0 = rf(limit, start, organizationalUnitId, sort, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.ListRulesResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(*int64, *string, *string, *string, *string) *apiutils.APIError); ok {
		r1 = rf(limit, start, organizationalUnitId, sort, filter)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*apiutils.APIError)
		}
	}

	return r0, r1
}

// MockPolicyRuleClient_ListPolicyRules_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListPolicyRules'
type MockPolicyRuleClient_ListPolicyRules_Call struct {
	*mock.Call
}

// ListPolicyRules is a helper method to define mock.On call
//   - limit *int64
//   - start *string
//   - organizationalUnitId *string
//   - sort *string
//   - filter *string
func (_e *MockPolicyRuleClient_Expecter) ListPolicyRules(limit interface{}, start interface{}, organizationalUnitId interface{}, sort interface{}, filter interface{}) *MockPolicyRuleClient_ListPolicyRules_Call {
	return &MockPolicyRuleClient_ListPolicyRules_Call{Call: _e.mock.On("ListPolicyRules", limit, start, organizationalUnitId, sort, filter)}
}

func (_c *MockPolicyRuleClient_ListPolicyRules_Call) Run(run func(limit *int64, start *string, organizationalUnitId *string, sort *string, filter *string)) *MockPolicyRuleClient_ListPolicyRules_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*int64), args[1].(*string), args[2].(*string), args[3].(*string), args[4].(*string))
	})
	return _c
}

func (_c *MockPolicyRuleClient_ListPolicyRules_Call) Return(_a0 *models.ListRulesResponse, _a1 *apiutils.APIError) *MockPolicyRuleClient_ListPolicyRules_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockPolicyRuleClient_ListPolicyRules_Call) RunAndReturn(run func(*int64, *string, *string, *string, *string) (*models.ListRulesResponse, *apiutils.APIError)) *MockPolicyRuleClient_ListPolicyRules_Call {
	_c.Call.Return(run)
	return _c
}

// ReadPolicyRule provides a mock function with given fields: ruleId
func (_m *MockPolicyRuleClient) ReadPolicyRule(ruleId string) (*models.ReadRuleResponse, *apiutils.APIError) {
	ret := _m.Called(ruleId)

	if len(ret) == 0 {
		panic("no return value specified for ReadPolicyRule")
	}

	var r0 *models.ReadRuleResponse
	var r1 *apiutils.APIError
	if rf, ok := ret.Get(0).(func(string) (*models.ReadRuleResponse, *apiutils.APIError)); ok {
		return rf(ruleId)
	}
	if rf, ok := ret.Get(0).(func(string) *models.ReadRuleResponse); ok {
		r0 = rf(ruleId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.ReadRuleResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(string) *apiutils.APIError); ok {
		r1 = rf(ruleId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*apiutils.APIError)
		}
	}

	return r0, r1
}

// MockPolicyRuleClient_ReadPolicyRule_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ReadPolicyRule'
type MockPolicyRuleClient_ReadPolicyRule_Call struct {
	*mock.Call
}

// ReadPolicyRule is a helper method to define mock.On call
//   - ruleId string
func (_e *MockPolicyRuleClient_Expecter) ReadPolicyRule(ruleId interface{}) *MockPolicyRuleClient_ReadPolicyRule_Call {
	return &MockPolicyRuleClient_ReadPolicyRule_Call{Call: _e.mock.On("ReadPolicyRule", ruleId)}
}

func (_c *MockPolicyRuleClient_ReadPolicyRule_Call) Run(run func(ruleId string)) *MockPolicyRuleClient_ReadPolicyRule_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockPolicyRuleClient_ReadPolicyRule_Call) Return(_a0 *models.ReadRuleResponse, _a1 *apiutils.APIError) *MockPolicyRuleClient_ReadPolicyRule_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockPolicyRuleClient_ReadPolicyRule_Call) RunAndReturn(run func(string) (*models.ReadRuleResponse, *apiutils.APIError)) *MockPolicyRuleClient_ReadPolicyRule_Call {
	_c.Call.Return(run)
	return _c
}

// UpdatePolicyRule provides a mock function with given fields: ruleId, body
func (_m *MockPolicyRuleClient) UpdatePolicyRule(ruleId string, body *models.UpdatePolicyRuleV1Request) (*models.UpdateRuleResponse, *apiutils.APIError) {
	ret := _m.Called(ruleId, body)

	if len(ret) == 0 {
		panic("no return value specified for UpdatePolicyRule")
	}

	var r0 *models.UpdateRuleResponse
	var r1 *apiutils.APIError
	if rf, ok := ret.Get(0).(func(string, *models.UpdatePolicyRuleV1Request) (*models.UpdateRuleResponse, *apiutils.APIError)); ok {
		return rf(ruleId, body)
	}
	if rf, ok := ret.Get(0).(func(string, *models.UpdatePolicyRuleV1Request) *models.UpdateRuleResponse); ok {
		r0 = rf(ruleId, body)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.UpdateRuleResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(string, *models.UpdatePolicyRuleV1Request) *apiutils.APIError); ok {
		r1 = rf(ruleId, body)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*apiutils.APIError)
		}
	}

	return r0, r1
}

// MockPolicyRuleClient_UpdatePolicyRule_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdatePolicyRule'
type MockPolicyRuleClient_UpdatePolicyRule_Call struct {
	*mock.Call
}

// UpdatePolicyRule is a helper method to define mock.On call
//   - ruleId string
//   - body *models.UpdatePolicyRuleV1Request
func (_e *MockPolicyRuleClient_Expecter) UpdatePolicyRule(ruleId interface{}, body interface{}) *MockPolicyRuleClient_UpdatePolicyRule_Call {
	return &MockPolicyRuleClient_UpdatePolicyRule_Call{Call: _e.mock.On("UpdatePolicyRule", ruleId, body)}
}

func (_c *MockPolicyRuleClient_UpdatePolicyRule_Call) Run(run func(ruleId string, body *models.UpdatePolicyRuleV1Request)) *MockPolicyRuleClient_UpdatePolicyRule_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(*models.UpdatePolicyRuleV1Request))
	})
	return _c
}

func (_c *MockPolicyRuleClient_UpdatePolicyRule_Call) Return(_a0 *models.UpdateRuleResponse, _a1 *apiutils.APIError) *MockPolicyRuleClient_UpdatePolicyRule_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockPolicyRuleClient_UpdatePolicyRule_Call) RunAndReturn(run func(string, *models.UpdatePolicyRuleV1Request) (*models.UpdateRuleResponse, *apiutils.APIError)) *MockPolicyRuleClient_UpdatePolicyRule_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockPolicyRuleClient creates a new instance of MockPolicyRuleClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockPolicyRuleClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockPolicyRuleClient {
	mock := &MockPolicyRuleClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
