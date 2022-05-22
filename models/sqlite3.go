package models

import (
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitSqlite3() {
	if DB != nil {
		return
	}

	db, err := gorm.Open(sqlite.Open("cloud.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	DB = db
	InitTables()
}

func InitTables() {
	DB.AutoMigrate(
		&PlayListItem{},
		&AudioItem{},
	)
}
