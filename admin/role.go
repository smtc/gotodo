package admin

import (
	"net/http"

	"github.com/smtc/goutils"
	"github.com/smtc/todolist/models"
	"github.com/zenazn/goji/web"
)

func RoleList(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		h     = goutils.HttpHandler(c, w, r)
		roles *[]models.Role
		err   error
	)

	roles, err = models.GetRoleList()
	if err != nil {
		h.RenderError(err.Error())
		return
	}

	h.RenderPage(roles)
}

func RoleEntity(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		h    = goutils.HttpHandler(c, w, r)
		id   = h.Param.GetInt64("id", 0)
		role = models.Role{Id: id}
	)

	if id == 0 {
		h.RenderJson(nil, 0)
	} else {
		role.Refresh()
		h.RenderJson(role, 1)
	}
}

func RoleSave(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		h    = goutils.HttpHandler(c, w, r)
		role models.Role
		err  error
	)

	println(1111)
	err = h.FormatBody(&role)
	if err != nil {
		h.RenderError(err.Error())
		return
	}

	println(2222)

	rr, err := goutils.ToJsonOnly(role)
	println(rr)
	err = role.Save()
	if err != nil {
		h.RenderError(err.Error())
		return
	}

	h.RenderJson(nil, 1)
}

func RoleTypes(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		h = goutils.HttpHandler(c, w, r)
	)

	h.RenderJson(models.GetRoleTypes(), 1)
}
