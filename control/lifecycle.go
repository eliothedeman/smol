package control

import "github.com/eliothedeman/smol/unit"

type Lifecycle struct{}

func (l *Lifecycle) Init(ctx unit.Ctx) {}

func (l *Lifecycle) Handle(ctx unit.Ctx, from unit.UnitRef, message any) error {
	return nil
}
