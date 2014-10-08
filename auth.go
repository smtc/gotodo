package main

import (
	"net/http"
	"time"

	"github.com/smtc/gocache"
	"github.com/smtc/gotodo/models"
	"github.com/smtc/goutils"
	"github.com/zenazn/goji/web"
)

const (
	AUTHENTICATION = "Authentication"
)

func getAuth(w http.ResponseWriter, r *http.Request, level int) (*models.User, bool) {
	var (
		key   = r.Header.Get(AUTHENTICATION)
		cache = gocache.GetCache()
	)
	if key == "" {
		w.Header().Set("Cache-Control", "no-cache")
		http.Redirect(w, r, "/mustlogin", 301)
		return nil, false
	}

	var (
		user *models.User
		err  error
	)

	i, suc := cache.Get(key)
	if !suc {
		w.Header().Set("Cache-Control", "no-cache")
		http.Redirect(w, r, "/mustlogin", 301)
		return nil, false
	}

	id := i.(int64)
	user, err = models.GetUser(id)
	if err != nil {
		return user, false
	}

	if level >= 0 {
		suc = user.Level <= level
	}

	if !suc {
		goutils.HttpHandler(web.C{}, w, r).RenderError("没有足够权限。")
	}

	return user, suc
}

func MustLogin(c web.C, w http.ResponseWriter, r *http.Request) {
	h := goutils.HttpHandler(c, w, r)
	h.RenderError("尚未登录，或者登录超时。")
}

func LoginPage(c web.C, w http.ResponseWriter, r *http.Request) {
	h := goutils.HttpHandler(c, w, r)
	h.RenderHtml("/login.html")
}

func Login(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		h               = goutils.HttpHandler(c, w, r)
		cache           = gocache.GetCache()
		email, password string
		err             error
		pm              = map[string]string{}
	)

	err = h.FormatBody(&pm)
	if err != nil {
		h.RenderError(err.Error())
		return
	}
	email = pm["email"]
	password = pm["password"]

	user, ok := models.UserLogin(email, password, r)
	if ok == false {
		h.RenderError("用户名或密码不正确")
	} else {
		key := goutils.ObjectId() + goutils.RandomString(16)
		cache.Set(key, user.Id, time.Hour*8)
		w.Header().Set(AUTHENTICATION, key)
		h.RenderJson(nil, 1, "")
	}
}

func Logout(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		h     = goutils.HttpHandler(c, w, r)
		cache = gocache.GetCache()
	)

	_ = h
	_ = cache
}
