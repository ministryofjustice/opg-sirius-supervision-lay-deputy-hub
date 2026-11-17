package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/util"
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

			proofOfDeathReceived := r.PostFormValue("proof-of-death")
			dateDeathCertificateReceieved := r.PostFormValue("date-proof-of-death-received")

			//Manually reset to null when submitted with proof of death false
			if proofOfDeathReceived == "false" {
				dateDeathCertificateReceieved = ""
			}

			deputyDeathForm := sirius.DeputyDeathRecord{
				DeputyId:                      app.DeputyDetails.ID,
				DateOfDeath:                   r.PostFormValue("date-of-death"),
				ProofOfDeathReceived:          proofOfDeathReceived,
				DateDeathCertificateReceieved: dateDeathCertificateReceieved,
				DateNotified:                  r.PostFormValue("date-notified"),
				NotifiedBy:                    r.PostFormValue("notified-by"),
				NotificationMethod:            r.PostFormValue("notification-method"),
				Notes:                         r.PostFormValue("notes"),
			}

			logger.Info("deputyDeathForm", deputyDeathForm)
			err := client.CreateDeputyDeathRecord(ctx, deputyDeathForm, logger)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars.Errors = util.RenameErrors(verr.Errors)
				vars.DeputyDeathForm = deputyDeathForm
				return tmpl.ExecuteTemplate(w, "page", vars)
			} else if err != nil {
				return err
			}

			return Redirect(fmt.Sprintf("/%d?success=recordDeath", app.DeputyId()))
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
