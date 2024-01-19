package customdate

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestDefaultDate(t *testing.T) {

	type Person struct {
		Name    string      `json:"name"`
		Created DefaultDate `json:"created"`
	}

	t.Run("marshal RFC3339 format should work", func(t *testing.T) {
		now := time.Now()

		val := &Person{
			Name:    "mason",
			Created: DefaultDate(now),
		}

		data, err := json.Marshal(val)
		if err != nil {
			t.Errorf("expected json.marshal to work: %v", err)
		}

		fmt.Printf("The time data: %s", string(data))

	})

}
