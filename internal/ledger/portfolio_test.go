package ledger

import (
	"errors"
	"slices"
	"testing"
)

func TestOperationsReturnsCopy(t *testing.T) {
	portfolio := NewPortfolio("portfolio-001", "Основной", nil)

	err := portfolio.AddOperation(Operation{
		ID:       "operation-001",
		Ticker:   "SBER",
		Quantity: 10,
	})
	if err != nil {
		t.Fatalf("AddOperation() error = %v", err)
	}

	operations := portfolio.Operations()
	operations[0].Quantity = 999

	got := portfolio.Operations()[0].Quantity
	const want = 10

	if got != want {
		t.Errorf("Quantity = %d, want %d", got, want)
	}
}

func TestNewPortfolioCopiesOperations(t *testing.T) {
	operations := []Operation{
		{
			ID:       "operation-001",
			Ticker:   "SBER",
			Quantity: 10,
		},
	}

	portfolio := NewPortfolio("portfolio-001", "Основной", operations)
	operations[0].Quantity = 999

	got := portfolio.Operations()[0].Quantity
	const want = 10

	if got != want {
		t.Errorf("Quantity = %d, want %d", got, want)
	}
}

func TestPositionsReturnsIndependentMap(t *testing.T) {
	portfolio := NewPortfolio("portfolio-001", "Основной", nil)

	err := portfolio.AddOperation(Operation{
		ID:       "operation-001",
		Ticker:   "SBER",
		Quantity: 10,
	})
	if err != nil {
		t.Fatalf("AddOperation() error = %v", err)
	}

	positions := portfolio.Positions()
	positions["SBER"] = 999

	got := portfolio.Positions()["SBER"]
	const want = 10

	if got != want {
		t.Errorf("Quantity = %d, want %d", got, want)
	}
}

func TestAddOperationRejectsOversell(t *testing.T) {
	portfolio := NewPortfolio("portfolio-001", "Основной", nil)

	err := portfolio.AddOperation(Operation{
		ID:       "operation-001",
		Ticker:   "SBER",
		Quantity: 10,
	})
	if err != nil {
		t.Fatalf("AddOperation() error = %v", err)
	}

	before := portfolio.Operations()

	err = portfolio.AddOperation(Operation{
		ID:       "operation-002",
		Ticker:   "SBER",
		Quantity: -15,
	})

	if err == nil {
		t.Errorf("expected error, got nil")
	}

	if !errors.Is(err, ErrInsufficientPosition) {
		t.Fatalf("error = %v, want ErrInsufficientPosition", err)
	}

	after := portfolio.Operations()

	if !slices.Equal(before, after) {
		t.Errorf("before %v, after %v", before, after)
	}
}

func TestTryAddDuplicateOperation(t *testing.T) {
	portfolio := NewPortfolio("portfolio-001", "Основной", nil)

	err := portfolio.AddOperation(Operation{
		ID:       "operation-001",
		Ticker:   "SBER",
		Quantity: 10,
	})
	if err != nil {
		t.Fatalf("AddOperation() error = %v", err)
	}

	before := portfolio.Operations()

	err = portfolio.AddOperation(Operation{
		ID:       "operation-001",
		Ticker:   "SBER",
		Quantity: 5,
	})

	if err == nil {
		t.Errorf("expected error, got nil")
	}

	if !errors.Is(err, ErrOperationAlreadyAdded) {
		t.Fatalf("error = %v, want ErrOperationAlreadyAdded", err)
	}

	after := portfolio.Operations()

	if !slices.Equal(before, after) {
		t.Errorf("before %v, after %v", before, after)
	}
}

func TestAddOperationRejectsZeroQuantity(t *testing.T) {
	portfolio := NewPortfolio("portfolio-001", "Основной", nil)

	before := portfolio.Operations()

	err := portfolio.AddOperation(Operation{
		ID:       "operation-001",
		Ticker:   "SBER",
		Quantity: 0,
	})

	if err == nil {
		t.Errorf("expected error, got nil")
	}

	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *ValidationError", err)
	}

	if validationErr.Reason != "must not be zero" {
		t.Errorf("error = %v, want 'must not be zero'", validationErr.Reason)
	}

	if validationErr.Field != "quantity" {
		t.Errorf("error = %v, want 'quantity'", validationErr.Field)
	}

	after := portfolio.Operations()

	if !slices.Equal(before, after) {
		t.Errorf("before %v, after %v", before, after)
	}
}
