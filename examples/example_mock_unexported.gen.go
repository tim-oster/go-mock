// Code generated by go-mock. DO NOT EDIT.
// version: go-mock v0.1.1
// source: example.go
// interface: Unexported
// flags: keepctx=false unexported=true
package examples

import (
	"context"
	mock "github.com/stretchr/testify/mock"
	"testing"
)

var _ Unexported = (*mockUnexported)(nil)

type mockUnexported struct {
	mock.Mock
}

func (x *mockUnexported) Normal(b_ bool) (int, error) {
	args := x.Called(b_)
	if len(args) > 0 {
		if t, ok := args.Get(0).(mockUnexported_Normal_ReturnFunc); ok {
			return t(b_)
		}
	}
	var r0 int
	if v := args.Get(0); v != nil {
		r0 = v.(int)
	}
	var r1 error
	if v := args.Get(1); v != nil {
		r1 = v.(error)
	}
	return r0, r1
}

type mockUnexported_Normal struct {
	*mock.Call
}

type mockUnexported_Normal_ReturnFunc func(b bool) (int, error)

func (c *mockUnexported_Normal) Return(arg0 int, arg1 error) *mock.Call {
	return c.Call.Return(arg0, arg1)
}

func (c *mockUnexported_Normal) ReturnFn(fn mockUnexported_Normal_ReturnFunc) *mock.Call {
	return c.Call.Return(fn)
}

func (x *mockUnexported) On_Normal(b_ bool) *mockUnexported_Normal {
	return &mockUnexported_Normal{Call: x.On("Normal", b_)}
}

func (x *mockUnexported) On_Normal_Any() *mockUnexported_Normal {
	return &mockUnexported_Normal{Call: x.On("Normal", mock.Anything)}
}

func (x *mockUnexported) On_Normal_Interface(b_ interface{}) *mockUnexported_Normal {
	return &mockUnexported_Normal{Call: x.On("Normal", b_)}
}

func (x *mockUnexported) Assert_Normal_Called(t *testing.T, b_ bool) bool {
	return x.AssertCalled(t, "Normal", b_)
}

func (x *mockUnexported) Assert_Normal_NumberOfCalls(t *testing.T, expectedCalls int) bool {
	return x.AssertNumberOfCalls(t, "Normal", expectedCalls)
}

func (x *mockUnexported) Assert_Normal_NotCalled(t *testing.T, b_ bool) bool {
	return x.AssertNotCalled(t, "Normal", b_)
}

func (x *mockUnexported) RemoveCtx(ctx_ context.Context) {
	args := x.Called()
	if len(args) > 0 {
		if t, ok := args.Get(0).(mockUnexported_RemoveCtx_ReturnFunc); ok {
			t(ctx_)
		}
	}
}

type mockUnexported_RemoveCtx struct {
	*mock.Call
}

type mockUnexported_RemoveCtx_ReturnFunc func(ctx context.Context)

func (c *mockUnexported_RemoveCtx) Return() *mock.Call {
	return c.Call.Return()
}

func (c *mockUnexported_RemoveCtx) ReturnFn(fn mockUnexported_RemoveCtx_ReturnFunc) *mock.Call {
	return c.Call.Return(fn)
}

func (x *mockUnexported) On_RemoveCtx() *mockUnexported_RemoveCtx {
	return &mockUnexported_RemoveCtx{Call: x.On("RemoveCtx")}
}

func (x *mockUnexported) Assert_RemoveCtx_Called(t *testing.T) bool {
	return x.AssertCalled(t, "RemoveCtx")
}

func (x *mockUnexported) Assert_RemoveCtx_NumberOfCalls(t *testing.T, expectedCalls int) bool {
	return x.AssertNumberOfCalls(t, "RemoveCtx", expectedCalls)
}

func (x *mockUnexported) Assert_RemoveCtx_NotCalled(t *testing.T) bool {
	return x.AssertNotCalled(t, "RemoveCtx")
}
