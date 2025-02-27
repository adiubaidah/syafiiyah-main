-- name: CreateEmployee :one
INSERT INTO
    "employee" (
        "nip",
        "name",
        "gender",
        "photo",
        "occupation_id",
        "user_id"
    )
VALUES
    (
        @nip,
        @name,
        @gender,
        sqlc.narg(photo) :: text,
        @occupation_id,
        sqlc.narg(user_id) :: integer
    ) RETURNING *;

-- name: CountEmployees :one

SELECT
    COUNT(*)
FROM
    employee
    LEFT JOIN "user" ON employee.user_id = "user".id
    LEFT JOIN employee_occupation ON employee.occupation_id = employee_occupation.id
WHERE
    (
        sqlc.narg(q)::text IS NULL
        OR employee.name ILIKE '%%' || sqlc.narg(q)::text || '%%'
        OR employee.nip ILIKE '%%' || sqlc.narg(q)::text || '%%'
    )
    AND (
        sqlc.narg(occupation_id)::int IS NULL
        OR employee.occupation_id = sqlc.narg(occupation_id)::int
    )
    AND (
        sqlc.narg(has_user)::boolean IS NULL
        OR (
            sqlc.narg(has_user)::boolean = TRUE
            AND "user".id IS NOT NULL
        )
        OR (
            sqlc.narg(has_user)::boolean = FALSE
            AND "user".id IS NULL
        )
    );

-- name: UpdateEmployee :one
UPDATE
    "employee"
SET
    "nip" = sqlc.narg(nip),
    "name" = COALESCE(sqlc.narg(name), name),
    "gender" = COALESCE(sqlc.narg(gender) :: gender_type, gender),
    "photo" = sqlc.narg(photo),
    "occupation_id" = COALESCE(sqlc.narg(occupation_id), occupation_id),
    "user_id" = sqlc.narg(user_id)
WHERE
    "id" = @id RETURNING *;

-- name: GetEmployeeByID :one
SELECT
    "employee".*,
    "user"."id" AS "userId",
    "user"."username" AS "userUsername"
FROM
    "employee"
    LEFT JOIN "user" ON "employee"."user_id" = "user"."id"
WHERE
    "employee"."id" = @id;

-- name: GetEmployeeByUserID :one
SELECT
    "employee".*
FROM
    "employee"
WHERE
    "user_id" = @user_id;

-- name: DeleteEmployee :one
DELETE FROM
    "employee"
WHERE
    "id" = @id RETURNING *;