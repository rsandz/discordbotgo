package llm

import (
	"context"
	"errors"
	"testing"

	"rsandz/bearlawyergo/internal/message"

	"io"
	"log/slog"

	"github.com/tmc/langchaingo/llms"
)

type mockLLM struct {
	GenerateContentFunc func(ctx context.Context, messages []llms.MessageContent, options ...llms.CallOption) (*llms.ContentResponse, error)
	CallFunc            func(ctx context.Context, prompt string, options ...llms.CallOption) (string, error)
}

func (m *mockLLM) GenerateContent(ctx context.Context, messages []llms.MessageContent, options ...llms.CallOption) (*llms.ContentResponse, error) {
	if m.GenerateContentFunc != nil {
		return m.GenerateContentFunc(ctx, messages, options...)
	}
	// Default behavior
	return &llms.ContentResponse{
		Choices: []*llms.ContentChoice{
			{Content: "mock response"},
		},
	}, nil
}

func (m *mockLLM) Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
	if m.CallFunc != nil {
		return m.CallFunc(ctx, prompt, options...)
	}
	return "mock response", nil
}

func TestLLMHandler_Handle(t *testing.T) {
	tests := []struct {
		name         string
		inputContent string
		mockResponse string
		mockError    error
		expectedResp string
		expectError  bool
	}{
		{
			name:         "Success",
			inputContent: "Hello",
			mockResponse: "I am a bear lawyer",
			expectedResp: "I am a bear lawyer",
		},
		{
			name:         "LLM Error",
			inputContent: "Hello",
			mockError:    errors.New("llm error"),
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockLLM{
				GenerateContentFunc: func(ctx context.Context, messages []llms.MessageContent, options ...llms.CallOption) (*llms.ContentResponse, error) {
					if tt.mockError != nil {
						return nil, tt.mockError
					}
					return &llms.ContentResponse{
						Choices: []*llms.ContentChoice{
							{Content: tt.mockResponse},
						},
					}, nil
				},
			}

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			h, err := NewLLMHandler(mock, logger)
			if err != nil {
				t.Fatalf("NewLLMHandler failed: %v", err)
			}

			resp := &message.Response{}
			err = h.Handle(context.Background(), &message.Request{RequestMessage: message.Message{Content: tt.inputContent}}, resp)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if resp.ResponseMessage.Content != tt.expectedResp {
				t.Errorf("expected content %q, got %q", tt.expectedResp, resp.ResponseMessage.Content)
			}
		})
	}
}

func TestLLMHandler_CanHandle(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	h, _ := NewLLMHandler(&mockLLM{}, logger)
	if !h.CanHandle(context.Background(), &message.Request{RequestMessage: message.Message{Content: "Hello"}}) {
		t.Error("CanHandle should always return true")
	}
}
