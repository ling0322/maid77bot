package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/net/proxy"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"maid77bot"
	"net/http"
)

// createHttpClientSocks5 creates a http client with socks5 proxy
func createHttpClientSocks5(address string) (*http.Client, error) {
	dialer, err := proxy.SOCKS5("tcp", address, nil, proxy.Direct)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	httpTransport := &http.Transport{Dial: dialer.Dial}
	httpClient := &http.Client{Transport: httpTransport}

	return httpClient, nil
}

// createHttpClient creates the httpClient for bot according to config
func createHttpClient(c *maid77bot.Config) (*http.Client, error) {
	if c.Socks5Proxy != "" {
		return createHttpClientSocks5(c.Socks5Proxy)
	} else {
		return &http.Client{}, nil
	}
}

func main() {
	maid77, err := maid77bot.NewBot("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	httpClient, err := createHttpClient(maid77.Config)
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPIWithClient(
		maid77.Config.TelegramToken,
		httpClient)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		fromMsg := &maid77bot.UserMessage{
			UserId: fmt.Sprint(update.Message.Chat.ID),
			Text:   update.Message.Text,
		}
		reply, err := maid77.Reply(fromMsg, context.Background())
		if err != nil {
			log.Println(err)
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply.Text)
		bot.Send(msg)
	}
}
