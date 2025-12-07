package memory

import (
	"sync"

	"rsandz/bearlawyergo/internal/message"
)

// ChatHistory is a thread-safe store for chat messages.
type ChatHistory struct {
	mu       sync.RWMutex
	messages []message.Message
}

// NewChatHistory creates a new ChatHistory instance.
func NewChatHistory() *ChatHistory {
	return &ChatHistory{
		messages: make([]message.Message, 0),
	}
}

// Add adds a message to the history.
func (h *ChatHistory) Add(msg message.Message) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.messages = append(h.messages, msg)
}

// Messages returns a copy of the message history.
func (h *ChatHistory) Messages() []message.Message {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// Return a copy to prevent race conditions if the caller modifies the slice
	msgs := make([]message.Message, len(h.messages))
	copy(msgs, h.messages)
	return msgs
}

// Clear clears the message history.
func (h *ChatHistory) Clear() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.messages = make([]message.Message, 0)
}
