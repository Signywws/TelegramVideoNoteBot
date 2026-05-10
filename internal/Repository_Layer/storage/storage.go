package storage

import (
	"context"
	"io"
	"time"
)

type FileStorage interface {
	Save(ctx context.Context, userID int64, category string, reader io.Reader) (string, error)
	Get(ctx context.Context, path string) (io.Reader, error)
}

type FileRecord struct {
	UserID    int64     `bson:"user_id"`
	ChatID    int64     `bson:"chat_id"`
	MessageID int       `bson:"message_id"`
	Original  string    `bson:"original_path"`
	VideoNote string    `bson:"video_note_path"`
	CreatedAt time.Time `bson:"created_at"`
}

type FileRepository interface {
	InsertRecord(ctx context.Context, record *FileRecord) error
	Close(ctx context.Context) error
}
