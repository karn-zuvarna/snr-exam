package utils

import "gorm.io/gorm"

func RollbackOnError(tx *gorm.DB, err error) bool {
	if err != nil {
		tx.Rollback()
	}
	return err != nil
}
