package server

import (
	"net/http"
)

type deputyHubVars struct {
	AppVars
}

func renderTemplateForDeputyHub(tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {

		app.PageName = "Deputy details"

		return tmpl.ExecuteTemplate(w, "page", deputyHubVars{
			AppVars: app,
		})
	}
}
