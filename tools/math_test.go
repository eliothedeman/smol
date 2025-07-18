package tools

import (
	"context"
	"testing"
)

// Reuse mock types from registers_test.go

func TestMathInit(t *testing.T) {
	math := &Math{}
	ctx := &mockCtx{Context: context.Background()}
	math.Init(ctx)
}

func TestMathBasicOperations(t *testing.T) {
	math := NewMath()

	// Test basic operations
	if result := math.Add(2, 3); result != 5 {
		t.Errorf("Expected 5, got %f", result)
	}

	if result := math.Subtract(5, 3); result != 2 {
		t.Errorf("Expected 2, got %f", result)
	}

	if result := math.Multiply(3, 4); result != 12 {
		t.Errorf("Expected 12, got %f", result)
	}

	if result, err := math.Divide(10, 2); err != nil || result != 5 {
		t.Errorf("Expected 5, got %f, err: %v", result, err)
	}

	if result := math.Power(2, 3); result != 8 {
		t.Errorf("Expected 8, got %f", result)
	}

	if result, err := math.Sqrt(9); err != nil || result != 3 {
		t.Errorf("Expected 3, got %f, err: %v", result, err)
	}

	if result := math.Sin(1.5707963267948966); result != 1 {
		t.Errorf("Expected 1, got %f", result)
	}

	if result := math.Cos(0); result != 1 {
		t.Errorf("Expected 1, got %f", result)
	}

	if result := math.Tan(0); result != 0 {
		t.Errorf("Expected 0, got %f", result)
	}

	if result, err := math.Log(2.718281828459045); err != nil || result != 1 {
		t.Errorf("Expected 1, got %f, err: %v", result, err)
	}

	if result := math.Round(3.14159, 2); result != 3.14 {
		t.Errorf("Expected 3.14, got %f", result)
	}
}

func TestMathErrorCases(t *testing.T) {
	math := NewMath()

	// Test division by zero
	if _, err := math.Divide(5, 0); err == nil {
		t.Error("Expected division by zero error")
	}

	// Test negative square root
	if _, err := math.Sqrt(-1); err == nil {
		t.Error("Expected negative sqrt error")
	}

	// Test invalid log
	if _, err := math.Log(-1); err == nil {
		t.Error("Expected invalid log error")
	}

	// Test invalid number parsing
	if _, err := math.ParseNumber("invalid"); err == nil {
		t.Error("Expected parse error")
	}
}

func TestMathHandleStringCommands(t *testing.T) {
	math := &Math{}
	ctx := &mockCtx{Context: context.Background()}
	math.Init(ctx)

	from := &testMessageHandler{}

	// Test add command
	err := math.Handle(ctx, from, "add 2 3 4")
	if err != nil {
		t.Errorf("Handle add command failed: %v", err)
	}

	// Test subtract command
	err = math.Handle(ctx, from, "sub 10 3")
	if err != nil {
		t.Errorf("Handle subtract command failed: %v", err)
	}

	// Test multiply command
	err = math.Handle(ctx, from, "mul 3 4")
	if err != nil {
		t.Errorf("Handle multiply command failed: %v", err)
	}

	// Test divide command
	err = math.Handle(ctx, from, "div 10 2")
	if err != nil {
		t.Errorf("Handle divide command failed: %v", err)
	}

	// Test power command
	err = math.Handle(ctx, from, "pow 2 3")
	if err != nil {
		t.Errorf("Handle power command failed: %v", err)
	}

	// Test sqrt command
	err = math.Handle(ctx, from, "sqrt 9")
	if err != nil {
		t.Errorf("Handle sqrt command failed: %v", err)
	}

	// Test help command
	err = math.Handle(ctx, from, "help")
	if err != nil {
		t.Errorf("Handle help command failed: %v", err)
	}
}

