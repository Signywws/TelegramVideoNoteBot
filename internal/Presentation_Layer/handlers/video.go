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
	// fileUnqID := msg.Video.FileUniqID
	size := msg.Video.Size
	log.Printf("Video from %d: file_id=%s, duration=%d", chatID, fileID, msg.Video.Duration)
	if msg.Video.Duration > 60 {
		h.client.SendMessage(chatID, "Видео слишком длинное. Пожалуйста, пришлите ролик до 60 секунд.")
		return
	}
	if size > 524288000 { // 500mb
		h.client.SendMessage(chatID, "Видео слишком большое. Пожалуйста, пришлите видое размером до 500mb.")
		return
	}

	h.pool.Submit(func() {
		h.client.SendMessage(chatID, "Видео получено. Обработка...")
		log.Printf("Processing video %s for chat %d", video.FileID, chatID)
		if err := h.processor.Process(context.Background(), int64(chatID), video, msg.MessageID); err != nil {
			log.Printf("Error processing video: %v", err)
			h.client.SendMessage(chatID, "Не удалось обработать видео. Возможно, формат не поддерживается.")
			return
		}
	})

}
