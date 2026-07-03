package banking

import (
	"errors"
	"strings"
	"testing"
)

func TestWithdrawSuccess(t *testing.T) {
	a := &Account{Balance: 100}
	if err := a.Withdraw(30); err != nil {
		t.Fatalf("Withdraw(30) unexpected error: %v", err)
	}
	if a.Balance != 70 {
		t.Errorf("Balance = %d, want 70", a.Balance)
	}
}

func TestWithdrawInsufficient(t *testing.T) {
	a := &Account{Balance: 50}
	err := a.Withdraw(80)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrInsufficientFunds) {
		t.Errorf("errors.Is(err, ErrInsufficientFunds) = false; err = %v", err)
	}
	if !IsInsufficientFunds(err) {
		t.Errorf("IsInsufficientFunds = false; err = %v", err)
	}
	// The wrapped error should still add context (the amount).
	if !strings.Contains(err.Error(), "80") {
		t.Errorf("wrapped error %q should mention the amount 80", err.Error())
	}
	if a.Balance != 50 {
		t.Errorf("Balance changed on failed withdraw: %d", a.Balance)
	}
}

func TestWithdrawValidation(t *testing.T) {
	a := &Account{Balance: 50}
	err := a.Withdraw(-5)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("errors.As did not find *ValidationError; err = %v", err)
	}
	if ve.Field != "amount" {
		t.Errorf("ValidationError.Field = %q, want %q", ve.Field, "amount")
	}
	if got := ve.Error(); got != "invalid amount: must be positive" {
		t.Errorf("Error() = %q", got)
	}
	field, ok := FieldInError(err)
	if !ok || field != "amount" {
		t.Errorf("FieldInError = (%q, %v), want (amount, true)", field, ok)
	}
}

func TestFieldInErrorAbsent(t *testing.T) {
	_, ok := FieldInError(errors.New("some other error"))
	if ok {
		t.Error("FieldInError should be false for a non-validation error")
	}
}
