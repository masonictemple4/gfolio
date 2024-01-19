package filestore

import (
	"context"
	"strings"
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

func GetRootPath(i Filestore) string {
	var res string
	switch i.(type) {
	case *InternalStore:
		res = i.(*InternalStore).root
	case *GCPStore:
		res = i.(*GCPStore).root
	}

	res = strings.TrimPrefix(res, "./")

	return res
}
