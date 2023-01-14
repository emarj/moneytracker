package datetime

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

const DateTimeFormat = "2006-01-02T15:04:05.999Z"

var nullBytes = []byte("null")

type DateTime struct {
	sql.NullTime
}

func Now() DateTime {
	return FromTime(time.Now())
}

func FromTime(t time.Time) DateTime {
	return DateTime{sql.NullTime{
		Time:  t,
		Valid: true,
	}}
}

func (t DateTime) Plus(d time.Duration) DateTime {
	return FromTime(t.Time.Add(d))
}

func (t DateTime) Minus(d time.Duration) DateTime {
	return t.Plus(-1 * d)
}

//JSON and SQL methods

func (t *DateTime) Scan(v interface{}) error {

	var s string
	switch z := v.(type) {
	case []byte:
		s = string(z)
	case string:
		s = z
	default:
		return errors.New("cannot convert time to string")
	}

	vt, err := time.Parse(DateTimeFormat, s)
	if err != nil {
		return err
	}
	t.Valid = true
	t.Time = vt
	return nil
}

func (t DateTime) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return driver.Value(t.String()), nil
}

func (t DateTime) String() string {
	if !t.Valid {
		return ""
	}
	return t.Time.UTC().Format(DateTimeFormat)
}

func (t *DateTime) UnmarshalJSON(data []byte) error {
	t.Valid = false
	if bytes.Equal(data, nullBytes) {
		return nil
	}

	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	vt, err := time.Parse(DateTimeFormat, str)
	if err != nil {
		return err
	}

	t.Valid = true
	t.Time = vt
	return nil
}

func (t DateTime) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return nullBytes, nil
	}
	return json.Marshal(t.String())
}
