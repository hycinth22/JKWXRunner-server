package database

import (
	"errors"
	"github.com/inkedawn/JKWXRunner-server/config"
	"github.com/jinzhu/gorm"
	"log"
	"os"
	"path/filepath"
	"time"

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
	var sqlLogger *log.Logger
	// initialize logger
	logger = log.New(os.Stdout, "[Database]", log.LstdFlags)
	// open log file && initialize sql logger
	{
		sqlLogDir, err := filepath.Abs(config.DBLogDir)
		if err != nil {
			logger.Println(sqlLogDir, err)
		}
		sqlLogFileName := filepath.Join(sqlLogDir, time.Now().Format(config.DBLogFileName))
		// ensure directory existed
		if err := os.MkdirAll(sqlLogDir, 0777); err != nil {
			logger.Println(sqlLogDir, err)
		}
		if sqlLogFile, err := os.Create(sqlLogFileName); err != nil {
			logger.Println(err)
			logger.Println("Warning: Failed to open log file. Database log will be shown only in database module log (usually StdOut)")
			sqlLogger = logger
		} else {
			sqlLogger = log.New(sqlLogFile, "", log.LstdFlags)
		}
	}
	// read database source config
	logger.Println("dsn: ", config.DSN)
	//  connect database
	db, err = gorm.Open("sqlite3", config.DSN) // db, err = gorm.Open("mysql", "root:root@/jkwxFucker?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	// configure database connection
	db.SetLogger(sqlLogger)
	//noinspection GoBoolExpressions
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
