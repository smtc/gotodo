package main

import (
	"flag"
	"fmt"
	"net/http"
	"regexp"

	"github.com/guotie/config"
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
	//deferinit.InitAll()

	models.InitDB()
	run()
}

func run() {
	// static files
	goji.Get("/assets/*", http.FileServer(http.Dir("./")))

	// route
	goji.Get("/", IndexHandler)

	goji.Get("/login", LoginPage)
	goji.Get("/logout", Logout)
	goji.Get("/mustlogin", MustLogin)
	goji.Post("/login", Login)

	goji.Get("/user/", UserList)
	goji.Get("/user/roles", UserRoles)
	goji.Get("/user/select", UserSelect)
	goji.Get("/user/info", UserInfo)
	goji.Post("/user/", UserSave)
	goji.Delete("/user/", UserDelete)

	goji.Get("/project/", ProjectList)
	goji.Get("/project/select", ProjectSelect)
	goji.Post("/project/", ProjectSave)
	goji.Delete("/project/", ProjectDelete)
	goji.Get("/project/level", ProjectLevel)

	goji.Get("/task/", TaskList)
	goji.Post("/task/", TaskSave)
	goji.Post("/task/refresh", TaskRefresh)
	goji.Delete("/task/", TaskDelete)

	goji.Get(regexp.MustCompile(`^/(?P<model>.+)\.(?P<fn>.+):(?P<param>.+)$`), TemplateHandler)
	goji.Get(regexp.MustCompile(`^/(?P<model>.+)\.(?P<fn>.+)$`), TemplateHandler)
	goji.Get(regexp.MustCompile(`^/(?P<model>.+)_(?P<fn>.+)$`), TemplateHandler)

	goji.NotFound(goutils.NotFound)

	goji.Serve()
}

func IndexHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	h := goutils.HttpHandler(c, w, r)
	h.RenderHtml("/index.html")
}

/*
模板页暂时以 model.fn:param 分级
*/
func TemplateHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	temp := fmt.Sprintf("/%s_%s.html", c.URLParams["model"], c.URLParams["fn"])
	goutils.Render(w).RenderHtml(temp)
}
