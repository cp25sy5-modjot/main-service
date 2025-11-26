package localfs

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	"github.com/cp25sy5-modjot/main-service/internal/storage"
)

type LocalStorage struct {
	baseDir string
}

var _ storage.Storage = (*LocalStorage)(nil)

func NewLocalStorage(baseDir string) (*LocalStorage, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, err
	}
	return &LocalStorage{baseDir: baseDir}, nil
}

// Save stores bytes and returns a relative path like "userID/2025/11/<uuid>.png".
func (s *LocalStorage) Save(ctx context.Context, userID string, data []byte, ext string) (string, error) {
	now := time.Now()
	fileName := fmt.Sprintf("%s.%s", uuid.New().String(), ext)
	relPath := filepath.Join(userID, fmt.Sprintf("%d", now.Year()), fmt.Sprintf("%02d", now.Month()), fileName)
	fullPath := filepath.Join(s.baseDir, relPath)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return "", err
	}
	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return "", err
	}
	return filepath.ToSlash(relPath), nil
}

func (s *LocalStorage) Load(ctx context.Context, path string) ([]byte, error) {
	fullPath := filepath.Join(s.baseDir, path)
	return os.ReadFile(fullPath)
}

func (s *LocalStorage) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(s.baseDir, path)
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
