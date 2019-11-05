package database

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/inkedawn/JKWXRunner-server/config"

	// _ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type DB = gorm.DB

type TX interface {
	Where(query interface{}, args ...interface{}) *gorm.DB
	Or(query interface{}, args ...interface{}) *gorm.DB
	Not(query interface{}, args ...interface{}) *gorm.DB
	Limit(limit interface{}) *gorm.DB
	Offset(offset interface{}) *gorm.DB
	Order(value interface{}, reorder ...bool) *gorm.DB
	Select(query interface{}, args ...interface{}) *gorm.DB
	Omit(columns ...string) *gorm.DB
	Group(query string) *gorm.DB
	Having(query interface{}, values ...interface{}) *gorm.DB
	Joins(query string, args ...interface{}) *gorm.DB
	Scopes(funcs ...func(*gorm.DB) *gorm.DB) *gorm.DB
	Unscoped() *gorm.DB
	Attrs(attrs ...interface{}) *gorm.DB
	Assign(attrs ...interface{}) *gorm.DB
	First(out interface{}, where ...interface{}) *gorm.DB
	Take(out interface{}, where ...interface{}) *gorm.DB
	Last(out interface{}, where ...interface{}) *gorm.DB
	Find(out interface{}, where ...interface{}) *gorm.DB
	Preloads(out interface{}) *gorm.DB
	Scan(dest interface{}) *gorm.DB
	Row() *sql.Row
	Rows() (*sql.Rows, error)
	ScanRows(rows *sql.Rows, result interface{}) error
	Pluck(column string, value interface{}) *gorm.DB
	Count(value interface{}) *gorm.DB
	Related(value interface{}, foreignKeys ...string) *gorm.DB
	FirstOrInit(out interface{}, where ...interface{}) *gorm.DB
	FirstOrCreate(out interface{}, where ...interface{}) *gorm.DB
	Update(attrs ...interface{}) *gorm.DB
	Updates(values interface{}, ignoreProtectedAttrs ...bool) *gorm.DB
	UpdateColumn(attrs ...interface{}) *gorm.DB
	UpdateColumns(values interface{}) *gorm.DB
	Save(value interface{}) *gorm.DB
	Create(value interface{}) *gorm.DB
	Delete(value interface{}, where ...interface{}) *gorm.DB
	Raw(sql string, values ...interface{}) *gorm.DB
	Exec(sql string, values ...interface{}) *gorm.DB
	Model(value interface{}) *gorm.DB
	Table(name string) *gorm.DB
	NewRecord(value interface{}) bool
	RecordNotFound() bool
	CreateTable(models ...interface{}) *gorm.DB
	DropTable(values ...interface{}) *gorm.DB
	DropTableIfExists(values ...interface{}) *gorm.DB
	HasTable(value interface{}) bool
	AutoMigrate(values ...interface{}) *gorm.DB
	ModifyColumn(column string, typ string) *gorm.DB
	DropColumn(column string) *gorm.DB
	AddIndex(indexName string, columns ...string) *gorm.DB
	AddUniqueIndex(indexName string, columns ...string) *gorm.DB
	RemoveIndex(indexName string) *gorm.DB
	AddForeignKey(field string, dest string, onDelete string, onUpdate string) *gorm.DB
	RemoveForeignKey(field string, dest string) *gorm.DB
	Association(column string) *gorm.Association
	Preload(column string, conditions ...interface{}) *gorm.DB
	Set(name string, value interface{}) *gorm.DB
	InstantSet(name string, value interface{}) *gorm.DB
	Get(name string) (value interface{}, ok bool)
	AddError(err error) error
	GetErrors() []error
}

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
