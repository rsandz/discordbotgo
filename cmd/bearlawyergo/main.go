package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"rsandz/bearlawyergo/internal/cli"
	"rsandz/bearlawyergo/internal/discord"
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
	// Parse flags
	useDiscord := flag.Bool("discord", false, "Run as Discord bot")
	flag.Parse()

	if *useDiscord {
		token := os.Getenv("DISCORD_TOKEN")
		if token == "" {
			logger.Error("DISCORD_TOKEN environment variable not set")
			os.Exit(1)
		}

		bot, err := discord.NewBot(token, orch, logger)
		if err != nil {
			logger.Error("Failed to create Discord bot", "error", err)
			os.Exit(1)
		}

		if err := bot.Start(); err != nil {
			logger.Error("Failed to start Discord bot", "error", err)
			os.Exit(1)
		}
		defer bot.Close()

		logger.Info("Discord bot running.")
		fmt.Println("Running. Press CTRL-C to exit.")
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop

		logger.Info("Shutting down Discord bot...")
	} else {
		repl := cli.NewREPL(orch, logger)
		repl.Start(ctx)
	}
}
