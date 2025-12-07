package memory

import (
	"rsandz/bearlawyergo/internal/message"
	"testing"
)

func TestChatHistory(t *testing.T) {
	history := NewChatHistory()

	if len(history.Messages()) != 0 {
		t.Errorf("Expected empty history, got %d messages", len(history.Messages()))
	}

	history.Add(message.Message{Role: message.UserRole, Content: "Hello"})
	if len(history.Messages()) != 1 {
		t.Errorf("Expected 1 message, got %d", len(history.Messages()))
	}
	if history.Messages()[0].Role != message.UserRole {
		t.Errorf("Expected user role, got %s", history.Messages()[0].Role)
	}

	history.Add(message.Message{Role: message.BotRole, Content: "Hi there"})
	if len(history.Messages()) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(history.Messages()))
	}
	if history.Messages()[1].Role != message.BotRole {
		t.Errorf("Expected bot role, got %s", history.Messages()[1].Role)
	}

	history.Clear()
	if len(history.Messages()) != 0 {
		t.Errorf("Expected empty history after clear, got %d", len(history.Messages()))
	}
}
