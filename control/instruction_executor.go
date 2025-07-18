package control

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/eliothedeman/smol/unit"
)

type CommandType string

const (
	CmdExecute CommandType = "execute"
	CmdQuery   CommandType = "query"
	CmdSet     CommandType = "set"
	CmdGet     CommandType = "get"
	CmdList    CommandType = "list"
	CmdHelp    CommandType = "help"
)

type Instruction struct {
	Type    CommandType
	Action  string
	Args    []string
	Context map[string]any
}

type CommandResult struct {
	Success bool
	Output  string
	Error   error
	Data    any
}

type InstructionExecutor struct {
	mu       sync.RWMutex
	commands map[CommandType]CommandHandler
	ctx      unit.Ctx
}

type CommandHandler func(ctx context.Context, args []string) (CommandResult, error)

func NewInstructionExecutor() *InstructionExecutor {
	ie := &InstructionExecutor{
		commands: make(map[CommandType]CommandHandler),
	}
	ie.registerDefaultCommands()
	return ie
}

func (ie *InstructionExecutor) Init(ctx unit.Ctx) {
	ie.mu.Lock()
	defer ie.mu.Unlock()
	ie.ctx = ctx
}

func (ie *InstructionExecutor) Handle(ctx unit.Ctx, from unit.UnitRef, message any) error {
	switch msg := message.(type) {
	case string:
		return ie.handleStringCommand(ctx, from, msg)
	case Instruction:
		return ie.handleInstruction(ctx, from, msg)
	case map[string]any:
		return ie.handleMapCommand(ctx, from, msg)
	default:
		return fmt.Errorf("unsupported message type: %T", message)
	}
}

func (ie *InstructionExecutor) handleStringCommand(ctx unit.Ctx, from unit.UnitRef, cmd string) error {
	instruction, err := ie.parseCommand(cmd)
	if err != nil {
		return fmt.Errorf("failed to parse command: %w", err)
	}
	return ie.handleInstruction(ctx, from, instruction)
}

func (ie *InstructionExecutor) handleInstruction(ctx unit.Ctx, from unit.UnitRef, instruction Instruction) error {
	handler, exists := ie.commands[instruction.Type]
	if !exists {
		return fmt.Errorf("unknown command type: %s", instruction.Type)
	}

	result, err := handler(ctx, instruction.Args)
	if err != nil {
		result = CommandResult{
			Success: false,
			Error:   err,
			Output:  fmt.Sprintf("Error: %v", err),
		}
	}

	from.Send(result)
	return nil
}

func (ie *InstructionExecutor) handleMapCommand(ctx unit.Ctx, from unit.UnitRef, cmd map[string]any) error {
	instruction := Instruction{
		Context: make(map[string]any),
	}

	if cmdType, ok := cmd["type"].(string); ok {
		instruction.Type = CommandType(cmdType)
	}
	if action, ok := cmd["action"].(string); ok {
		instruction.Action = action
	}
	if args, ok := cmd["args"].([]string); ok {
		instruction.Args = args
	} else if args, ok := cmd["args"].([]any); ok {
		instruction.Args = make([]string, len(args))
		for i, arg := range args {
			if str, ok := arg.(string); ok {
				instruction.Args[i] = str
			}
		}
	}

	return ie.handleInstruction(ctx, from, instruction)
}

func (ie *InstructionExecutor) parseCommand(input string) (Instruction, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return Instruction{}, fmt.Errorf("empty command")
	}

	parts := strings.Fields(input)
	if len(parts) == 0 {
		return Instruction{}, fmt.Errorf("invalid command format")
	}

	cmdType := CommandType(strings.ToLower(parts[0]))
	action := ""
	args := []string{}

	if len(parts) > 1 {
		action = parts[1]
	}
	if len(parts) > 2 {
		args = parts[2:]
	}

	return Instruction{
		Type:    cmdType,
		Action:  action,
		Args:    args,
		Context: make(map[string]any),
	}, nil
}

func (ie *InstructionExecutor) RegisterCommand(cmdType CommandType, handler CommandHandler) {
	ie.mu.Lock()
	defer ie.mu.Unlock()
	ie.commands[cmdType] = handler
}

func (ie *InstructionExecutor) registerDefaultCommands() {
	ie.RegisterCommand(CmdHelp, ie.handleHelp)
	ie.RegisterCommand(CmdList, ie.handleList)
	ie.RegisterCommand(CmdQuery, ie.handleQuery)
}

func (ie *InstructionExecutor) handleHelp(ctx context.Context, args []string) (CommandResult, error) {
	helpText := `Available commands:
  help [command] - Show help information
  list [type] - List available units or resources
  query <target> [args...] - Query information from a unit
  execute <script> - Execute a script or command
  set <key> <value> - Set a configuration value
  get <key> - Get a configuration value`

	return CommandResult{
		Success: true,
		Output:  helpText,
		Data:    helpText,
	}, nil
}

func (ie *InstructionExecutor) handleList(ctx context.Context, args []string) (CommandResult, error) {
	if ie.ctx == nil {
		return CommandResult{
			Success: false,
			Error:   fmt.Errorf("context not initialized"),
		}, nil
	}

	units := ie.ctx.Units()
	var unitNames []string
	for _, unit := range units {
		unitNames = append(unitNames, unit.Name)
	}

	output := fmt.Sprintf("Available units: %v", unitNames)
	return CommandResult{
		Success: true,
		Output:  output,
		Data:    unitNames,
	}, nil
}

func (ie *InstructionExecutor) handleQuery(ctx context.Context, args []string) (CommandResult, error) {
	if len(args) == 0 {
		return CommandResult{
			Success: false,
			Error:   fmt.Errorf("query requires a target"),
		}, nil
	}

	target := args[0]
	if ie.ctx == nil {
		return CommandResult{
			Success: false,
			Error:   fmt.Errorf("context not initialized"),
		}, nil
	}

	units := ie.ctx.Units()
	for _, unit := range units {
		if unit.Name == target {
			return CommandResult{
				Success: true,
				Output:  fmt.Sprintf("Found unit: %s", target),
				Data:    unit,
			}, nil
		}
	}

	return CommandResult{
		Success: false,
		Error:   fmt.Errorf("unit not found: %s", target),
	}, nil
}
