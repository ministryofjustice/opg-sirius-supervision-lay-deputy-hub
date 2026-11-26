package server

import (
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/shared"
	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/util"
)

type deputyDeathRecordVars struct {
	AppVars
	DeputyDeathRecord sirius.DeputyDeathRecord
}

func renderTemplateForDeputyDeathGet(client LayDeputyHubClient, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)
		app.PageName = fmt.Sprintf("Record death of %s %s", app.DeputyDetails.DeputyFirstName, app.DeputyDetails.DeputySurname)

		vars := deputyDeathRecordVars{
			AppVars: app,
		}

		deathRecord, err := client.GetDeputyDeathRecord(ctx, app.DeputyID())
		if err != nil {
			return err
		}
		vars.DeputyDeathRecord = deathRecord

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}

func renderTemplateForDeputyDeathPost(client LayDeputyHubClient, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)
		app.PageName = fmt.Sprintf("Record death of %s %s", app.DeputyDetails.DeputyFirstName, app.DeputyDetails.DeputySurname)

		vars := deputyDeathRecordVars{
			AppVars: app,
		}

		proofOfDeathReceived := r.PostFormValue("proof-of-death")

		deputyDeathRecordForm := sirius.DeputyDeathRecord{
			DeputyID:             app.DeputyID(),
			DateOfDeath:          shared.NewDate(r.PostFormValue("date-of-death")),
			ProofOfDeathReceived: proofOfDeathReceived,
			DateNotified:         shared.NewDate(r.PostFormValue("date-notified")),
			NotifiedBy: sirius.RefData{
				Handle: r.PostFormValue("notified-by"),
			},
			NotificationMethod: r.PostFormValue("notification-method"),
			Notes:              r.PostFormValue("notes"),
		}

		//Only add element when form submitted with proof of death true
		if proofOfDeathReceived == "true" {
			deputyDeathRecordForm.DateDeathCertificateReceieved = shared.NewDate(r.PostFormValue("date-proof-of-death-received"))
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
	}
}
