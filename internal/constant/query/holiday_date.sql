-- name: CreateHolidayDates :copyfrom
INSERT INTO
    "holiday_date" ("date", "holiday_id")
VALUES
    (@date, @holiday_id);

-- name: DeleteHolidayDateByHolidayId :exec
DELETE FROM
    "holiday_date"
WHERE
    "holiday_id" = @holiday_id;