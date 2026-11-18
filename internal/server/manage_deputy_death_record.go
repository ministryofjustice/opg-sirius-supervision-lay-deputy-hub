package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/util"
)

type deputyDeathVars struct {
	AppVars
	DeputyDeathRecord sirius.DeputyDeathRecord
}

func renderTemplateForDeputyDeath(client LayDeputyHubClient, tmpl Template, logger *slog.Logger) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)
		app.PageName = fmt.Sprintf("Record death of %s %s", app.DeputyDetails.DeputyFirstName, app.DeputyDetails.DeputySurname)

		vars := deputyDeathVars{
			AppVars: app,
		}

		switch r.Method {
		case http.MethodGet:
			deathRecord, err := client.GetDeputyDeathRecord(ctx, app.DeputyDetails.ID)
			if err != nil {
				return err
			}
			vars.DeputyDeathRecord = deathRecord
			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			proofOfDeathReceived := r.PostFormValue("proof-of-death")
			proofOfDeathReceivedBool, err := strconv.ParseBool(proofOfDeathReceived)
			dateDeathCertificateReceieved := r.PostFormValue("date-proof-of-death-received")

			//Manually reset to null when form submitted with proof of death false
			if proofOfDeathReceivedBool == false {
				dateDeathCertificateReceieved = ""
			}

			deputyDeathRecordForm := sirius.DeputyDeathRecord{
				DeputyId:                      app.DeputyDetails.ID,
				DateOfDeath:                   r.PostFormValue("date-of-death"),
				ProofOfDeathReceived:          proofOfDeathReceivedBool,
				DateDeathCertificateReceieved: dateDeathCertificateReceieved,
				DateNotified:                  r.PostFormValue("date-notified"),
				NotifiedBy: sirius.RefData{
					Handle: r.PostFormValue("notified-by"),
				},
				NotificationMethod: r.PostFormValue("notification-method"),
				Notes:              r.PostFormValue("notes"),
			}

			logger.Info("deputyDeathForm", deputyDeathRecordForm)
			err = client.UpdateDeputyDeathRecord(ctx, deputyDeathRecordForm)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars.Errors = util.RenameErrors(verr.Errors)
				vars.DeputyDeathRecord = deputyDeathRecordForm
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
