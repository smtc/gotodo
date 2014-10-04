package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Task struct {
	Id        int64  `json:"id"`
	ObjectId  string `sql:"size:16" json:"object_id"`
	ParentId  string `sql:"size:16" json:"parent_id"`
	Path      string `sql:"size:200" json:"path"`
	ProjectId int64  `json:"project_id"`

	Name      string `sql:"size:128" json:"name"`
	Des       string `sql:"size:5000" json:"des"`
	User      int64  `json:"user"`
	Level     int    `json:"level"`
	Status    int    `json:"status"`
	SubNum    int    `json:"sub_num"`
	Progress  int    `json:"progress"`
	CreatedBy int64  `json:"created_by"`
	UpdatedBy int64  `json:"created_by"`

	CreatedAt int64 `json:"created_at"`
	StartAt   int64 `json:"start_at"`
	FinishAt  int64 `json:"finish_at"`
	UpdatedAt int64 `json:"updated_at"`
	Deadline  int64 `json:"deadline"`

	UserText    string `sql:"-" json:"user_text"`
	StatusText  string `sql:"-" json:"status_text"`
	CreatedText string `sql:"-" json:"created_text"`
	UpdatedText string `sql:"-" json:"updated_text"`
	SubTask     []Task `sql:"-" json:"sub_task"`
}

func getTaskDB() *gorm.DB {
	return GetDB(DEFAULT_DB)
}

func TaskList(where string) ([]Task, error) {
	var (
		db    = getTaskDB()
		err   error
		tasks []Task
	)
	err = db.Where(where).Order("level desc").Order("deadline").Find(&tasks).Error
	return tasks, err
}

func (t *Task) Refresh() error {
	db := getTaskDB()
	return db.First(t, t.Id).Error
}

func (t *Task) Save() error {
	var (
		db     = getTaskDB()
		err    error
		old    Task
		parent Task
	)

	if t.Id == 0 {
		t.ObjectId = ""
		t.CreatedAt = time.Now().Unix()
	} else {
		old.Id = t.Id
		err = old.Refresh()
		if err != nil {
			return err
		}
		t.CreatedAt = old.CreatedAt
		//t.ObjectId = old.ObjectId
	}
	t.UpdatedAt = time.Now().Unix()

	if t.ParentId == "" {
		t.Path = t.ObjectId
	} else {
		err = db.Where("object_id = ?", t.ParentId).First(&parent).Error
		if err != nil {
			return err
		}
		t.Path = parent.Path + "," + t.ObjectId
	}

	err = db.Save(t).Error
	return err
}
