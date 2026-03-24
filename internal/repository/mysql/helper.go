package mysql

import (
	"errors"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"

	"github.com/template-be-commerce/pkg/apperrors"
)

const mysqlErrDuplicateEntry uint16 = 1062

// isUniqueViolation returns true when err is a MySQL duplicate-entry error.
func isUniqueViolation(err error) bool {
	var me *mysql.MySQLError
	return errors.As(err, &me) && me.Number == mysqlErrDuplicateEntry
}

// handleGORMError converts common GORM errors to apperrors sentinels.
func handleGORMError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperrors.ErrNotFound
	}
	if isUniqueViolation(err) {
		return apperrors.ErrEmailExists // caller may override with a more specific error
	}
	return err
}
