// Code generated by mockery. DO NOT EDIT.

package sdkclients

import (
	apiutils "github.com/clumio-code/clumio-go-sdk/api_utils"
	mock "github.com/stretchr/testify/mock"

	models "github.com/clumio-code/clumio-go-sdk/models"
)

// MockUserClient is an autogenerated mock type for the UserClient type
type MockUserClient struct {
	mock.Mock
}

type MockUserClient_Expecter struct {
	mock *mock.Mock
}

func (_m *MockUserClient) EXPECT() *MockUserClient_Expecter {
	return &MockUserClient_Expecter{mock: &_m.Mock}
}

// ChangePassword provides a mock function with given fields: body
func (_m *MockUserClient) ChangePassword(body *models.ChangePasswordV2Request) (*models.ChangePasswordResponse, *apiutils.APIError) {
	ret := _m.Called(body)

	if len(ret) == 0 {
		panic("no return value specified for ChangePassword")
	}

	var r0 *models.ChangePasswordResponse
	var r1 *apiutils.APIError
	if rf, ok := ret.Get(0).(func(*models.ChangePasswordV2Request) (*models.ChangePasswordResponse, *apiutils.APIError)); ok {
		return rf(body)
	}
	if rf, ok := ret.Get(0).(func(*models.ChangePasswordV2Request) *models.ChangePasswordResponse); ok {
		r0 = rf(body)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.ChangePasswordResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.ChangePasswordV2Request) *apiutils.APIError); ok {
		r1 = rf(body)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*apiutils.APIError)
		}
	}

	return r0, r1
}

// MockUserClient_ChangePassword_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ChangePassword'
type MockUserClient_ChangePassword_Call struct {
	*mock.Call
}

// ChangePassword is a helper method to define mock.On call
//   - body *models.ChangePasswordV2Request
func (_e *MockUserClient_Expecter) ChangePassword(body interface{}) *MockUserClient_ChangePassword_Call {
	return &MockUserClient_ChangePassword_Call{Call: _e.mock.On("ChangePassword", body)}
}

func (_c *MockUserClient_ChangePassword_Call) Run(run func(body *models.ChangePasswordV2Request)) *MockUserClient_ChangePassword_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*models.ChangePasswordV2Request))
	})
	return _c
}

func (_c *MockUserClient_ChangePassword_Call) Return(_a0 *models.ChangePasswordResponse, _a1 *apiutils.APIError) *MockUserClient_ChangePassword_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockUserClient_ChangePassword_Call) RunAndReturn(run func(*models.ChangePasswordV2Request) (*models.ChangePasswordResponse, *apiutils.APIError)) *MockUserClient_ChangePassword_Call {
	_c.Call.Return(run)
	return _c
}

// CreateUser provides a mock function with given fields: body
func (_m *MockUserClient) CreateUser(body *models.CreateUserV2Request) (*models.CreateUserResponse, *apiutils.APIError) {
	ret := _m.Called(body)

	if len(ret) == 0 {
		panic("no return value specified for CreateUser")
	}

	var r0 *models.CreateUserResponse
	var r1 *apiutils.APIError
	if rf, ok := ret.Get(0).(func(*models.CreateUserV2Request) (*models.CreateUserResponse, *apiutils.APIError)); ok {
		return rf(body)
	}
	if rf, ok := ret.Get(0).(func(*models.CreateUserV2Request) *models.CreateUserResponse); ok {
		r0 = rf(body)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.CreateUserResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.CreateUserV2Request) *apiutils.APIError); ok {
		r1 = rf(body)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*apiutils.APIError)
		}
	}

	return r0, r1
}

// MockUserClient_CreateUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateUser'
type MockUserClient_CreateUser_Call struct {
	*mock.Call
}

// CreateUser is a helper method to define mock.On call
//   - body *models.CreateUserV2Request
func (_e *MockUserClient_Expecter) CreateUser(body interface{}) *MockUserClient_CreateUser_Call {
	return &MockUserClient_CreateUser_Call{Call: _e.mock.On("CreateUser", body)}
}

func (_c *MockUserClient_CreateUser_Call) Run(run func(body *models.CreateUserV2Request)) *MockUserClient_CreateUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*models.CreateUserV2Request))
	})
	return _c
}

