package util

import (
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// ParseTime parses a time string in "HH:MM:SS" format and returns a time.Time object with the current date.
func ParseTime(timeString string) (time.Time, error) {
	return time.Parse("15:04:05", timeString)
}

func ParseHHMMWithCurrentDate(timeString string) (time.Time, error) {
	if timeString == "" {
		return time.Time{}, errors.New("time string is empty")
	}

	parsedTime, err := time.Parse("15:04", timeString)
	if err != nil {
		return time.Time{}, err
	}

	now := time.Now()
	currentDate := now.Format("2006-01-02")
	fullTimeString := fmt.Sprintf("%s %s", currentDate, parsedTime.Format("15:04"))
	fullTime, err := time.Parse("2006-01-02 15:04", fullTimeString)
	if err != nil {
		return time.Time{}, err
	}

	return fullTime, nil
}

// ParseDate parses a date string in "YYYY-MM-DD" format and returns a time.Time object.
func ParseDate(dateString string) (time.Time, error) {
	return time.Parse("2006-01-02", dateString)
}

// ConvertToPgxTime converts a time.Time object to pgtype.Time with microseconds precision.
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
