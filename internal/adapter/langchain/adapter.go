package langchain

import (
	"rsandz/bearlawyergo/internal/message"

	"github.com/tmc/langchaingo/llms"
)

// ToLLMMessage converts an internal message to a LangChain message content.
func ToLLMMessage(msg message.Message) llms.MessageContent {
	role := llms.ChatMessageTypeAI
	if msg.Role == message.UserRole {
		role = llms.ChatMessageTypeHuman
	}

	return llms.MessageContent{
		Role:  role,
		Parts: []llms.ContentPart{llms.TextContent{Text: msg.Content}},
	}
}

// ToLLMMessages converts a slice of internal messages to a slice of LangChain message contents.
func ToLLMMessages(messages []message.Message) []llms.MessageContent {
	if messages == nil {
		return nil
	}

	result := make([]llms.MessageContent, len(messages))
	for i, msg := range messages {
		result[i] = ToLLMMessage(msg)
	}
	return result
}
