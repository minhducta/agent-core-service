---
applyTo: "**/*.go"
---

# Go Development Instructions

Follow idiomatic Go: [Effective Go](https://go.dev/doc/effective_go) + [Code Review Comments](https://go.dev/wiki/CodeReviewComments).

## Code Style
- Write simple, clear, idiomatic Go
- Keep happy path left-aligned (return early, reduce nesting)
- Use `gofmt` and `goimports`
- No emoji in code/comments
- Each `.go` file has exactly ONE `package` declaration

## Naming
- `mixedCaps` not underscores
- Exported names start with capital
- Avoid stuttering (`http.Server` not `http.HTTPServer`)

## Error Handling
- Check errors immediately
- Wrap with `fmt.Errorf("context: %w", err)`
- Use `errors.Is` / `errors.As`
- Error messages: lowercase, no trailing punctuation

## Testing
- Table-driven tests: slice of anonymous structs + `t.Run`
- Name: `Test_functionName_scenario`
- Use `testify/assert` and `testify/require`
- Target >80% coverage
- `t.Helper()` for helper functions
