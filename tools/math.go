package tools

import (
	"fmt"
	"math"
	"strconv"
)

type Math struct{}

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

var (
	ErrDivisionByZero = fmt.Errorf("division by zero")
	ErrNegativeSqrt   = fmt.Errorf("square root of negative number")
	ErrInvalidLog     = fmt.Errorf("logarithm of non-positive number")
)
