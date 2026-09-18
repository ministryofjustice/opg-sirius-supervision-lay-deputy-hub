package server

import (
	"net/http"
)

type deputyHubVars struct {
	AppVars
}

func renderTemplateForDeputyHub(tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		app.PageName = "Deputy details"

		vars := deputyHubVars{
			AppVars: app,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
