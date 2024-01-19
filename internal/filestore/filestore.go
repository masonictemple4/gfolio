package filestore

import (
	"context"
)

// Filestore handles the storage of static files.
// Wrapper around io.Reader & io.Writer interfaces.
// Write(p []byte) (n int, err error)
// Read(p []byte) (n int, err error)
type Filestore interface {
	Read(ctx context.Context, path string) ([]byte, error)
	Write(ctx context.Context, object string, p []byte) (n int64, err error)
	Delete(ctx context.Context, path string) error
}
