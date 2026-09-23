package server

import (
	"net/http"
)

type deputyHubEventVars struct {
	AppVars
}

func renderTemplateForDeputyHubEvents(tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		app.PageName = "Timeline"

		return tmpl.ExecuteTemplate(w, "page", deputyHubEventVars{
			AppVars: app,
		})
	}
}
