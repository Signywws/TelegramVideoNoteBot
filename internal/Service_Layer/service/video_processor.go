package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"videonotebot/internal/Presentation_Layer/clients"
	"videonotebot/internal/Repository_Layer/storage"
)

type VideoProcessor struct {
	client    *clients.Client
	converter *Converter
	fileStore storage.FileStorage
	repo      storage.FileRepository
}

func NewVideoProcessor(client *clients.Client, converter *Converter, storage storage.FileStorage, repo storage.FileRepository) *VideoProcessor {

	return &VideoProcessor{
		client:    client,
		converter: converter,
		fileStore: storage,
		repo:      repo,
	}
}

func (p *VideoProcessor) Process(ctx context.Context, chatID int64, video *clients.Video, messageID int) error {
	fileInfo, err := p.client.GetFile(video.FileID) // получаем информацию о файле
	if err != nil {
		return fmt.Errorf("get file: %w", err)
	}

	tmpFile, err := os.CreateTemp("", "video_original_*.mp4")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name()) // удаляем после конвертации

	log.Printf("Downloading file %s...", fileInfo.FilePath)

	if err := p.client.DownloadFile(fileInfo.FilePath, tmpFile); err != nil {
		return fmt.Errorf("download: %w", err)
	}

	defer tmpFile.Close()

	// Сохраняем оригинал в постоянное хранилище (передаём открытый файл)
	// Переоткрываем, т.к. tmpFile уже закрыт
	origFile, err := os.Open(tmpFile.Name())
	if err != nil {
		return fmt.Errorf("reopen temp: %w", err)
	}
	origPath, err := p.fileStore.Save(ctx, chatID, "original", origFile)
	origFile.Close()
	if err != nil {
		return fmt.Errorf("save original: %w", err)
	}
	log.Printf("Original saved to %s", origPath)
	// Конвертируем в Video Note
	// Создаём временный файл для результата
	noteTmp, err := os.CreateTemp("", "video_note_*.mp4")
	if err != nil {
		return fmt.Errorf("create temp note: %w", err)
	}
	noteTmpPath := noteTmp.Name()
	noteTmp.Close()
	defer os.Remove(noteTmpPath)

	log.Printf("Converting %s -> %s", tmpFile.Name(), noteTmpPath)
	if err := p.converter.ConvertToVideoNote(ctx, tmpFile.Name(), noteTmpPath); err != nil {
		return fmt.Errorf("convert: %w", err)
	}

	// Сохраняем кружочек в сторадж
	noteFile, err := os.Open(noteTmpPath)
	if err != nil {
		return fmt.Errorf("open note: %w", err)
	}
	notePath, err := p.fileStore.Save(ctx, chatID, "video_note", noteFile)
	noteFile.Close()
	if err != nil {
		return fmt.Errorf("save note: %w", err)
	}
	log.Printf("Video note saved to %s", notePath)

	if err := p.repo.InsertRecord(context.Background(), &storage.FileRecord{
		UserID:    chatID,
		ChatID:    chatID,
		MessageID: messageID, // нужно прокинуть из Handler
		Original:  origPath,
		VideoNote: notePath,
		CreatedAt: time.Now(),
	}); err != nil {
		log.Printf("Failed to insert mongo record: %v", err)
	}

	// Отправляем Video Note
	log.Printf("Sending video note to chat %d...", chatID)
	if err := p.client.SendVideoNote(chatID, notePath, video.Duration); err != nil {
		return fmt.Errorf("send video note: %w", err)
	}

	log.Printf("Successfully processed video for chat %d", chatID)
	return nil

}
