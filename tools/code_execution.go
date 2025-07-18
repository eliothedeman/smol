package tools

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type CodeExecution struct {
	workDir string
	timeout time.Duration
}

type ExecutionRequest struct {
	Code     string            `json:"code"`
	Language string            `json:"language"`
	Env      map[string]string `json:"env,omitempty"`
}

type ExecutionResult struct {
	Output   string `json:"output"`
	Error    string `json:"error,omitempty"`
	ExitCode int    `json:"exit_code"`
	Duration string `json:"duration"`
}

func NewCodeExecution(workDir string) *CodeExecution {
	return &CodeExecution{
		workDir: workDir,
		timeout: 30 * time.Second,
	}
}

func (ce *CodeExecution) Execute(req ExecutionRequest) (*ExecutionResult, error) {
	start := time.Now()

	switch strings.ToLower(req.Language) {
	case "go":
		return ce.executeGo(req, start)
	case "python":
		return ce.executePython(req, start)
	case "bash":
		return ce.executeBash(req, start)
	default:
		return nil, fmt.Errorf("unsupported language: %s", req.Language)
	}
}

func (ce *CodeExecution) executeGo(req ExecutionRequest, start time.Time) (*ExecutionResult, error) {
	file, err := os.CreateTemp(ce.workDir, "code_*.go")
	if err != nil {
		return nil, err
	}
	defer os.Remove(file.Name())

	if _, err := file.WriteString(req.Code); err != nil {
		return nil, err
	}
	file.Close()

	ctx, cancel := context.WithTimeout(context.Background(), ce.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "run", file.Name())
	cmd.Env = ce.buildEnv(req.Env)

	output, err := cmd.CombinedOutput()
	duration := time.Since(start)

	result := &ExecutionResult{
		Output:   string(output),
		Duration: duration.String(),
	}

	if err != nil {
		result.Error = err.Error()
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		}
	}

	return result, nil
}

func (ce *CodeExecution) executePython(req ExecutionRequest, start time.Time) (*ExecutionResult, error) {
	file, err := os.CreateTemp(ce.workDir, "code_*.py")
	if err != nil {
		return nil, err
	}
	defer os.Remove(file.Name())

	if _, err := file.WriteString(req.Code); err != nil {
		return nil, err
	}
	file.Close()

	ctx, cancel := context.WithTimeout(context.Background(), ce.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "python3", file.Name())
	cmd.Env = ce.buildEnv(req.Env)

	output, err := cmd.CombinedOutput()
	duration := time.Since(start)

	result := &ExecutionResult{
		Output:   string(output),
		Duration: duration.String(),
	}

	if err != nil {
		result.Error = err.Error()
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		}
	}

	return result, nil
}

func (ce *CodeExecution) executeBash(req ExecutionRequest, start time.Time) (*ExecutionResult, error) {
	file, err := os.CreateTemp(ce.workDir, "code_*.sh")
	if err != nil {
		return nil, err
	}
	defer os.Remove(file.Name())

	if _, err := file.WriteString(req.Code); err != nil {
		return nil, err
	}
	file.Close()

	ctx, cancel := context.WithTimeout(context.Background(), ce.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "bash", file.Name())
	cmd.Env = ce.buildEnv(req.Env)

	output, err := cmd.CombinedOutput()
	duration := time.Since(start)

	result := &ExecutionResult{
		Output:   string(output),
		Duration: duration.String(),
	}

	if err != nil {
		result.Error = err.Error()
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		}
	}

	return result, nil
}

func (ce *CodeExecution) buildEnv(env map[string]string) []string {
	baseEnv := os.Environ()
	for k, v := range env {
		baseEnv = append(baseEnv, fmt.Sprintf("%s=%s", k, v))
	}
	return baseEnv
}
