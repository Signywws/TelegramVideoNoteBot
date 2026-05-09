package handlers

import (
	"context"
	"log"
	"videonotebot/internal/Presentation_Layer/clients"
	"videonotebot/internal/Service_Layer/service"
	"videonotebot/internal/pool"
)

const getFileURL = "https://api.telegram.org/file/bot"

type VideoHandler struct {
	client    *clients.Client
	pool      *pool.Pool
	processor *service.VideoProcessor
}

func NewVideoHandler(client *clients.Client, pool *pool.Pool, processor *service.VideoProcessor) *VideoHandler {
	return &VideoHandler{
		client:    client,
		pool:      pool,
		processor: processor,
	}
}

func (h *VideoHandler) Handle(msg *clients.Message) {
	video := msg.Video
	chatID := msg.Chat.ChatID
	fileID := msg.Video.FileID
	log.Printf("Video from %d: file_id=%s, duration=%d", chatID, fileID, msg.Video.Duration)
	if msg.Video.Duration > 60 {
		h.client.SendMessage(chatID, "Видео слишком длинное. Пожалуйста, пришлите ролик до 60 секунд.")
		return
	}

	h.client.SendMessage(chatID, "Видео получено. Обработка...")
	h.pool.Submit(func() {
		log.Printf("Processing video %s for chat %d", video.FileID, chatID)
		if err := h.processor.Process(context.Background(), int64(chatID), video); err != nil {
			log.Printf("Error processing video: %v", err)
			h.client.SendMessage(chatID, "Не удалось обработать видео. Возможно, формат не поддерживается.")
			return
		}
	})

}
