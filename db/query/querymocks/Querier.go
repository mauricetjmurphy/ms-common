// Code generated by mockery v1.0.0. DO NOT EDIT.

package querymocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// Querier is an autogenerated mock type for the Querier type
type Querier struct {
	mock.Mock
}

// Find provides a mock function with given fields: _a0, _a1
func (_m *Querier) Find(_a0 context.Context, _a1 interface{}) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindAll provides a mock function with given fields: _a0, _a1
func (_m *Querier) FindAll(_a0 context.Context, _a1 interface{}) (int64, error) {
	ret := _m.Called(_a0, _a1)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) int64); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, interface{}) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}