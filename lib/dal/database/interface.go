package database

import "gorm.io/gorm"

type GormGetSettable interface {
	GormGettable
}

type GormGettable interface {
	GetConn() *gorm.DB
	SetConn(*gorm.DB)
}
