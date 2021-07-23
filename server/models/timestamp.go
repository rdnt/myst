package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Timestamp wraps the Time struct offering UNIX timestamp encoding/decoding
type Timestamp struct {
	time.Time
}

var _ json.Marshaler = (*Timestamp)(nil)
var _ json.Unmarshaler = (*Timestamp)(nil)

// MarshalJSON implements the json.Marshaler interface
func (t *Timestamp) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("%d", t.Time.Unix())
	return []byte(stamp), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (t *Timestamp) UnmarshalJSON(data []byte) (err error) {
	ts, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	t.Time = time.Unix(ts, 0)
	return
}
