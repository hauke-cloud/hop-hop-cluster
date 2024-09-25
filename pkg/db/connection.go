package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	domain "github.com/hauke-cloud/hop-hop-cluster/pkg/domain"
)

func ConnectDatabase() (*gorm.DB, error) {
	db, dbErr := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	db.AutoMigrate(&domain.Cluster{})

	return db, dbErr
}
