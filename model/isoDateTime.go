package model

import (
	"fmt"
	"time"
)

// ISODateTime represent a moment in time that serializes
// as an ISO (RFC3339 with nanoseconds) Date Time text representation.
type ISODateTime time.Time

// String implements fmt.Stringer.
func (dt ISODateTime) String() string {
	return time.Time(dt).Format(time.RFC3339Nano)
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (dt *ISODateTime) UnmarshalText(text []byte) error {
	parsed, err := time.Parse(time.RFC3339Nano, string(text))
	if err != nil {
		return fmt.Errorf("error parsing ISO8601DateTime: %v", err)
	}
	*dt = ISODateTime(parsed)
	return nil
}
