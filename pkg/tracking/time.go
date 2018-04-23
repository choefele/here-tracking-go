package tracking

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(data []byte) error {
	millis, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	*t = Time{time.Unix(0, millis*int64(time.Millisecond))}
	return nil
}

func (t Time) MarshalJSON() ([]byte, error) {
	if t.Unix() < 0 {
		return nil, errors.New("Time must be after 1 January 1970 00:00:00 UTC")
	}
	time := fmt.Sprintf("%v", t.UnixNano()/int64(time.Millisecond))
	return []byte(time), nil
}
