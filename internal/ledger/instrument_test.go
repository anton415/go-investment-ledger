package ledger

import (
	"errors"
	"testing"
)

func TestInstrumentCreation(t *testing.T) {
	instrument, err := NewInstrument("instrument-001", "secid-001")
	if err != nil {
		t.Fatalf("NewInstrument() error = %v", err)
	}

	if instrument.ID != "instrument-001" {
		t.Errorf("Instrument ID = %s, want %s", instrument.ID, "instrument-001")
	}
	if instrument.SECID != "secid-001" {
		t.Errorf("Instrument SECID = %s, want %s", instrument.SECID, "secid-001")
	}
}

func TestInstrumentCreationRejectsInvalidFields(t *testing.T) {
	tests := []struct {
		name       string
		instrument Instrument
		wantField  string
		wantReason string
	}{
		{
			name: "empty SECID",
			instrument: Instrument{
				ID:    "instrument-001",
				SECID: "",
			},
			wantField:  "secid",
			wantReason: "must not be empty",
		},
		{
			name: "empty ID",
			instrument: Instrument{
				ID:    "",
				SECID: "secid-001",
			},
			wantField:  "id",
			wantReason: "must not be empty",
		},
		{
			name:       "empty ID and SECID",
			instrument: Instrument{},
			wantField:  "id",
			wantReason: "must not be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instrument, err := NewInstrument(tt.instrument.ID, tt.instrument.SECID)

			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if instrument != (Instrument{}) {
				t.Errorf("Instrument = %+v, want zero value", instrument)
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
		})
	}
}
