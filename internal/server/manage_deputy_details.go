package server

import "net/http"

//type ListClientsVars struct {
//	AppVars
//}

func renderTemplateForManageDeputyDetailsForm(tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {

		app.PageName = "Manage deputy contact details"

		return tmpl.ExecuteTemplate(w, "page", ListClientsVars{
			AppVars: app,
		})
	}
}
