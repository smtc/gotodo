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
	CreatedAt int64  `json:"created_at"`
	UpdatedBy int64  `json:"updated_by"`
	UpdatedAt int64  `json:"updated_at"`
	Level     int    `json:"level"`
	Ongoing   int    `json:"ongoing"`
	Finished  int    `json:"finished"`
	Expired   int    `json:"expired"`
	Des       string `sql:"size:512" json:"des"`

	ChiefText string `sql:"-" json:"chief_text"`
	UsersText string `sql:"-" json:"users_text"`
	LevelText string `sql:"-" json:"level_text"`
}

func getProjectDB() *gorm.DB {
	return GetDB(DEFAULT_DB)
}

func (p *Project) Refresh() error {
	db := getProjectDB()
	return db.First(p, p.Id).Error
}

func GetProjectList(page, size int, where string) (int, *[]Project, error) {
	var (
		db       = getProjectDB()
		projects []Project
		total    int
		err      error
	)

	err = db.Model(&Project{}).Where(where).Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	if total == 0 {
		return 0, &[]Project{}, nil
	}

	err = db.Where(where).Order("edit_at desc").Offset((page - 1) * size).Limit(size).Find(&projects).Error
	for i := 0; i < len(projects); i++ {
		projects[i].LevelText = LEVELS[projects[i].Level]
	}

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
	p.UpdatedAt = time.Now().Unix()

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
