package ledger

import (
	"errors"
	"fmt"
	"slices"
)

type PortfolioID string
type OperationID string

type Operation struct {
	ID           OperationID
	InstrumentID InstrumentID
	// Quantity задаёт изменение позиции: положительное значение увеличивает,
	// отрицательное — уменьшает количество инструмента.
	Quantity int
}

type Portfolio struct {
	ID         PortfolioID
	Name       string
	operations []Operation
}

type ValidationError struct {
	Field  string
	Reason string
}

type HistoryIntegrityError struct {
	Index       int
	OperationID OperationID
	Err         error
}

func (e *HistoryIntegrityError) Error() string {
	return fmt.Sprintf(
		"invalid history operation[%d] %q: %v",
		e.Index,
		e.OperationID,
		e.Err,
	)
}

func (e *HistoryIntegrityError) Unwrap() error {
	return e.Err
}

func (e ValidationError) Error() string {
	return fmt.Sprintf(
		"field %q, reason %q",
		e.Field,
		e.Reason,
	)
}

var (
	ErrInsufficientPosition  = errors.New("недостаточное количество инструмента в позиции")
	ErrOperationAlreadyAdded = errors.New("операция уже добавлена")
)

func NewPortfolio(id PortfolioID, name string) (Portfolio, error) {
	if id == "" {
		return Portfolio{}, fmt.Errorf(
			"create portfolio: %w",
			&ValidationError{
				Field:  "id",
				Reason: "must not be empty",
			},
		)
	}
	return Portfolio{
		ID:         id,
		Name:       name,
		operations: []Operation{},
	}, nil
}

func RestorePortfolio(
	id PortfolioID,
	name string,
	operations []Operation,
) (Portfolio, error) {
	p, err := NewPortfolio(id, name)
	if err != nil {
		return Portfolio{}, fmt.Errorf(
			"restore portfolio: %w",
			err,
		)
	}

	for i, operation := range operations {
		if err := p.AddOperation(operation); err != nil {
			return Portfolio{}, fmt.Errorf(
				"restore portfolio %q: %w",
				id,
				&HistoryIntegrityError{
					Index:       i,
					OperationID: operation.ID,
					Err:         err,
				},
			)
		}
	}

	return p, nil
}

// AddOperation добавляет операцию, только если её ID уникален,
// а итоговая позиция по инструменту не становится отрицательной.
func (p *Portfolio) AddOperation(operation Operation) error {
	if operation.ID == "" {
		return fmt.Errorf(
			"add operation to portfolio %q: %w",
			p.ID,
			&ValidationError{
				Field:  "id",
				Reason: "must not be empty",
			},
		)
	}

	if operation.InstrumentID == "" {
		return fmt.Errorf(
			"add operation %q to portfolio %q: %w",
			operation.ID,
			p.ID,
			&ValidationError{
				Field:  "instrument_id",
				Reason: "must not be empty",
			},
		)
	}

	if operation.Quantity == 0 {
		return fmt.Errorf(
			"add operation %q to portfolio %q: %w",
			operation.ID,
			p.ID,
			&ValidationError{
				Field:  "quantity",
				Reason: "must not be zero",
			},
		)
	}

	if containsOperation(p.operations, operation.ID) {
		return fmt.Errorf(
			"add operation %q to portfolio %q: %w",
			operation.ID,
			p.ID,
			ErrOperationAlreadyAdded,
		)
	}

	if p.position(operation.InstrumentID)+operation.Quantity < 0 {
		return fmt.Errorf(
			"add operation %q to portfolio %q: %w",
			operation.ID,
			p.ID,
			ErrInsufficientPosition,
		)
	}

	p.operations = append(p.operations, operation)
	return nil
}

func (p Portfolio) Summary() (string, int) {
	return string(p.ID) + ": " + p.Name, len(p.Positions())
}

// Positions пересчитывает открытые позиции из журнала операций,
// который остаётся единственным источником истины.
func (p Portfolio) Positions() map[InstrumentID]int {
	positions := make(map[InstrumentID]int)
	for _, operation := range p.operations {
		positions[operation.InstrumentID] += operation.Quantity
		if positions[operation.InstrumentID] == 0 {
			delete(positions, operation.InstrumentID)
		}
	}
	return positions
}

// Operations возвращает копию среза операций, чтобы предотвратить
// изменение внутреннего состояния портфеля извне.
func (p Portfolio) Operations() []Operation {
	return slices.Clone(p.operations)
}

func (p Portfolio) position(instrumentID InstrumentID) int {
	return p.Positions()[instrumentID]
}

func containsOperation(operations []Operation, id OperationID) bool {
	for _, operation := range operations {
		if operation.ID == id {
			return true
		}
	}
	return false
}
