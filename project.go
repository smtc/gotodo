package main

import (
	"net/http"

	"github.com/smtc/goutils"
	"github.com/smtc/todolist/models"
	"github.com/zenazn/goji/web"
)

func ProjectList(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		h        = goutils.HttpHandler(c, w, r)
		projects *[]models.Project
		err      error
		page     int
		size     int
		total    int
		where    string
	)

	page, size = h.GetPageSize()
	if name := h.Query.GetString("f.name", ""); name != "" {
		where = "`name` like '%" + name + "%'"
	}
	total, projects, err = models.GetProjectList(page, size, where)
	if err != nil {
		h.RenderError(err.Error())
		return
	}

	list, err := goutils.ToMapList(projects, []string{}, goutils.FilterModeExclude)
	if err != nil {
		h.RenderError(err.Error())
		return
	}
	_ = list

	h.RenderPage(projects, total)
}

func ProjectEntity(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		h       = goutils.HttpHandler(c, w, r)
		id      = h.Param.GetInt64("id", 0)
		project = models.Project{Id: id}
	)

	if id == 0 {
		h.RenderJson(nil, 0)
	} else {
		project.Refresh()
		h.RenderJson(&project, 1)
	}
}

func ProjectDelete(c web.C, w http.ResponseWriter, r *http.Request) {
}
