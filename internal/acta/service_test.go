package acta

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
)

func TestServiceHandlesRecord(t *testing.T) {
	t.Parallel()

	recorder, _ := testRecorder(t)
	service := NewService(recorder)

	encoded, err := json.Marshal(validRecordParams())
	if err != nil {
		t.Fatalf("marshal record params: %v", err)
	}

	response := service.Handle(context.Background(), Request{
		Protocol:  ProtocolVersion,
		Operation: "record",
		Params:    encoded,
	})
	if !response.OK || response.Error != nil {
		t.Fatalf("Handle() response = %#v", response)
	}

	var result RecordResult
	if err := json.Unmarshal(response.Result, &result); err != nil {
		t.Fatalf("decode record result: %v", err)
	}

	if !result.Recorded || result.ID == "" {
		t.Fatalf("record result = %#v", result)
	}
}

func TestServiceRejectsInvalidRequest(t *testing.T) {
	t.Parallel()

	recorder, _ := testRecorder(t)
	service := NewService(recorder)

	tests := []struct {
		name      string
		operation string
		params    json.RawMessage
		wantCode  ErrorCode
	}{
		{
			name:      "unknown operation",
			operation: "missing",
			params:    json.RawMessage(`{}`),
			wantCode:  ErrorUnknownOperation,
		},
		{
			name:      "malformed params",
			operation: "record",
			params:    json.RawMessage(`{"unknown":true}`),
			wantCode:  ErrorInvalidRequest,
		},
		{
			name:      "invalid event",
			operation: "record",
			params:    json.RawMessage(`{}`),
			wantCode:  ErrorInvalidEvent,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			response := service.Handle(context.Background(), Request{
				Protocol:  ProtocolVersion,
				Operation: test.operation,
				Params:    test.params,
			})

			if response.OK || response.Error == nil || response.Error.Code != test.wantCode {
				t.Fatalf("Handle() response = %#v, want code %q", response, test.wantCode)
			}
		})
	}
}

func TestRecordErrorResponse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		err      error
		wantCode ErrorCode
	}{
		{name: "invalid", err: ErrEventInvalid, wantCode: ErrorInvalidEvent},
		{name: "corrupt", err: ErrJournalCorrupt, wantCode: ErrorJournalCorrupt},
		{name: "unavailable", err: ErrJournalNotLoaded, wantCode: ErrorJournalUnavailable},
		{name: "other", err: errors.New("storage failure"), wantCode: ErrorJournalUnavailable},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			response := recordErrorResponse(test.err)
			if response.Error == nil || response.Error.Code != test.wantCode {
				t.Fatalf("recordErrorResponse() = %#v", response)
			}
		})
	}
}
