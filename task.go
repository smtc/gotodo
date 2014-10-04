package main

import (
	"net/http"

	"github.com/smtc/gotodo/models"
	"github.com/smtc/goutils"
	"github.com/zenazn/goji/web"
)

func TaskList(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		h    = goutils.HttpHandler(c, w, r)
		list []models.Task
		err  error
	)

	list, err = models.TaskList("")
	if err != nil {
		list = []models.Task{}
	}

	h.RenderPage(list, 0)
}
