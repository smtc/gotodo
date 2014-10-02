package main

import (
	"net/http"

	"github.com/smtc/gotodo/models"
	"github.com/smtc/goutils"
	"github.com/zenazn/goji/web"
)

func UserList(w http.ResponseWriter, r *http.Request) {
	users, _ := models.UserList(0, 20, nil)
	list, _ := goutils.ToMapList(users, []string{"email", "name", "roles"}, goutils.FilterModeInclude)
	h := goutils.HttpHandler(web.C{}, w, r)
	h.RenderPage(list, 20)
}

func UserEntity(c web.C, w http.ResponseWriter, r *http.Request) {
}
