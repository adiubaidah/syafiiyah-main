-- name: CreateUser :one
INSERT INTO
    "user" ("role", "username", "password")
VALUES
    (
        @role :: role_type,
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
        OR "role" = sqlc.narg(role)
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

-- name: GetUser :one
SELECT
    "user"."id",
    "user"."role",
    "user"."username",
    "user"."password"
FROM
    "user"
WHERE
    (
        sqlc.narg(id)::integer IS NOT NULL
        AND "id" = sqlc.narg(id)::integer
    )
    OR (
        sqlc.narg(username)::text IS NOT NULL
        AND "username" = sqlc.narg(username)::text
    )
LIMIT
    1;

-- name: UpdateUser :one
UPDATE
    "user"
SET
    "role" = COALESCE(sqlc.narg(role)::role_type, "role"),
    "username" = COALESCE(sqlc.narg(username), "username"),
    "password" = COALESCE(sqlc.narg(password), "password")
WHERE
    "id" = @id RETURNING *;

-- name: DeleteUser :one
DELETE FROM
    "user"
WHERE
    "id" = @id RETURNING *;