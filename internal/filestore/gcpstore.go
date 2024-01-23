package filestore

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
)

// verify the Filestore interface is implemented

var _ Filestore = (*GCPStore)(nil)

type GCPStore struct {
	ProgressEnabled bool
	ChunkSize       int64
	root            string
}

// NewGCPStore creates a new gcp store and returns a pointer to it.
//
// Note: Passing a value here overrides env ASSET_DIR
//
// You can leave assetRoot empty or fill it with os.Getenv("ASSET_DIR")
// when implenting. That or replace it with the directory name for your
// public files..
func NewGCPStore(assetRoot string, progress bool, chunk int64) *GCPStore {
	return &GCPStore{
		ProgressEnabled: progress,
		ChunkSize:       chunk,
		root:            assetRoot,
	}
}

// This is essentially for when you're downloading data, or possibly using it for
// some sort of updates. Should not be used to send back to the client or to store
// into the dbs.
func (g *GCPStore) Read(ctx context.Context, path string) ([]byte, error) {
	baseurl := os.Getenv("BUCKET_BASE_URL")
	bucket := os.Getenv("STORAGE_BUCKET")

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	// TODO: Find filename
	buf := new(bytes.Buffer)

	path = removeLocalPath(path, os.Getenv("WORKDIR"))

	println("baseurl: ", baseurl)
	objPath := strings.TrimPrefix(path, baseurl)
	println("bucket: ", bucket)
	println("objPath: ", objPath)

	gcpRdr, err := client.Bucket(bucket).Object(objPath).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer gcpRdr.Close()

	if _, err := io.Copy(buf, gcpRdr); err != nil {
		return nil, fmt.Errorf("io.Copy: %w", err)
	}
	return buf.Bytes(), nil
}

// Path corresponds to the path inside of storage bucket.
func (g *GCPStore) Write(ctx context.Context, object string, p []byte) (n int64, err error) {
	object = removeLocalPath(object, os.Getenv("WORKDIR"))
	bucket := os.Getenv("STORAGE_BUCKET")

	client, err := storage.NewClient(ctx)
	if err != nil {
		return 0, err
	}

	defer client.Close()

	buf := bytes.NewBuffer(p)

	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	gcpWtr := client.Bucket(bucket).Object(object).NewWriter(ctx)
	gcpWtr.ChunkSize = int(g.ChunkSize)
	if g.ProgressEnabled {
		gcpWtr.ProgressFunc = func(i int64) {
			// TODO: Enable redis ore something.
			println()
			fmt.Printf("Progress %d/%d\r", i, int64(buf.Len()))
		}
	}

	if n, err = io.Copy(gcpWtr, buf); err != nil {
		return 0, fmt.Errorf("io.Copy: %w", err)
	}

	// Data can continue to be added to the file until the writer is closed.
	if err := gcpWtr.Close(); err != nil {
		return 0, fmt.Errorf("Writer.Close: %w", err)
	}

	return n, nil
}

func (g *GCPStore) Delete(ctx context.Context, path string) error {
	// TODO Fill this out.
	// path = removeLocalPath(path, os.Getenv("WORKDIR"))
	return nil
}

func removeLocalPath(s, prefix string) string {
	if strings.Contains(s, prefix) {
		s = strings.TrimPrefix(s, prefix)
		if strings.HasPrefix(s, "/") {
			s = strings.TrimPrefix(s, "/")
		}
	}
	return s
}
