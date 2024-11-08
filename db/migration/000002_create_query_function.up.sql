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
    gender gender,
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
            santri.gender::gender,
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
