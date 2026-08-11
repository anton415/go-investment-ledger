package main

import (
	"errors"
	"fmt"
)

type PortfolioID string
type OperationID string
type Ticker string

const demoPortfolioID PortfolioID = "portfolio-001"

type Operation struct {
	ID     OperationID
	Ticker Ticker
	// Quantity задаёт изменение позиции: положительное значение увеличивает,
	// отрицательное — уменьшает количество инструмента.
	Quantity int
}

type Portfolio struct {
	ID         PortfolioID
	Name       string
	Operations []Operation
}

var (
	ErrInsufficientPosition  = errors.New("недостаточное количество инструмента в позиции")
	ErrOperationAlreadyAdded = errors.New("операция уже добавлена")
)

func main() {
	portfolio := createPortfolio(demoPortfolioID, "Основной портфель", nil)

	operations := []Operation{
		{ID: "operation-001", Ticker: "SBER", Quantity: 10},
		{ID: "operation-002", Ticker: "YNDX", Quantity: 3},
		{ID: "operation-003", Ticker: "SBER", Quantity: -4},
		{ID: "operation-003", Ticker: "SBER", Quantity: -4},
		{ID: "operation-004", Ticker: "YNDX", Quantity: -5},
	}

	for _, operation := range operations {
		err := portfolio.addOperation(operation)
		switch err {
		case nil:
			fmt.Printf("%s: %+d %s — добавлена\n", operation.ID, operation.Quantity, operation.Ticker)
		case ErrInsufficientPosition:
			fmt.Printf("%s — %v\n", operation.ID, err)
		case ErrOperationAlreadyAdded:
			fmt.Printf("%s — %v\n", operation.ID, err)
		default:
			fmt.Printf("%s — неизвестная ошибка: %v\n", operation.ID, err)
		}
	}

	positions := portfolio.buildPositions()
	for _, ticker := range []Ticker{"SBER", "YNDX", "MOEX"} {
		quantity, found := positions[ticker]
		fmt.Printf("position=%s, quantity=%d, found=%t\n", ticker, quantity, found)
	}

	portfolios := map[PortfolioID]Portfolio{portfolio.ID: portfolio}
	savedPortfolio, found := portfolios[demoPortfolioID]
	fmt.Printf("portfolio=%q, found=%t\n", savedPortfolio.Name, found)

	missingPortfolio, found := portfolios[PortfolioID("portfolio-999")]
	fmt.Printf("portfolio=%q, found=%t\n", missingPortfolio.Name, found)

	rawID := string(portfolio.ID)
	fmt.Println(rawID)

	summary, positionCount := portfolio.buildPortfolioSummary()
	fmt.Printf("summary=%q, positions=%d\n", summary, positionCount)

	if positionCount > 0 {
		fmt.Println("portfolio has open positions")
	} else {
		fmt.Println("portfolio is empty")
	}

	var diversification string
	switch {
	case positionCount == 0:
		diversification = "empty"
	case positionCount == 1:
		diversification = "single instrument"
	default:
		diversification = "multiple instruments"
	}
	fmt.Println(diversification)
}

func (p Portfolio) buildPortfolioSummary() (string, int) {
	return string(p.ID) + ": " + p.Name, len(p.buildPositions())
}

// buildPositions пересчитывает открытые позиции из журнала операций,
// который остаётся единственным источником истины.
func (p Portfolio) buildPositions() map[Ticker]int {
	positions := make(map[Ticker]int)
	for _, operation := range p.Operations {
		positions[operation.Ticker] += operation.Quantity
		if positions[operation.Ticker] == 0 {
			delete(positions, operation.Ticker)
		}
	}
	return positions
}

func (p Portfolio) position(ticker Ticker) int {
	return p.buildPositions()[ticker]
}

func containsOperation(operations []Operation, id OperationID) bool {
	for _, operation := range operations {
		if operation.ID == id {
			return true
		}
	}
	return false
}

// addOperation добавляет операцию, только если её ID уникален,
// а итоговая позиция по инструменту не становится отрицательной.
func (p *Portfolio) addOperation(operation Operation) error {
	if containsOperation(p.Operations, operation.ID) {
		return ErrOperationAlreadyAdded
	}
	if p.position(operation.Ticker)+operation.Quantity < 0 {
		return ErrInsufficientPosition
	}
	p.Operations = append(p.Operations, operation)
	return nil
}

func createPortfolio(id PortfolioID, name string, operations []Operation) Portfolio {
	return Portfolio{
		ID:         id,
		Name:       name,
		Operations: operations,
	}
}
