package main

import (
	"net/http"
	"time"

	"github.com/smtc/gocache"
	"github.com/smtc/gotodo/models"
	"github.com/smtc/goutils"
	"github.com/zenazn/goji/web"
)

func LoginPage(c web.C, w http.ResponseWriter, r *http.Request) {
	h := goutils.HttpHandler(c, w, r)
	h.RenderHtml("/login.html")
}

func Login(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		h                     = goutils.HttpHandler(c, w, r)
		cache                 = gocache.GetCache()
		email, password, etag string
		err                   error
		pm                    = map[string]string{}
	)

	err = h.FormatBody(&pm)
	if err != nil {
		h.RenderError(err.Error())
		return
	}
	email = pm["email"]
	password = pm["password"]
	etag = pm["etag"]

	if etag == "" {
		h.RenderError("登录异常.")
		return
	}

	user, ok := models.UserLogin(email, password, r)
	if ok == false {
		h.RenderError("用户名或密码不正确")
	} else {
		cache.Set(etag, user.Id, time.Hour*8)
		h.RenderJson(nil, 1, "")
	}
}
