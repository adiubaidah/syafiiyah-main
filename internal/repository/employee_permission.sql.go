// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: employee_permission.sql

package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createEmployeePermission = `-- name: CreateEmployeePermission :one
INSERT INTO
    "employee_permission" (
        employee_id,
        start_permission,
        end_permission,
        "type",
        excuse
    )
VALUES
    (
        $1,
        $2,
        $3,
        $4 :: santri_permission_type,
        $5
    ) RETURNING id, employee_id, type, start_permission, end_permission, excuse
`

type CreateEmployeePermissionParams struct {
	EmployeeID      int32              `db:"employee_id"`
	StartPermission pgtype.Timestamptz `db:"start_permission"`
	EndPermission   pgtype.Timestamptz `db:"end_permission"`
	Type            interface{}        `db:"type"`
	Excuse          string             `db:"excuse"`
}

func (q *Queries) CreateEmployeePermission(ctx context.Context, arg CreateEmployeePermissionParams) (EmployeePermission, error) {
	row := q.db.QueryRow(ctx, createEmployeePermission,
		arg.EmployeeID,
		arg.StartPermission,
		arg.EndPermission,
		arg.Type,
		arg.Excuse,
	)
	var i EmployeePermission
	err := row.Scan(
		&i.ID,
		&i.EmployeeID,
		&i.Type,
		&i.StartPermission,
		&i.EndPermission,
		&i.Excuse,
	)
	return i, err
}

const deleteEmployeePermission = `-- name: DeleteEmployeePermission :one
DELETE FROM
    "employee_permission"
WHERE
    "id" = $1 RETURNING id, employee_id, type, start_permission, end_permission, excuse
`

func (q *Queries) DeleteEmployeePermission(ctx context.Context, id int32) (EmployeePermission, error) {
	row := q.db.QueryRow(ctx, deleteEmployeePermission, id)
	var i EmployeePermission
	err := row.Scan(
		&i.ID,
		&i.EmployeeID,
		&i.Type,
		&i.StartPermission,
		&i.EndPermission,
		&i.Excuse,
	)
	return i, err
}

const getEmployeePermission = `-- name: GetEmployeePermission :one
SELECT
    employee_permission.id, employee_permission.employee_id, employee_permission.type, employee_permission.start_permission, employee_permission.end_permission, employee_permission.excuse,
    "employee"."name" AS "employee_name"
FROM
    "employee_permission"
    INNER JOIN "employee" ON "employee_permission"."employee_id" = "employee"."id"
WHERE
    "employee_permission"."id" = $1
`

type GetEmployeePermissionRow struct {
	ID              int32              `db:"id"`
	EmployeeID      int32              `db:"employee_id"`
	Type            PermissionType     `db:"type"`
	StartPermission pgtype.Timestamptz `db:"start_permission"`
	EndPermission   pgtype.Timestamptz `db:"end_permission"`
	Excuse          string             `db:"excuse"`
	EmployeeName    string             `db:"employee_name"`
}

func (q *Queries) GetEmployeePermission(ctx context.Context, id int32) (GetEmployeePermissionRow, error) {
	row := q.db.QueryRow(ctx, getEmployeePermission, id)
	var i GetEmployeePermissionRow
	err := row.Scan(
		&i.ID,
		&i.EmployeeID,
		&i.Type,
		&i.StartPermission,
		&i.EndPermission,
		&i.Excuse,
		&i.EmployeeName,
	)
	return i, err
}

const listEmployeePermissions = `-- name: ListEmployeePermissions :many
SELECT
    employee_permission.id, employee_permission.employee_id, employee_permission.type, employee_permission.start_permission, employee_permission.end_permission, employee_permission.excuse,
    "employee"."name" AS "employee_name"
FROM
    "employee_permission"
    INNER JOIN "employee" ON "employee_permission"."employee_id" = "employee"."id"
WHERE
    ($1 :: text IS NULL
    OR "employee"."name" ILIKE '%' || $1 || '%')
    AND (
        $2 :: integer IS NULL
        OR "employee_id" = $2 :: integer
    )
    AND (
        $3 :: santri_permission_type IS NULL
        OR "type" = $3 :: santri_permission_type
    )
    AND (
        $4 :: timestamptz IS NULL
        OR "start_permission" >= $4 :: timestamptz
    )
    AND (
        $5 :: timestamptz IS NULL
        OR "end_permission" <= $5 :: timestamptz
    )
LIMIT
    $7 OFFSET $6
`

type ListEmployeePermissionsParams struct {
	Q            pgtype.Text        `db:"q"`
	EmployeeID   pgtype.Int4        `db:"employee_id"`
	Type         interface{}        `db:"type"`
	FromDate     pgtype.Timestamptz `db:"from_date"`
	EndDate      pgtype.Timestamptz `db:"end_date"`
	OffsetNumber int32              `db:"offset_number"`
	LimitNumber  int32              `db:"limit_number"`
}

type ListEmployeePermissionsRow struct {
	ID              int32              `db:"id"`
	EmployeeID      int32              `db:"employee_id"`
	Type            PermissionType     `db:"type"`
	StartPermission pgtype.Timestamptz `db:"start_permission"`
	EndPermission   pgtype.Timestamptz `db:"end_permission"`
	Excuse          string             `db:"excuse"`
	EmployeeName    string             `db:"employee_name"`
}

func (q *Queries) ListEmployeePermissions(ctx context.Context, arg ListEmployeePermissionsParams) ([]ListEmployeePermissionsRow, error) {
	rows, err := q.db.Query(ctx, listEmployeePermissions,
		arg.Q,
		arg.EmployeeID,
		arg.Type,
		arg.FromDate,
		arg.EndDate,
		arg.OffsetNumber,
		arg.LimitNumber,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListEmployeePermissionsRow{}
	for rows.Next() {
		var i ListEmployeePermissionsRow
		if err := rows.Scan(
			&i.ID,
			&i.EmployeeID,
			&i.Type,
			&i.StartPermission,
			&i.EndPermission,
			&i.Excuse,
			&i.EmployeeName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateEmployeePermission = `-- name: UpdateEmployeePermission :one
UPDATE
    "employee_permission"
SET
    "employee_id" = COALESCE($1, employee_id),
    "start_permission" = COALESCE($2, start_permission),
    "end_permission" = $3,
    "excuse" = COALESCE($4, excuse)
WHERE
    "id" = $5 RETURNING id, employee_id, type, start_permission, end_permission, excuse
`

type UpdateEmployeePermissionParams struct {
	EmployeeID      pgtype.Int4        `db:"employee_id"`
	StartPermission pgtype.Timestamptz `db:"start_permission"`
	EndPermission   pgtype.Timestamptz `db:"end_permission"`
	Excuse          pgtype.Text        `db:"excuse"`
	ID              int32              `db:"id"`
}

func (q *Queries) UpdateEmployeePermission(ctx context.Context, arg UpdateEmployeePermissionParams) (EmployeePermission, error) {
	row := q.db.QueryRow(ctx, updateEmployeePermission,
		arg.EmployeeID,
		arg.StartPermission,
		arg.EndPermission,
		arg.Excuse,
		arg.ID,
	)
	var i EmployeePermission
	err := row.Scan(
		&i.ID,
		&i.EmployeeID,
		&i.Type,
		&i.StartPermission,
		&i.EndPermission,
		&i.Excuse,
	)
	return i, err
}
