package persistence

import "context"

type QuerierV2 interface {
	Querier
	ListSantri(ctx context.Context, arg ListSantriParams) ([]ListSantriRow, error)
}

var _ QuerierV2 = (*Queries)(nil)
