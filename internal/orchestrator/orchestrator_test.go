package orchestrator

import (
	"context"
	"io"
	"log/slog"
	"rsandz/bearlawyergo/internal/message"
	"testing"
)

type mockHandler struct {
	canHandle      bool
	shouldContinue bool
	called         bool
}

func (h *mockHandler) Handle(ctx context.Context, m *message.Request, response *message.Response) error {
	h.called = true
	response.ResponseMessage = message.Message{Content: "mock response"}
	response.ShouldContinueHandling = h.shouldContinue
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
				&mockHandler{canHandle: false, shouldContinue: true},
			},
			wantErr:     true,
			expectedErr: ErrNoHandlerFound,
		},
		{
			name: "handler can handle",
			handlers: []Handler{
				&mockHandler{canHandle: true, shouldContinue: true},
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

func TestOrchestrator_Handle_MultipleHandlers(t *testing.T) {
	h1 := &mockHandler{canHandle: true, shouldContinue: true}
	h2 := &mockHandler{canHandle: true, shouldContinue: false}
	h3 := &mockHandler{canHandle: true, shouldContinue: true}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	o := NewOrchestrator([]Handler{h1, h2, h3}, logger)

	_, err := o.Handle(context.Background(), &message.Request{})
	if err != nil {
		t.Errorf("Handle() unexpected error = %v", err)
	}

	if !h1.called {
		t.Error("Handler 1 should have been called")
	}
	if !h2.called {
		t.Error("Handler 2 should have been called")
	}
	if h3.called {
		t.Error("Handler 3 should NOT have been called")
	}
}
