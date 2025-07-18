package unit

import "context"

type UnitDesc struct {
	Name  string
	Proxy Unit
}

type Ctx interface {
	context.Context
	Units() []UnitDesc
	Spawn(name string, f UnitFactory) UnitRef
	Self() UnitRef
	Subscribe(other Unit)
	Unsubscribe(other Unit)
}

type UnitFactory func() Unit

type UnitRef interface {
	Name() string
	Send(msg any)
	Stop()
}

type Unit interface {
	Init(ctx Ctx)
	Handle(ctx Ctx, from UnitRef, message any) error
}
