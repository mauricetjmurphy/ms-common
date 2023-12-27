// Code generated by mockery v2.14.0. DO NOT EDIT.

package secretsmocks

import (
	context "context"

	secrets "github.com/NBCUniversal/gvs-ms-common/clients/aws/secrets"
	mock "github.com/stretchr/testify/mock"
)

// Client is an autogenerated mock type for the Client type
type Client struct {
	mock.Mock
}

// GetAzureSecret provides a mock function with given fields: ctx, secretID
func (_m *Client) GetAzureSecret(ctx context.Context, secretID string) (*secrets.AzureAPISecret, error) {
	ret := _m.Called(ctx, secretID)

	var r0 *secrets.AzureAPISecret
	if rf, ok := ret.Get(0).(func(context.Context, string) *secrets.AzureAPISecret); ok {
		r0 = rf(ctx, secretID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*secrets.AzureAPISecret)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, secretID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSecret provides a mock function with given fields: ctx, secretID
func (_m *Client) GetSecret(ctx context.Context, secretID string) (*secrets.Secret, error) {
	ret := _m.Called(ctx, secretID)

	var r0 *secrets.Secret
	if rf, ok := ret.Get(0).(func(context.Context, string) *secrets.Secret); ok {
		r0 = rf(ctx, secretID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*secrets.Secret)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, secretID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetValue provides a mock function with given fields: ctx, secretID, target
func (_m *Client) GetValue(ctx context.Context, secretID string, target interface{}) error {
	ret := _m.Called(ctx, secretID, target)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, interface{}) error); ok {
		r0 = rf(ctx, secretID, target)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewClient creates a new instance of Client. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewClient(t mockConstructorTestingTNewClient) *Client {
	mock := &Client{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
