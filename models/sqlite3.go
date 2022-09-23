package models

import (
	"fmt"
	"path"

	"gitee.com/tzxhy/dlna-home/constants"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitSqlite3() {
	if DB != nil {
		return
	}
	dbPath := path.Join(constants.StorageRoot, "cloud.db")
	fmt.Println("dbPath: ", dbPath)
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
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
