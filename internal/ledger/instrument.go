package ledger

import "fmt"

type InstrumentID string
type SECID string

type Instrument struct {
	ID    InstrumentID
	SECID SECID
}

func NewInstrument(id InstrumentID, secid SECID) (Instrument, error) {
	if id == "" {
		return Instrument{}, fmt.Errorf(
			"create instrument: %w",
			&ValidationError{
				Field:  "id",
				Reason: "must not be empty",
			},
		)
	}

	if secid == "" {
		return Instrument{}, fmt.Errorf(
			"create instrument: %w",
			&ValidationError{
				Field:  "secid",
				Reason: "must not be empty",
			},
		)
	}

	return Instrument{
		ID:    id,
		SECID: secid,
	}, nil
}
