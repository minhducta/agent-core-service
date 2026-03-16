---
applyTo: "**/*_test.go"
---

# Golang Unit Testing Instructions

All test files in this project MUST strictly adhere to the following rules.

## 1. General Testing Principles

- **Always Test New Code:** Whenever new functions, methods, handlers, or logic are generated, the corresponding unit tests MUST be generated immediately.
- **High Coverage:** Aim for > 80% code coverage. Ensure edge cases, happy paths, and error states are all covered.
- **Table-Driven Testing:** All unit tests MUST use the Golang **Table-Driven Testing** pattern (slice of anonymous structs).
- **Assertions:** Use `github.com/stretchr/testify/assert` and `github.com/stretchr/testify/require` for assertions instead of standard `if err != nil` checks.
- **Test Naming:** Name tests descriptively using `Test_TypeName_MethodName` or `Test_FunctionName` for package-level functions.
- **Subtests:** Use `t.Run(tt.name, ...)` inside table-driven loops for better organization and output.

## 2. Mocking Strategy & Libraries

| Dependency | Mock Library | Notes |
|---|---|---|
| Domain interfaces (repos) | `github.com/stretchr/testify/mock` | Generate mock structs that implement `domain.XxxRepository` |
| SQL Database | `github.com/DATA-DOG/go-sqlmock` | Mock `database/sql` connections; define `ExpectQuery`, `ExpectExec` |
| Redis | `github.com/go-redis/redismock/v9` | Mock Redis commands without a real server |
| Fiber handlers | `fiber.App.Test(req)` + `net/http/httptest` | Create `httptest.NewRequest`, pass to Fiber app, assert status and body |
| Cache (`pkg/cache`) | Pass `nil` or use `github.com/go-redis/redismock/v9` | Usecases guard cache calls with `if uc.cache != nil` |
| Logger (`pkg/logger`) | Use a real `*logger.Logger` with `"error"` level | Logger has no side effects worth mocking; use a real instance |
| Kafka producer | Pass `nil` for `*kafka.Producer` | Usecases guard Kafka calls with `if uc.producer != nil` |

## 3. Test File Organization

- Test files live **next to** the file they test: `foo.go` → `foo_test.go`.
- Test files use the **same package** as the file under test (white-box testing).
- Mock implementations of domain interfaces live in `internal/mocks/repository_mock.go`.
- Common test helpers go as unexported helpers at the top of the test file.

## 4. Layer-Specific Testing Rules

### Repository Tests (`internal/repository/*_test.go`)

- Use `go-sqlmock` to mock the `*sqlx.DB` inside `database.DB`.
- Create a `database.DB` wrapper around the mocked `*sql.DB` via `sqlx.NewDb(mockDB, "postgres")`.
- Test each repository method: happy path, not-found, SQL error.
- Verify all expectations are met with `mock.ExpectationsWereMet()`.
- DO NOT connect to a real database.

### Usecase Tests (`internal/usecase/*_test.go`)

- Mock the `domain.XxxRepository` interface using `testify/mock`.
- Pass `nil` for `*cache.Cache` — usecases guard cache access with nil checks.
- Pass `nil` for `*kafka.Producer` — usecases guard Kafka access with nil checks.
- Use a real `*logger.Logger` (level `"error"` to suppress noise).
- Verify that the correct repository methods are called with expected arguments.
- Test business rule validation.

### Handler Tests (`internal/handler/*_test.go`)

- Create a `fiber.New()` app, register the single route under test.
- Inject a mocked usecase (via `testify/mock`).
- Build requests with `httptest.NewRequest`, set headers (`Content-Type`, `Authorization`).
- Call `app.Test(req, -1)` (disable timeout).
- Assert HTTP status code and JSON response body.
- Verify `mockUsecase.AssertExpectations(t)`.

### Middleware Tests (`internal/middleware/*_test.go`)

- Test each middleware (CORS, Recovery, RequestLogger, Auth) independently.
- For auth middleware: test with valid API key, invalid API key, and missing API key.
- `RequestLogger` must handle `nil` logger gracefully.

### Domain Tests (`internal/domain/*_test.go`)

- Test pure helper methods if any exist.
- No mocks needed — domain has no external dependencies.

## 5. Code Examples

### Table-Driven Test Pattern
```go
func Test_FunctionName(t *testing.T) {
    tests := []struct {
        name     string
        input    InputType
        expected OutputType
        wantErr  bool
    }{
        {
            name:     "happy path",
            input:    validInput,
            expected: expectedOutput,
            wantErr:  false,
        },
        {
            name:     "error case",
            input:    invalidInput,
            expected: zeroValue,
            wantErr:  true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := FunctionName(tt.input)
            if tt.wantErr {
                require.Error(t, err)
                return
            }
            require.NoError(t, err)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### Mock Repository Pattern (testify/mock)
```go
type MockBotRepository struct {
    mock.Mock
}

func (m *MockBotRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Bot, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.Bot), args.Error(1)
}
```

### Fiber Handler Test Pattern
```go
func Test_BotHandler_GetProfile(t *testing.T) {
    mockUsecase := new(MockBotUsecase)
    handler := NewBotHandler(mockUsecase)

    app := fiber.New()
    app.Get("/v1/me", handler.GetProfile)

    req := httptest.NewRequest(http.MethodGet, "/v1/me", nil)
    req.Header.Set("Content-Type", "application/json")

    resp, err := app.Test(req, -1)
    require.NoError(t, err)
    assert.Equal(t, fiber.StatusOK, resp.StatusCode)
    mockUsecase.AssertExpectations(t)
}
```

## 6. Running Tests

```bash
# Run all tests
make test

# Run tests for a specific package
go test ./internal/usecase/... -v -cover

# Run a specific test
go test ./internal/usecase/... -run Test_BotUsecase_GetProfile -v

# Run with coverage report
make test-coverage
```
