// Code generated by go-mock. DO NOT EDIT.
// version: go-mock v0.1.1
// source: example.go
// interface: ImportEdgeCases
// flags: keepctx=false unexported=false
package examples

import (
	mock "github.com/stretchr/testify/mock"
	somepkg "github.com/tim-oster/go-mock/examples/some-pkg"
	"testing"
)

var _ ImportEdgeCases = (*MockImportEdgeCases)(nil)

type MockImportEdgeCases struct {
	mock.Mock
}

func (x *MockImportEdgeCases) PathVsName(t_ somepkg.SomeType) {
	args := x.Called(t_)
	if len(args) > 0 {
		if t, ok := args.Get(0).(mockImportEdgeCases_PathVsName_ReturnFunc); ok {
			t(t_)
		}
	}
}

type mockImportEdgeCases_PathVsName struct {
	*mock.Call
}

type mockImportEdgeCases_PathVsName_ReturnFunc func(t somepkg.SomeType)

func (c *mockImportEdgeCases_PathVsName) Return() *mock.Call {
	return c.Call.Return()
}

func (c *mockImportEdgeCases_PathVsName) ReturnFn(fn mockImportEdgeCases_PathVsName_ReturnFunc) *mock.Call {
	return c.Call.Return(fn)
}

func (x *MockImportEdgeCases) On_PathVsName(t_ somepkg.SomeType) *mockImportEdgeCases_PathVsName {
	return &mockImportEdgeCases_PathVsName{Call: x.On("PathVsName", t_)}
}

func (x *MockImportEdgeCases) On_PathVsName_Any() *mockImportEdgeCases_PathVsName {
	return &mockImportEdgeCases_PathVsName{Call: x.On("PathVsName", mock.Anything)}
}

func (x *MockImportEdgeCases) On_PathVsName_Interface(t_ any) *mockImportEdgeCases_PathVsName {
	return &mockImportEdgeCases_PathVsName{Call: x.On("PathVsName", t_)}
}

func (x *MockImportEdgeCases) Assert_PathVsName_Called(t *testing.T, t_ somepkg.SomeType) bool {
	return x.AssertCalled(t, "PathVsName", t_)
}

func (x *MockImportEdgeCases) Assert_PathVsName_NumberOfCalls(t *testing.T, expectedCalls int) bool {
	return x.AssertNumberOfCalls(t, "PathVsName", expectedCalls)
}

func (x *MockImportEdgeCases) Assert_PathVsName_NotCalled(t *testing.T, t_ somepkg.SomeType) bool {
	return x.AssertNotCalled(t, "PathVsName", t_)
}

func (x *MockImportEdgeCases) ImportAlias(t_ somepkg.SomeType) {
	args := x.Called(t_)
	if len(args) > 0 {
		if t, ok := args.Get(0).(mockImportEdgeCases_ImportAlias_ReturnFunc); ok {
			t(t_)
		}
	}
}

type mockImportEdgeCases_ImportAlias struct {
	*mock.Call
}

type mockImportEdgeCases_ImportAlias_ReturnFunc func(t somepkg.SomeType)

func (c *mockImportEdgeCases_ImportAlias) Return() *mock.Call {
	return c.Call.Return()
}

func (c *mockImportEdgeCases_ImportAlias) ReturnFn(fn mockImportEdgeCases_ImportAlias_ReturnFunc) *mock.Call {
	return c.Call.Return(fn)
}

func (x *MockImportEdgeCases) On_ImportAlias(t_ somepkg.SomeType) *mockImportEdgeCases_ImportAlias {
	return &mockImportEdgeCases_ImportAlias{Call: x.On("ImportAlias", t_)}
}

func (x *MockImportEdgeCases) On_ImportAlias_Any() *mockImportEdgeCases_ImportAlias {
	return &mockImportEdgeCases_ImportAlias{Call: x.On("ImportAlias", mock.Anything)}
}

func (x *MockImportEdgeCases) On_ImportAlias_Interface(t_ any) *mockImportEdgeCases_ImportAlias {
	return &mockImportEdgeCases_ImportAlias{Call: x.On("ImportAlias", t_)}
}

func (x *MockImportEdgeCases) Assert_ImportAlias_Called(t *testing.T, t_ somepkg.SomeType) bool {
	return x.AssertCalled(t, "ImportAlias", t_)
}

func (x *MockImportEdgeCases) Assert_ImportAlias_NumberOfCalls(t *testing.T, expectedCalls int) bool {
	return x.AssertNumberOfCalls(t, "ImportAlias", expectedCalls)
}

func (x *MockImportEdgeCases) Assert_ImportAlias_NotCalled(t *testing.T, t_ somepkg.SomeType) bool {
	return x.AssertNotCalled(t, "ImportAlias", t_)
}
