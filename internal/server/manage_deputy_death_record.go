package server

import (
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/util"
)

type deputyDeathRecordVars struct {
	AppVars
	DeputyDeathRecord sirius.DeputyDeathRecord
}

func renderTemplateForDeputyDeath(client LayDeputyHubClient, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)
		app.PageName = fmt.Sprintf("Record death of %s %s", app.DeputyDetails.DeputyFirstName, app.DeputyDetails.DeputySurname)

		vars := deputyDeathRecordVars{
			AppVars: app,
		}

		switch r.Method {
		case http.MethodGet:
			deathRecord, err := client.GetDeputyDeathRecord(ctx, app.DeputyID())
			if err != nil {
				return err
			}
			vars.DeputyDeathRecord = deathRecord
			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			proofOfDeathReceived := r.PostFormValue("proof-of-death")
			dateDeathCertificateReceieved := r.PostFormValue("date-proof-of-death-received")

			//Manually reset to null when form submitted with proof of death false
			if proofOfDeathReceived == "false" {
				dateDeathCertificateReceieved = ""
			}

			deputyDeathRecordForm := sirius.DeputyDeathRecord{
				DeputyID:                      app.DeputyID(),
				DateOfDeath:                   r.PostFormValue("date-of-death"),
				ProofOfDeathReceived:          proofOfDeathReceived,
				DateDeathCertificateReceieved: dateDeathCertificateReceieved,
				DateNotified:                  r.PostFormValue("date-notified"),
				NotifiedBy: sirius.RefData{
					Handle: r.PostFormValue("notified-by"),
				},
				NotificationMethod: r.PostFormValue("notification-method"),
				Notes:              r.PostFormValue("notes"),
			}

			err := client.UpdateDeputyDeathRecord(ctx, deputyDeathRecordForm)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars.Errors = util.RenameErrors(verr.Errors)
				vars.DeputyDeathRecord = deputyDeathRecordForm
				return tmpl.ExecuteTemplate(w, "page", vars)
			} else if err != nil {
				return err
			}

			return Redirect(fmt.Sprintf("/%d?success=recordDeath", app.DeputyID()))
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
