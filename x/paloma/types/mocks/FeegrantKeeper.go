// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"

	feegrant "cosmossdk.io/x/feegrant"
	cosmos_sdktypes "github.com/cosmos/cosmos-sdk/types"

	mock "github.com/stretchr/testify/mock"
)

// FeegrantKeeper is an autogenerated mock type for the FeegrantKeeper type
type FeegrantKeeper struct {
	mock.Mock
}

// AllowancesByGranter provides a mock function with given fields: ctx, req
func (_m *FeegrantKeeper) AllowancesByGranter(ctx context.Context, req *feegrant.QueryAllowancesByGranterRequest) (*feegrant.QueryAllowancesByGranterResponse, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for AllowancesByGranter")
	}

	var r0 *feegrant.QueryAllowancesByGranterResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *feegrant.QueryAllowancesByGranterRequest) (*feegrant.QueryAllowancesByGranterResponse, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *feegrant.QueryAllowancesByGranterRequest) *feegrant.QueryAllowancesByGranterResponse); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*feegrant.QueryAllowancesByGranterResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *feegrant.QueryAllowancesByGranterRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GrantAllowance provides a mock function with given fields: ctx, granter, grantee, feeAllowance
func (_m *FeegrantKeeper) GrantAllowance(ctx context.Context, granter cosmos_sdktypes.AccAddress, grantee cosmos_sdktypes.AccAddress, feeAllowance feegrant.FeeAllowanceI) error {
	ret := _m.Called(ctx, granter, grantee, feeAllowance)

	if len(ret) == 0 {
		panic("no return value specified for GrantAllowance")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, cosmos_sdktypes.AccAddress, cosmos_sdktypes.AccAddress, feegrant.FeeAllowanceI) error); ok {
		r0 = rf(ctx, granter, grantee, feeAllowance)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewFeegrantKeeper creates a new instance of FeegrantKeeper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFeegrantKeeper(t interface {
	mock.TestingT
	Cleanup(func())
}) *FeegrantKeeper {
	mock := &FeegrantKeeper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
