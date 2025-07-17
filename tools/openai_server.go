package tools

import "github.com/eliothedeman/smol/unit"

type OpenAIServer struct{}

func (o *OpenAIServer) Init(ctx unit.Ctx) {}

func (o *OpenAIServer) Handle(ctx unit.Ctx, from unit.UnitRef, message any) error {
	return nil
}
