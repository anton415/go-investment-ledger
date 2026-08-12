package ledger

import (
	"errors"
	"slices"
)

type PortfolioID string
type OperationID string
type Ticker string

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
	operations []Operation
}

var (
	ErrInsufficientPosition  = errors.New("недостаточное количество инструмента в позиции")
	ErrOperationAlreadyAdded = errors.New("операция уже добавлена")
)

func NewPortfolio(id PortfolioID, name string, operations []Operation) Portfolio {
	return Portfolio{
		ID:         id,
		Name:       name,
		operations: slices.Clone(operations),
	}
}

// AddOperation добавляет операцию, только если её ID уникален,
// а итоговая позиция по инструменту не становится отрицательной.
func (p *Portfolio) AddOperation(operation Operation) error {
	if containsOperation(p.operations, operation.ID) {
		return ErrOperationAlreadyAdded
	}
	if p.position(operation.Ticker)+operation.Quantity < 0 {
		return ErrInsufficientPosition
	}
	p.operations = append(p.operations, operation)
	return nil
}

func (p Portfolio) Summary() (string, int) {
	return string(p.ID) + ": " + p.Name, len(p.Positions())
}

// Positions пересчитывает открытые позиции из журнала операций,
// который остаётся единственным источником истины.
func (p Portfolio) Positions() map[Ticker]int {
	positions := make(map[Ticker]int)
	for _, operation := range p.operations {
		positions[operation.Ticker] += operation.Quantity
		if positions[operation.Ticker] == 0 {
			delete(positions, operation.Ticker)
		}
	}
	return positions
}

// Operations возвращает копию среза операций, чтобы предотвратить
// изменение внутреннего состояния портфеля извне.
func (p Portfolio) Operations() []Operation {
	return slices.Clone(p.operations)
}

func (p Portfolio) position(ticker Ticker) int {
	return p.Positions()[ticker]
}

func containsOperation(operations []Operation, id OperationID) bool {
	for _, operation := range operations {
		if operation.ID == id {
			return true
		}
	}
	return false
}
