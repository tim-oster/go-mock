package examples

import (
	"context"
	"iter"

	context_alias "context"
	alias "github.com/tim-oster/go-mock/examples/some-pkg"
	different_name "github.com/tim-oster/go-mock/examples/some-pkg"
	"golang.org/x/exp/constraints"
)

//go:generate go-mock ComplexTypes
type ComplexTypes interface {
	Normal(b bool) (int, error)
	RemoveCtx(ctx context.Context)
	RemoveCtxByType(context context.Context)
	RemoveCtxAlias(context context_alias.Context)
	NamedReturns() (i int, e error)
	UnnamedParams(context.Context, int)
	SharedTypeParams(a, b string, c int)
	AnonymousInterface(i interface {
		TestMethod() (bool, error)
	})
	AnonymousStruct(s struct {
		testVar bool
	})
	Any(any) any
	Channels(i chan<- int, o <-chan int) (chan bool, error)
	Variadic(i int, i2 ...int) (bool, error)
}

//go:generate go-mock -unexported Unexported
type Unexported interface {
	Normal(b bool) (int, error)
	RemoveCtx(ctx context.Context)
}

//go:generate go-mock -keepctx KeepCtx
type KeepCtx interface {
	Normal(b bool) (int, error)
	KeepCtx(ctx context.Context)
}

//go:generate go-mock ImportEdgeCases
type ImportEdgeCases interface {
	PathVsName(t different_name.SomeType)
	ImportAlias(t alias.SomeType)
}

//go:generate go-mock Generics
type Generics[T, R any, O constraints.Ordered, E ~[]R] interface {
	WithGeneric(typed T) (T, error)
	NestedGeneric(iter.Seq[T])
	WithoutGeneric(i int)
}
