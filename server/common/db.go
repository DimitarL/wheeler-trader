package common

import (
	"github.com/jackc/pgconn"
)

const fkConstraintErrorCode = "23503"
const keyDuplicateErrorCode = "23505"

func IsForeignKeyError(err error) bool {
	return isDbError(err, fkConstraintErrorCode)
}

func IsDuplicateKeyError(err error) bool {
	return isDbError(err, keyDuplicateErrorCode)
}

func isDbError(err error, code string) bool {
	if err == nil {
		return false
	}

	if err, ok := err.(*pgconn.PgError); ok {
		return err.Code == code
	} else {
		return false
	}
}
