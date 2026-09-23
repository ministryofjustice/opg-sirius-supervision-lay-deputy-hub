package server

import (
	"net/http"
)

type deputyHubEventVars struct {
	AppVars
}

func renderTemplateForDeputyHubEvents(tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {

		app.PageName = "Timeline"

		return tmpl.ExecuteTemplate(w, "page", deputyHubEventVars{
			AppVars: app,
		})
	}
}
