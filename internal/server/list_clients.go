package server

import "net/http"

type ListClientsVars struct {
	AppVars
}

func renderTemplateForClientTab(tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			return StatusError(http.StatusMethodNotAllowed)
		}

		app.PageName = "Clients"

		return tmpl.ExecuteTemplate(w, "page", ListClientsVars{
			AppVars: app,
		})
	}
}
