package filestore

import (
	"bufio"
	"context"
	"encoding/binary"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var _ Filestore = (*InternalStore)(nil)

var (
	enc = binary.BigEndian
)

const (
	lenWidth = 8
)

type InternalStore struct {
	*os.File
	mu   sync.Mutex
	buf  *bufio.Writer
	size uint64
	root string
}

// NewInternalStore creates a new internal store
// and returns a pointer to it.
//
// Important: Requires current working directory.
//
// Note: Passing a value here overrides env ASSET_DIR
// and if that is not set either, we default to "assets".
//
// "/" is not an eligible path for internal/local storage.
// However, this can be used in remote storages as the root.
//
// You can leave assetRoot empty or fill it with os.Getenv("ASSET_DIR")
// when implenting. That or replace it with the directory name for your
// public files..
func NewInternalStore(assetRoot string) (*InternalStore, error) {
	if assetRoot == "" {
		if val := os.Getenv("ASSET_DIR"); val != "" && val != "/" {
			assetRoot = val
		} else {
			assetRoot = "assets"
		}
	}

	wd := os.Getenv("WORKDIR")

	if !strings.Contains(wd, assetRoot) {
		if strings.HasSuffix(wd, "/") {
			assetRoot = wd + strings.TrimPrefix(assetRoot, "/")
		} else {
			assetRoot = wd + "/" + strings.TrimPrefix(assetRoot, "/")
		}
	}

	if err := os.MkdirAll(assetRoot, 0755); err != nil {
		return nil, err
	}

	return &InternalStore{
		root: assetRoot,
	}, nil
}

func (i *InternalStore) init(path string, flag int) error {

	path = filepath.Join(os.Getenv("WORKDIR"), path)

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	f, err := os.OpenFile(path, flag, 0755)
	if err != nil {
		return err
	}

	fi, err := os.Stat(f.Name())
	if err != nil {
		return err
	}

	i.File = f
	i.buf = bufio.NewWriter(f)

	i.size = uint64(fi.Size())

	return nil
}

func (i *InternalStore) Write(ctx context.Context, object string, p []byte) (n int64, err error) {
	i.mu.Lock()

	// NOTE: Will replace the existing data.
	// As soon as it opens, so only call write if necessary.
	if err := i.init(object, os.O_RDWR|os.O_CREATE|os.O_TRUNC); err != nil {
		return 0, err
	}

	w, err := i.buf.Write(p)
	if err != nil {
		return 0, err
	}

	// finWidth := w + lenWidth
	i.size += uint64(w)

	// need to pass the lock here.
	defer i.reset()
	defer i.mu.Unlock()

	return int64(w), nil
}

func (i *InternalStore) Read(ctx context.Context, path string) ([]byte, error) {
	i.mu.Lock()

	// I'm back and forth on this one
	// the reset will not get called until
	// after we unlock from reading.
	// This means that all other methods that interact
	// with the store will need to call reset to prevent
	// unwanted behavior.
	defer i.reset()
	defer i.mu.Unlock()

	if err := i.init(path, os.O_RDONLY); err != nil {
		return nil, err
	}

	data, err := io.ReadAll(i.File)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (i *InternalStore) Delete(ctx context.Context, path string) error {
	i.mu.Lock()

	defer i.reset()
	defer i.mu.Unlock()

	if err := os.Remove(path); err != nil {
		return err
	}

	return nil
}

func (i *InternalStore) reset() error {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.buf = nil
	i.size = 0

	return i.File.Close()
}
