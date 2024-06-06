package postgresdb

import (
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"
)

const (
	ForeignKeyViolation = "23503"
	UniqueViolation     = "23505"
)

// var ErrRecordNotFound = sql.ErrNoRows
var ErrRecordNotFound = pgx.ErrNoRows

var ErrUniqueViolation = &pq.Error{
	Code: UniqueViolation,
}

func ErrorCode(err error) string {
	var pgErr *pq.Error
	if errors.As(err, &pgErr) {
		return pgErr.Code.Name()
	}
	return ""
}
