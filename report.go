package main

import (
	"net/http"

	"github.com/smtc/gotodo/models"
	"github.com/smtc/goutils"
	"github.com/zenazn/goji/web"
)

func ReportList(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		h          = goutils.HttpHandler(c, w, r)
		task_id    = h.Param.GetInt64("task_id", 0)
		project_id = h.Param.GetInt64("project_id", 0)
		list       []models.Report
		err        error
	)

	user, suc := getAuth(w, r, models.ROLE_MEMBER)
	if !suc {
		return
	}

	if task_id == 0 || project_id == 0 {
		h.RenderError("task not found")
		return
	}

	project, err := models.GetProject(project_id)
	if err != nil {
		h.RenderError(err.Error())
		return
	}

	if !(user.IsAdmin() || project.HasMember(user.Id)) {
		h.RenderError("没有足够权限")
		return
	}

	list, err = models.ReportList(task_id)
	if err != nil {
		list = []models.Report{}
	}

	h.RenderPage(list, 0)
}

func ReportSave(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		h      = goutils.HttpHandler(c, w, r)
		report models.Report
		err    error
	)

	err = h.FormatBody(&report)
	if err != nil {
		h.RenderError(err.Error())
		return
	}

	user, suc := getAuth(w, r, models.ROLE_MEMBER)
	if !suc {
		return
	}

	err = report.Save(user)
	if err != nil {
		h.RenderError(err.Error())
		return
	}

	h.RenderJson(report, 1, "")

}
