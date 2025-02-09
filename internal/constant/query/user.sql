-- name: CreateUser :one
INSERT INTO
    "user" ("role","email", "username", "password")
VALUES
    (
        @role :: role_type,
        @email :: text,
        @username :: text,
        @password :: text
    ) RETURNING *;

-- name: CountUsers :one
SELECT
    COUNT(*) AS "count"
FROM
    "user"
    LEFT JOIN "parent" ON "user"."id" = "parent"."user_id"
    LEFT JOIN "employee" ON "user"."id" = "employee"."user_id"
WHERE
    (
        sqlc.narg(q) :: text IS NULL
        OR "username" ILIKE '%' || sqlc.narg(q) || '%'
    )
    AND (
        sqlc.narg(role) :: role_type  IS NULL
        OR "role" = sqlc.narg(role) :: role_type
    )
    AND (
        sqlc.narg(has_owner) :: boolean IS NULL
        OR (
            sqlc.narg(has_owner) = TRUE
            AND (
                parent.id IS NOT NULL
                OR employee.id IS NOT NULL
            )
        )
        OR (
            sqlc.narg(has_owner) = FALSE
            AND parent.id IS NULL
            AND employee.id IS NULL
        )
    );

-- name: GetUserById :one
SELECT
    "user".*,
    CASE
        WHEN "parent"."id" IS NOT NULL THEN "parent"."id"
        WHEN "employee"."id" IS NOT NULL THEN "employee"."id"
        ELSE NULL
    END AS "owner_id"
FROM
    "user"
LEFT JOIN "parent" ON "user"."id" = "parent"."user_id"
LEFT JOIN "employee" ON "user"."id" = "employee"."user_id"
WHERE
    (
        sqlc.narg(id)::integer IS NOT NULL
        AND "user"."id" = sqlc.narg(id)::integer
    )
LIMIT
    1;

-- name: GetUserByUsername :one
SELECT
    "user".*,
    CASE
        WHEN "parent"."id" IS NOT NULL THEN "parent"."id"
        WHEN "employee"."id" IS NOT NULL THEN "employee"."id"
        ELSE NULL
    END AS "owner_id"
FROM
    "user"
LEFT JOIN "parent" ON "user"."id" = "parent"."user_id"
LEFT JOIN "employee" ON "user"."id" = "employee"."user_id"
WHERE
    (
        sqlc.narg(username)::text IS NOT NULL
        AND "user"."username" = sqlc.narg(username)::text
    )
LIMIT
    1;

-- name: GetUserByEmail :one
SELECT
    "user".*,
    CASE
        WHEN "parent"."id" IS NOT NULL THEN "parent"."id"
        WHEN "employee"."id" IS NOT NULL THEN "employee"."id"
        ELSE NULL
    END AS "owner_id"
FROM
    "user"
LEFT JOIN "parent" ON "user"."id" = "parent"."user_id"
LEFT JOIN "employee" ON "user"."id" = "employee"."user_id"
WHERE
    (
        sqlc.narg(email)::text IS NOT NULL
        AND "user"."email" = sqlc.narg(email)::text
    )
LIMIT
    1;


-- name: UpdateUser :one
UPDATE
    "user"
SET
    "role" = COALESCE(sqlc.narg(role)::role_type, "role"),
    "email" = COALESCE(sqlc.narg(email), "email"),
    "username" = COALESCE(sqlc.narg(username), "username"),
    "password" = COALESCE(sqlc.narg(password), "password")
WHERE
    "id" = @id RETURNING *;

-- name: DeleteUser :one
DELETE FROM
    "user"
WHERE
    "id" = @id RETURNING *;