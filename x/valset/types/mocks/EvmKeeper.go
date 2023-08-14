// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	types "github.com/cosmos/cosmos-sdk/types"
	mock "github.com/stretchr/testify/mock"
)

// EvmKeeper is an autogenerated mock type for the EvmKeeper type
type EvmKeeper struct {
	mock.Mock
}

// MissingChains provides a mock function with given fields: ctx, chainReferenceIDs
func (_m *EvmKeeper) MissingChains(ctx types.Context, chainReferenceIDs []string) ([]string, error) {
	ret := _m.Called(ctx, chainReferenceIDs)

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, []string) ([]string, error)); ok {
		return rf(ctx, chainReferenceIDs)
	}
	if rf, ok := ret.Get(0).(func(types.Context, []string) []string); ok {
		r0 = rf(ctx, chainReferenceIDs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(types.Context, []string) error); ok {
		r1 = rf(ctx, chainReferenceIDs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewEvmKeeper interface {
	mock.TestingT
	Cleanup(func())
}

// NewEvmKeeper creates a new instance of EvmKeeper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewEvmKeeper(t mockConstructorTestingTNewEvmKeeper) *EvmKeeper {
	mock := &EvmKeeper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}