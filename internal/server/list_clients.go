package server

import "net/http"

type ListClientsVars struct {
	AppVars
}

func renderTemplateForClientTab(tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {

		app.PageName = "Clients"

		return tmpl.ExecuteTemplate(w, "page", ListClientsVars{
			AppVars: app,
		})
	}
}
