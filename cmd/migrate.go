package cmd

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/masonictemple4/masonictempl/internal/filestore"
	"github.com/spf13/cobra"
	"google.golang.org/api/iterator"
)

type migrateArg int

const (
	src migrateArg = iota
	dst
)

var migrateToRemoteCmd = &cobra.Command{
	Use:   "migrate [src] [dst]",
	Short: "Migrate assets from one dest to another.",
	Long: `Migrate local assets to remote. For example:
masonictempl assets migrate local remote or masonictempl remote local. (Eligible values for src and dst are:
local remote)
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		assetsCmd.PersistentPreRun(assetsCmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			log.Fatal("masonictempl assets migrate [src] [dst]")
		}

		if err := runMigrate(cmd.Context(), args[src], args[dst]); err != nil {
			log.Fatalf("[assets migrate]: %v\n", err)
		}

	},
}

// TODO: Add ability to compare local and remote to remove
// remote assets that no longer exist locally.
func runMigrate(ctx context.Context, src, dst string) error {
	if src == dst {
		return errors.New("src and dst cannot be the same")
	}

	switch {
	case src == dst:
		return errors.New("src and dst cannot be the same")
	case src == "local":
		if err := migrateFromLocalToRemote(ctx); err != nil {
			return err
		}
	case src == "remote":
		if err := migrateFromRemoteToLocal(ctx); err != nil {
			return err
		}
	default:
		return fmt.Errorf("dst must be either local or remote, got %s", dst)
	}

	return nil
}

func migrateFromLocalToRemote(ctx context.Context) error {

	remote, err := filestore.NewFilestore("gcp", "")
	if err != nil {
		return err
	}

	local, err := filestore.NewFilestore("internal", "")
	if err != nil {
		return err
	}

	fmt.Printf("remote: %s\n", filestore.GetRootPath(remote))
	fmt.Printf("local: %s\n", filestore.GetRootPath(local))

	filepath.WalkDir(filestore.GetRootPath(local), func(path string, file fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !file.IsDir() {
			// only need to write on file.
			// Gcp will take care of missing dirs.
			fmt.Println("file: ", path)
			fData, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			wtnCnt, err := remote.Write(ctx, path, fData)
			if err != nil {
				return err
			}

			fmt.Printf("wrote %d bytes to %s\n", wtnCnt, strings.TrimPrefix(path, os.Getenv("WORKDIR")))
		}

		return nil
	})

	return nil
}

// TODO: Add ability to compare remote and local to remove
// local assets that no longer exist in remote..
func migrateFromRemoteToLocal(ctx context.Context) error {
	bucket := os.Getenv("STORAGE_BUCKET")

	remote, err := filestore.NewFilestore("gcp", "")
	if err != nil {
		return err
	}

	local, err := filestore.NewFilestore("internal", "")
	if err != nil {
		return err
	}

	fmt.Printf("remote: %s\n", filestore.GetRootPath(remote))
	fmt.Printf("local: %s\n", filestore.GetRootPath(local))

	// bucket := "bucket-name"
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	it := client.Bucket(bucket).Objects(ctx, nil)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("Bucket(%q).Objects: %w", bucket, err)
		}

		data, err := remote.Read(ctx, attrs.Name)
		if err != nil {
			return err
		}

		wtnCnt, err := local.Write(ctx, attrs.Name, data)
		if err != nil {
			return err
		}
		localPath := filepath.Join(attrs.Name, os.Getenv("WORKDIR"))

		fmt.Printf("wrote %d bytes to %s\n", wtnCnt, localPath)
	}
	return nil

}
