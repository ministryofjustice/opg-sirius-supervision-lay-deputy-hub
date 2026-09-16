package server

import (
	"net/http"
)

//type LayDeputyHubInformation interface {
//	GetDeputyDetails(sirius.Context, int) (sirius.DeputyDetails, error)
//}

type deputyHubVars struct {
	AppVars
}

func renderTemplateForDeputyHub(tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		//ctx := getContext(r)

		//_, err := client.GetDeputyDetails(ctx, app.DeputyId())
		//if err != nil {
		//	return err
		//}

		app.PageName = "Deputy details"

		vars := deputyHubVars{
			AppVars: app,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
