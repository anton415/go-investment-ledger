package ledger

import (
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
