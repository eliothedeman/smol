package unit

import (
	"context"
	"fmt"
	"time"
)

type UnitRef string

type Ctx struct {
	Context context.Context
	ID      UnitRef
	Send    func(to UnitRef, msg any) error
	Log     func(format string, args ...any)
}

type Unit interface {
	Init(ctx Ctx)
	Handle(ctx Ctx, from UnitRef, message any) error
}

type Message struct {
	From    UnitRef
	To      UnitRef
	Type    string
	Payload any
	Time    time.Time
}

type Registry struct {
	units map[UnitRef]Unit
}

func NewRegistry() *Registry {
	return &Registry{
		units: make(map[UnitRef]Unit),
	}
}

func (r *Registry) Register(id UnitRef, unit Unit) {
	r.units[id] = unit
}

func (r *Registry) Get(id UnitRef) (Unit, bool) {
	unit, exists := r.units[id]
	return unit, exists
}

func (r *Registry) Start() error {
	for id, unit := range r.units {
		ctx := Ctx{
			Context: context.Background(),
			ID:      id,
			Send:    r.sendMessage,
			Log:     r.log,
		}
		unit.Init(ctx)
	}
	return nil
}

func (r *Registry) sendMessage(to UnitRef, msg any) error {
	unit, exists := r.units[to]
	if !exists {
		return fmt.Errorf("unit %s not found", to)
	}

	ctx := Ctx{
		Context: context.Background(),
		ID:      to,
		Send:    r.sendMessage,
		Log:     r.log,
	}

	return unit.Handle(ctx, "system", msg)
}

func (r *Registry) log(format string, args ...any) {
	fmt.Printf("[UNIT] "+format+"\n", args...)
}
