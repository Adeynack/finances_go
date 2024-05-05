package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type MdIntDate struct {
	Date time.Time
}

var (
	_ json.Unmarshaler = &MdIntDate{}
)

// UnmarshalJSON implements json.Unmarshaler.
func (m *MdIntDate) UnmarshalJSON(data []byte) error {
	var intValue int
	if err := json.Unmarshal(data, &intValue); err != nil {
		return fmt.Errorf("unable to unmarshal Moneydance integer date: %w", err)
	}
	year := intValue / 10_000
	month := time.Month((intValue / 100) % 100)
	day := intValue % 100

	m.Date = time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	return nil
}

func (m *MdIntDate) Year() int {
	return m.Date.Year()
}

func (m *MdIntDate) Month() time.Month {
	return m.Date.Month()
}

func (m *MdIntDate) Day() int {
	return m.Date.Day()
}
