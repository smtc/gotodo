package models

import (
	"flag"

	"github.com/guotie/config"
)

func init() {
	configFn := flag.String("config", "../test.json", "config file path")
	config.ReadCfg(*configFn)

	dropTables()

	InitDB()
}

func dropTables() {
	db := GetDB(DEFAULT_DB)
	db.DropTableIfExists(User{})
	db.DropTableIfExists(Role{})
	db.DropTableIfExists(Project{})
	db.DropTableIfExists(Task{})
}
