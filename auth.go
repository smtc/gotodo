package main

import (
	"net/http"

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
		h        = goutils.HttpHandler(c, w, r)
		email    string
		password string
		err      error
		ep       = map[string]string{}
	)

	err = h.FormatBody(&ep)
	if err != nil {
		h.RenderError(err.Error())
		return
	}
	email = ep["email"]
	password = ep["password"]

	k, ok := models.UserLogin(email, password, r)
	if ok == false {
		h.RenderError("用户名密码不正确")
	} else {
		h.RenderJson(k, 1, "")
	}
}
