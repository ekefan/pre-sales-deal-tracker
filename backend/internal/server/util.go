package server 

import (
	"time"
	"encoding/json"
)

// Custom type for Unix timestamp
type UnixTime struct {
	Time  time.Time
	Valid bool
}


// UnmarshalJSON for UnixTime
func (ut *UnixTime) UnmarshalJSON(data []byte) error {
	var timestamp int64
	if err := json.Unmarshal(data, &timestamp); err != nil {
		return err
	}

	if timestamp == 0 {
		ut.Valid = false
		return nil
	}

	ut.Time = time.Unix(timestamp, 0).UTC()
	ut.Valid = true
	return nil
}