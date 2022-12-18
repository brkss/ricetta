// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/brkss/vanillefraise2/db/sqlc (interfaces: Store)

// Package mockdb is a generated GoMock package.
package mockdb

import (
	context "context"
	reflect "reflect"

	db "github.com/brkss/vanillefraise2/db/sqlc"
	gomock "github.com/golang/mock/gomock"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// AssignRecipeToCategory mocks base method.
func (m *MockStore) AssignRecipeToCategory(arg0 context.Context, arg1 db.AssignRecipeToCategoryParams) (db.RecipeCategoryRecipe, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AssignRecipeToCategory", arg0, arg1)
	ret0, _ := ret[0].(db.RecipeCategoryRecipe)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AssignRecipeToCategory indicates an expected call of AssignRecipeToCategory.
func (mr *MockStoreMockRecorder) AssignRecipeToCategory(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AssignRecipeToCategory", reflect.TypeOf((*MockStore)(nil).AssignRecipeToCategory), arg0, arg1)
}

// CreateRecipe mocks base method.
func (m *MockStore) CreateRecipe(arg0 context.Context, arg1 db.CreateRecipeParams) (db.Recipe, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRecipe", arg0, arg1)
	ret0, _ := ret[0].(db.Recipe)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateRecipe indicates an expected call of CreateRecipe.
func (mr *MockStoreMockRecorder) CreateRecipe(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRecipe", reflect.TypeOf((*MockStore)(nil).CreateRecipe), arg0, arg1)
}

// CreateRecipeCategory mocks base method.
func (m *MockStore) CreateRecipeCategory(arg0 context.Context, arg1 db.CreateRecipeCategoryParams) (db.RecipeCategory, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRecipeCategory", arg0, arg1)
	ret0, _ := ret[0].(db.RecipeCategory)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateRecipeCategory indicates an expected call of CreateRecipeCategory.
func (mr *MockStoreMockRecorder) CreateRecipeCategory(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRecipeCategory", reflect.TypeOf((*MockStore)(nil).CreateRecipeCategory), arg0, arg1)
}

// CreateUser mocks base method.
func (m *MockStore) CreateUser(arg0 context.Context, arg1 db.CreateUserParams) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockStoreMockRecorder) CreateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockStore)(nil).CreateUser), arg0, arg1)
}

// CreateUserInfo mocks base method.
func (m *MockStore) CreateUserInfo(arg0 context.Context, arg1 db.CreateUserInfoParams) (db.UserInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUserInfo", arg0, arg1)
	ret0, _ := ret[0].(db.UserInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUserInfo indicates an expected call of CreateUserInfo.
func (mr *MockStoreMockRecorder) CreateUserInfo(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUserInfo", reflect.TypeOf((*MockStore)(nil).CreateUserInfo), arg0, arg1)
}
