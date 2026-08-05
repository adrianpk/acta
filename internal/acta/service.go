package acta

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
)

// Service dispatches validated protocol operations to application workflows.
type Service struct {
	recorder *Recorder
}

// NewService constructs a service without performing IO.
func NewService(recorder *Recorder) *Service {
	return &Service{recorder: recorder}
}

// Handle executes one local service request.
func (service *Service) Handle(ctx context.Context, request Request) Response {
	switch request.Operation {
	case "record":
		return service.handleRecord(ctx, request.Params)
	default:
		return ErrorResponse(ErrorUnknownOperation, "unknown operation")
	}
}

func (service *Service) handleRecord(ctx context.Context, encoded json.RawMessage) Response {
	var params RecordParams

	decoder := json.NewDecoder(bytes.NewReader(encoded))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&params); err != nil {
		return ErrorResponse(ErrorInvalidRequest, "decode record params")
	}

	if err := requireJSONEnd(decoder); err != nil {
		return ErrorResponse(ErrorInvalidRequest, "record params contain trailing JSON")
	}

	result, err := service.recorder.Record(ctx, params)
	if err != nil {
		return recordErrorResponse(err)
	}

	return SuccessResponse(result)
}

func recordErrorResponse(err error) Response {
	switch {
	case errors.Is(err, ErrEventInvalid):
		return ErrorResponse(ErrorInvalidEvent, err.Error())
	case errors.Is(err, ErrJournalCorrupt):
		return ErrorResponse(ErrorJournalCorrupt, err.Error())
	case errors.Is(err, ErrJournalNotLoaded):
		return ErrorResponse(ErrorJournalUnavailable, "journal is unavailable")
	default:
		return ErrorResponse(ErrorJournalUnavailable, "journal operation failed")
	}
}
