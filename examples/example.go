package examples

import "context"

//go:generate go-mock ComplexTypes
type ComplexTypes interface {
	Normal(b bool) (int, error)
	RemoveCtx(ctx context.Context)
	AnonymousInterface(i interface {
		TestMethod() (bool, error)
	})
	AnonymousStruct(s struct {
		testVar bool
	})
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
