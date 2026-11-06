package server

import (
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/util"
)

type DeputyDetailsInformation interface {
	UpdateDeputyDetails(sirius.Context, int, sirius.DeputyDetailsForm) error
}

type ManageDeputyDetailsVars struct {
	AppVars
	Form sirius.DeputyDetailsForm
}

func renderTemplateForManageDeputyDetailsForm(client DeputyDetailsInformation, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)
		app.PageName = "Manage deputy contact details"
		form := sirius.DeputyDetailsForm{
			DeputyFirstName:        app.DeputyDetails.DeputyFirstName,
			DeputyMiddleName:       app.DeputyDetails.DeputyMiddleName,
			DeputySurname:          app.DeputyDetails.DeputySurname,
			DeputyOtherNames:       app.DeputyDetails.DeputyOtherNames,
			DeputySubType:          sirius.DeputySubType{SubType: "PERSON"},
			DeputyDateOfBirth:      app.DeputyDetails.DeputyDateOfBirth,
			Email:                  app.DeputyDetails.Email,
			MobileTelephoneNumber:  app.DeputyDetails.MobileTelephoneNumber,
			DaytimeTelephoneNumber: app.DeputyDetails.DaytimeTelephoneNumber,
			EveningTelephoneNumber: app.DeputyDetails.EveningTelephoneNumber,
			AddressLine1:           app.DeputyDetails.AddressLine1,
			AddressLine2:           app.DeputyDetails.AddressLine2,
			AddressLine3:           app.DeputyDetails.AddressLine3,
			Town:                   app.DeputyDetails.Town,
			County:                 app.DeputyDetails.County,
			Country:                app.DeputyDetails.Country,
			Postcode:               app.DeputyDetails.Postcode,
			IsAirmailRequired:      app.DeputyDetails.IsAirmailRequired,
		}
		vars := ManageDeputyDetailsVars{AppVars: app, Form: form}

		switch r.Method {
		case http.MethodGet:
			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			form := sirius.DeputyDetailsForm{
				DeputyFirstName:        r.PostFormValue("deputy-first-name"),
				DeputyMiddleName:       r.PostFormValue("deputy-middle-name"),
				DeputySurname:          r.PostFormValue("deputy-last-name"),
				DeputyOtherNames:       r.PostFormValue("deputy-other-names"),
				DeputySubType:          sirius.DeputySubType{SubType: "PERSON"},
				DeputyDateOfBirth:      r.PostFormValue("deputy-date-of-birth"),
				MobileTelephoneNumber:  r.PostFormValue("deputy-mobile-number"),
				DaytimeTelephoneNumber: r.PostFormValue("deputy-phone-number"),
				EveningTelephoneNumber: r.PostFormValue("deputy-evening-number"),
				Email:                  r.PostFormValue("deputy-email"),
				AddressLine1:           r.PostFormValue("address-line-1"),
				AddressLine2:           r.PostFormValue("address-line-2"),
				AddressLine3:           r.PostFormValue("address-line-3"),
				Town:                   r.PostFormValue("town"),
				County:                 r.PostFormValue("county"),
				Postcode:               r.PostFormValue("postcode"),
				Country:                r.PostFormValue("country"),
				IsAirmailRequired:      r.PostFormValue("airmail-required") == "true",
			}

			err := client.UpdateDeputyDetails(ctx, app.DeputyId(), form)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars.Errors = util.RenameErrors(verr.Errors)
				vars.Form = form
				return tmpl.ExecuteTemplate(w, "page", vars)
			} else if err != nil {
				return err
			}

			return Redirect(fmt.Sprintf("/%d?success=deputyDetails", app.DeputyId()))
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
