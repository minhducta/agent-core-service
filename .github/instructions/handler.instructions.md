---
applyTo: "internal/handler/**"
---

# Handler Layer Instructions

Fiber HTTP handlers — parse requests, call usecase, return JSON.

## Rules

1. **Never import `repository/` directly** — only `usecase/` interfaces.
2. **Parse and validate input** before calling usecase.
3. **Map domain errors to HTTP status codes** in `helpers.go`.
4. **Extract `botID` from `fiber.Locals("botID")`** — set by auth middleware.
5. **Response format**: `{"data": ...}` for success, `{"error": {"code": ..., "message": ...}}` for errors.

## Route Registration

All routes registered in `router.go`. Group protected routes under auth middleware.

## Testing

- Create `fiber.New()` app, register route, inject mocked usecase.
- Test via `app.Test(httptest.NewRequest(...))`.
- Test HTTP status codes + response body.
