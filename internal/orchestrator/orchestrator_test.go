package orchestrator

import (
	"context"
	"io"
	"log/slog"
	"rsandz/bearlawyergo/internal/message"
	"testing"
)

type mockHandler struct {
	canHandle bool
}

func (h *mockHandler) Handle(ctx context.Context, m *message.Request, response *message.Response) error {
	response.ResponseMessage = message.Message{Content: "mock response"}
	return nil
}

func (h *mockHandler) CanHandle(ctx context.Context, m *message.Request) bool {
	return h.canHandle
}

func TestOrchestrator_Handle(t *testing.T) {
	tests := []struct {
		name        string
		handlers    []Handler
		wantErr     bool
		expectedErr error
	}{
		{
			name:        "no handlers",
			handlers:    []Handler{},
			wantErr:     true,
			expectedErr: ErrNoHandlerFound,
		},
		{
			name: "handler cannot handle",
			handlers: []Handler{
				&mockHandler{canHandle: false},
			},
			wantErr:     true,
			expectedErr: ErrNoHandlerFound,
		},
		{
			name: "handler can handle",
			handlers: []Handler{
				&mockHandler{canHandle: true},
			},
			wantErr:     false,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			o := NewOrchestrator(tt.handlers, logger)

			_, err := o.Handle(context.Background(), &message.Request{})
			if (err != nil) != tt.wantErr {
				t.Errorf("Handle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.expectedErr != nil && err != tt.expectedErr {
				t.Errorf("Handle() error = %v, expectedErr %v", err, tt.expectedErr)
			}
		})
	}
}