func (_c *MockUserClient_CreateUser_Call) Return(_a0 *models.CreateUserResponse, _a1 *apiutils.APIError) *MockUserClient_CreateUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockUserClient_CreateUser_Call) RunAndReturn(run func(*models.CreateUserV2Request) (*models.CreateUserResponse, *apiutils.APIError)) *MockUserClient_CreateUser_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteUser provides a mock function with given fields: userId
func (_m *MockUserClient) DeleteUser(userId int64) (interface{}, *apiutils.APIError) {
	ret := _m.Called(userId)

	if len(ret) == 0 {
		panic("no return value specified for DeleteUser")
	}

	var r0 interface{}
	var r1 *apiutils.APIError
	if rf, ok := ret.Get(0).(func(int64) (interface{}, *apiutils.APIError)); ok {
		return rf(userId)
	}
	if rf, ok := ret.Get(0).(func(int64) interface{}); ok {
		r0 = rf(userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	if rf, ok := ret.Get(1).(func(int64) *apiutils.APIError); ok {
		r1 = rf(userId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*apiutils.APIError)
		}
	}

	return r0, r1
}

// MockUserClient_DeleteUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteUser'
type MockUserClient_DeleteUser_Call struct {
	*mock.Call
}

// DeleteUser is a helper method to define mock.On call
//   - userId int64
func (_e *MockUserClient_Expecter) DeleteUser(userId interface{}) *MockUserClient_DeleteUser_Call {
	return &MockUserClient_DeleteUser_Call{Call: _e.mock.On("DeleteUser", userId)}
}

func (_c *MockUserClient_DeleteUser_Call) Run(run func(userId int64)) *MockUserClient_DeleteUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int64))
	})
	return _c
}

func (_c *MockUserClient_DeleteUser_Call) Return(_a0 interface{}, _a1 *apiutils.APIError) *MockUserClient_DeleteUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockUserClient_DeleteUser_Call) RunAndReturn(run func(int64) (interface{}, *apiutils.APIError)) *MockUserClient_DeleteUser_Call {
	_c.Call.Return(run)
	return _c
}

// ListUsers provides a mock function with given fields: limit, start, filter
func (_m *MockUserClient) ListUsers(limit *int64, start *string, filter *string) (*models.ListUsersResponse, *apiutils.APIError) {
	ret := _m.Called(limit, start, filter)

	if len(ret) == 0 {
		panic("no return value specified for ListUsers")
	}

	var r0 *models.ListUsersResponse
	var r1 *apiutils.APIError
	if rf, ok := ret.Get(0).(func(*int64, *string, *string) (*models.ListUsersResponse, *apiutils.APIError)); ok {
		return rf(limit, start, filter)
	}
	if rf, ok := ret.Get(0).(func(*int64, *string, *string) *models.ListUsersResponse); ok {
		r0 = rf(limit, start, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.ListUsersResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(*int64, *string, *string) *apiutils.APIError); ok {
		r1 = rf(limit, start, filter)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*apiutils.APIError)
		}
	}

	return r0, r1
}

// MockUserClient_ListUsers_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListUsers'
type MockUserClient_ListUsers_Call struct {
	*mock.Call
}

// ListUsers is a helper method to define mock.On call
//   - limit *int64
//   - start *string
//   - filter *string
func (_e *MockUserClient_Expecter) ListUsers(limit interface{}, start interface{}, filter interface{}) *MockUserClient_ListUsers_Call {
	return &MockUserClient_ListUsers_Call{Call: _e.mock.On("ListUsers", limit, start, filter)}
}

func (_c *MockUserClient_ListUsers_Call) Run(run func(limit *int64, start *string, filter *string)) *MockUserClient_ListUsers_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*int64), args[1].(*string), args[2].(*string))
	})
	return _c
}

func (_c *MockUserClient_ListUsers_Call) Return(_a0 *models.ListUsersResponse, _a1 *apiutils.APIError) *MockUserClient_ListUsers_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockUserClient_ListUsers_Call) RunAndReturn(run func(*int64, *string, *string) (*models.ListUsersResponse, *apiutils.APIError)) *MockUserClient_ListUsers_Call {
	_c.Call.Return(run)
	return _c
}

// ReadUser provides a mock function with given fields: userId
func (_m *MockUserClient) ReadUser(userId int64) (*models.ReadUserResponse, *apiutils.APIError) {
	ret := _m.Called(userId)

	if len(ret) == 0 {
		panic("no return value specified for ReadUser")
	}

	var r0 *models.ReadUserResponse
	var r1 *apiutils.APIError
	if rf, ok := ret.Get(0).(func(int64) (*models.ReadUserResponse, *apiutils.APIError)); ok {
		return rf(userId)
	}
	if rf, ok := ret.Get(0).(func(int64) *models.ReadUserResponse); ok {
		r0 = rf(userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.ReadUserResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(int64) *apiutils.APIError); ok {
		r1 = rf(userId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*apiutils.APIError)
		}
	}

	return r0, r1
}

// MockUserClient_ReadUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ReadUser'
type MockUserClient_ReadUser_Call struct {
	*mock.Call
}

// ReadUser is a helper method to define mock.On call
//   - userId int64
func (_e *MockUserClient_Expecter) ReadUser(userId interface{}) *MockUserClient_ReadUser_Call {
	return &MockUserClient_ReadUser_Call{Call: _e.mock.On("ReadUser", userId)}
}

