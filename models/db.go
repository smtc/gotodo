package models

import (
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/guotie/config"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

// ===========================================================

type table_database int

const (
	DEFAULT_DB table_database = iota
	ACCOUNT_DB
)

var (
	dbs = map[string]*gorm.DB{}
)

func getSchema(model table_database) string {
	db := ""
	switch model {
	case DEFAULT_DB:
		db = ""
	}
	return db
}

func GetDB(model table_database) *gorm.DB {
	return getDB(getSchema(model))
}

func getDB(dbname string) *gorm.DB {
	//return &db
	if dbname == "" {
		dbname = config.GetStringDefault("dbname", "")
	}

	db := dbs[dbname]
	if db == nil {
		newDB, err := opendb(dbname, "", "")
		if err != nil {
			panic(err)
		}
		db = &newDB
		dbs[dbname] = db
	}
	return db
}

// 建立数据库连接
func opendb(dbname, dbuser, dbpass string) (gorm.DB, error) {
	var (
		dbtype, dsn string
		db          gorm.DB
		err         error
	)

	if dbuser == "" {
		dbuser = config.GetStringDefault("dbuser", "")
	}
	if dbpass == "" {
		dbpass = config.GetStringDefault("dbpass", "")
	}

	dbtype = strings.ToLower(config.GetStringDefault("dbtype", "mysql"))
	if dbtype == "mysql" {
		dsn = fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
			dbuser,
			dbpass,
			config.GetStringDefault("dbproto", "tcp"),
			config.GetStringDefault("dbhost", "127.0.0.1"),
			config.GetIntDefault("dbport", 3306),
			dbname,
		)
		//dsn += "&loc=Asia%2FShanghai"
	} else if dbtype == "pg" || dbtype == "postgres" || dbtype == "postgresql" {
		dbtype = "postgres"
		dsn = fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
			dbuser,
			dbpass,
			config.GetStringDefault("dbhost", "127.0.0.1"),
			config.GetIntDefault("dbport", 5432),
			dbname)
	}

	//println(dbtype, dsn)
	db, err = gorm.Open(dbtype, dsn)
	if err != nil {
		log.Println(err.Error())
		return db, err
	}

	err = db.DB().Ping()
	if err != nil {
		log.Println(err.Error())
		return db, err
	}

	return db, nil
}

/*
// 供外部调用的API
func OpenDB(dbname, dbuser, dbpass string) {
	db, err := opendb(dbname, dbuser, dbpass)
	if err != nil {
		panic(err)
	}
	dbs[dbname] = &db
}

func CloseDB(dbname string) {
	GetDB(dbname).DB().Close()
}
*/
