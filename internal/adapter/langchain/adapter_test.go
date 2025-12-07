package langchain

import (
	"rsandz/bearlawyergo/internal/message"
	"testing"

	"github.com/tmc/langchaingo/llms"
)

func TestToLLMMessage(t *testing.T) {
	tests := []struct {
		name     string
		input    message.Message
		expected llms.ChatMessageType
	}{
		{
			name:     "User Message",
			input:    message.Message{Role: message.UserRole, Content: "Hello"},
			expected: llms.ChatMessageTypeHuman,
		},
		{
			name:     "Bot Message",
			input:    message.Message{Role: message.BotRole, Content: "Hi"},
			expected: llms.ChatMessageTypeAI,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToLLMMessage(tt.input)
			if result.Role != tt.expected {
				t.Errorf("Expected role %s, got %s", tt.expected, result.Role)
			}
			content := result.Parts[0].(llms.TextContent).Text
			if content != tt.input.Content {
				t.Errorf("Expected content %s, got %s", tt.input.Content, content)
			}
		})
	}
}
