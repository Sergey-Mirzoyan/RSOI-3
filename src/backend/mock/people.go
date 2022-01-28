// Code generated by MockGen. DO NOT EDIT.
// Source: ..\usecases\uc_interface\people.go

// Package mock is a generated GoMock package.
package mock

import (
	models "lab1/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIPersonsUsecase is a mock of IPersonsUsecase interface.
type MockIPersonsUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockIPersonsUsecaseMockRecorder
}

// MockIPersonsUsecaseMockRecorder is the mock recorder for MockIPersonsUsecase.
type MockIPersonsUsecaseMockRecorder struct {
	mock *MockIPersonsUsecase
}

// NewMockIPersonsUsecase creates a new mock instance.
func NewMockIPersonsUsecase(ctrl *gomock.Controller) *MockIPersonsUsecase {
	mock := &MockIPersonsUsecase{ctrl: ctrl}
	mock.recorder = &MockIPersonsUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIPersonsUsecase) EXPECT() *MockIPersonsUsecaseMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockIPersonsUsecase) Create(person *models.InputPerson) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", person)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockIPersonsUsecaseMockRecorder) Create(person interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockIPersonsUsecase)(nil).Create), person)
}

// Delete mocks base method.
func (m *MockIPersonsUsecase) Delete(id int) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", id)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockIPersonsUsecaseMockRecorder) Delete(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockIPersonsUsecase)(nil).Delete), id)
}

// GetAll mocks base method.
func (m *MockIPersonsUsecase) GetAll() ([]*models.Person, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll")
	ret0, _ := ret[0].([]*models.Person)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockIPersonsUsecaseMockRecorder) GetAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockIPersonsUsecase)(nil).GetAll))
}

// GetById mocks base method.
func (m *MockIPersonsUsecase) GetById(id int) (*models.Person, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", id)
	ret0, _ := ret[0].(*models.Person)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockIPersonsUsecaseMockRecorder) GetById(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockIPersonsUsecase)(nil).GetById), id)
}

// Update mocks base method.
func (m *MockIPersonsUsecase) Update(id int, newInfo *models.InputPerson) (models.Person, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", id, newInfo)
	ret0, _ := ret[0].(models.Person)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Update indicates an expected call of Update.
func (mr *MockIPersonsUsecaseMockRecorder) Update(id, newInfo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIPersonsUsecase)(nil).Update), id, newInfo)
}