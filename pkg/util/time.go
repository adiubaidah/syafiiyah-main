package util

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

var loc *time.Location

func ConvertToPgxTime(parsedTime time.Time) pgtype.Time {

	microseconds := int64(parsedTime.Hour()*3600+parsedTime.Minute()*60+parsedTime.Second()) * 1e6
	return pgtype.Time{
		Microseconds: microseconds,
		Valid:        true,
	}
}
