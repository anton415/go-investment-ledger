package ledger

import (
	"errors"
	"slices"
	"testing"
)

func TestOperationsReturnsCopy(t *testing.T) {
	portfolio, errPortfolio := NewPortfolio("portfolio-001", "Основной")
	if errPortfolio != nil {
		t.Fatalf("NewPortfolio() error = %v", errPortfolio)
	}

	err := portfolio.AddOperation(Operation{
		ID:           "operation-001",
		InstrumentID: "instrument-001",
		Quantity:     10,
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

func TestRestorePortfolioCopiesOperations(t *testing.T) {
	operations := []Operation{
		{
			ID:           "operation-001",
			InstrumentID: "instrument-001",
			Quantity:     10,
		},
	}

	portfolio, err := RestorePortfolio(
		"portfolio-001",
		"Основной",
		operations,
	)
	if err != nil {
		t.Fatalf("RestorePortfolio() error = %v", err)
	}

	operations[0].Quantity = 999

	got := portfolio.Operations()[0].Quantity
	const want = 10

	if got != want {
		t.Errorf("Quantity = %d, want %d", got, want)
	}
}

func TestRestorePortfolioRejectsInvalidHistory(t *testing.T) {
	operations := []Operation{
		{
			ID:           "operation-001",
			InstrumentID: "instrument-001",
			Quantity:     10,
		},
		{
			ID:           "operation-002",
			InstrumentID: "instrument-001",
			Quantity:     -15,
		},
	}

	portfolio, err := RestorePortfolio(
		"portfolio-001",
		"Основной",
		operations,
	)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var historyErr *HistoryIntegrityError
	if !errors.As(err, &historyErr) {
		t.Fatalf("error = %v, want *HistoryIntegrityError", err)
	}

	if historyErr.Index != 1 {
		t.Errorf("Index = %d, want 1", historyErr.Index)
	}
	if historyErr.OperationID != "operation-002" {
		t.Errorf(
			"OperationID = %q, want %q",
			historyErr.OperationID,
			"operation-002",
		)
	}

	if !errors.Is(err, ErrInsufficientPosition) {
		t.Errorf("error = %v, want ErrInsufficientPosition", err)
	}

	if portfolio.ID != "" {
		t.Errorf("ID = %q, want empty", portfolio.ID)
	}
	if portfolio.Name != "" {
		t.Errorf("Name = %q, want empty", portfolio.Name)
	}
	if got := len(portfolio.Operations()); got != 0 {
		t.Errorf("len(Operations()) = %d, want 0", got)
	}
}

func TestPositionsReturnsIndependentMap(t *testing.T) {
	portfolio, errPortfolio := NewPortfolio("portfolio-001", "Основной")

	if errPortfolio != nil {
		t.Fatalf("NewPortfolio() error = %v", errPortfolio)
	}

	err := portfolio.AddOperation(Operation{
		ID:           "operation-001",
		InstrumentID: "instrument-001",
		Quantity:     10,
	})
	if err != nil {
		t.Fatalf("AddOperation() error = %v", err)
	}

	positions := portfolio.Positions()
	positions["instrument-001"] = 999

	got := portfolio.Positions()["instrument-001"]
	const want = 10

	if got != want {
		t.Errorf("Quantity = %d, want %d", got, want)
	}
}

func TestAddOperationRejectsOversell(t *testing.T) {
	portfolio, errPortfolio := NewPortfolio("portfolio-001", "Основной")

	if errPortfolio != nil {
		t.Fatalf("NewPortfolio() error = %v", errPortfolio)
	}

	err := portfolio.AddOperation(Operation{
		ID:           "operation-001",
		InstrumentID: "instrument-001",
		Quantity:     10,
	})
	if err != nil {
		t.Fatalf("AddOperation() error = %v", err)
	}

	before := portfolio.Operations()

	err = portfolio.AddOperation(Operation{
		ID:           "operation-002",
		InstrumentID: "instrument-001",
		Quantity:     -15,
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
	portfolio, errPortfolio := NewPortfolio("portfolio-001", "Основной")

	if errPortfolio != nil {
		t.Fatalf("NewPortfolio() error = %v", errPortfolio)
	}

	err := portfolio.AddOperation(Operation{
		ID:           "operation-001",
		InstrumentID: "instrument-001",
		Quantity:     10,
	})
	if err != nil {
		t.Fatalf("AddOperation() error = %v", err)
	}

	before := portfolio.Operations()

	err = portfolio.AddOperation(Operation{
		ID:           "operation-001",
		InstrumentID: "instrument-001",
		Quantity:     5,
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

func TestAddOperationRejectsInvalidFields(t *testing.T) {
	tests := []struct {
		name       string
		operation  Operation
		wantField  string
		wantReason string
	}{
		{
			name: "empty ID",
			operation: Operation{
				InstrumentID: "instrument-001",
				Quantity:     10,
			},
			wantField:  "id",
			wantReason: "must not be empty",
		},
		{
			name: "empty instrument ID",
			operation: Operation{
				ID:       "operation-001",
				Quantity: 10,
			},
			wantField:  "instrument_id",
			wantReason: "must not be empty",
		},
		{
			name: "zero quantity",
			operation: Operation{
				ID:           "operation-001",
				InstrumentID: "instrument-001",
			},
			wantField:  "quantity",
			wantReason: "must not be zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			portfolio, errPortfolio := NewPortfolio("portfolio-001", "Основной")

			if errPortfolio != nil {
				t.Fatalf("NewPortfolio() error = %v", errPortfolio)
			}
			before := portfolio.Operations()

			err := portfolio.AddOperation(tt.operation)
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			var validationErr *ValidationError
			if !errors.As(err, &validationErr) {
				t.Fatalf("error = %v, want *ValidationError", err)
			}

			if validationErr.Field != tt.wantField {
				t.Errorf("Field = %q, want %q", validationErr.Field, tt.wantField)
			}
			if validationErr.Reason != tt.wantReason {
				t.Errorf("Reason = %q, want %q", validationErr.Reason, tt.wantReason)
			}

			after := portfolio.Operations()
			if !slices.Equal(before, after) {
				t.Errorf("operations before = %v, after = %v", before, after)
			}
		})
	}
}

func TestAddOperationValidatesFieldsBeforeState(t *testing.T) {
	portfolio, errPortfolio := NewPortfolio("portfolio-001", "Основной")

	if errPortfolio != nil {
		t.Fatalf("NewPortfolio() error = %v", errPortfolio)
	}

	err := portfolio.AddOperation(Operation{
		ID:           "operation-001",
		InstrumentID: "instrument-001",
		Quantity:     10,
	})

	if err != nil {
		t.Fatalf("AddOperation() error = %v", err)
	}

	before := portfolio.Operations()

	err = portfolio.AddOperation(Operation{
		ID:           "operation-001",
		InstrumentID: "",
		Quantity:     5,
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *ValidationError", err)
	}

	if validationErr.Field != "instrument_id" {
		t.Errorf("Field = %q, want %q", validationErr.Field, "instrument_id")
	}
	if validationErr.Reason != "must not be empty" {
		t.Errorf("Reason = %q, want %q", validationErr.Reason, "must not be empty")
	}

	if errors.Is(err, ErrOperationAlreadyAdded) {
		t.Fatalf("error = %v, must not match ErrOperationAlreadyAdded", err)
	}

	after := portfolio.Operations()
	if !slices.Equal(before, after) {
		t.Errorf("operations before = %v, after = %v", before, after)
	}
}

func TestNewPortfolioRejectsEmptyID(t *testing.T) {
	portfolio, err := NewPortfolio("", "Основной")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *ValidationError", err)
	}

	if validationErr.Field != "id" {
		t.Errorf("Field = %q, want %q", validationErr.Field, "id")
	}
	if validationErr.Reason != "must not be empty" {
		t.Errorf("Reason = %q, want %q", validationErr.Reason, "must not be empty")
	}

	if portfolio.ID != "" {
		t.Errorf("ID = %q, want empty", portfolio.ID)
	}
	if portfolio.Name != "" {
		t.Errorf("Name = %q, want empty", portfolio.Name)
	}
	if got := len(portfolio.Operations()); got != 0 {
		t.Errorf("len(Operations()) = %d, want 0", got)
	}
}
