package main

import (
	"fmt"
	"net/http"

	"github.com/smtc/gotodo/models"
	"github.com/smtc/goutils"
	"github.com/zenazn/goji/web"
)

type taskList struct {
	Id    int64          `json:"id"`
	Name  string         `json:"name"`
	Level int            `json:"level"`
	Tasks []*models.Task `json:"tasks"`
}

func TaskList(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		h        = goutils.HttpHandler(c, w, r)
		projects map[int64]models.Project
		dict     map[int64]*models.Task
		tasks    = map[int64]*taskList{}
		err      error
	)

	user, suc := getAuth(w, r, models.ROLE_MEMBER)
	if !suc {
		return
	}

	projects, err = models.GetProjectCache()
	if err != nil {
		h.RenderError(err.Error())
		return
	}

	dict, err = models.TaskList(fmt.Sprintf("`status`!='%s' and `status`!='%s'", models.TASK_STATUS_STOPED, models.TASK_STATUS_FINISHED))
	if err == nil {
		var p models.Project
		for _, t := range dict {
			p = projects[t.ProjectId]
			if user.IsAdmin() || p.HasMember(user.Id) || t.User == user.Id {
				t.Editable = user.IsAdmin() || p.Chief == user.Id

				tl, suc := tasks[p.Id]
				if !suc {
					tl = &taskList{
						Id:    p.Id,
						Name:  p.Name,
						Level: p.Level,
						Tasks: []*models.Task{},
					}
					tasks[p.Id] = tl
				}

				tl.Tasks = append(tl.Tasks, t)
				/*
					if t.ParentId == 0 {
						tasks = append(tasks, t)
					} else {
						if parent, suc := dict[t.ParentId]; suc {
							parent.SubTask = append(parent.SubTask, t)
						}
					}
				*/
			}
		}
	}

	newTasks := []*taskList{}
	for _, t := range tasks {
		newTasks = append(newTasks, t)
	}

	h.RenderPage(newTasks, 0)
}

func TaskFinish(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		h   = goutils.HttpHandler(c, w, r)
		id  int64
		err error
	)

	err = h.FormatBody(&id)
	if err != nil {
		h.RenderError(err.Error())
		return
	}

	user, suc := getAuth(w, r, models.ROLE_MEMBER)
	if !suc {
		return
	}

	task, err := models.GetTask(id)
	if err != nil {
		h.RenderError(err.Error())
		return
	}

	if !task.HasAuth(user, true) {
		h.RenderError("没有足够的权限")
		return
	}

	if task.Status != models.TASK_STATUS_TESTING {
		h.RenderError("error status")
		return
	}

	task.UpdatedBy = user.Id
	task.Status = models.TASK_STATUS_FINISHED
	err = task.Save()
	if err != nil {
		h.RenderError(err.Error())
		return
	}

	h.RenderJson(task.ProjectId, 1, "")
}

func TaskSave(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		h    = goutils.HttpHandler(c, w, r)
		task models.Task
		err  error
	)

	err = h.FormatBody(&task)
	if err != nil {
		h.RenderError(err.Error())
		return
	}

	user, suc := getAuth(w, r, models.ROLE_MEMBER)
	if !suc {
		return
	}

	if !task.HasAuth(user, true) {
		h.RenderError("没有足够的权限")
		return
	}

	task.UpdatedBy = user.Id
	err = task.Save()
	if err != nil {
		h.RenderError(err.Error())
		return
	}

	task.Editable = true
	h.RenderJson(task, 1, "")
}

func TaskDelete(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		h   = goutils.HttpHandler(c, w, r)
		id  int64
		err error
	)
	err = h.FormatBody(&id)
	if err != nil {
		h.RenderError(err.Error())
		return
	}

	task, err := models.GetTask(id)
	if err != nil {
		h.RenderError(err.Error())
		return
	}

	user, suc := getAuth(w, r, models.ROLE_MEMBER)
	if !suc {
		return
	}

	if !task.HasAuth(user, true) {
		h.RenderError("没有足够的权限")
		return
	}

	suc, err = task.Delete()
	if err != nil {
		h.RenderError(err.Error())
		return
	}

	h.RenderJson(nil, 1, "")
}

func TaskRefresh(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		h    = goutils.HttpHandler(c, w, r)
		id   int64
		err  error
		task *models.Task
	)
	err = h.FormatBody(&id)
	if err != nil {
		h.RenderError(err.Error())
		return
	}

	task, err = models.GetTask(id)
	if err != nil {
		h.RenderError(err.Error())
		return
	}

	user, suc := getAuth(w, r, models.ROLE_MEMBER)
	if !suc {
		return
	}

	if !task.HasAuth(user, true) {
		h.RenderError("没有足够的权限")
		return
	}

	task.Status = models.TASK_STATUS_PROGRESS
	err = task.Save()
	if err != nil {
		h.RenderError(err.Error())
		return
	}

	h.RenderJson(task, 1, "")
}

func TaskWeightSelect(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		h    = goutils.HttpHandler(c, w, r)
		list = []models.TextValue{}
	)

	for i := 0; i < 7; i++ {
		list = append(list, models.TextValue{Text: fmt.Sprintf("%d", i+1), Value: i})
	}

	h.RenderJson(list, 1, "")
}
