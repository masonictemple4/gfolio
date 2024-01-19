package utils

import "encoding/json"

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
