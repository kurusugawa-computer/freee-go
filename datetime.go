package freee

import (
	"encoding/json"
	"time"
)

type DateTime time.Time
type Date time.Time
type Time time.Time

func NewDateTime(year int, month time.Month, day int, hour int, min int, sec int) *DateTime {
	d := DateTime(time.Date(year, month, day, hour, min, sec, 0, time.UTC))
	return &d
}

func (d *DateTime) MarshalJSON() ([]byte, error) {
	s := ""
	if d != nil {
		t := time.Time(*d)
		s = t.Format("2006-01-02 15:04:05")
	}
	return json.Marshal(s)
}

func (d *DateTime) String() string {
	s := ""
	if d != nil {
		t := time.Time(*d)
		s = t.Format("2006-01-02 15:04:05")
	}
	return s
}

func NewDate(year int, month time.Month, day int) *Date {
	d := Date(time.Date(year, month, day, 0, 0, 0, 0, time.UTC))
	return &d
}

func (d *Date) MarshalJSON() ([]byte, error) {
	s := ""
	if d != nil {
		t := time.Time(*d)
		s = t.Format("2006-01-02")
	}
	return json.Marshal(s)
}

func (d *Date) String() string {
	s := ""
	if d != nil {
		t := time.Time(*d)
		s = t.Format("2006-01-02")
	}
	return s
}

func NewTime(hour int, min int, sec int) *DateTime {
	d := DateTime(time.Date(0, 0, 0, hour, min, sec, 0, time.UTC))
	return &d
}

func (t *Time) MarshalJSON() ([]byte, error) {
	s := ""
	if t != nil {
		u := time.Time(*t)
		s = u.Format("15:04:05")
	}
	return json.Marshal(s)
}

func (d *Time) String() string {
	s := ""
	if d != nil {
		t := time.Time(*d)
		s = t.Format("15:04:05")
	}
	return s
}
