package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/sirius"
)

type LayDeputyHubInformation interface {
	GetDeputyClients(sirius.Context, sirius.ClientListParams) (sirius.ClientList, error)
}

type deputyHubVars struct {
	AppVars
	ActiveClientCount int
}

func renderTemplateForDeputyHub(client LayDeputyHubInformation, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {

		ctx := getContext(r)

		params := sirius.ClientListParams{
			DeputyID: app.DeputyID(),
		}

		clientList, err := client.GetDeputyClients(ctx, params)
		if err != nil {
			return err
		}

		app.PageName = "Deputy details"

		var successMessage string

		switch r.URL.Query().Get("success") {
		case "deputyDetails":
			successMessage = "Deputy details updated"
		case "deleteDeathRecord":
			successMessage = "Deputy death record deleted"
		case "recordDeath":
			successMessage = "Deputy death record updated"
		default:
			successMessage = ""
		}
		app.SuccessMessage = successMessage

		vars := deputyHubVars{
			AppVars:           app,
			ActiveClientCount: clientList.Metadata.TotalActiveClients,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
