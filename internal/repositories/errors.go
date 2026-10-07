package repositories

import (
	"errors"

	"github.com/lib/pq"
)

// Codes d'erreur PostgreSQL : https://www.postgresql.org/docs/current/errcodes-appendix.html
const (
	foreignKeyViolationCode = "23503"
	uniqueViolationCode     = "23505"
)

// IsUniqueViolation : la valeur existe déjà (ex : email ou pseudo déjà pris)
func IsUniqueViolation(err error) bool {
	return hasPostgresCode(err, uniqueViolationCode)
}

// IsForeignKeyViolation : la ligne est encore utilisée ailleurs (ex : un utilisateur qui a des commandes)
func IsForeignKeyViolation(err error) bool {
	return hasPostgresCode(err, foreignKeyViolationCode)
}

func hasPostgresCode(err error, code string) bool {
	var postgresError *pq.Error
	return errors.As(err, &postgresError) && string(postgresError.Code) == code
}
