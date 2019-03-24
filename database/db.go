package database

import (
	"errors"
	"github.com/inkedawn/JKWXFucker-server/config"
	"github.com/jinzhu/gorm"
	"io"
	"log"
	"os"

	// _ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type DB = gorm.DB

var (
	db     *DB
	logger *log.Logger
)

func init() {
	var err error
	var logOutput io.Writer = os.Stdout
	// initialize logger
	logger = log.New(logOutput, "[Database]", log.LstdFlags)
	// read database source config
	logger.Println("dsn: ", config.DSN)
	//  connect database
	db, err = gorm.Open("sqlite3", config.DSN) // db, err = gorm.Open("mysql", "root:root@/jkwxFucker?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	// configure database connection
	db.LogMode(config.DBLogMode)
	db.DB().SetMaxIdleConns(config.DBMaxIdleConn)
	db.DB().SetMaxOpenConns(config.DBMaxOpenConn)
	db.SingularTable(true)
	// test database connection
	if err := db.DB().Ping(); err != nil {
		panic(errors.New("Ping Failed. " + err.Error()))
	}
}

func GetDB() *DB {
	return db
}

func IsRecordNotFoundError(err error) bool {
	return gorm.IsRecordNotFoundError(err)
}
