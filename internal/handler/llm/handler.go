package llm

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"rsandz/bearlawyergo/internal/adapter/langchain"
	"rsandz/bearlawyergo/internal/config"
	"rsandz/bearlawyergo/internal/message"

	"github.com/tmc/langchaingo/llms"
)

type LLMHandler struct {
	llm          llms.Model
	logger       *slog.Logger
	systemPrompt string
}

func NewLLMHandler(llm llms.Model, logger *slog.Logger) (*LLMHandler, error) {
	prompts, err := config.LoadPrompts()
	if err != nil {
		return nil, fmt.Errorf("failed to load prompts: %w", err)
	}

	return &LLMHandler{
		llm:          llm,
		logger:       logger,
		systemPrompt: prompts.SystemPrompt,
	}, nil
}

func (h *LLMHandler) Handle(ctx context.Context, msg *message.Request, response *message.Response) error {
	h.logger.InfoContext(ctx, "LLMHandler processing message")

	messages := h.buildContextWindow(msg.RequestMessage.Content, msg.History)

	completion, err := h.inferCompletion(ctx, messages)
	if err != nil {
		h.logger.ErrorContext(ctx, "Failed to generate completion", "error", err)
		return fmt.Errorf("failed to generate completion: %w", err)
	}

	h.logger.InfoContext(ctx, "LLMHandler completed request")
	response.ResponseMessage = message.Message{
		Content: completion,
	}
	return nil
}

func (h *LLMHandler) CanHandle(ctx context.Context, msg *message.Request) bool {
	return true
}

func (h *LLMHandler) buildContextWindow(userInput string, history []message.Message) []llms.MessageContent {
	messages := []llms.MessageContent{
		{
			Role:  llms.ChatMessageTypeSystem,
			Parts: []llms.ContentPart{llms.TextContent{Text: h.systemPrompt}},
		},
	}

	messages = append(messages, langchain.ToLLMMessages(history)...)

	messages = append(messages, llms.MessageContent{
		Role:  llms.ChatMessageTypeHuman,
		Parts: []llms.ContentPart{llms.TextContent{Text: userInput}},
	})

	return messages
}

func (h *LLMHandler) inferCompletion(ctx context.Context, messages []llms.MessageContent) (string, error) {
	resp, err := h.llm.GenerateContent(ctx, messages)
	if err != nil {
		return "", err
	}

	choices := resp.Choices
	if len(choices) < 1 {
		return "", errors.New("empty response from model")
	}
	c1 := choices[0]
	return c1.Content, nil
}
