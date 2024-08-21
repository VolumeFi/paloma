// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	types "github.com/palomachain/paloma/x/valset/types"
)

// ValsetKeeper is an autogenerated mock type for the ValsetKeeper type
type ValsetKeeper struct {
	mock.Mock
}

// GetCurrentSnapshot provides a mock function with given fields: ctx
func (_m *ValsetKeeper) GetCurrentSnapshot(ctx context.Context) (*types.Snapshot, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetCurrentSnapshot")
	}

	var r0 *types.Snapshot
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*types.Snapshot, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *types.Snapshot); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Snapshot)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewValsetKeeper creates a new instance of ValsetKeeper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewValsetKeeper(t interface {
	mock.TestingT
	Cleanup(func())
},
) *ValsetKeeper {
	mock := &ValsetKeeper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}