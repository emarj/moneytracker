package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type DateTime struct{ time.Time }

func (t *DateTime) Scan(v interface{}) error {

	var s string
	switch z := v.(type) {
	case []uint8:
		s = string(z)
	case string:
		s = z
	default:
		return errors.New("cannot convert time to string")
	}

	vt, err := time.Parse("2006-01-02T15:04:05", s)
	if err != nil {
		return err
	}
	t.Time = vt
	return nil
}

func (t DateTime) Value() (driver.Value, error) {
	return driver.Value(t.Format("2006-01-02T15:04:05")), nil
}

func (t *DateTime) UnmarshalJSON(json []byte) error {
	str := string(json[1 : len(json)-1])
	vt, err := time.Parse("2006-01-02T15:04:05", str)
	if err != nil {
		vt, err = time.Parse("2006-01-02", str)
		if err != nil {
			return err
		}
	}
	t.Time = vt
	return nil
}

func (t DateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Format("2006-01-02T15:04:05"))
}

type Date struct{ time.Time }

func (t *Date) Scan(v interface{}) error {

	var s string
	switch z := v.(type) {
	case []uint8:
		s = string(z)
	case string:
		s = z
	default:
		return errors.New("cannot convert time to string")
	}

	vt, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	t.Time = vt
	return nil
}

func (t Date) Value() (driver.Value, error) {
	return driver.Value(t.Format("2006-01-02")), nil
}

func (t *Date) UnmarshalJSON(json []byte) error {
	str := string(json[1 : len(json)-1])
	vt, err := time.Parse("2006-01-02", str)
	if err != nil {
		return err
	}
	t.Time = vt
	return nil
}

func (t Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Format("2006-01-02"))
}
