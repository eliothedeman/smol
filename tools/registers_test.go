package tools

import (
	"context"
	"testing"

	"github.com/eliothedeman/smol/unit"
)

type mockCtx struct {
	context.Context
	units []unit.UnitDesc
}

func (m *mockCtx) Units() []unit.UnitDesc {
	return m.units
}

func (m *mockCtx) Spawn(name string, f unit.UnitFactory) unit.UnitRef {
	return &mockUnitRef{name: name}
}

func (m *mockCtx) Self() unit.UnitRef {
	return &mockUnitRef{name: "test"}
}

func (m *mockCtx) Subscribe(other unit.Unit)   {}
func (m *mockCtx) Unsubscribe(other unit.Unit) {}

type mockUnitRef struct {
	name string
}

func (m *mockUnitRef) Name() string {
	return m.name
}

func (m *mockUnitRef) Send(msg any) {}
func (m *mockUnitRef) Stop()        {}

type testMessageHandler struct {
	lastMessage any
	name        string
}

func (t *testMessageHandler) Name() string {
	return "test"
}

func (t *testMessageHandler) Send(msg any) {
	t.lastMessage = msg
}

func (t *testMessageHandler) Stop() {}

func TestRegistersInit(t *testing.T) {
	registers := &Registers{}
	ctx := &mockCtx{Context: context.Background()}
	registers.Init(ctx)
}

func TestRegistersSetGet(t *testing.T) {
	registers := NewRegisters()
	ctx := &mockCtx{Context: context.Background()}
	registers.Init(ctx)

	// Test setting and getting values
	registers.Set("test", 42)
	if val, exists := registers.Get("test"); !exists || val != 42 {
		t.Errorf("Expected 42, got %v", val)
	}

	// Test type-specific getters
	registers.Set("int_val", 100)
	if val, exists := registers.GetInt("int_val"); !exists || val != 100 {
		t.Errorf("Expected 100, got %v", val)
	}

	registers.Set("float_val", 3.14)
	if val, exists := registers.GetFloat("float_val"); !exists || val != 3.14 {
		t.Errorf("Expected 3.14, got %v", val)
	}

	registers.Set("string_val", "hello")
	if val, exists := registers.GetString("string_val"); !exists || val != "hello" {
		t.Errorf("Expected 'hello', got %v", val)
	}
}

func TestRegistersList(t *testing.T) {
	registers := NewRegisters()
	registers.Set("a", 1)
	registers.Set("b", 2)

	list := registers.List()
	if len(list) != 2 {
		t.Errorf("Expected 2 items, got %d", len(list))
	}
	if list["a"] != 1 || list["b"] != 2 {
		t.Errorf("Unexpected values in list: %v", list)
	}
}

func TestRegistersClear(t *testing.T) {
	registers := NewRegisters()
	registers.Set("test", 42)
	registers.Clear()

	if _, exists := registers.Get("test"); exists {
		t.Error("Expected register to be cleared")
	}
}

func TestRegistersHandleStringCommands(t *testing.T) {
	registers := NewRegisters()
	ctx := &mockCtx{Context: context.Background()}
	registers.Init(ctx)

	from := &testMessageHandler{}

	// Test set command
	err := registers.Handle(ctx, from, "set test 42")
	if err != nil {
		t.Errorf("Handle set command failed: %v", err)
	}

	// Test get command
	err = registers.Handle(ctx, from, "get test")
	if err != nil {
		t.Errorf("Handle get command failed: %v", err)
	}

	// Test list command
	err = registers.Handle(ctx, from, "list")
	if err != nil {
		t.Errorf("Handle list command failed: %v", err)
	}

	// Test clear command
	err = registers.Handle(ctx, from, "clear")
	if err != nil {
		t.Errorf("Handle clear command failed: %v", err)
	}
}

func TestRegistersHandleMapCommands(t *testing.T) {
	registers := NewRegisters()
	ctx := &mockCtx{Context: context.Background()}
	registers.Init(ctx)

	from := &testMessageHandler{}

	// Test set action
	err := registers.Handle(ctx, from, map[string]interface{}{
		"action": "set",
		"name":   "test",
		"value":  42,
	})
	if err != nil {
		t.Errorf("Handle set action failed: %v", err)
	}

	// Test get action
	err = registers.Handle(ctx, from, map[string]interface{}{
		"action": "get",
		"name":   "test",
	})
	if err != nil {
		t.Errorf("Handle get action failed: %v", err)
	}
}

func TestRegistersParseValue(t *testing.T) {
	registers := &Registers{}

	tests := []struct {
		input    string
		expected interface{}
	}{
		{"42", 42},
		{"3.14", 3.14},
		{"true", true},
		{"false", false},
		{"hello", "hello"},
		{"123abc", "123abc"},
	}

	for _, test := range tests {
		result := registers.parseValue(test.input)
		if result != test.expected {
			t.Errorf("Expected %v (%T), got %v (%T)", test.expected, test.expected, result, result)
		}
	}
}
