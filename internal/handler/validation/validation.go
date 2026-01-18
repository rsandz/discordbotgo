package validation

import (
	"context"
	"fmt"
	"rsandz/bearlawyergo/internal/message"
)

const MaxMessageLength = 500

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// Execute validation on the message.
// If the message is invalid, set the response message to the validation error and prevent further handling.
// Never returns an error.
func (h *Handler) Handle(ctx context.Context, msg *message.Request, response *message.Response) error {
	if !validateMessageLength(msg) {
		return failValidation(response, fmt.Sprintf("Your message is too long. Please keep it under %d characters.", MaxMessageLength))
	}
	if !validateMessageExists(msg) {
		return failValidation(response, "Please provide a message.")
	}
	return nil
}

// Always handles all messages
func (h *Handler) CanHandle(ctx context.Context, msg *message.Request) bool {
	return true
}

func failValidation(response *message.Response, message string) error {
	response.ResponseMessage.Content = message
	response.ShouldContinueHandling = false
	return nil
}

func validateMessageExists(msg *message.Request) bool {
	return msg.RequestMessage.Content != ""
}

func validateMessageLength(msg *message.Request) bool {
	return len(msg.RequestMessage.Content) <= MaxMessageLength
}
