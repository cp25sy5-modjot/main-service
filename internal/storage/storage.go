package storage

import "context"

type Storage interface {
	Save(ctx context.Context, userID string, data []byte, ext string) (string, error)
	Load(ctx context.Context, path string) ([]byte, error)
	Delete(ctx context.Context, path string) error
}
