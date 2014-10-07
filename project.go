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
		list     []models.Project
		projects []models.Project
		p        models.Project
		user     *models.User
		err      error
		suc      bool
	)

	if user, suc = getAuth(w, r, models.ROLE_MEMBER); !suc {
		return
	}

	projects, err = models.GetProjectList()
	if err != nil {
		h.RenderError(err.Error())
		return
	}

	list = []models.Project{}
	for i := 0; i < len(projects); i++ {
		p = projects[i]
		if user.Level <= models.ROLE_MANAGER || p.Visibility == 0 || p.HasMember(user.Id) {
			list = append(list, p)
		}
	}

	h.RenderPage(list, 0)
}

func ProjectDelete(c web.C, w http.ResponseWriter, r *http.Request) {
	if _, suc := getAuth(w, r, models.ROLE_MANAGER); !suc {
		return
	}
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
	user, suc := getAuth(w, r, models.ROLE_MANAGER)
	if !suc {
		return
	}

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

	project.UpdatedBy = user.Id
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
