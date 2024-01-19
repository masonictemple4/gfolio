package filestore

import (
	"bufio"
	"context"
	"encoding/binary"
	"io"
	"os"
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
}

func NewInternalStore(path string) (*InternalStore, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}

	size := uint64(fi.Size())

	return &InternalStore{
		File: f,
		buf:  bufio.NewWriter(f),
		size: size,
	}, nil
}

func (i *InternalStore) Write(ctx context.Context, object string, p []byte) (n int64, err error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	w, err := i.buf.Write(p)
	if err != nil {
		return 0, err
	}

	// finWidth := w + lenWidth
	i.size += uint64(w)

	return int64(w), nil
}

func (i *InternalStore) Read(ctx context.Context, path string) ([]byte, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if err := i.buf.Flush(); err != nil {
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

func (i *InternalStore) Size(ctx context.Context, path string) (uint64, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	return i.size, nil
}

func (i *InternalStore) Close() error {
	i.mu.Lock()
	defer i.mu.Unlock()

	if err := i.buf.Flush(); err != nil {
		return err
	}

	return i.File.Close()
}
