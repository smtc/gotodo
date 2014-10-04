package main

import (
	"net/http"

	"github.com/smtc/gotodo/models"
	"github.com/smtc/goutils"
	"github.com/zenazn/goji/web"
)

func TaskList(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		h    = goutils.HttpHandler(c, w, r)
		list []models.Task
		err  error
	)

	list, err = models.TaskList("")
	if err != nil {
		list = []models.Task{}
	}

	h.RenderPage(list, 0)
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

	task, err := models.TaskDelete(id)
	if err != nil {
		h.RenderError(err.Error())
		return
	}

	h.RenderJson(task, 1, "")
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

	task, err = models.TaskRefresh(id)
	if err != nil {
		h.RenderError(err.Error())
		return
	}

	h.RenderJson(task, 1, "")
}
