package util

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ConvertToPgxTime(timeStr string) (pgtype.Time, error) {
	parsedTime, err := time.Parse("15:04", timeStr)
	if err != nil {
		return pgtype.Time{}, err
	}

	microseconds := int64(parsedTime.Hour()*3600+parsedTime.Minute()*60) * 1e6
	return pgtype.Time{
		Microseconds: microseconds,
		Valid:        true,
	}, nil
}
