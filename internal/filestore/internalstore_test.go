package filestore

import (
	"os"
	"testing"
)

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
