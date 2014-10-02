package models

func InitDB() {
	db := GetDB(DEFAULT_DB)
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Role{})
	db.AutoMigrate(&Project{})
	db.AutoMigrate(&Task{})
	db.Model(&Task{}).AddIndex("idx_object_id", "object_id")
	db.Model(&Task{}).AddIndex("idx_parent_id", "parent_id")
}
