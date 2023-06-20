// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	gorm "gorm.io/gorm"
)

// FilterOption is an autogenerated mock type for the FilterOption type
type FilterOption struct {
	mock.Mock
}

// Execute provides a mock function with given fields: _a0
func (_m *FilterOption) Execute(_a0 *gorm.DB) *gorm.DB {
	ret := _m.Called(_a0)

	var r0 *gorm.DB
	if rf, ok := ret.Get(0).(func(*gorm.DB) *gorm.DB); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gorm.DB)
		}
	}

	return r0
}

type mockConstructorTestingTNewFilterOption interface {
	mock.TestingT
	Cleanup(func())
}

// NewFilterOption creates a new instance of FilterOption. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewFilterOption(t mockConstructorTestingTNewFilterOption) *FilterOption {
	mock := &FilterOption{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
