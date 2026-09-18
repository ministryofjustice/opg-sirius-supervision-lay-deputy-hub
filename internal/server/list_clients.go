package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/urlbuilder"
)

type ListClientsVars struct {
	ListPage
}

func (lcv ListClientsVars) CreateUrlBuilder() urlbuilder.UrlBuilder {
	return urlbuilder.UrlBuilder{
		OriginalPath: "clients",
	}
}

func renderTemplateForClientTab(tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			return StatusError(http.StatusMethodNotAllowed)
		}

		var vars ListClientsVars

		app.PageName = "Clients"
		vars.AppVars = app
		vars.UrlBuilder = vars.CreateUrlBuilder()

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
