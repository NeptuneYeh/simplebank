package postgresdb

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SQLStore provides all functions to execute db queries and transactions
type SQLStore struct {
	*Queries
	//db *sql.DB
	connPool *pgxpool.Pool
}

var TxKey = struct {
}{}

// Store provides all functions to execute db queries and transactions
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
	VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error)
}

// NewStore creates a new Store
func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}
