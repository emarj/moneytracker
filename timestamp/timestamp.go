package timestamp

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Trailing zeros are essential otherwise Go will remove them
// const TimestampFormat = "2006-01-02T15:04:05.999Z"
const TimestampFormat = "2006-01-02T15:04:05.000Z"

var nullBytes = []byte("null")

type Timestamp struct {
	sql.NullTime
}

func Now() Timestamp {
	return FromTime(time.Now())
}

func FromTime(t time.Time) Timestamp {
	return Timestamp{sql.NullTime{
		Time:  t,
		Valid: true,
	}}
}

func (t Timestamp) Plus(d time.Duration) Timestamp {
	return FromTime(t.Time.Add(d))
}

func (t Timestamp) Minus(d time.Duration) Timestamp {
	return t.Plus(-1 * d)
}

//JSON and SQL methods

func (t *Timestamp) Scan(v interface{}) error {

	var s string
	switch z := v.(type) {
	case []byte:
		s = string(z)
	case string:
		s = z
	default:
		return errors.New("cannot convert time to string")
	}

	vt, err := time.Parse(TimestampFormat, s)
	if err != nil {
		return err
	}
	t.Valid = true
	t.Time = vt
	return nil
}

func (t Timestamp) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return driver.Value(t.String()), nil
}

func (t Timestamp) String() string {
	if !t.Valid {
		return ""
	}
	return t.Time.UTC().Format(TimestampFormat)
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	t.Valid = false
	if bytes.Equal(data, nullBytes) {
		return nil
	}

	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	vt, err := time.Parse(TimestampFormat, str)
	if err != nil {
		return err
	}

	t.Valid = true
	t.Time = vt
	return nil
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return nullBytes, nil
	}
	return json.Marshal(t.String())
}
