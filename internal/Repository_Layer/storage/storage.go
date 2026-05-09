package storage

import (
	"context"
	"io"
)

type FileStorage interface {
	Save(ctx context.Context, userID int64, category string, reader io.Reader) (string, error)
	Get(ctx context.Context, path string) (io.Reader, error)
	// Delete(ctx context.Context, path string) error
}
