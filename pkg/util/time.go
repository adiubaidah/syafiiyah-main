package util

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ParseTime(timeString string) (time.Time, error) {
	return time.Parse("15:04:05", timeString)
}

func ConvertToPgxTime(parsedTime time.Time) pgtype.Time {

	microseconds := int64(parsedTime.Hour()*3600+parsedTime.Minute()*60+parsedTime.Second()) * 1e6
	return pgtype.Time{
		Microseconds: microseconds,
		Valid:        true,
	}
}

func ConvertToTime(pgxTime pgtype.Time) string {
	seconds := pgxTime.Microseconds / 1e6
	nanoseconds := (pgxTime.Microseconds % 1e6) * 1e3
	t := time.Unix(seconds, nanoseconds)
	return t.Format("15:04:05")
}
