// The custom date package allows us to setup
// the default date format we'll be using across
// our blog application.
// This gives us control over the format we receive date information in
// and the format we return it in.
package customdate

import (
	"encoding/json"
	"strings"
	"time"
)

var _ json.Marshaler = (*DefaultDate)(nil)
var _ json.Unmarshaler = (*DefaultDate)(nil)

// DefaultDate implements both he Marshaler and
// Unmarshaler interfaces. Allowing us to control
// the format in which we receive and return
// time.
type DefaultDate time.Time // RFC3339

func (d *DefaultDate) UnmarshalJSON(b []byte) error {

	value := strings.ToUpper(strings.Trim(string(b), `"`))
	if value == "" || value == "null" {
		return nil
	}

	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return err
	}

	*d = DefaultDate(t)

	return nil
}

func (d *DefaultDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(*d))
}
