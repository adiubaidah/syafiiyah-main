package persistence

import "github.com/jackc/pgx/v5/pgxpool"

type Store interface {
	//Querier // Uncomment this line to include Querier interface
	QuerierV2 // Advance version of Querier interface
}

type SQLStore struct {
	connPool *pgxpool.Pool
	*Queries
}

func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}
