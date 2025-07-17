package tools

import (
	"fmt"
	"sync"
)

type Register struct {
	name  string
	value interface{}
}

type Registers struct {
	registers map[string]interface{}
	mu        sync.RWMutex
}

func NewRegisters() *Registers {
	return &Registers{
		registers: make(map[string]interface{}),
	}
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

func (r *Registers) GetInt(name string) (int, bool) {
	val, exists := r.Get(name)
	if !exists {
		return 0, false
	}

	if intVal, ok := val.(int); ok {
		return intVal, true
	}
	return 0, false
}

func (r *Registers) GetFloat(name string) (float64, bool) {
	val, exists := r.Get(name)
	if !exists {
		return 0, false
	}

	if floatVal, ok := val.(float64); ok {
		return floatVal, true
	}
	return 0, false
}

func (r *Registers) GetString(name string) (string, bool) {
	val, exists := r.Get(name)
	if !exists {
		return "", false
	}

	if strVal, ok := val.(string); ok {
		return strVal, true
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

func (r *Registers) Info() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return fmt.Sprintf("Registers: %d items", len(r.registers))
}
