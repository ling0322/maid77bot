package maid77bot

import (
	"context"
)

// Bot provides the interface of maid77 bot
type Bot struct {
	dialogCtrl *DialogController

	// Config of this bot
	Config *Config
}

// NewBot creates a new instance of Bot from config file
func NewBot(configFile string) (*Bot, error) {
	config, err := LoadConfig(configFile)
	if err != nil {
		return nil, err
	}

	dialogCtrl, err := NewDialogController(config)
	if err != nil {
		return nil, err
	}

	return &Bot{
		dialogCtrl: dialogCtrl,
		Config:     config,
	}, nil
}

// Reply replies a user message
func (b *Bot) Reply(m *UserMessage, ctx context.Context) (*ReplyMessage, error) {
	return b.dialogCtrl.Reply(m, ctx)
}
