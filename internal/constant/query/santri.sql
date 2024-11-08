-- name: CreateSantri :one
INSERT INTO
    "santri" (
        "nis",
        "name",
        "gender",
        "is_active",
        "generation",
        "photo",
        "occupation_id",
        "parent_id"
    )
VALUES
    (
        @nis,
        @name,
        @gender,
        @is_active,
        @generation,
        sqlc.narg(photo) :: text,
        @occupation_id,
        @parent_id
    ) RETURNING *;

-- ListSantri :many
SELECT
    "id",
    "name",
    "nis",
    "generation",
    "parent_id",
    "parent_name",
    "parent_whatsapp_number",
    "occupation_id",
    "occupation_name"
FROM
    list_santri(
        sqlc.narg(q),
        sqlc.narg(occupation_id),
        sqlc.narg(generation),
        @limit_number,
        @offset_number,
        sqlc.narg(order_by) :: santri_order_by
    ) AS "list_santri";

-- name: CountSantri :one
SELECT
    COUNT(*) AS "count"
FROM
    "santri"
    LEFT JOIN "parent" ON "santri".parent_id = "parent".id
    LEFT JOIN santri_occupation ON "santri".occupation_id = santri_occupation.id
WHERE
    (
        sqlc.narg(q)::text IS NULL
        OR "santri".name ILIKE '%' || sqlc.narg(q)::text || '%'
        OR "santri".nis ILIKE '%' || sqlc.narg(q)::text || '%'
    )
    AND (
        sqlc.narg(occupation_id)::integer IS NULL
        OR "santri".occupation_id = sqlc.narg(occupation_id)
    )
    AND (
        sqlc.narg(generation)::integer IS NULL
        OR "santri".generation = sqlc.narg(generation)
    )
    AND (
        sqlc.narg(is_active)::boolean IS NULL
        OR "santri".is_active = sqlc.narg(is_active)::boolean
    );

-- name: GetSantri :one
SELECT
    "santri".*,
    "parent"."id" AS "parent_id",
    "parent"."name" AS "parent_name",
    "parent"."whatsapp_number" AS "parent_whatsapp_number",
    "parent"."address" AS "parent_address",
    "santri_occupation"."name" AS "occupation_name"
FROM
    "santri"
    LEFT JOIN "parent" ON "santri"."parent_id" = "parent"."id"
    LEFT JOIN "santri_occupation" ON "santri"."occupation_id" = "santri_occupation"."id"
WHERE
    "santri"."id" = @id;

-- name: UpdateSantri :one
UPDATE
    "santri"
SET
    "nis" = @nis,
    "name" = @name,
    "generation" = @generation,
    "is_active" = @is_active :: boolean,
    "gender" = @gender::gender,
    "photo" = sqlc.narg(photo) :: text,
    "occupation_id" = @occupation_id,
    "parent_id" = sqlc.narg(parent_id) :: integer
WHERE
    "id" = @id RETURNING *;

-- name: DeleteSantri :one
DELETE FROM
    "santri"
WHERE
    "id" = @id RETURNING *;