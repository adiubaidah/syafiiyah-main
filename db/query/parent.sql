-- name: QueryParentsAsc :many
SELECT
    "parent".*,
    "user"."id" AS "userId",
    "user"."username" AS "userUsername"
FROM
    "parent"
    LEFT JOIN "user" ON "parent"."user_id" = "user"."id"
WHERE
    (
        sqlc.narg(q) :: text IS NULL
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

-- name: QueryParentsDesc :many
SELECT
    "parent".*,
    "user"."id" AS "userId",
    "user"."username" AS "userUsername"
FROM
    "parent"
    LEFT JOIN "user" ON "parent"."user_id" = "user"."id"
WHERE
    (
        sqlc.narg(q) :: text IS NULL
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

-- name: CountParents :one
SELECT
    COUNT(*) AS "count"
FROM
    "parent"
WHERE
    (
        sqlc.narg(q) :: text IS NULL
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
    );

-- name: CreateParent :one
INSERT INTO
    "parent" (
        "name",
        "address",
        "gender",
        "no_wa",
        "photo",
        "user_id"
    )
VALUES
    (
        @name,
        @address,
        @gender,
        @no_wa,
        sqlc.narg(photo),
        sqlc.narg(user_id)
    ) RETURNING *;

-- name: UpdateParent :one
UPDATE
    "parent"
SET
    "name" = @name,
    "address" = @address,
    "gender" = @gender,
    "no_wa" = @no_wa,
    "photo" = sqlc.narg(photo),
    "user_id" = sqlc.narg(user_id)
WHERE
    "id" = @id RETURNING *;

-- name: GetParent :one
SELECT
    "parent".*,
    "user"."id" AS "userId",
    "user"."username" AS "userUsername"
FROM
    "parent"
    LEFT JOIN "user" ON "parent"."user_id" = "user"."id"
WHERE
    "parent"."id" = @id;

-- name: DeleteParent :one
DELETE FROM
    "parent"
WHERE
    "id" = @id RETURNING *;