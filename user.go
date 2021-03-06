package main

import (
	"net/http"

	"github.com/smtc/gocache"
	"github.com/smtc/gotodo/models"
	"github.com/smtc/goutils"
	"github.com/zenazn/goji/web"
)

func UserList(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		h = goutils.HttpHandler(c, w, r)
	)
	users, err := models.GetAllUsers()
	if err != nil {
		h.RenderError(err.Error())
		return
	}
	h.RenderPage(users, len(users))
}

func UserInfo(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		h     = goutils.HttpHandler(c, w, r)
		cache = gocache.GetCache()
		key   string
		id    int64
		user  *models.User
		err   error
	)

	key = r.Header.Get(AUTHENTICATION)
	if key != "" {
		_id, suc := cache.Get(key)
		if suc {
			id, suc = _id.(int64)
			if suc {
				user, err = models.GetUser(id)
				if err == nil {
					h.RenderJson(user, 1, "")
					return
				}
			}
		}
	}

	h.RenderError("没有登录")
}

/*
func UserEntity(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		h    = goutils.HttpHandler(c, w, r)
		id   = h.Param.GetInt64("id", 0)
		user *models.User
		err  error
	)

	if id == 0 {
		h.RenderJson(nil, 0, "")
	}

	user, err = models.GetUser(id)
	if err != nil {
		h.RenderError(err.Error())
		return
	}

	h.RenderJson(user, 1, "")
}
*/

func UserSave(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		h    = goutils.HttpHandler(c, w, r)
		user models.User
		err  error
	)

	err = h.FormatBody(&user)
	if err != nil {
		h.RenderError(err.Error())
		return
	}

	err = user.Save()
	if err != nil {
		h.RenderError(err.Error())
		return
	}

	h.RenderJson(user, 1, "")
}

func UserDelete(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		h   = goutils.HttpHandler(c, w, r)
		err error
		id  int64
	)

	err = h.FormatBody(&id)
	if err != nil {
		h.RenderError(err.Error())
		return
	}

	user, err := models.UserDelete(id)
	if err != nil {
		h.RenderJson(user, 1, err.Error())
		return
	}

	h.RenderJson(user, 1, "")
}

func UserSelect(c web.C, w http.ResponseWriter, r *http.Request) {
	h := goutils.HttpHandler(c, w, r)
	h.RenderJson(models.GetUserSelectData(), 1, "")
}

func UserRoles(c web.C, w http.ResponseWriter, r *http.Request) {
	h := goutils.HttpHandler(c, w, r)
	h.RenderJson(models.GetRoles(), 1, "")
}
