package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"videonotebot/internal/config"
	"videonotebot/internal/polling"
	"videonotebot/internal/telegram"
)

const (
	tgBotHost = "api.telegram.org"
)

func main() {
	cfg := config.Load()
	if cfg.BotToken == "" {
		log.Fatal("TOKEN is not set")
	}

	client := telegram.NewClient(cfg.BotToken)

	poller := polling.New(client, cfg.PollTimeout)

	// Контекст который завершается после прерывания например Ctrl + C
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	poller.Start(ctx)

	log.Println("Bot Stopped")

}
