package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/guotie/config"
	"github.com/guotie/deferinit"
	"github.com/smtc/goutils"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"

	"github.com/smtc/gotodo/models"
)

var (
	configFn = flag.String("config", "./config.json", "config file path")
)

func main() {
	config.ReadCfg(*configFn)
	deferinit.InitAll()

	models.InitDB()
	run()
}

func run() {
	// static files
	goji.Get("/assets/*", http.FileServer(http.Dir("./")))

	// route
	goji.Get("/", indexHandler)
	goji.Get("/menu", menuHandler)

	goji.Get("/user/", UserList)
	goji.Get("/user/:id", UserEntity)

	goji.Get("/role/", RoleList)
	goji.Delete("/role/", RoleDelete)
	goji.Get("/role/types", RoleTypes)
	goji.Get("/role/:id", RoleEntity)
	goji.Post("/role/:id", RoleSave)

	goji.Get("/project/", ProjectList)
	goji.Post("/project/", ProjectSave)
	goji.Delete("/project/", ProjectDelete)
	goji.Get("/project/level", ProjectLevel)
	goji.Get("/project/:id", ProjectEntity)

	goji.Get(regexp.MustCompile(`^/(?P<model>.+)\.(?P<fn>.+):(?P<param>.+)$`), templateHandler)
	goji.Get(regexp.MustCompile(`^/(?P<model>.+)\.(?P<fn>.+)$`), templateHandler)
	goji.Get(regexp.MustCompile(`^/(?P<model>.+)_(?P<fn>.+)$`), templateHandler)

	goji.NotFound(goutils.NotFound)

	goji.Serve()
}

func indexHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	h := goutils.HttpHandler(c, w, r)
	h.RenderHtml("/main.html")
}

/*
模板页暂时以 model.fn:param 分级
*/
func templateHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	temp := fmt.Sprintf("/%s_%s.html", c.URLParams["model"], c.URLParams["fn"])
	goutils.Render(w).RenderHtml(temp)
}

func menuHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	buf, err := ioutil.ReadFile("./menu.json")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(buf)
}
