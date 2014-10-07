package main

import (
	"net/http"

	"github.com/smtc/gotodo/models"
	"github.com/smtc/goutils"
	"github.com/zenazn/goji/web"
)

func TaskList(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		h        = goutils.HttpHandler(c, w, r)
		projects map[int64]models.Project
		list     []models.Task
		tasks    = []models.Task{}
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

	list, err = models.TaskList("")
	if err == nil {
		var p models.Project
		for _, t := range list {
			p = projects[t.ProjectId]
			if user.IsAdmin() || p.HasMember(user.Id) || t.User == user.Id {
				t.Editable = user.IsAdmin() || p.Chief == user.Id
				tasks = append(tasks, t)
			}
		}
	}

	h.RenderPage(tasks, 0)
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
