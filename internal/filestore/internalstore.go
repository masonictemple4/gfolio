package filestore

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
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

	if !strings.Contains("/home/mason/personal/masonictempl", root) {
		root = "/home/mason/personal/masonictempl/" + strings.TrimPrefix(root, "/")
	}

	return &InternalStore{
		root: root,
	}, nil
}

func GetRootPath(i *InternalStore) string {
	urlFriendly := strings.TrimPrefix(i.root, "./")

	return urlFriendly
}

func (i *InternalStore) init(path string, flag int) error {

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	println("[Filestore (INIT)] Opening: ", path)
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

	fmt.Println("[Filestore] Wrote ", i.size, " bytes to buffer.")

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

	println("[Filestore (READ)] Calling init from path: ", path)
	if err := i.init(path, os.O_RDONLY); err != nil {
		return nil, err
	}

	data, err := io.ReadAll(i.File)
	if err != nil {
		return nil, err
	}

	fmt.Printf("[Filestore (READ)] Data: %s\n", string(data))

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
