package persistence

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type ListSantriParams struct {
	Q            pgtype.Text       `db:"q"`
	OccupationID pgtype.Int4       `db:"occupation_id"`
	Generation   pgtype.Int4       `db:"generation"`
	IsActive     pgtype.Bool       `db:"is_active"`
	LimitNumber  int32             `db:"limit_number"`
	OffsetNumber int32             `db:"offset_number"`
	OrderBy      NullSantriOrderBy `db:"order_by"`
}

const listSantri = `-- name: ListSantri :many
SELECT
    "list_santri"."id",
    "list_santri"."name",
	"list_santri"."gender",
    "list_santri"."nis",
    "list_santri"."generation",
	"list_santri"."is_active",
	"list_santri"."photo",
    "list_santri"."parent_id" AS "parent_id",
    "list_santri"."parent_name" AS "parent_name",
    "list_santri"."parent_whatsapp_number" AS "parent_whatsapp_number",
    "list_santri"."occupation_id",
    "list_santri"."occupation_name"
FROM
    list_santri(
        $1,
        $2,
		$3,
        $4,
        $5,
        $6,
        $7
    )
`

type ListSantriRow struct {
	ID             int32       `db:"id" json:"id"`
	Name           string      `db:"name" json:"name"`
	Gender         GenderType  `db:"gender" json:"gender"`
	Nis            pgtype.Text `db:"nis" json:"nis"`
	Generation     int32       `db:"generation" json:"generation"`
	IsActive       pgtype.Bool `db:"is_active" json:"is_active"`
	Photo          pgtype.Text `db:"photo"`
	ParentID       pgtype.Int4 `db:"parent_id" json:"parent_id"`
	ParentName     pgtype.Text `db:"parent_name" json:"parent_name"`
	ParentWhatsapp pgtype.Text `db:"parent_whatsapp_number" json:"parent_whatsapp_number"`
	OccupationID   pgtype.Int4 `db:"occupation_id" json:"occupation_id"`
	OccupationName pgtype.Text `db:"occupation_name" json:"occupation_name"`
}

func (q *Queries) ListSantri(ctx context.Context, arg ListSantriParams) ([]ListSantriRow, error) {
	rows, err := q.db.Query(ctx, listSantri,
		arg.Q,
		arg.OccupationID,
		arg.Generation,
		arg.IsActive,
		arg.LimitNumber,
		arg.OffsetNumber,
		arg.OrderBy,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListSantriRow{}
	for rows.Next() {
		var i ListSantriRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Gender,
			&i.Nis,
			&i.Generation,
			&i.IsActive,
			&i.Photo,
			&i.ParentID,
			&i.ParentName,
			&i.ParentWhatsapp,
			&i.OccupationID,
			&i.OccupationName,
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

type ListUserParams struct {
	Q            pgtype.Text     `db:"q"`
	Role         NullRoleType    `db:"role"`
	HasOwner     pgtype.Bool     `db:"has_owner"`
	LimitNumber  int32           `db:"limit_number"`
	OffsetNumber int32           `db:"offset_number"`
	OrderBy      NullUserOrderBy `db:"order_by"`
}

const listUser = `-- name: ListUsers :many
SELECT
	"list_user"."id",
	"list_user"."username",
	"list_user"."role",
	"list_user"."id_owner",
	"list_user"."name_owner"
FROM list_user(
	$1,
	$2,
	$3,
	$4,
	$5,
	$6
)`

type ListUserRow struct {
	ID        int32       `db:"id"`
	Username  string      `db:"username"`
	Role      RoleType    `db:"role"`
	IDOwner   pgtype.Int4 `db:"id_owner"`
	NameOwner pgtype.Text `db:"name_owner"`
}

func (q *Queries) ListUsers(ctx context.Context, arg ListUserParams) ([]ListUserRow, error) {
	rows, err := q.db.Query(ctx, listUser,
		arg.Q,
		arg.Role,
		arg.HasOwner,
		arg.LimitNumber,
		arg.OffsetNumber,
		arg.OrderBy,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListUserRow{}
	for rows.Next() {
		var i ListUserRow
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.Role,
			&i.IDOwner,
			&i.NameOwner,
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

type ListParentParams struct {
	Q            pgtype.Text       `db:"q"`
	HasUser      pgtype.Bool       `db:"has_user"`
	LimitNumber  int32             `db:"limit_number"`
	OffsetNumber int32             `db:"offset_number"`
	OrderBy      NullParentOrderBy `db:"order_by"`
}

const listParent = `-- name: ListParent :many
SELECT
	"list_parent"."id",
	"list_parent"."name",
	"list_parent"."gender",
	"list_parent"."address",
	"list_parent"."whatsapp_number",
	"list_parent"."photo",
	"list_parent"."user_id",
	"list_parent"."username"
FROM list_parent(
	$1,
	$2,
	$3,
	$4,
	$5
)`

type ListParentRow struct {
	ID             int32       `db:"id"`
	Name           string      `db:"name"`
	Gender         GenderType  `db:"gender"`
	Address        string      `db:"address"`
	WhatsappNumber pgtype.Text `db:"whatsapp_number"`
	Photo          pgtype.Text `db:"photo"`
	UserID         pgtype.Int4 `db:"user_id"`
	Username       pgtype.Text `db:"username"`
}

func (q *Queries) ListParents(ctx context.Context, arg ListParentParams) ([]ListParentRow, error) {
	rows, err := q.db.Query(ctx, listParent,
		arg.Q,
		arg.HasUser,
		arg.LimitNumber,
		arg.OffsetNumber,
		arg.OrderBy,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListParentRow{}
	for rows.Next() {
		var i ListParentRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Gender,
			&i.Address,
			&i.WhatsappNumber,
			&i.Photo,
			&i.UserID,
			&i.Username,
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
