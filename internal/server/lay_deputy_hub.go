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
			DeputyId: app.DeputyId(),
		}

		clientList, err := client.GetDeputyClients(ctx, params)
		if err != nil {
			return err
		}

		app.PageName = "Deputy details"

		vars := deputyHubVars{
			AppVars:           app,
			ActiveClientCount: clientList.Metadata.TotalActiveClients,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
