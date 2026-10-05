package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/sirius"
)

type deputyTimelineVars struct {
	AppVars
	Timeline []sirius.TimelineEvent
}

func renderTemplateForDeputyTimeline(client LayDeputyHubClient, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {

		ctx := getContext(r)

		timeline, err := client.GetDeputyTimeline(ctx, app.DeputyDetails.ID)
		if err != nil {
			return err
		}

		app.PageName = "Timeline"

		return tmpl.ExecuteTemplate(w, "page", deputyTimelineVars{
			AppVars:  app,
			Timeline: timeline,
		})
	}
}