func TestMathHandleMapCommands(t *testing.T) {
	math := &Math{}
	ctx := &mockCtx{Context: context.Background()}
	math.Init(ctx)

	from := &testMessageHandler{}

	// Test add action
	err := math.Handle(ctx, from, map[string]interface{}{
		"action":  "add",
		"numbers": []interface{}{2.0, 3.0, 4.0},
	})
	if err != nil {
		t.Errorf("Handle add action failed: %v", err)
	}

	// Test subtract action
	err = math.Handle(ctx, from, map[string]interface{}{
		"action": "subtract",
		"a":      10.0,
		"b":      3.0,
	})
	if err != nil {
		t.Errorf("Handle subtract action failed: %v", err)
	}

	// Test divide action
	err = math.Handle(ctx, from, map[string]interface{}{
		"action": "divide",
		"a":      10.0,
		"b":      2.0,
	})
	if err != nil {
		t.Errorf("Handle divide action failed: %v", err)
	}
}

func TestMathParseNumbers(t *testing.T) {
	math := NewMath()

	// Test parse single number
	if num, err := math.parseNumber("42"); err != nil || num != 42 {
		t.Errorf("Expected 42, got %f, err: %v", num, err)
	}

	// Test parse multiple numbers
	numbers, err := math.parseNumbers([]string{"1", "2", "3"})
	if err != nil {
		t.Errorf("Parse numbers failed: %v", err)
	}
	if len(numbers) != 3 || numbers[0] != 1 || numbers[1] != 2 || numbers[2] != 3 {
		t.Errorf("Unexpected numbers: %v", numbers)
	}

	// Test parse invalid number
	if _, err := math.parseNumbers([]string{"1", "invalid", "3"}); err == nil {
		t.Error("Expected parse error")
	}
}

func TestMathSumAndProduct(t *testing.T) {
	math := NewMath()

	numbers := []float64{1, 2, 3, 4, 5}

	if sum := math.sum(numbers); sum != 15 {
		t.Errorf("Expected 15, got %f", sum)
	}

	if product := math.product(numbers); product != 120 {
		t.Errorf("Expected 120, got %f", product)
	}
}

func TestMathExtractNumbers(t *testing.T) {
	math := NewMath()

	// Test extract single number
	if num, err := math.extractNumber(42.0); err != nil || num != 42 {
		t.Errorf("Expected 42, got %f, err: %v", num, err)
	}

	// Test extract from string
	if num, err := math.extractNumber("3.14"); err != nil || num != 3.14 {
		t.Errorf("Expected 3.14, got %f, err: %v", num, err)
	}

	// Test extract numbers array
	numbers, err := math.extractNumbers([]interface{}{1.0, 2.0, 3.0})
	if err != nil {
		t.Errorf("Extract numbers failed: %v", err)
	}
	if len(numbers) != 3 || numbers[0] != 1 || numbers[1] != 2 || numbers[2] != 3 {
		t.Errorf("Unexpected numbers: %v", numbers)
	}

	// Test extract invalid type
	if _, err := math.extractNumber("invalid"); err == nil {
		t.Error("Expected extract error")
	}
}

func TestMathErrorHandling(t *testing.T) {
	math := &Math{}
	ctx := &mockCtx{Context: context.Background()}
	math.Init(ctx)

	from := &testMessageHandler{}

	// Test invalid command
	err := math.Handle(ctx, from, "invalid 1 2")
	if err == nil {
		t.Error("Expected invalid command error")
	}

	// Test invalid numbers
	err = math.Handle(ctx, from, "add 1 invalid")
	if err == nil {
		t.Error("Expected invalid number error")
	}

	// Test division by zero
	err = math.Handle(ctx, from, "div 1 0")
	if err == nil {
		t.Error("Expected division by zero error")
	}

	// Test missing action
	err = math.Handle(ctx, from, map[string]interface{}{
		"numbers": []interface{}{1.0, 2.0},
	})
	if err == nil {
		t.Error("Expected missing action error")
	}
}
