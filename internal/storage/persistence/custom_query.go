package persistence

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type ListSantriParams struct {
	Q            pgtype.Text       `db:"q" json:"q"`
	OccupationID pgtype.Int4       `db:"occupation_id" json:"occupation_id"`
	Generation   pgtype.Int4       `db:"generation" json:"generation"`
	IsActive     pgtype.Bool       `db:"is_active" json:"is_actived"`
	LimitNumber  int32             `db:"limit_number" json:"limit_number"`
	OffsetNumber int32             `db:"offset_number" json:"offset_number"`
	OrderBy      NullSantriOrderBy `db:"order_by" json:"order_by"`
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
	Gender         Gender      `db:"gender" json:"gender"`
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
