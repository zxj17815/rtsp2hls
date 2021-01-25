package models

import (
	"database/sql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB
var SqlDB *sql.DB

func init() {
	var (
		err error
	)
	DB, err = gorm.Open(sqlite.Open("db/rtsp2hls.sqlite"), &gorm.Config{})

	// 获取通用数据库对象 sql.DB ，然后使用其提供的功能
	SqlDB, err := DB.DB()

	log.Println(SqlDB.Stats())

	if err != nil {
		log.Println(err)
	}
}
