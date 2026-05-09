package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type LocalStotage struct {
	basePath string
}

func NewFileStore(basePath string) (*LocalStotage, error) {
	absPath, err := filepath.Abs(basePath)
	if err != nil {
		return nil, fmt.Errorf("abs path: %w", err)
	}
	path := filepath.Dir(absPath)
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, fmt.Errorf("create base dir: %w", err)
	}
	return &LocalStotage{basePath: absPath}, nil
}

func (l *LocalStotage) Save(_ context.Context, userID int64, category string, reader io.Reader) (string, error) {
	dir := filepath.Join(l.basePath, fmt.Sprintf("%d", userID))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("create user dir: %w", err)
	}
	fileName := uuid.NewString() + ".mp4"
	pathToFile := filepath.Join(dir, fileName)
	file, err := os.Create(pathToFile)
	if err != nil {
		return "", fmt.Errorf("create file err")
	}
	defer file.Close()
	if _, err := io.Copy(file, reader); err != nil {
		return "", fmt.Errorf("write data: %w", err)
	}
	return pathToFile, nil
}

func (l *LocalStotage) Get(_ context.Context, path string) (io.Reader, error) {
	absPath, err := filepath.Abs(l.basePath)
	if err != nil {
		return nil, fmt.Errorf("abs path: %w", err)
	}
	fmt.Println("get path: ", absPath)
	pathToFile := filepath.Join(absPath, path)
	fmt.Println("Path To FIle: ", pathToFile)
	return os.Open(pathToFile)
}
