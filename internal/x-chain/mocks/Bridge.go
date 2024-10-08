// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"

	xchain "github.com/palomachain/paloma/v2/internal/x-chain"
	mock "github.com/stretchr/testify/mock"
)

// Bridge is an autogenerated mock type for the Bridge type
type Bridge struct {
	mock.Mock
}

// ExecuteJob provides a mock function with given fields: ctx, jcfg
func (_m *Bridge) ExecuteJob(ctx context.Context, jcfg *xchain.JobConfiguration) (uint64, error) {
	ret := _m.Called(ctx, jcfg)

	if len(ret) == 0 {
		panic("no return value specified for ExecuteJob")
	}

	var r0 uint64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *xchain.JobConfiguration) (uint64, error)); ok {
		return rf(ctx, jcfg)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *xchain.JobConfiguration) uint64); ok {
		r0 = rf(ctx, jcfg)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *xchain.JobConfiguration) error); ok {
		r1 = rf(ctx, jcfg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// VerifyJob provides a mock function with given fields: ctx, definition, payload, refID
func (_m *Bridge) VerifyJob(ctx context.Context, definition []byte, payload []byte, refID string) error {
	ret := _m.Called(ctx, definition, payload, refID)

	if len(ret) == 0 {
		panic("no return value specified for VerifyJob")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []byte, []byte, string) error); ok {
		r0 = rf(ctx, definition, payload, refID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// XChainReferenceIDs provides a mock function with given fields: _a0
func (_m *Bridge) XChainReferenceIDs(_a0 context.Context) []string {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for XChainReferenceIDs")
	}

	var r0 []string
	if rf, ok := ret.Get(0).(func(context.Context) []string); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// XChainType provides a mock function with given fields:
func (_m *Bridge) XChainType() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for XChainType")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// NewBridge creates a new instance of Bridge. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewBridge(t interface {
	mock.TestingT
	Cleanup(func())
},
) *Bridge {
	mock := &Bridge{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
