package errors

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

func IsForeignKeyViolation(err error, keyName string) bool {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) && pgErr.Code == "23503" && pgErr.ConstraintName == keyName {
		return true
	}

	return false
}
