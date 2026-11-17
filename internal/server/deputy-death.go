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

		vars := deputyDeathVars{
			AppVars:         app,
			DeputyDeathForm: sirius.DeputyDeathRecord{},
		}

		vars.Errors = sirius.ValidationErrors{}

		switch r.Method {
		case http.MethodGet:
			return tmpl.ExecuteTemplate(w, "page", vars)

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

			if verr, ok := err.(sirius.ValidationError); ok {
				vars.Errors = verr.Errors
				return tmpl.ExecuteTemplate(w, "page", vars)
			}
			if err != nil {
				return err
			}

			return Redirect(fmt.Sprintf("/%d?success=recordDeath", app.DeputyId()))
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
