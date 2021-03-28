package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
)

type Bot struct {
	s *discordgo.Session
}

func NewBot(token string) (*Bot, error) {
	b := &Bot{}

	// Create Discord Session
	s, err := discordgo.New(fmt.Sprintf("Bot %s", token))
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Fatal("failed to create bot", zap.Error(err))
	}

	// Add event handlers
	s.AddHandler(b.ready)

	b.s = s
	return b, err
}

func (b *Bot) Start() error {
	return b.s.Open()
}

func (b *Bot) Stop() error {
	return b.s.Close()
}

func (*Bot) ready(session *discordgo.Session, event *discordgo.Ready) {
	zap.L().Info("bot is ready", zap.Int("guilds", len(event.Guilds)))
}
