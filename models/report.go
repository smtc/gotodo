package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

type Report struct {
	Id       int64  `json:"id"`
	TaskId   int64  `json:"task_id"`
	Progress int    `json:"progress"`
	Des      string `sql:"size:500" json:"des"`
	ReportBy int64  `json:"report_by"`
	ReportAt int64  `json:"report_at"`

	ReportName string `sql:"-" json:"report_name"`
}

func getReportDB() *gorm.DB {
	return GetDB(DEFAULT_DB)
}

func ReportList(taskId int64) ([]Report, error) {
	var (
		db      = getReportDB()
		reports []Report
		err     error
	)

	err = db.Where("task_id=?", taskId).Find(&reports).Error
	if err == nil {
		var r *Report
		for i := 0; i < len(reports); i++ {
			r = &reports[i]
			r.ReportName = GetUserName(r.ReportBy)
		}
	}
	return reports, err
}

func (r *Report) Save(user *User) error {
	var (
		db   = getTaskDB()
		err  error
		task *Task
	)

	task, err = GetTask(r.TaskId)
	if err != nil {
		return err
	}

	if !task.HasAuth(user, false) {
		return fmt.Errorf("没有足够权限")
	}

	r.ReportBy = user.Id
	r.ReportAt = time.Now().Unix()
	r.ReportName = GetUserName(r.ReportBy)
	err = db.Save(r).Error
	if err != nil {
		return err
	}

	if r.Progress > 0 && r.Progress < 100 {
		task.Progress = r.Progress
		task.Status = TASK_STATUS_PROGRESS
		return task.Save()
	} else if r.Progress == 100 {
		task.Progress = r.Progress
		task.Status = TASK_STATUS_TESTING
		return task.Save()
	}

	return nil
}
