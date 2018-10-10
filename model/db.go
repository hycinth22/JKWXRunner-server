package model

import (
	"errors"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"log"
	// _ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var DSN string
var db *gorm.DB

const dsnConfig = "database.dsn"

func init() {
	var err error
	dsnConfigData, err := ioutil.ReadFile(dsnConfig)
	if err != nil {
		panic(err)
	}
	DSN = string(dsnConfigData)
	log.Println("DSN: ", DSN)
	db, err = gorm.Open("sqlite3", DSN) // db, err = gorm.Open("mysql", "root:root@/jkwxFucker?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	db.DB().SetMaxIdleConns(2)
	db.DB().SetMaxOpenConns(10)
	db.SingularTable(true)
	if !db.HasTable(&Account{}) {
		if err := db.CreateTable(&Account{}).Error; err != nil {
			panic(errors.New("CreateTable Failed. " + err.Error()))
		}
	}
	if !db.HasTable(&Ticket{}) {
		if err := db.CreateTable(&Ticket{}).Error; err != nil {
			panic(errors.New("CreateTable Failed. " + err.Error()))
		}
	}
	if !db.HasTable(&SessionStore{}) {
		if err := db.CreateTable(&SessionStore{}).Error; err != nil {
			panic(errors.New("CreateTable Failed. " + err.Error()))
		}
	}
	if !db.HasTable(&AccountLog{}) {
		if err := db.CreateTable(&AccountLog{}).Error; err != nil {
			panic(errors.New("CreateTable Failed. " + err.Error()))
		}
	}
	if err := db.DB().Ping(); err != nil {
		panic(errors.New("Ping Failed. " + err.Error()))
	}
	db.AutoMigrate(&Account{})
}
