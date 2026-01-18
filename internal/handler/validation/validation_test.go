package validation

import (
	"context"
	"rsandz/bearlawyergo/internal/message"
	"strings"
	"testing"
)

func TestValidationHandler(t *testing.T) {
	tests := []struct {
		name                     string
		inputRequest             *message.Request
		expectedPassesValidation bool
	}{
		{
			name:                     "Success",
			inputRequest:             buildRequestForString("Hello"),
			expectedPassesValidation: true,
		},
		{
			name:                     "Message too long",
			inputRequest:             buildRequestForString(strings.Repeat("a", 501)),
			expectedPassesValidation: false,
		},
		{
			name:                     "No message",
			inputRequest:             buildRequestForString(""),
			expectedPassesValidation: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vh := NewHandler()
			msg := tt.inputRequest
			response := &message.Response{ShouldContinueHandling: true}
			if err := vh.Handle(context.Background(), msg, response); err != nil {
				t.Errorf("Handle() error = %v", err)
			}

			if response.ShouldContinueHandling != tt.expectedPassesValidation {
				t.Errorf("expected continue handling %v, got %v", tt.expectedPassesValidation, response.ShouldContinueHandling)
			}
		})
	}
}

func buildRequestForString(content string) *message.Request {
	return &message.Request{
		RequestMessage: message.Message{
			Content: content,
		},
	}
}
