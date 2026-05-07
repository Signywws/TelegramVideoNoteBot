package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken    string
	PollTimeout int // секунды ожидания клиента
}

func Load() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	timeout, err := strconv.Atoi(os.Getenv("POLLTIMEOUT"))
	if err != nil {
		log.Fatal(err)
	}
	return &Config{
		BotToken:    os.Getenv("TOKEN"),
		PollTimeout: timeout,
	}
}
