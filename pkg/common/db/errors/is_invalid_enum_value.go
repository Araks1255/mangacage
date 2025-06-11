package errors

import (
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

func IsInvalidEnumValue(err error, enumType string) bool {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) && pgErr.Code == "22P02" && strings.Contains(pgErr.Message, enumType) { // Тут только через драйвер не получается, потому что поле DataTypeName, служащее для подобных случаев, при ошибке enum не заполняется
		return true
	}

	return false
}
