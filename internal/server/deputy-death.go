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

			proofOfDeathReceived := r.PostFormValue("proofofdeath")
			dateDeathCertificateReceieved := r.PostFormValue("dateproofofdeathreceived")

			//Manually reset to null when submitted with proof of death false
			if proofOfDeathReceived == "false" {
				logger.Info("pod received", proofOfDeathReceived)
				dateDeathCertificateReceieved = ""
			}
			logger.Info("hpodr after", proofOfDeathReceived)

			deputyDeathForm := sirius.DeputyDeathRecord{
				DeputyId:                      app.DeputyDetails.ID,
				DateOfDeath:                   r.PostFormValue("dateofdeath"),
				ProofOfDeathReceived:          proofOfDeathReceived,
				DateDeathCertificateReceieved: dateDeathCertificateReceieved,
				DateNotified:                  r.PostFormValue("datenotified"),
				NotifiedBy:                    r.PostFormValue("notified-by"),
				NotificationMethod:            r.PostFormValue("notification-method"),
				Notes:                         r.PostFormValue("notes"),
			}

			logger.Info("deputyDeathForm", deputyDeathForm)
			err := client.CreateDeputyDeathRecord(ctx, deputyDeathForm, logger)

			if verr, ok := err.(sirius.ValidationError); ok {
				logger.Info("has a validation error", verr)

				vars.Errors = util.RenameErrors(verr.Errors)
				vars.DeputyDeathForm = deputyDeathForm
				logger.Info("ready to return", verr)

				return tmpl.ExecuteTemplate(w, "page", vars)
			} else if err != nil {
				logger.Info("non validation death", err)

				return err
			}

			return Redirect(fmt.Sprintf("/%d?success=recordDeath", app.DeputyId()))
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
