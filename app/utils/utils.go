package app_utils

import (
	"fmt"
	"strconv"
	"time"
)

func StrToInt64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func MsToTime(ms string) (time.Time, error) {
	msInt, err := strconv.ParseInt(ms, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(0, msInt*int64(time.Millisecond)), nil
}

func ToUTCDateString(datetime time.Time) string {
	// return fmt.Sprintf("%s+00:00", datetime.UTC().Format("2006-01-02T15:04:05"))
	y, m, d := datetime.Date()
	datetime = time.Date(y, m, d, 0, 0, 0, 0, time.Local)
	return fmt.Sprintf("%s+00:00", datetime.UTC().Format("2006-01-02T15:04:05"))
}
