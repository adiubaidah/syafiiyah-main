package persistence

import "context"

type QuerierV2 interface {
	Querier
	ListSantri(ctx context.Context, arg ListSantriParams) ([]ListSantriRow, error)
	ListUsers(ctx context.Context, arg ListUserParams) ([]ListUserRow, error)
	ListParents(ctx context.Context, arg ListParentParams) ([]ListParentRow, error)
}

var _ QuerierV2 = (*Queries)(nil)
