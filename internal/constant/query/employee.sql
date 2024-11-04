-- name: ListEmployeesAsc :many
SELECT
    "employee".*,
    "user"."id" AS "userId",
    "user"."username" AS "userUsername"
FROM
    "employee"
    LEFT JOIN "user" ON "employee"."user_id" = "user"."id"
WHERE
    (
        sqlc.narg(q)::text IS NULL
        OR "name" ILIKE '%' || sqlc.narg(q) || '%'
    )
    AND (
        @has_user :: smallint IS NULL
        OR (
            @has_user = 1
            AND "user_id" IS NOT NULL
        )
        OR (
            @has_user = 0
            AND "user_id" IS NULL
        )
        OR (@has_user = -1)
    )
ORDER BY
    "name" ASC
LIMIT
    @limit_number OFFSET @offset_number;

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
        sqlc.narg(photo)::text,
        @occupation_id,
        sqlc.narg(user_id)::integer
    ) RETURNING *;

-- name: UpdateEmployee :one
UPDATE
    "employee"
SET
    "nip" = @nip,
    "name" = @name,
    "gender" = @gender,
    "photo" = sqlc.narg(photo),
    "occupation_id" = @occupation_id,
    "user_id" = sqlc.narg(user_id)
WHERE
    "id" = @id RETURNING *;

-- name: GetEmployee :one
SELECT
    "employee".*,
    "user"."id" AS "userId",
    "user"."username" AS "userUsername"
FROM
    "employee"
    LEFT JOIN "user" ON "employee"."user_id" = "user"."id"
WHERE
    "employee"."id" = @id;

-- name: DeleteEmployee :one
DELETE FROM
    "employee"
WHERE
    "id" = @id
RETURNING *;