func (_c *MockUserClient_ReadUser_Call) Run(run func(userId int64)) *MockUserClient_ReadUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int64))
	})
	return _c
}

func (_c *MockUserClient_ReadUser_Call) Return(_a0 *models.ReadUserResponse, _a1 *apiutils.APIError) *MockUserClient_ReadUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockUserClient_ReadUser_Call) RunAndReturn(run func(int64) (*models.ReadUserResponse, *apiutils.APIError)) *MockUserClient_ReadUser_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateUser provides a mock function with given fields: userId, body
func (_m *MockUserClient) UpdateUser(userId int64, body *models.UpdateUserV2Request) (*models.UpdateUserResponse, *apiutils.APIError) {
	ret := _m.Called(userId, body)

	if len(ret) == 0 {
		panic("no return value specified for UpdateUser")
	}

	var r0 *models.UpdateUserResponse
	var r1 *apiutils.APIError
	if rf, ok := ret.Get(0).(func(int64, *models.UpdateUserV2Request) (*models.UpdateUserResponse, *apiutils.APIError)); ok {
		return rf(userId, body)
	}
	if rf, ok := ret.Get(0).(func(int64, *models.UpdateUserV2Request) *models.UpdateUserResponse); ok {
		r0 = rf(userId, body)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.UpdateUserResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(int64, *models.UpdateUserV2Request) *apiutils.APIError); ok {
		r1 = rf(userId, body)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*apiutils.APIError)
		}
	}

	return r0, r1
}

// MockUserClient_UpdateUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateUser'
type MockUserClient_UpdateUser_Call struct {
	*mock.Call
}

// UpdateUser is a helper method to define mock.On call
//   - userId int64
//   - body *models.UpdateUserV2Request
func (_e *MockUserClient_Expecter) UpdateUser(userId interface{}, body interface{}) *MockUserClient_UpdateUser_Call {
	return &MockUserClient_UpdateUser_Call{Call: _e.mock.On("UpdateUser", userId, body)}
}

func (_c *MockUserClient_UpdateUser_Call) Run(run func(userId int64, body *models.UpdateUserV2Request)) *MockUserClient_UpdateUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int64), args[1].(*models.UpdateUserV2Request))
	})
	return _c
}

func (_c *MockUserClient_UpdateUser_Call) Return(_a0 *models.UpdateUserResponse, _a1 *apiutils.APIError) *MockUserClient_UpdateUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockUserClient_UpdateUser_Call) RunAndReturn(run func(int64, *models.UpdateUserV2Request) (*models.UpdateUserResponse, *apiutils.APIError)) *MockUserClient_UpdateUser_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateUserProfile provides a mock function with given fields: body
func (_m *MockUserClient) UpdateUserProfile(body *models.UpdateUserProfileV2Request) (*models.EditProfileResponse, *apiutils.APIError) {
	ret := _m.Called(body)

	if len(ret) == 0 {
		panic("no return value specified for UpdateUserProfile")
	}

	var r0 *models.EditProfileResponse
	var r1 *apiutils.APIError
	if rf, ok := ret.Get(0).(func(*models.UpdateUserProfileV2Request) (*models.EditProfileResponse, *apiutils.APIError)); ok {
		return rf(body)
	}
	if rf, ok := ret.Get(0).(func(*models.UpdateUserProfileV2Request) *models.EditProfileResponse); ok {
		r0 = rf(body)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.EditProfileResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.UpdateUserProfileV2Request) *apiutils.APIError); ok {
		r1 = rf(body)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*apiutils.APIError)
		}
	}

	return r0, r1
}

// MockUserClient_UpdateUserProfile_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateUserProfile'
type MockUserClient_UpdateUserProfile_Call struct {
	*mock.Call
}

// UpdateUserProfile is a helper method to define mock.On call
//   - body *models.UpdateUserProfileV2Request
func (_e *MockUserClient_Expecter) UpdateUserProfile(body interface{}) *MockUserClient_UpdateUserProfile_Call {
	return &MockUserClient_UpdateUserProfile_Call{Call: _e.mock.On("UpdateUserProfile", body)}
}

func (_c *MockUserClient_UpdateUserProfile_Call) Run(run func(body *models.UpdateUserProfileV2Request)) *MockUserClient_UpdateUserProfile_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*models.UpdateUserProfileV2Request))
	})
	return _c
}

func (_c *MockUserClient_UpdateUserProfile_Call) Return(_a0 *models.EditProfileResponse, _a1 *apiutils.APIError) *MockUserClient_UpdateUserProfile_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockUserClient_UpdateUserProfile_Call) RunAndReturn(run func(*models.UpdateUserProfileV2Request) (*models.EditProfileResponse, *apiutils.APIError)) *MockUserClient_UpdateUserProfile_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockUserClient creates a new instance of MockUserClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockUserClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockUserClient {
	mock := &MockUserClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
