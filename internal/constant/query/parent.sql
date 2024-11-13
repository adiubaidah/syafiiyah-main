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
        sqlc.narg(has_user)::boolean IS NULL
        OR (sqlc.narg(has_user) = TRUE AND "parent"."user_id" IS NOT NULL)
        OR (sqlc.narg(has_user) = FALSE AND "parent"."user_id" IS NULL)
        )
    ;

-- name: CreateParent :one
INSERT INTO
    "parent" (
        "name",
        "address",
        "gender",
        "whatsapp_number",
        "photo",
        "user_id"
    )
VALUES
    (
        @name,
        @address,
        @gender,
        @whatsapp_number,
        sqlc.narg(photo),
        sqlc.narg(user_id)
    ) RETURNING *;

-- name: UpdateParent :one
UPDATE
    "parent"
SET
    "name" = COALESCE(sqlc.narg(name), name),
    "address" = COALESCE(sqlc.narg(address), address),
    "gender" = COALESCE(sqlc.narg(gender)::gender, gender),
    "whatsapp_number" = COALESCE(sqlc.narg(whatsapp_number), whatsapp_number),
    "photo" = COALESCE(sqlc.narg(photo), photo),
    "user_id" = sqlc.narg(user_id)
WHERE
    "id" = @id RETURNING *;

-- name: GetParent :one
SELECT
    "parent".*,
    "user"."id" AS "user_id",
    "user"."username" AS "user_username"
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