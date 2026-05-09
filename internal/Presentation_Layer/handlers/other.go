package handlers

import (
	"fmt"
	"log"
	"videonotebot/internal/Presentation_Layer/clients"
)

type OtherHandler struct {
	client *clients.Client
}

func NewOtherHandler(client *clients.Client) *OtherHandler {
	return &OtherHandler{
		client: client,
	}
}

func (h *OtherHandler) Handle(msg *clients.Message) {
	chatID := msg.Chat.ChatID
	log.Println("Unsupported content from %d", chatID)

	reply := "Я понимаю только видео. Отправьте видео, и я сделаю из него кружочек."
	if err := h.client.SendMessage(chatID, reply); err != nil {
		fmt.Errorf("OtherHandler error: %w", err)
	}
}
