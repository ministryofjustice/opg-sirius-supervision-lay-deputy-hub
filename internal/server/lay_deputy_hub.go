package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/sirius"
)

type LayDeputyHubInformation interface {
	GetDeputyDetails(sirius.Context, int) (sirius.DeputyDetails, error)
}

type deputyHubVars struct {
	AppVars
}

func renderTemplateForDeputyHub(client LayDeputyHubInformation, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)

		var selectedOrderStatuses []string
		selectedOrderStatuses = append(selectedOrderStatuses, "ACTIVE")

		_, err := client.GetDeputyDetails(ctx, app.DeputyId())
		if err != nil {
			return err
		}

		app.PageName = "Deputy details"

		vars := deputyHubVars{
			AppVars: app,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
