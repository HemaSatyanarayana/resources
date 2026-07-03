# 06 — Errors

Go treats errors as **values**, not exceptions. You return them, inspect them, wrap them, and unwrap them. Master the modern (`errors.Is`/`As`/`%w`) toolkit and you'll write robust, debuggable Go.

## Concepts

- **`error` is just an interface:** `interface { Error() string }`.
- **Sentinel errors** are exported `var Err... = errors.New(...)` values. Callers compare with **`errors.Is(err, ErrX)`** — never `err == ErrX` once wrapping is involved.
- **Wrapping** with `fmt.Errorf("context: %w", err)` adds context while preserving the original for `errors.Is`/`As`.
- **Custom error types** carry structured data. Recover them with **`errors.As(err, &target)`**.
- **Convention:** the last return value is the error; return `nil` on success; don't wrap with `%w` unless callers need to unwrap.

## Your task

Implement everything in [`errors.go`](errors.go). You'll need `import "fmt"` for the wrapping.

| Function | Skill |
|----------|-------|
| `ValidationError.Error` | Satisfy the `error` interface |
| `Account.Withdraw` | Return custom + wrapped sentinel errors |
| `IsInsufficientFunds` | `errors.Is` through a wrap chain |
| `FieldInError` | `errors.As` to extract a typed error |

## Run

```bash
go test -v ./exercises/06-errors/
```

## Hints

- Return the validation error as `&ValidationError{Field: "amount", Reason: "must be positive"}`. Pointer receiver on `Error()` means the pointer is the `error`.
- `errors.As` needs a pointer to the target: `var ve *ValidationError; errors.As(err, &ve)`.
- Don't mutate `Balance` before you've confirmed the withdrawal is valid.

<details>
<summary>Reference solution</summary>

```go
package banking

import (
	"errors"
	"fmt"
)

func (e *ValidationError) Error() string {
	return fmt.Sprintf("invalid %s: %s", e.Field, e.Reason)
}

func (a *Account) Withdraw(amount int) error {
	if amount <= 0 {
		return &ValidationError{Field: "amount", Reason: "must be positive"}
	}
	if amount > a.Balance {
		return fmt.Errorf("withdraw %d: %w", amount, ErrInsufficientFunds)
	}
	a.Balance -= amount
	return nil
}

func IsInsufficientFunds(err error) bool {
	return errors.Is(err, ErrInsufficientFunds)
}

func FieldInError(err error) (string, bool) {
	var ve *ValidationError
	if errors.As(err, &ve) {
		return ve.Field, true
	}
	return "", false
}
```

</details>
