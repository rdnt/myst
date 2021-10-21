package timestamp

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"strconv"
	"time"
)

var (
	ErrInvalidType = fmt.Errorf("invalid type")
)

// Timestamp wraps the Time struct offering UNIX timestamp encoding/decoding
type Timestamp struct {
	time.Time
}

// New returns a timestamp with the current time
func New() Timestamp {
	return Timestamp{time.Now()}
}

// assert timestamp has custom json marshaler/unmarshaler on compile-time
var _ json.Marshaler = (*Timestamp)(nil)
var _ json.Unmarshaler = (*Timestamp)(nil)

// MarshalJSON implements the json.Marshaler interface
func (t Timestamp) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("%d", t.Time.UnixNano())
	return []byte(stamp), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (t *Timestamp) UnmarshalJSON(data []byte) (err error) {
	ts, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	t.Time = time.Unix(0, ts)
	return
}

// assert timestamp has custom bson value marshaler/unmarshaler on compile-time
var _ bson.ValueMarshaler = (*Timestamp)(nil)
var _ bson.ValueUnmarshaler = (*Timestamp)(nil)

// MarshalBSONValue implements the bson.ValueMarshaler interface
func (t Timestamp) MarshalBSONValue() (bsontype.Type, []byte, error) {
	typ, b, err := bson.MarshalValue(t.Time.UnixNano())
	if err != nil {
		return typ, nil, err
	}
	return typ, b, nil
}

// UnmarshalBSONValue implements the bson.ValueUnmarshaler interface
func (t *Timestamp) UnmarshalBSONValue(typ bsontype.Type, b []byte) error {
	if typ != bsontype.Int64 {
		return ErrInvalidType
	}

	i := binary.LittleEndian.Uint64(b)
	t.Time = time.Unix(0, int64(i))

	return nil
}
