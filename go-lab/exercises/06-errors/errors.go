// Package banking drills Go's error model: sentinel errors, custom error
// types, wrapping with %w, and errors.Is / errors.As.
package banking

import (
	"errors"
	"fmt"
)

// ErrInsufficientFunds is a sentinel error callers can compare against with
// errors.Is — even after it has been wrapped with additional context.
var ErrInsufficientFunds = errors.New("banking: insufficient funds")

// ValidationError is a custom error type carrying structured detail about which
// field failed and why. Callers can recover it with errors.As.
type ValidationError struct {
	Field  string
	Reason string
}

// Error makes ValidationError satisfy the error interface. Format it as
// "invalid <Field>: <Reason>", e.g. "invalid amount: must be positive".
func (e *ValidationError) Error() string {
	// panic("TODO: implement ValidationError.Error")
	return fmt.Sprintf("invalid %s: %s", e.Field, e.Reason)
}

// Account holds a balance in whole cents.
type Account struct {
	Balance int
}

// Withdraw subtracts amount from the balance. Rules:
//
//   - If amount <= 0, return a *ValidationError with Field "amount" and Reason
//     "must be positive". (Return it as an error.)
//   - If amount > Balance, return an error that WRAPS ErrInsufficientFunds with
//     context, e.g. fmt.Errorf("withdraw %d: %w", amount, ErrInsufficientFunds).
//   - Otherwise subtract amount from Balance and return nil.
func (a *Account) Withdraw(amount int) error {
	// panic("TODO: implement Account.Withdraw")
	if amount <= 0 {
		return &ValidationError{
			Field:  "amount",
			Reason: "must be positive",
		}
	} else if amount > a.Balance {
		return fmt.Errorf("withdraw %d: %w", amount, ErrInsufficientFunds)
	} else {
		a.Balance -= amount
		return nil
	}
}

// IsInsufficientFunds reports whether err (or anything it wraps) is
// ErrInsufficientFunds. Use errors.Is.
func IsInsufficientFunds(err error) bool {
	// panic("TODO: implement IsInsufficientFunds")
	return errors.Is(err, ErrInsufficientFunds)
}

// FieldInError returns the Field of a *ValidationError found anywhere in err's
// chain, plus true. If there is no ValidationError in the chain, it returns
// "", false. Use errors.As.
func FieldInError(err error) (string, bool) {
	// panic("TODO: implement FieldInError")
	var ve *ValidationError
	if errors.As(err, &ve) {
		return ve.Field, true
	}
	return "", false
}
