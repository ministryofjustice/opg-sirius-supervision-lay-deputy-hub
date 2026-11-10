package server

import (
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/util"
)

type AdditionalInformationInformation interface {
	UpdateAdditionalInformation(sirius.Context, int, sirius.AdditionalInformationForm) error
}

type ManageAdditionalInformationVars struct {
	AppVars
	Form sirius.AdditionalInformationForm
}

func renderTemplateForManageAdditionalInformationForm(client AdditionalInformationInformation, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)
		app.PageName = "Manage additional information"
		form := sirius.AdditionalInformationForm{
			InterpreterRequired:               app.DeputyDetails.InterpreterRequired,
			SpecialCorrespondenceRequirements: app.DeputyDetails.SpecialCorrespondenceRequirements,
			Notes:                             app.DeputyDetails.Notes,
		}

		vars := ManageAdditionalInformationVars{AppVars: app, Form: form}

		switch r.Method {
		case http.MethodGet:
			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			form := sirius.AdditionalInformationForm{
				InterpreterRequired: r.PostFormValue("interpreter-required"),
				SpecialCorrespondenceRequirements: sirius.SpecialCorrespondenceRequirementsType{
					AudioTape:                  r.PostFormValue("audio") == "on",
					LargePrint:                 r.PostFormValue("large-print") == "on",
					HearingImpaired:            r.PostFormValue("hearing-impaired") == "on",
					SpellingOfNameRequiresCare: r.PostFormValue("spelling-takes-care") == "on",
				},
				Notes:                 r.PostFormValue("notes"),
				CorrespondenceByWelsh: vars.DeputyDetails.CorrespondenceByWelsh,
				CorrespondenceByEmail: vars.DeputyDetails.CorrespondenceByEmail,
				CorrespondenceByPhone: vars.DeputyDetails.CorrespondenceByPhone,
				CorrespondenceByPost:  vars.DeputyDetails.CorrespondenceByPost,
			}

			err := client.UpdateAdditionalInformation(ctx, app.DeputyId(), form)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars.Errors = util.RenameErrors(verr.Errors)
				vars.Form = form
				return tmpl.ExecuteTemplate(w, "page", vars)
			} else if err != nil {
				return err
			}

			return Redirect(fmt.Sprintf("/%d?success=additionalInformation", app.DeputyId()))
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
