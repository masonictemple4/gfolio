package utils

import (
	"encoding/json"
	"os"
)

func Convert[S any, T any](source S, target T) error {
	tmp, err := json.Marshal(source)
	if err != nil {
		return err
	}
	err = json.Unmarshal(tmp, &target)
	if err != nil {
		return err
	}
	return nil
}

func FileExists(fn string) bool {
	fInfo, err := os.Stat(fn)
	return !os.IsNotExist(err) && !fInfo.IsDir()
}
