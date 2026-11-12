package server

import (
	"fmt"
	"log/slog"
	"net/http"
)

type deputyDeathVars struct {
	AppVars
	DeputyDeathForm deputyDeathForm
}

type deputyDeathForm struct {
	ProofOfDeath             string
	DateOfDeath              string
	DateProofOfDeathReceived string
	DateNotified             string
	NotifiedBy               string
	HowOpgNotified           string
	Notes                    string
}

func renderTemplateForDeputyDeath(client LayDeputyHubClient, tmpl Template, logger *slog.Logger) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		app.PageName = fmt.Sprintf("Record death of %s %s", app.DeputyDetails.DeputyFirstName, app.DeputyDetails.DeputySurname)

		switch r.Method {
		case http.MethodGet:
			return tmpl.ExecuteTemplate(w, "page", deputyDeathVars{
				AppVars:         app,
				DeputyDeathForm: deputyDeathForm{},
			})

		case http.MethodPost:
			//ctx := getContext(r)

			var proofOfDeath = r.PostFormValue("proofofdeath")
			var dateOfDeath = r.PostFormValue("dateofdeath")
			var dateOfDeathOptional = r.PostFormValue("dateofdeath-optional")
			var dateProofOfDeathReceived = r.PostFormValue("dateproofofdeathreceived")
			var dateNotified = r.PostFormValue("datenotified")
			var notifiedBy = r.PostFormValue("notified-by")
			var howOpgNotified = r.PostFormValue("how-opg-notified")
			var notes = r.PostFormValue("notes")

			if dateOfDeathOptional != "" {
				dateOfDeath = dateOfDeathOptional
			}

			deputyDeathForm := deputyDeathForm{
				ProofOfDeath:             proofOfDeath,
				DateOfDeath:              dateOfDeath,
				DateProofOfDeathReceived: dateProofOfDeathReceived,
				DateNotified:             dateNotified,
				NotifiedBy:               notifiedBy,
				HowOpgNotified:           howOpgNotified,
				Notes:                    notes,
			}

			logger.Info("deputyDeathForm", deputyDeathForm)
			//timeline, err := client.GetDeputyTimeline(ctx, app.DeputyDetails.ID)
			//if err != nil {
			//	return err
			//}
			return tmpl.ExecuteTemplate(w, "page", deputyDeathVars{
				AppVars: app,
			})
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
