package models

import (
	"github.com/jinzhu/gorm"
	"github.com/smtc/todolist/database"
)

// ===========================================================

type table_database int

const (
	DEFAULT_DB table_database = iota
	ACCOUNT_DB
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
	return database.GetDB(getSchema(model))
}

func InitDB() {
	db := GetDB(DEFAULT_DB)
	db.AutoMigrate(&Account{})
	db.AutoMigrate(&Role{})
	db.AutoMigrate(&Project{})
	db.AutoMigrate(&Task{})
	db.Model(&Task{}).AddIndex("idx_object_id", "object_id")
	db.Model(&Task{}).AddIndex("idx_parent_id", "parent_id")
}
