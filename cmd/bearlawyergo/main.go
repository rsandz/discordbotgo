package main

import (
	"context"
	"log"
	"os"

	"rsandz/bearlawyergo/internal/cli"
	llmHandler "rsandz/bearlawyergo/internal/handler/llm"
	"rsandz/bearlawyergo/internal/handler/validation"
	"rsandz/bearlawyergo/internal/logging"
	"rsandz/bearlawyergo/internal/orchestrator"

	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms/openai"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize structured logger
	logger, err := logging.NewLogger("INFO", "bearlawyer.log")
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	logger.Info("Starting Bear Lawyer")

	ctx := context.Background()

	llm, err := openai.New()
	if err != nil {
		logger.Error("Failed to create LLM", "error", err)
		os.Exit(1)
	}

	validationHandler := validation.NewHandler()
	llmHandler, err := llmHandler.NewLLMHandler(llm, logger)
	if err != nil {
		logger.Error("Failed to create LLM handler", "error", err)
		os.Exit(1)
	}

	handlers := []orchestrator.Handler{
		validationHandler,
		llmHandler,
	}

	orch := orchestrator.NewOrchestrator(handlers, logger)
	repl := cli.NewREPL(orch, logger)

	repl.Start(ctx)
}
