-- Santri query function

CREATE TYPE santri_order_by AS ENUM (
    'asc:name',
    'asc:nis',
    'asc:generation',
    'desc:name',
    'desc:nis',
    'desc:generation'
);

CREATE OR REPLACE FUNCTION list_santri(
    q TEXT,
    occupation_id_param INTEGER,
    generation_param INTEGER,
    is_active_param BOOLEAN,
    limit_number INTEGER,
    offset_number INTEGER,
    order_by santri_order_by
) RETURNS TABLE (
    id INTEGER,
    name TEXT,
    gender gender_type,
    nis TEXT,
    generation INTEGER,
    is_active BOOLEAN,
    photo TEXT,
    parent_id INTEGER,
    parent_name TEXT,
    parent_whatsapp_number TEXT,
    occupation_id INTEGER,
    occupation_name TEXT
) AS $$
DECLARE
    order_column TEXT := 'name';
    order_direction TEXT := 'ASC';
BEGIN
    IF order_by = 'asc:name' THEN
        order_column := 'name';
        order_direction := 'ASC';
    ELSIF order_by = 'asc:nis' THEN
        order_column := 'nis';
        order_direction := 'ASC';
    ELSIF order_by = 'asc:generation' THEN
        order_column := 'generation';
        order_direction := 'ASC';
    ELSIF order_by = 'desc:name' THEN
        order_column := 'name';
        order_direction := 'DESC';
    ELSIF order_by = 'desc:nis' THEN
        order_column := 'nis';
        order_direction := 'DESC';
    ELSIF order_by = 'desc:generation' THEN
        order_column := 'generation';
        order_direction := 'DESC';
    END IF;

    RETURN QUERY EXECUTE format(
        $query$
        SELECT
            santri.id,
            santri.name::text,
            santri.gender::gender_type,
            santri.nis::text,
            santri.generation,
            santri.is_active::boolean,
            santri.photo::text,
            parent.id AS parent_id,
            parent.name::text AS parent_name,
            parent.whatsapp_number::text AS parent_whatsapp_number,
            santri_occupation.id AS occupation_id,
            santri_occupation.name::text AS occupation_name
        FROM
            santri
            LEFT JOIN parent ON santri.parent_id = parent.id
            LEFT JOIN santri_occupation ON santri.occupation_id = santri_occupation.id
        WHERE
            ($1 IS NULL OR santri.name ILIKE '%%' || $1 || '%%' OR santri.nis ILIKE '%%' || $1 || '%%')
            AND ($2 IS NULL OR santri.occupation_id = $2)
            AND ($3 IS NULL OR santri.generation = $3)
            AND ($4 IS NULL OR santri.is_active = $4)
        ORDER BY santri.%I %s
        LIMIT $5 OFFSET $6
        $query$,
        order_column,
        order_direction
    )
    USING q, occupation_id_param, generation_param, is_active_param, limit_number, offset_number;
END;
$$ LANGUAGE plpgsql;

-- User query function

CREATE TYPE user_order_by AS ENUM (
    'asc:username',
    'desc:username',
    'asc:name',
    'desc:name'
);

CREATE OR REPLACE FUNCTION list_user(
    q TEXT,
    role_param role_type,
    has_owner BOOLEAN,
    limit_number INTEGER,
    offset_number INTEGER,
    order_by user_order_by
) RETURNS TABLE (
    id INTEGER,
    username TEXT,
    role role_type,
    id_owner INTEGER,
    name_owner TEXT
) AS $$
DECLARE
    order_column TEXT := 'username';
    order_direction TEXT := 'ASC';
BEGIN
    -- Tentukan kolom pengurutan berdasarkan nilai dari `order_by`
    IF order_by = 'asc:username' THEN
        order_column := 'username';
        order_direction := 'ASC';
    ELSIF order_by = 'desc:username' THEN
        order_column := 'username';
        order_direction := 'DESC';
    ELSIF order_by = 'asc:name' THEN
        order_column := 'COALESCE(parent.name, employee.name)';
        order_direction := 'ASC';
    ELSIF order_by = 'desc:name' THEN
        order_column := 'COALESCE(parent.name, employee.name)';
        order_direction := 'DESC';
    END IF;

    RETURN QUERY EXECUTE format(
        $query$
        SELECT
            "user".id,
            "user".username::text,
            "user".role::role_type,
            COALESCE(parent.id, employee.id) AS id_owner,
            COALESCE(parent.name::text, employee.name::text) AS name_owner
        FROM
            "user"
            LEFT JOIN parent ON "user".id = parent.user_id
            LEFT JOIN employee ON "user".id = employee.user_id
        WHERE
            ($1 IS NULL OR "user".username ILIKE '%%' || $1 || '%%')
            AND ($2 IS NULL OR "user"."role" = $2)
            AND (
                $3 IS NULL
                OR ($3 = TRUE AND (parent.id IS NOT NULL OR employee.id IS NOT NULL))
                OR ($3 = FALSE AND parent.id IS NULL AND employee.id IS NULL)
                )
        ORDER BY %s %s
        LIMIT $4 OFFSET $5
        $query$,
        order_column,
        order_direction
    )
    USING q,role_param,has_owner, limit_number, offset_number;
END;
$$ LANGUAGE plpgsql;

-- OR ($3 = TRUE AND ("employee".id IS NOT NULL OR "parent".id IS NOT NULL))
-- OR ($3 = FALSE AND "employee".id IS NULL AND "parent".id IS NULL)

-- Parent query function

CREATE TYPE parent_order_by AS ENUM (
    'asc:name',
    'desc:name'
);

CREATE OR REPLACE FUNCTION list_parent(
    q TEXT,
    has_user BOOLEAN,
    limit_number INTEGER,
    offset_number INTEGER,
    order_by parent_order_by
) RETURNS TABLE (
    id INTEGER,
    name TEXT,
    address TEXT,
    gender gender_type,
    whatsapp_number TEXT,
    photo TEXT,
    user_id INTEGER,
    username TEXT
) AS $$
DECLARE
    order_column TEXT := 'name';
    order_direction TEXT := 'ASC';
BEGIN
    IF order_by = 'asc:name' THEN
        order_column := 'name';
        order_direction := 'ASC';
    ELSIF order_by = 'desc:name' THEN
        order_column := 'name';
        order_direction := 'DESC';
    END IF;

    RETURN QUERY EXECUTE format(
        $query$
        SELECT
            parent.id,
            parent.name::text,
            parent.address::text,
            parent.gender::gender_type,
            parent.whatsapp_number::text,
            parent.photo::text,
            "user".id AS user_id,
            "user".username::text
        FROM
            parent
            LEFT JOIN "user" ON parent.user_id = "user".id
        WHERE
            ($1 IS NULL OR parent.name ILIKE '%%' || $1 || '%%')
            AND (
                $2 IS NULL
                OR ($2 = TRUE AND "user".id IS NOT NULL)
                OR ($2 = FALSE AND "user".id IS NULL)
            )
        ORDER BY %I %s
        LIMIT $3 OFFSET $4
        $query$,
        order_column,
        order_direction
    )
    USING q, has_user, limit_number, offset_number;
END;
$$ LANGUAGE plpgsql;