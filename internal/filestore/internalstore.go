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

func NewInternalStore(root string) (*InternalStore, error) {
	if err := os.MkdirAll(root, 0755); err != nil {
		return nil, err
	}

	return &InternalStore{
		root: root,
	}, nil
}

func (i *InternalStore) init(object string) error {
	path := filepath.Join(i.root, object)
	if strings.Contains(object, i.root) {
		path = object
	}

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	fi, err := os.Stat(f.Name())
	if err != nil {
		return err
	}

	i.size = uint64(fi.Size())

	return nil
}

func (i *InternalStore) Write(ctx context.Context, object string, p []byte) (n int64, err error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if err := i.init(object); err != nil {
		return 0, err
	}

	w, err := i.buf.Write(p)
	if err != nil {
		return 0, err
	}

	// finWidth := w + lenWidth
	i.size += uint64(w)

	i.reset()

	return int64(w), nil
}

func (i *InternalStore) Read(ctx context.Context, path string) ([]byte, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.reset()

	if err := i.init(path); err != nil {
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
	defer i.mu.Unlock()

	if err := os.Remove(path); err != nil {
		return err
	}

	return nil
}

func (i *InternalStore) reset() error {
	i.mu.Lock()
	defer i.mu.Unlock()

	if err := i.buf.Flush(); err != nil {
		return err
	}

	i.buf = nil
	i.size = 0

	return i.File.Close()
}
