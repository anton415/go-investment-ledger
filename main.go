package main

import (
	"fmt"

	"github.com/anton415/go-investment-ledger/internal/ledger"
)

type operationRecorder interface {
	AddOperation(ledger.Operation) error
}

const demoPortfolioID ledger.PortfolioID = "portfolio-001"

func main() {
	portfolio := ledger.NewPortfolio(demoPortfolioID, "Основной портфель", nil)
	var recorder operationRecorder = &portfolio

	operations := []ledger.Operation{
		{ID: "operation-001", Ticker: "SBER", Quantity: 10},
		{ID: "operation-002", Ticker: "YNDX", Quantity: 3},
		{ID: "operation-003", Ticker: "SBER", Quantity: -4},
		{ID: "operation-003", Ticker: "SBER", Quantity: -4},
		{ID: "operation-004", Ticker: "YNDX", Quantity: -5},
	}

	for _, operation := range operations {
		err := recorder.AddOperation(operation)
		switch err {
		case nil:
			fmt.Printf("%s: %+d %s — добавлена\n", operation.ID, operation.Quantity, operation.Ticker)
		case ledger.ErrInsufficientPosition:
			fmt.Printf("%s — %v\n", operation.ID, err)
		case ledger.ErrOperationAlreadyAdded:
			fmt.Printf("%s — %v\n", operation.ID, err)
		default:
			fmt.Printf("%s — неизвестная ошибка: %v\n", operation.ID, err)
		}
	}

	positions := portfolio.Positions()

	findPosition := func(ticker ledger.Ticker) (int, bool) {
		quantity, found := positions[ticker]
		return quantity, found
	}

	trackedTickers := [...]ledger.Ticker{"SBER", "YNDX", "MOEX"}

	for _, ticker := range trackedTickers {
		quantity, found := findPosition(ticker)
		fmt.Printf("position=%s, quantity=%d, found=%t\n", ticker, quantity, found)
	}

	portfolios := map[ledger.PortfolioID]ledger.Portfolio{portfolio.ID: portfolio}
	savedPortfolio, found := portfolios[demoPortfolioID]
	fmt.Printf("portfolio=%q, found=%t\n", savedPortfolio.Name, found)

	missingPortfolio, found := portfolios[ledger.PortfolioID("portfolio-999")]
	fmt.Printf("portfolio=%q, found=%t\n", missingPortfolio.Name, found)

	rawID := string(portfolio.ID)
	fmt.Println(rawID)

	summary, positionCount := portfolio.Summary()
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

	initial, hasInitial := firstRune(portfolio.Name)
	fmt.Printf("portfolio_initial=%q, found=%t\n", initial, hasInitial)
}

func firstRune(value string) (rune, bool) {
	for _, symbol := range value {
		return symbol, true
	}
	return 0, false
}
