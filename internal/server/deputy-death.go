package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/sirius"
)

type deputyDeathVars struct {
	AppVars
	DeputyDeathForm sirius.DeputyDeathRecord
}

func renderTemplateForDeputyDeath(client LayDeputyHubClient, tmpl Template, logger *slog.Logger) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		app.PageName = fmt.Sprintf("Record death of %s %s", app.DeputyDetails.DeputyFirstName, app.DeputyDetails.DeputySurname)

		switch r.Method {
		case http.MethodGet:
			return tmpl.ExecuteTemplate(w, "page", deputyDeathVars{
				AppVars:         app,
				DeputyDeathForm: sirius.DeputyDeathRecord{},
			})

		case http.MethodPost:
			ctx := getContext(r)

			var proofOfDeath = r.PostFormValue("proofofdeath")
			var dateOfDeath = r.PostFormValue("dateofdeath")
			var dateOfDeathOptional = r.PostFormValue("dateofdeath-optional")
			var dateProofOfDeathReceived = r.PostFormValue("dateproofofdeathreceived")
			var dateNotified = r.PostFormValue("datenotified")
			var notifiedBy = r.PostFormValue("notified-by")
			var notifiedByOther = r.PostFormValue("notified-by-other-text")
			var notificationMethod = r.PostFormValue("notification-method")
			var notes = r.PostFormValue("notes")

			if dateOfDeathOptional != "" {
				dateOfDeath = dateOfDeathOptional
			}

			deputyDeathForm := sirius.DeputyDeathRecord{
				DeputyId:                      app.DeputyDetails.ID,
				DateOfDeath:                   dateOfDeath,
				ProofOfDeathReceived:          proofOfDeath,
				DateDeathCertificateReceieved: dateProofOfDeathReceived,
				DateNotified:                  dateNotified,
				NotifiedBy:                    notifiedBy,
				NotifiedByOther:               notifiedByOther,
				NotificationMethod:            notificationMethod,
				Notes:                         notes,
			}

			logger.Info("deputyDeathForm", deputyDeathForm)
			err := client.CreateDeputyDeathRecord(ctx, deputyDeathForm)
			if err != nil {
				return err
			}
			return tmpl.ExecuteTemplate(w, "page", deputyDeathVars{
				AppVars: app,
			})
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
