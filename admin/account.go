package admin

import (
	"net/http"

	"github.com/smtc/goutils"
	"github.com/smtc/todolist/models"
	"github.com/zenazn/goji/web"
)

func AccountList(w http.ResponseWriter, r *http.Request) {
	models, _ := models.AccountList(0, 20, nil)
	list, _ := goutils.ToMapList(models, []string{"email", "name", "roles"}, goutils.FilterModeInclude)
	goutils.Render(w).RenderPage(list)
}

func AccountEntity(c web.C, w http.ResponseWriter, r *http.Request) {
}
