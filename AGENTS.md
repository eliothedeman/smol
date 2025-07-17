# AGENTS.md

## Build/Lint/Test
- Build: `go build`
- Lint: `go vet` and `gofmt -s -w .`
- Test all: `go test ./...`
- Single test: `go test -run TestName`

## Code Style
- Imports: Group stdlib, third-party, local; use `goimports`
- Formatting: Use `gofmt -s` (no trailing commas)
- Types: Prefer explicit types over `interface{}`
- Naming: CamelCase for functions/variables, snake_case for constants
- Error handling: Check errors explicitly, use `fmt.Errorf` for wrapping

## Tool calling rules
- You must always read a file before overwriting it. Use the Read tool first

