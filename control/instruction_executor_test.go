package control

import (
	"context"
	"testing"
	"time"

	"github.com/eliothedeman/smol/unit"
)

type testUnit struct{}

func (t *testUnit) Init(ctx unit.Ctx) {}
func (t *testUnit) Handle(ctx unit.Ctx, from unit.UnitRef, message any) error {
	return nil
}

func TestNewInstructionExecutor(t *testing.T) {
	ie := NewInstructionExecutor()
	if ie == nil {
		t.Fatal("Expected non-nil InstructionExecutor")
	}
}

func TestParseCommand(t *testing.T) {
	ie := NewInstructionExecutor()

	tests := []struct {
		name     string
		input    string
		expected Instruction
		wantErr  bool
	}{
		{
			name:  "simple command",
			input: "help",
			expected: Instruction{
				Type:   CmdHelp,
				Action: "",
				Args:   []string{},
			},
			wantErr: false,
		},
		{
			name:  "command with action",
			input: "query units",
			expected: Instruction{
				Type:   CmdQuery,
				Action: "units",
				Args:   []string{},
			},
			wantErr: false,
		},
		{
			name:  "command with args",
			input: "execute script.py arg1 arg2",
			expected: Instruction{
				Type:   CmdExecute,
				Action: "script.py",
				Args:   []string{"arg1", "arg2"},
			},
			wantErr: false,
		},
		{
			name:     "empty command",
			input:    "",
			expected: Instruction{},
			wantErr:  true,
		},
		{
			name:     "whitespace only",
			input:    "   ",
			expected: Instruction{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ie.parseCommand(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if result.Type != tt.expected.Type {
					t.Errorf("parseCommand() Type = %v, want %v", result.Type, tt.expected.Type)
				}
				if result.Action != tt.expected.Action {
					t.Errorf("parseCommand() Action = %v, want %v", result.Action, tt.expected.Action)
				}
				if len(result.Args) != len(tt.expected.Args) {
					t.Errorf("parseCommand() Args length = %v, want %v", len(result.Args), len(tt.expected.Args))
				}
			}
		})
	}
}

func TestHandleStringCommand(t *testing.T) {
	ie := NewInstructionExecutor()
	ctx := &mockCtx{
		Context: context.Background(),
		units: []unit.UnitDesc{
			{Name: "test1", Proxy: &testUnit{}},
			{Name: "test2", Proxy: &testUnit{}},
		},
	}

	ie.Init(ctx)

	from := &mockUnitRef{name: "test"}

	err := ie.Handle(ctx, from, "help")
	if err != nil {
		t.Errorf("Handle help command failed: %v", err)
	}

	err = ie.Handle(ctx, from, "list")
	if err != nil {
		t.Errorf("Handle list command failed: %v", err)
	}

	err = ie.Handle(ctx, from, "query test1")
	if err != nil {
		t.Errorf("Handle query command failed: %v", err)
	}
}

func TestHandleInstruction(t *testing.T) {
	ie := NewInstructionExecutor()
	ctx := &mockCtx{
		Context: context.Background(),
		units: []unit.UnitDesc{
			{Name: "test1", Proxy: &testUnit{}},
		},
	}

	ie.Init(ctx)
	from := &mockUnitRef{name: "test"}

	instruction := Instruction{
		Type:    CmdQuery,
		Action:  "test1",
		Args:    []string{},
		Context: make(map[string]any),
	}

	err := ie.Handle(ctx, from, instruction)
	if err != nil {
		t.Errorf("Handle instruction failed: %v", err)
	}
}

func TestHandleMapCommand(t *testing.T) {
	ie := NewInstructionExecutor()
	ctx := &mockCtx{
		Context: context.Background(),
		units: []unit.UnitDesc{
			{Name: "test1", Proxy: &testUnit{}},
		},
	}

	ie.Init(ctx)
	from := &mockUnitRef{name: "test"}

	cmd := map[string]any{
		"type":   "query",
		"action": "test1",
		"args":   []string{},
	}

	err := ie.Handle(ctx, from, cmd)
	if err != nil {
		t.Errorf("Handle map command failed: %v", err)
	}
}

func TestRegisterCommand(t *testing.T) {
	ie := NewInstructionExecutor()

	customHandler := func(ctx context.Context, args []string) (CommandResult, error) {
		return CommandResult{
			Success: true,
			Output:  "custom command executed",
		}, nil
	}

	ie.RegisterCommand("custom", customHandler)

	instruction := Instruction{
		Type:    CommandType("custom"),
		Action:  "test",
		Args:    []string{},
		Context: make(map[string]any),
	}

	ctx := &mockCtx{Context: context.Background()}
	ie.Init(ctx)
	from := &mockUnitRef{name: "test"}

	err := ie.Handle(ctx, from, instruction)
	if err != nil {
		t.Errorf("Handle custom command failed: %v", err)
	}
}

func TestUnknownCommand(t *testing.T) {
	ie := NewInstructionExecutor()
	ctx := &mockCtx{Context: context.Background()}
	ie.Init(ctx)
	from := &mockUnitRef{name: "test"}

	err := ie.Handle(ctx, from, "unknown command")
	if err == nil {
		t.Error("Expected error for unknown command, got nil")
	}
}

func TestConcurrentCommandHandling(t *testing.T) {
	ie := NewInstructionExecutor()
	ctx := &mockCtx{
		Context: context.Background(),
		units: []unit.UnitDesc{
			{Name: "test1", Proxy: &testUnit{}},
		},
	}

	ie.Init(ctx)

	from := &mockUnitRef{name: "test"}
	done := make(chan bool)

	for i := 0; i < 10; i++ {
		go func() {
			err := ie.Handle(ctx, from, "help")
			if err != nil {
				t.Errorf("Concurrent command failed: %v", err)
			}
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		select {
		case <-done:
		case <-time.After(1 * time.Second):
			t.Fatal("Timeout waiting for concurrent commands")
		}
	}
}
