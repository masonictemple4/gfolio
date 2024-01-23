package filestore

import (
	"context"
	"os"
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

// Returns the fully qualified path for the root of the filestore.
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

// NewFilestore returns a Filestore based on the method provided.
// If method is not provided, it defaults to internal.
//
// Valid methods are:
// internal, gcp
//
// Important:
//   - when using internal Requires current working directory.
//
// Note: Passing a value here overrides env ASSET_DIR
// and if that is not set either, we default to "assets".
//
// You can leave assetRoot empty or fill it with os.Getenv("ASSET_DIR")
// when implenting. That or replace it with the directory name for your
// public files..
func NewFilestore(method, assetRoot string) (Filestore, error) {
	var fs Filestore
	var err error

	if assetRoot == "" {
		if val := os.Getenv("ASSET_DIR"); val != "" && val != "/" {
			assetRoot = val
		} else {
			assetRoot = "assets"
		}
	}

	switch method {
	case "internal":
		fs, err = NewInternalStore(assetRoot)
	case "gcp":
		fs = NewGCPStore(assetRoot, false, 0)
	default:
		fs, err = NewInternalStore(assetRoot)
	}

	return fs, err
}
