package control

import "github.com/eliothedeman/smol/unit"

type InstructionExecutor struct{}

func (ie *InstructionExecutor) Init(ctx unit.Ctx) {}

func (ie *InstructionExecutor) Handle(ctx unit.Ctx, from unit.UnitRef, message any) error {
	return nil
}
