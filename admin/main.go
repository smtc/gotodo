package admin

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/smtc/goutils"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
)

func AdminMux() *web.Mux {
	mux := web.New()

	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Get("/admin/", indexHandler)
	mux.Get("/admin/menu", menuHandler)

	mux.Get("/admin/account/", AccountList)
	mux.Get("/admin/account/:id", AccountEntity)

	mux.Get("/admin/role/", RoleList)
	mux.Get("/admin/role/types", RoleTypes)
	mux.Get("/admin/role/:id", RoleEntity)
	mux.Post("/admin/role/:id", RoleSave)

	mux.Get(regexp.MustCompile(`^/admin/(?P<model>.+)\.(?P<fn>.+):(?P<param>.+)$`), templateHandler)
	mux.Get(regexp.MustCompile(`^/admin/(?P<model>.+)\.(?P<fn>.+)$`), templateHandler)

	mux.NotFound(goutils.NotFound)
	return mux
}

func indexHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	h := goutils.HttpHandler(c, w, r)
	h.RenderHtml("/admin/main.html")
}

/*
模板页暂时以 model.fn:param 分级
*/
func templateHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	temp := fmt.Sprintf("/admin/%s_%s.html", c.URLParams["model"], c.URLParams["fn"])
	goutils.Render(w).RenderHtml(temp)
}

func menuHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	buf, err := ioutil.ReadFile("./admin/menu.json")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(buf)
}
