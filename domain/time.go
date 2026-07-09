package domain

import "time"

const (
	DATE_FORMAT          = "2006-01-02T15:04:05"
	FILENAME_DATE_FORMAT = "2006_01_02_150405"
	TEXT_DATE_FORMAT     = "2006-01-02 15:04:05"
)

func NanosecondsToTime(nanoseconds int64) time.Time {
	converted := time.Unix(0, nanoseconds).UTC()
	return converted
}
