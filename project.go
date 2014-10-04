package main

import (
	"net/http"

	"github.com/smtc/gotodo/models"
	"github.com/smtc/goutils"
	"github.com/zenazn/goji/web"
)

func ProjectList(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		h        = goutils.HttpHandler(c, w, r)
		projects []models.Project
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

	h.RenderPage(projects, total)
}

func ProjectEntity(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		h       = goutils.HttpHandler(c, w, r)
		id      = h.Param.GetInt64("id", 0)
		project = models.Project{Id: id}
	)

	if id == 0 {
		h.RenderJson(nil, 0, "")
	} else {
		project.Refresh()
		h.RenderJson(&project, 1, "")
	}
}

func ProjectDelete(c web.C, w http.ResponseWriter, r *http.Request) {
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

	err = models.ProjectDelete(id)
	if err != nil {
		h.RenderError(err.Error())
		return
	}

	h.RenderJson(nil, 1, "")
}

func ProjectSave(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		h       = goutils.HttpHandler(c, w, r)
		project models.Project
		err     error
	)

	err = h.FormatBody(&project)
	if err != nil {
		h.RenderError(err.Error())
		return
	}

	err = project.Save()
	if err != nil {
		h.RenderError(err.Error())
		return
	}

	h.RenderJson(project, 1, "")
}

func ProjectLevel(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		h = goutils.HttpHandler(c, w, r)
	)
	levels := models.GetLevel()
	h.RenderJson(levels, 1, "")
}

func ProjectSelect(c web.C, w http.ResponseWriter, r *http.Request) {
	h := goutils.HttpHandler(c, w, r)
	ps, err := models.GetProjectCache()
	if err != nil {
		h.RenderError(err.Error())
		return
	}

	kvs := []models.TextValue{}
	for _, v := range ps {
		kvs = append(kvs, models.TextValue{Text: v.Name, Value: v.Id})
	}

	h.RenderJson(kvs, 1, "")
}
