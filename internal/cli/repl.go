package cli

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"rsandz/bearlawyergo/internal/memory"
	"rsandz/bearlawyergo/internal/message"
	"rsandz/bearlawyergo/internal/orchestrator"
)

const (
	cliUser    = "cli-user"
	cliBot     = "cli-bot"
	cliChannel = "cli-channel"
	maxHistory = 10
)

type REPL struct {
	orchestrator *orchestrator.Orchestrator

	logger *slog.Logger
}

func NewREPL(orchestrator *orchestrator.Orchestrator, logger *slog.Logger) *REPL {
	return &REPL{
		orchestrator: orchestrator,
		logger:       logger,
	}
}

func (r *REPL) Start(ctx context.Context) {
	reader := bufio.NewReader(os.Stdin)
	chatHistory := memory.NewChatHistory()
	fmt.Println("Bear Lawyer CLI started. Type 'quit' or 'exit' to leave.")

	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if text == "quit" || text == "exit" {
			break
		}

		if text == "" {
			continue
		}

		msg := &message.Message{
			Content: text,
			User:    cliUser,
			Role:    message.UserRole,
		}

		history := chatHistory.Messages()
		r.logger.Info("Processing message", "message", text, "history", history)

		// Add user message to history after resolve above to prevent duplicates.
		chatHistory.Add(*msg)

		request := &message.Request{
			RequestMessage: *msg,
			History:        history,
			Channel:        cliChannel,
		}

		resp, err := r.orchestrator.Handle(ctx, request)
		if err != nil {
			fmt.Printf("Error handling message: %v\n", err)
			continue
		}

		if resp != nil {
			fmt.Printf("Bear Lawyer: %s\n", resp.ResponseMessage.Content)
			chatHistory.Add(message.Message{
				Role:    message.BotRole,
				Content: resp.ResponseMessage.Content,
				User:    cliBot,
			})
		}
	}
}
