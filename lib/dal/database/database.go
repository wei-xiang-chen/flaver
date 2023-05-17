package database

import (
	"flaver/lib/dal/database/tools"

	"gorm.io/gorm"
)

func GetClientDB() *gorm.DB {
	return tools.GetClientDB()
}
