package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"videonotebot/internal/Presentation_Layer/clients"
	"videonotebot/internal/Presentation_Layer/dispatcher"
	"videonotebot/internal/Presentation_Layer/polling"
	"videonotebot/internal/Repository_Layer/storage"
	"videonotebot/internal/Service_Layer/service"
	"videonotebot/internal/config"
	"videonotebot/internal/pool"
)

const (
	tgBotHost = "api.telegram.org"
)

func main() {
	cfg := config.Load()
	if cfg.BotToken == "" {
		log.Fatal("TOKEN is not set")
	}

	client := clients.NewClient(cfg.BotToken)

	// инициализация сервисов

	// Хранилище
	fileStore, err := storage.NewFileStore("../media/")
	if err != nil {
		log.Fatalf("Failed to init file storage: %v", err)
	}
	converter := service.NewConverter()

	mongoURL := os.Getenv("MONGO_URI")
	if mongoURL == "" {
		mongoURL = "mongodb://localhost:27017/"
	}

	mongoRepo, err := storage.NewMongoRepo(context.Background(), mongoURL, "data")

	// Video Processor
	videoProcessor := service.NewVideoProcessor(client, converter, fileStore, mongoRepo)

	// Worker Pool
	pool := pool.NewPool(5)

	// Dispatcher
	dispatcher := dispatcher.NewDispatcher(client, pool, videoProcessor, fileStore)

	// Poller
	poller := polling.NewPoller(client, dispatcher, cfg.PollTimeout)

	// Контекст который завершается после прерывания например Ctrl + C
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	poller.Start(ctx)
	log.Println("Bot started. Press Ctrl+C to stop.")

	<-ctx.Done()
	log.Println("Shutting down...")
	pool.Shutdown()

	log.Println("Bot Stopped")

}
