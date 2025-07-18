package tools

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/eliothedeman/smol/unit"
)

type Math struct {
	ctx unit.Ctx
}

func NewMath() *Math {
	return &Math{}
}

func (m *Math) Add(a, b float64) float64 {
	return a + b
}

func (m *Math) Subtract(a, b float64) float64 {
	return a - b
}

func (m *Math) Multiply(a, b float64) float64 {
	return a * b
}

func (m *Math) Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivisionByZero
	}
	return a / b, nil
}

func (m *Math) Power(base, exponent float64) float64 {
	return math.Pow(base, exponent)
}

func (m *Math) Sqrt(x float64) (float64, error) {
	if x < 0 {
		return 0, ErrNegativeSqrt
	}
	return math.Sqrt(x), nil
}

func (m *Math) Sin(x float64) float64 {
	return math.Sin(x)
}

func (m *Math) Cos(x float64) float64 {
	return math.Cos(x)
}

func (m *Math) Tan(x float64) float64 {
	return math.Tan(x)
}

func (m *Math) Log(x float64) (float64, error) {
	if x <= 0 {
		return 0, ErrInvalidLog
	}
	return math.Log(x), nil
}

func (m *Math) ParseNumber(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func (m *Math) Round(x float64, decimals int) float64 {
	factor := math.Pow(10, float64(decimals))
	return math.Round(x*factor) / factor
}

func (m *Math) Init(ctx unit.Ctx) {
	m.ctx = ctx
}

func (m *Math) Handle(ctx unit.Ctx, from unit.UnitRef, message any) error {
	switch msg := message.(type) {
	case string:
		return m.handleStringCommand(ctx, from, msg)
	case map[string]interface{}:
		return m.handleMapCommand(ctx, from, msg)
	default:
		return fmt.Errorf("unsupported message type: %T", message)
	}
}

func (m *Math) handleStringCommand(ctx unit.Ctx, from unit.UnitRef, command string) error {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}

	cmd := strings.ToLower(parts[0])

	switch cmd {
	case "add", "sum":
		if len(parts) < 3 {
			return fmt.Errorf("add requires at least 2 numbers")
		}
		result, err := m.parseNumbers(parts[1:])
		if err != nil {
			return err
		}
		from.Send(fmt.Sprintf("%.6f", m.sum(result)))

	case "sub", "subtract":
		if len(parts) != 3 {
			return fmt.Errorf("subtract requires exactly 2 numbers")
		}
		a, b, err := m.parseTwoNumbers(parts[1], parts[2])
		if err != nil {
			return err
		}
		from.Send(fmt.Sprintf("%.6f", m.Subtract(a, b)))

	case "mul", "multiply":
		if len(parts) < 3 {
			return fmt.Errorf("multiply requires at least 2 numbers")
		}
		result, err := m.parseNumbers(parts[1:])
		if err != nil {
			return err
		}
		from.Send(fmt.Sprintf("%.6f", m.product(result)))

	case "div", "divide":
		if len(parts) != 3 {
			return fmt.Errorf("divide requires exactly 2 numbers")
		}
		a, b, err := m.parseTwoNumbers(parts[1], parts[2])
		if err != nil {
			return err
		}
		result, err := m.Divide(a, b)
		if err != nil {
			return err
		}
		from.Send(fmt.Sprintf("%.6f", result))

	case "pow", "power":
		if len(parts) != 3 {
			return fmt.Errorf("power requires exactly 2 numbers")
		}
		base, exp, err := m.parseTwoNumbers(parts[1], parts[2])
		if err != nil {
			return err
		}
		from.Send(fmt.Sprintf("%.6f", m.Power(base, exp)))

	case "sqrt":
		if len(parts) != 2 {
			return fmt.Errorf("sqrt requires exactly 1 number")
		}
		x, err := m.parseNumber(parts[1])
		if err != nil {
			return err
		}
		result, err := m.Sqrt(x)
		if err != nil {
			return err
		}
		from.Send(fmt.Sprintf("%.6f", result))

	case "sin":
		if len(parts) != 2 {
			return fmt.Errorf("sin requires exactly 1 number")
		}
		x, err := m.parseNumber(parts[1])
		if err != nil {
			return err
		}
		from.Send(fmt.Sprintf("%.6f", m.Sin(x)))

	case "cos":
		if len(parts) != 2 {
			return fmt.Errorf("cos requires exactly 1 number")
		}
		x, err := m.parseNumber(parts[1])
		if err != nil {
			return err
		}
		from.Send(fmt.Sprintf("%.6f", m.Cos(x)))

	case "tan":
		if len(parts) != 2 {
			return fmt.Errorf("tan requires exactly 1 number")
		}
		x, err := m.parseNumber(parts[1])
		if err != nil {
			return err
		}
		from.Send(fmt.Sprintf("%.6f", m.Tan(x)))

	case "log":
		if len(parts) != 2 {
			return fmt.Errorf("log requires exactly 1 number")
		}
		x, err := m.parseNumber(parts[1])
		if err != nil {
			return err
		}
		result, err := m.Log(x)
		if err != nil {
			return err
		}
		from.Send(fmt.Sprintf("%.6f", result))

	case "round":
		if len(parts) < 2 || len(parts) > 3 {
			return fmt.Errorf("round requires 1-2 numbers")
		}
		x, err := m.parseNumber(parts[1])
		if err != nil {
			return err
		}
		decimals := 0
		if len(parts) == 3 {
			if d, err := m.parseNumber(parts[2]); err == nil {
				decimals = int(d)
			}
		}
		from.Send(fmt.Sprintf("%.6f", m.Round(x, decimals)))

	case "help":
		from.Send(`Math commands:
  add <num1> <num2> [...] - Add numbers
  sub <a> <b> - Subtract b from a
  mul <num1> <num2> [...] - Multiply numbers
  div <a> <b> - Divide a by b
  pow <base> <exp> - Power operation
  sqrt <x> - Square root
  sin <x> - Sine function
  cos <x> - Cosine function
  tan <x> - Tangent function
  log <x> - Natural logarithm
  round <x> [decimals] - Round number`)

	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}

	return nil
}

