package discord

import (
	"context"
	"log/slog"
	"rsandz/bearlawyergo/internal/message"
	"rsandz/bearlawyergo/internal/orchestrator"
	"slices"

	discordgo "github.com/bwmarrin/discordgo"
)

type Bot struct {
	discord      *discordgo.Session
	orchestrator *orchestrator.Orchestrator

	logger *slog.Logger
}

func NewBot(token string, orchestrator *orchestrator.Orchestrator, logger *slog.Logger) (*Bot, error) {
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	bot := &Bot{discord: discord, orchestrator: orchestrator, logger: logger}
	discord.AddHandler(bot.handleMessage)

	return bot, nil
}

func (b *Bot) Start() error {
	return b.discord.Open()
}

func (b *Bot) Close() error {
	return b.discord.Close()
}

func (b *Bot) handleMessage(session *discordgo.Session, m *discordgo.MessageCreate) {
	ctx := context.Background()

	if !b.shouldRespond(m.Message) {
		return
	}

	b.logger.Info("Responding to Discord message", "user", m.Author.ID, "user_name", m.Author.Username, "content", m.Content)

	session.ChannelTyping(m.ChannelID)

	history := b.resolveHistory(m.ChannelID)
	msg := message.NewMessage(m.Author.Username, m.Content, message.UserRole)
	req := message.NewRequest(
		*msg,
		history,
		m.ChannelID,
	)

	resp, err := b.orchestrator.Handle(ctx, req)
	if err != nil {
		b.handleError(ctx, m.ChannelID, err)
		return
	}

	b.discord.ChannelMessageSend(m.ChannelID, resp.ResponseMessage.Content)
}

func (b *Bot) shouldRespond(m *discordgo.Message) bool {
	if m.Author.ID == b.discord.State.User.ID {
		b.logger.Debug("Received message from self", "content", m.Content)
		return false
	}
	return slices.ContainsFunc(m.Mentions, func(mention *discordgo.User) bool {
		return mention.ID == b.discord.State.User.ID
	})
}

func (b *Bot) resolveHistory(channelID string) []message.Message {
	var history []message.Message
	messages, _ := b.discord.ChannelMessages(channelID, 10, "", "", "")

	// Discord returns messages from newest to oldest. Iterate backwards to keep chronological order.
	for i := len(messages) - 1; i >= 0; i-- {
		dm := messages[i]
		role := message.UserRole
		if dm.Author.ID == b.discord.State.User.ID {
			role = message.BotRole
		}
		history = append(history, *message.NewMessage(dm.Author.Username, dm.Content, role))
	}
	return history
}

func (b *Bot) handleError(ctx context.Context, channelId string, err error) {
	b.discord.ChannelMessageSend(channelId, "Sorry! Something went wrong. Please try again later.")
	b.logger.InfoContext(ctx, "error handling message", "error", err, "channel_id", channelId)
}
