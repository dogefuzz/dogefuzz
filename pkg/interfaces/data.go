package interfaces

import "gorm.io/gorm"

type Connection interface {
	Migrate() error
	Clean() error
	GetDB() *gorm.DB
}
