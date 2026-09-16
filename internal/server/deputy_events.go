package server

import (
	"net/http"
)

type deputyHubEventVars struct {
	ListPage
}

func renderTemplateForDeputyHubEvents(tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		vars := deputyHubEventVars{}
		app.PageName = "Timeline"
		vars.AppVars = app

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
