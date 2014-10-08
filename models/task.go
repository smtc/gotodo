package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/smtc/goutils"
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
	Status    string `sql:"size:10" json:"status"`
	SubNum    int    `json:"sub_num"`
	Progress  int    `json:"progress"`
	CreatedBy int64  `json:"created_by"`
	UpdatedBy int64  `json:"created_by"`

	CreatedAt int64 `json:"created_at"`
	StartAt   int64 `json:"start_at"`
	FinishAt  int64 `json:"finish_at"`
	UpdatedAt int64 `json:"updated_at"`
	ReportAt  int64 `json:"report_at"`
	Deadline  int64 `json:"deadline"`

	UserName    string `sql:"-" json:"user_name"`
	StatusText  string `sql:"-" json:"status_text"`
	CreatedText string `sql:"-" json:"created_text"`
	UpdatedText string `sql:"-" json:"updated_text"`
	SubTask     []Task `sql:"-" json:"sub_task"`
	Editable    bool   `sql:"-" json:"editable"`
}

func getTaskDB() *gorm.DB {
	return GetDB(DEFAULT_DB)
}

func TaskList(where string, data ...interface{}) ([]Task, error) {
	var (
		db    = getTaskDB()
		err   error
		tasks []Task
		task  *Task
		count int
	)
	err = db.Where(where, data).Order("level desc").Order("deadline").Find(&tasks).Error

	count = len(tasks)
	for i := 0; i < count; i++ {
		task = &tasks[i]
		task.setName()
	}

	return tasks, err
}

func GetTask(id int64) (*Task, error) {
	var (
		db   = getTaskDB()
		task Task
		err  error
	)
	err = db.First(&task, id).Error
	return &task, err
}

func (t *Task) Save() error {
	var (
		db     = getTaskDB()
		err    error
		old    Task
		parent Task
	)

	if t.Id == 0 {
		t.ObjectId = goutils.ObjectId()
		t.CreatedAt = time.Now().Unix()
		t.Status = TASK_STATUS_CREATED
	} else {
		_ = old
	}
	t.UpdatedAt = time.Now().Unix()
	t.setName()

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

func (t *Task) setName() {
	t.UserName = GetUserName(t.User)
	t.StatusText = TASK_STATUS[t.Status]
}

// 如果任务还没开始-删除，如果已经进行，更改状态为-中止
func (t *Task) Delete() (bool, error) {
	db := getTaskDB()
	if t.Status == TASK_STATUS_FINISHED {
		return false, fmt.Errorf("任务已完成，不能删除。")
	}

	t.setName()
	if t.Status == TASK_STATUS_CREATED {
		return true, db.Delete(t).Error
	} else {
		t.Status = TASK_STATUS_STOPED
		return false, db.Save(t).Error
	}
}

// mustChief - 需要项目组长以上权限
func (t *Task) HasAuth(user *User, mustChief bool) bool {
	if user.IsAdmin() {
		return true
	}

	p, err := GetProject(t.ProjectId)
	if err != nil {
		return false
	}

	if mustChief {
		return p.Chief == user.Id
	} else {
		for _, uid := range p.Users {
			if uid == user.Id {
				return true
			}
		}
	}
	return false
}
