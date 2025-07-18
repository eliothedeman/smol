package tools

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// Registers manages named values
type Registers struct {
	mu        sync.RWMutex
	registers map[string]interface{}
}

// NewRegisters creates a new Registers instance
func NewRegisters() *Registers {
	return &Registers{
		registers: make(map[string]interface{}),
	}
}

// Handle processes register commands
func (r *Registers) Handle(cmd string, from MessageSender) error {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}

	action := parts[0]

	switch action {
	case "set":
		if len(parts) < 3 {
			return fmt.Errorf("set requires name and value")
		}
		name := parts[1]
		value := strings.Join(parts[2:], " ")
		parsed, err := r.parseValue(value)
		if err != nil {
			return fmt.Errorf("invalid value: %v", err)
		}
		r.Set(name, parsed)
		from.Send(fmt.Sprintf("%s = %v", name, parsed))

	case "get":
		if len(parts) < 2 {
			return fmt.Errorf("get requires name")
		}
		name := parts[1]
		if val, exists := r.Get(name); exists {
			from.Send(fmt.Sprintf("%s = %v", name, val))
		} else {
			from.Send(fmt.Sprintf("%s not found", name))
		}

	case "list":
		regs := r.List()
		if len(regs) == 0 {
			from.Send("No registers")
		} else {
			result := "Registers:\n"
			for name, val := range regs {
				result += fmt.Sprintf("  %s = %v\n", name, val)
			}
			from.Send(result)
		}

	case "clear":
		r.Clear()
		from.Send("All registers cleared")

	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}

	return nil
}

func (r *Registers) parseValue(s string) (interface{}, error) {
	// Try int
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return i, nil
	}
	// Try float
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f, nil
	}
	// Try bool
	if s == "true" {
		return true, nil
	}
	if s == "false" {
		return false, nil
	}
	// Return as string
	return s, nil
}

func (r *Registers) Set(name string, value interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.registers[name] = value
}

func (r *Registers) Get(name string) (interface{}, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	val, exists := r.registers[name]
	return val, exists
}

func (r *Registers) GetInt(name string) (int64, bool) {
	if val, exists := r.Get(name); exists {
		if i, ok := val.(int64); ok {
			return i, true
		}
		if f, ok := val.(float64); ok {
			return int64(f), true
		}
	}
	return 0, false
}

func (r *Registers) GetFloat(name string) (float64, bool) {
	if val, exists := r.Get(name); exists {
		if f, ok := val.(float64); ok {
			return f, true
		}
		if i, ok := val.(int64); ok {
			return float64(i), true
		}
	}
	return 0, false
}

func (r *Registers) GetString(name string) (string, bool) {
	if val, exists := r.Get(name); exists {
		if s, ok := val.(string); ok {
			return s, true
		}
		return fmt.Sprintf("%v", val), true
	}
	return "", false
}

func (r *Registers) List() map[string]interface{} {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make(map[string]interface{})
	for k, v := range r.registers {
		result[k] = v
	}
	return result
}

func (r *Registers) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.registers = make(map[string]interface{})
}
