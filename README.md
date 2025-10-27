# go-mock

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.18-blue.svg)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

A powerful Go code generation tool that automatically creates type-safe mock implementations from interface definitions. Designed to work seamlessly with [testify/mock](https://github.com/stretchr/testify), go-mock eliminates boilerplate and provides a fluent, type-safe API for setting up test expectations and assertions.

## Overview

Writing mocks manually is tedious and error-prone. go-mock analyzes your Go interfaces and generates comprehensive mock implementations with:

- **Type-safe method setup** - No more casting or type assertions
- **Fluent API** - Chain method calls for readable test setup
- **Smart parameter handling** - Automatic context removal, variadic parameters, generics support
- **Comprehensive assertions** - Verify calls, arguments, and invocation counts
- **Complex type support** - Handles channels, anonymous interfaces, function types, and more

## Features

- ✅ **Automatic code generation** via `//go:generate` directives
- ✅ **Type-safe mock implementations** with compile-time interface checks
- ✅ **Fluent expectation API** - `On_Method()`, `On_Method_Any()`, `On_Method_Interface()`
- ✅ **Flexible return values** - `Return()` for fixed values, `ReturnFn()` for custom logic
- ✅ **Built-in assertions** - `Assert_Method_Called()`, `Assert_Method_NumberOfCalls()`, etc.
- ✅ **Smart context handling** - Automatically removes `context.Context` from mock setup (optional)
- ✅ **Generics support** (Go 1.18+) - Full support for type parameters and constraints
- ✅ **Complex Go types** - Arrays, slices, maps, channels, structs, interfaces, function types
- ✅ **Variadic parameters** - Proper handling of `...T` syntax
- ✅ **Import resolution** - Automatic package imports and qualification

## Installation

```bash
go install github.com/tim-oster/go-mock@latest
```

Requires Go 1.18 or later.

## Quick Start

### 1. Define your interface

```go
package myservice // file: myservice.go

import "context"

type UserRepository interface {
    GetUser(ctx context.Context, id string) (*User, error)
    CreateUser(ctx context.Context, user *User) error
}
```

### 2. Add a go:generate directive

```go
//go:generate go-mock UserRepository
```

### 3. Generate the mock

```bash
go generate ./...
```

This creates `myservice_mock_userrepository.gen.go` with a full mock implementation.

### 4. Use the mock in your tests

```go
package myservice_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "myservice"
)

func TestUserService(t *testing.T) {
    // Create mock
    mockRepo := &myservice.MockUserRepository{}

    // Setup expectation
    expectedUser := &myservice.User{ID: "123", Name: "John"}
    mockRepo.On_GetUser("123").Return(expectedUser, nil)

    // Use mock in your code
    service := myservice.NewUserService(mockRepo)
    user, err := service.FetchUser(context.Background(), "123")

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, expectedUser, user)
    mockRepo.AssertExpectations(t)
}
```

## Usage

### Basic Command

```go
//go:generate go-mock InterfaceName
```

### Command-Line Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-keepctx` | Keep `context.Context` parameters in mock setup methods | `false` |
| `-unexported` | Generate unexported mock struct (e.g., `mockInterface` instead of `MockInterface`) | `false` |

### Interface Renaming

Rename the generated mock struct:

```go
//go:generate go-mock UserRepository:UserRepo
```

This generates `MockUserRepo` instead of `MockUserRepository`.

### Examples

```go
// Standard interface
//go:generate go-mock UserRepository

// Keep context in mock setup
//go:generate go-mock -keepctx APIClient

// Generate unexported mock for internal testing
//go:generate go-mock -unexported internalService

// Rename the generated mock
//go:generate go-mock UserRepository:Repo
```

## Generated Code Overview

For each interface method, go-mock generates:

### Expectation Setup Methods

```go
// Setup with specific arguments (context auto-removed by default)
func (m *MockInterface) On_MethodName(param string) *MockInterface_MethodName_Call

// Setup with wildcard arguments (uses mock.Anything internally)
func (m *MockInterface) On_MethodName_Any() *MockInterface_MethodName_Call

// Setup with interface{} parameters for maximum flexibility
func (m *MockInterface) On_MethodName_Interface(param interface{}) *MockInterface_MethodName_Call
```

### Return Value Helpers

```go
// Set fixed return values
func (c *MockInterface_MethodName_Call) Return(s string, err error)

// Provide custom logic with a function
func (c *MockInterface_MethodName_Call) ReturnFn(fn func(string) (string, error))
```

### Assertion Methods

```go
// Assert method was called with specific arguments
func (m *MockInterface) Assert_MethodName_Called(t *testing.T, param string)

// Assert method was called N times
func (m *MockInterface) Assert_MethodName_NumberOfCalls(t *testing.T, n int)

// Assert method was never called
func (m *MockInterface) Assert_MethodName_NotCalled(t *testing.T)
```

## Advanced Features

### Context Handling

By default, go-mock removes `context.Context` as the first parameter from mock setup methods, since context is rarely important for mocking behavior:

```go
// Interface method
GetUser(ctx context.Context, id string) (*User, error)

// Generated setup method (context removed)
On_GetUser(id string) *MockUserRepository_GetUser_Call
```

To preserve context in setup methods, use `-keepctx`:

```go
//go:generate go-mock -keepctx UserRepository

// Generated setup method (context kept)
On_GetUser(ctx context.Context, id string) *MockUserRepository_GetUser_Call
```

## How It Works

1. **AST Parsing** - go-mock uses Go's `go/parser` to analyze your source files
2. **Interface Extraction** - Finds and filters interface definitions matching your pattern
3. **Code Generation** - Uses [jennifer](https://github.com/dave/jennifer) to generate type-safe mock code
4. **Output** - Creates `{source}_mock_{interface}.gen.go` files

## Contributing

Contributions are welcome! Please feel free to submit issues, feature requests, or pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
