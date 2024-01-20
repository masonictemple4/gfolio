package filestore

import (
	"os"
	"strings"
	"testing"
)

const curPath = "internal/filestore"

func TestPathing(t *testing.T) {
	t.Run("[Internal store] Freehand path expiriments", func(t *testing.T) {
		workingDir := os.Getenv("PWD")

		t.Logf("Working dir Pre-Correction: %s", workingDir)
		workingDir = strings.TrimSuffix(workingDir, curPath)

		t.Logf("Working dir Post-Correction: %s", workingDir)

		assetDir := os.Getenv("ASSET_DIR")

		if assetDir == "" {
			assetDir = "assets"
		}

		var root string
		t.Logf("Before: %s", assetDir)
		if !strings.Contains(workingDir, assetDir) {
			root = workingDir + strings.TrimPrefix(assetDir, "/")
		}
		t.Logf("After: %s", root)

	})
	t.Run("[Internal store] New Internal Store", func(t *testing.T) {
		err := os.Setenv("ASSET_DIR", "share")
		if err != nil {
			t.Errorf("there was an error setting the env variable: %v", err)
		}

		assetDir := os.Getenv("ASSET_DIR")

		if assetDir == "" {
			assetDir = "assets"
		}

		store, err := NewInternalStore(assetDir)
		if err != nil {
			t.Errorf("there was an error creating the internal store: %v", err)
		}

		// Remember to replace the curPath because we're running inside
		// the package directory.
		t.Logf("Store root Raw: %s", store.root)
		t.Logf("Store root: %s", strings.Replace(store.root, curPath+"/", "", 1))

	})

}

func TestInternalStore(t *testing.T) {

	t.Run("Test open", func(t *testing.T) {
		fp := "../../assets/blogs/test.md"
		size := 476

		ogFile, err := os.OpenFile(fp, os.O_RDONLY, 0755)
		if err != nil {
			t.Errorf("there was an error opening the file: %v", err)
		}

		ogFile.Close()

		ogStat, err := os.Stat(fp)
		if err != nil {
			t.Errorf("there was an error getting the file stats: %v", err)
		}

		if ogStat.Size() != int64(size) {
			t.Errorf("the file size was not correct, expected %d, got %d", size, ogStat.Size())
		}

		// 		_, err = NewInternalStore("./assets/blogs")
		// 		if err != nil {
		// 			t.Errorf("there was an error opening the store: %v", err)
		// 		}

	})

}
