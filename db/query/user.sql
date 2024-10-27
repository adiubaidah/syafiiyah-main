-- name: CreateUser :one
INSERT INTO
    "user" ("role", "username", "password")
VALUES
    (@role::user_role, @username::text, @password::text) RETURNING *;

-- name: QueryUsersAscUsername :many
SELECT
    "user"."id",
    "user"."username",
    "user"."role",
    COALESCE("parent"."id", 0) AS "parentID",
    "parent"."name" AS "parentName",
    COALESCE("employee"."id", 0) AS "employeeID",
    "employee"."name" AS "employeeName"
FROM
    "user"
    LEFT JOIN "parent" ON "user"."id" = "parent"."user_id"
    LEFT JOIN "employee" ON "user"."id" = "employee"."user_id"
WHERE
    (
        sqlc.narg(q)::text IS NULL
        OR "user"."username" ILIKE '%' || sqlc.narg(q) || '%'
    )
    AND (
        sqlc.narg(role)::user_role IS NULL
        OR "user"."role" = sqlc.narg(role)
    )
    AND (
        @has_relation::smallint IS NULL
        OR (
            @has_relation = 1
            AND "parent"."id" IS NOT NULL OR "employee"."id" IS NOT NULL
        )
        OR (
            @has_relation = 0
            AND "parent"."id" IS NULL AND "employee"."id" IS NULL
        )
        OR (@has_relation = -1)
    )
ORDER BY
    "user"."username" ASC
LIMIT
    @limit_number OFFSET @offset_number;

-- name: QueryUsersDescUsername :many
SELECT
    "user"."id",
    "user"."username",
    "user"."role",
    "parent"."id" AS "parentID",
    "parent"."name" AS "parentName",
    "employee"."id" AS "employeeID",
    "employee"."name" AS "employeeName"
FROM
    "user"
    LEFT JOIN "parent" ON "user"."id" = "parent"."user_id"
    LEFT JOIN "employee" ON "user"."id" = "employee"."user_id"
WHERE
    (
        sqlc.narg(q)::text IS NULL
        OR "user"."username" ILIKE '%' || sqlc.narg(q) || '%'
    )
    AND (
        sqlc.narg(role)::text IS NULL
        OR "user"."role" = sqlc.narg(role)
    )
    AND (
        @has_relation::smallint IS NULL
        OR (
            @has_relation = 1
            AND "parent"."id" IS NOT NULL OR "employee"."id" IS NOT NULL
        )
        OR (
            @has_relation = 0
            AND "parent"."id" IS NULL AND "employee"."id" IS NULL
        )
        OR (@has_relation = -1)
    )
ORDER BY
    "user"."username" DESC
LIMIT
    @limit_number OFFSET @offset_number;

-- name: CountUsers :one
SELECT
    COUNT(*) AS "count"
FROM
    "user"
WHERE
    (
        sqlc.narg(q)::text IS NULL
        OR "username" ILIKE '%' || sqlc.narg(q) || '%'
    )
    AND (
        sqlc.narg(role)::text IS NULL
        OR "role" = sqlc.narg(role)
    )
    AND (
        @has_relation::smallint IS NULL
        OR (
            @has_relation = 1
            AND "parent_id" IS NOT NULL OR "employee_id" IS NOT NULL
        )
        OR (
            @has_relation = 0
            AND "parent_id" IS NULL AND "employee_id" IS NULL
        )
        OR (@has_relation = -1)
    );

-- name: GetUserByID :one
SELECT
    "user"."id",
    "user"."role",
    "user"."username"
FROM
    "user"
WHERE
    "id" = @id;

-- name: UpdateUser :one
UPDATE
    "user"
SET
    "role" = @role,
    "username" = @username,
    "password" = sqlc.narg(password)
WHERE
    "id" = @id
RETURNING *;

-- name: DeleteUser :one
DELETE FROM
    "user"
WHERE
    "id" = @id
RETURNING *;