func (m *Math) handleMapCommand(ctx unit.Ctx, from unit.UnitRef, command map[string]interface{}) error {
	action, ok := command["action"].(string)
	if !ok {
		return fmt.Errorf("missing action field")
	}

	switch action {
	case "add":
		numbers, err := m.extractNumbers(command["numbers"])
		if err != nil {
			return err
		}
		from.Send(m.sum(numbers))

	case "subtract":
		a, b, err := m.extractTwoNumbers(command["a"], command["b"])
		if err != nil {
			return err
		}
		from.Send(m.Subtract(a, b))

	case "multiply":
		numbers, err := m.extractNumbers(command["numbers"])
		if err != nil {
			return err
		}
		from.Send(m.product(numbers))

	case "divide":
		a, b, err := m.extractTwoNumbers(command["a"], command["b"])
		if err != nil {
			return err
		}
		result, err := m.Divide(a, b)
		if err != nil {
			return err
		}
		from.Send(result)

	case "power":
		base, exp, err := m.extractTwoNumbers(command["base"], command["exponent"])
		if err != nil {
			return err
		}
		from.Send(m.Power(base, exp))

	case "sqrt":
		x, err := m.extractNumber(command["x"])
		if err != nil {
			return err
		}
		result, err := m.Sqrt(x)
		if err != nil {
			return err
		}
		from.Send(result)

	case "sin":
		x, err := m.extractNumber(command["x"])
		if err != nil {
			return err
		}
		from.Send(m.Sin(x))

	case "cos":
		x, err := m.extractNumber(command["x"])
		if err != nil {
			return err
		}
		from.Send(m.Cos(x))

	case "tan":
		x, err := m.extractNumber(command["x"])
		if err != nil {
			return err
		}
		from.Send(m.Tan(x))

	case "log":
		x, err := m.extractNumber(command["x"])
		if err != nil {
			return err
		}
		result, err := m.Log(x)
		if err != nil {
			return err
		}
		from.Send(result)

	case "round":
		x, err := m.extractNumber(command["x"])
		if err != nil {
			return err
		}
		decimals := 0
		if d, ok := command["decimals"].(float64); ok {
			decimals = int(d)
		}
		from.Send(m.Round(x, decimals))

	default:
		return fmt.Errorf("unknown action: %s", action)
	}

	return nil
}

func (m *Math) parseNumber(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func (m *Math) parseNumbers(strings []string) ([]float64, error) {
	numbers := make([]float64, len(strings))
	for i, s := range strings {
		num, err := m.parseNumber(s)
		if err != nil {
			return nil, fmt.Errorf("invalid number: %s", s)
		}
		numbers[i] = num
	}
	return numbers, nil
}

func (m *Math) parseTwoNumbers(a, b string) (float64, float64, error) {
	num1, err := m.parseNumber(a)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid first number: %s", a)
	}
	num2, err := m.parseNumber(b)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid second number: %s", b)
	}
	return num1, num2, nil
}

func (m *Math) extractNumber(val interface{}) (float64, error) {
	switch v := val.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case string:
		return m.parseNumber(v)
	default:
		return 0, fmt.Errorf("invalid number type: %T", val)
	}
}

func (m *Math) extractTwoNumbers(a, b interface{}) (float64, float64, error) {
	num1, err := m.extractNumber(a)
	if err != nil {
		return 0, 0, err
	}
	num2, err := m.extractNumber(b)
	if err != nil {
		return 0, 0, err
	}
	return num1, num2, nil
}

func (m *Math) extractNumbers(val interface{}) ([]float64, error) {
	switch v := val.(type) {
	case []interface{}:
		numbers := make([]float64, len(v))
		for i, item := range v {
			num, err := m.extractNumber(item)
			if err != nil {
				return nil, err
			}
			numbers[i] = num
		}
		return numbers, nil
	case []float64:
		return v, nil
	default:
		return nil, fmt.Errorf("invalid numbers format: %T", val)
	}
}

func (m *Math) sum(numbers []float64) float64 {
	sum := 0.0
	for _, num := range numbers {
		sum += num
	}
	return sum
}

func (m *Math) product(numbers []float64) float64 {
	product := 1.0
	for _, num := range numbers {
		product *= num
	}
	return product
}

var (
	ErrDivisionByZero = fmt.Errorf("division by zero")
	ErrNegativeSqrt   = fmt.Errorf("square root of negative number")
	ErrInvalidLog     = fmt.Errorf("logarithm of non-positive number")
)
