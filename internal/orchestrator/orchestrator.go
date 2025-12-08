package orchestrator

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"rsandz/bearlawyergo/internal/logging"
	"rsandz/bearlawyergo/internal/message"
)

var (
	ErrNoHandlerFound = errors.New("no handler found for message")
)

type Handler interface {
	Handle(ctx context.Context, message *message.Request, response *message.Response) error
	CanHandle(ctx context.Context, message *message.Request) bool
}

type Orchestrator struct {
	handlers []Handler
	logger   *slog.Logger
}

func NewOrchestrator(handlers []Handler, logger *slog.Logger) *Orchestrator {
	return &Orchestrator{
		handlers: handlers,
		logger:   logger,
	}
}

func (orchestrator *Orchestrator) Handle(ctx context.Context, msg *message.Request) (*message.Response, error) {
	traceID, err := generateTraceID()
	response := &message.Response{ShouldContinueHandling: true}

	if err != nil {
		orchestrator.logger.Error("Failed to generate trace ID", "error", err)
		// Fallback to no trace ID or some default? Let's proceed with empty or error
	}

	// Inject trace ID into context
	ctx = context.WithValue(ctx, logging.TraceIDKey, traceID)

	orchestrator.logger.InfoContext(ctx, "Orchestrator received message", "content", msg.RequestMessage.Content)
	messageWasHandled := false
	for _, h := range orchestrator.handlers {
		if h.CanHandle(ctx, msg) {
			orchestrator.logger.InfoContext(ctx, "Handler found for message", "handler_type", fmt.Sprintf("%T", h))
			if err := h.Handle(ctx, msg, response); err != nil {
				orchestrator.logger.ErrorContext(ctx, "Handler failed to handle message", "error", err)
				return nil, err
			}
			messageWasHandled = true

			if !response.ShouldContinueHandling {
				orchestrator.logger.InfoContext(ctx, "Halting request handling early")
				break
			}
		}
	}

	if !messageWasHandled {
		orchestrator.logger.WarnContext(ctx, "No handler found for message", "content", msg.RequestMessage.Content)
		return nil, ErrNoHandlerFound
	}
	return response, nil
}

func generateTraceID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to read random bytes: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}
