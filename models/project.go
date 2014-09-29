package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

type Project struct {
	Id        int64  `json:"id"`
	Name      string `sql:"size:100" json:"name"`
	Chief     int64  `json:"chief"`
	Users     string `sql:"size:10000" json:"users"`
	CreatedAt int64  `json:"created_at" format:"datetime"`
	EditBy    int64  `json:"edit_by"`
	EditAt    int64  `json:"edit_at" format:"datetime"`
	Status    string `sql:"size:10" json:"status"`

	ChiefText  string `sql:"-" json:"chief_text"`
	UsersText  string `sql:"-" json:"users_text"`
	StatusText string `sql:"-" json:"status_text"`
}

func getProjectDB() *gorm.DB {
	return GetDB(DEFAULT_DB)
}

func (p *Project) Refresh() error {
	db := getProjectDB()
	return db.First(p, p.Id).Error
}

func GetProjectList(page, size int, where string, data ...interface{}) (int, *[]Project, error) {
	var (
		db       = getProjectDB()
		projects []Project
		total    int
		err      error
	)

	err = db.Model(&Project{}).Where(where, data).Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	err = db.Where(where, data).Offset((page - 1) * size).Limit(size).Find(&projects).Error
	return total, &projects, err
}

func (p *Project) Save() error {
	var (
		db  = getProjectDB()
		old Project
		err error
	)

	if p.Id == 0 {
		p.CreatedAt = time.Now().Unix()
	} else {
		err = db.First(&old, p.Id).Error
		if err != nil {
			return err
		}
	}

	return db.Save(p).Error
}

func ProjectDelete(id int64) error {
	var (
		db    = getProjectDB()
		count int
		err   error
	)

	err = db.Model(&Task{}).Where("project_id = ?", id).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("Project has %d tasks, can't delete.", count)
	}

	return db.Where("id = ?", id).Delete(&Project{}).Error
}
