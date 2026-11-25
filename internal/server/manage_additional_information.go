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

func renderTemplateForManageAdditionalInformationGet(client AdditionalInformationInformation, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		app.PageName = "Manage additional information"
		form := sirius.AdditionalInformationForm{
			InterpreterRequired:               app.DeputyDetails.InterpreterRequired,
			SpecialCorrespondenceRequirements: app.DeputyDetails.SpecialCorrespondenceRequirements,
			CorrespondenceByWelsh:             app.DeputyDetails.CorrespondenceByWelsh,
			Notes:                             app.DeputyDetails.Notes,
		}

		vars := ManageAdditionalInformationVars{AppVars: app, Form: form}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}

func renderTemplateForManageAdditionalInformationPost(client AdditionalInformationInformation, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)
		app.PageName = "Manage additional information"

		form := sirius.AdditionalInformationForm{
			InterpreterRequired: r.PostFormValue("interpreter-required"),
			SpecialCorrespondenceRequirements: sirius.SpecialCorrespondenceRequirementsType{
				AudioTape:       r.PostFormValue("audio") == "on",
				LargePrint:      r.PostFormValue("large-print") == "on",
				HearingImpaired: r.PostFormValue("hearing-impaired") == "on",
			},
			Notes:                 r.PostFormValue("notes"),
			CorrespondenceByWelsh: r.PostFormValue("welsh") == "on",
			CorrespondenceByEmail: app.DeputyDetails.CorrespondenceByEmail,
			CorrespondenceByPhone: app.DeputyDetails.CorrespondenceByPhone,
			CorrespondenceByPost:  app.DeputyDetails.CorrespondenceByPost,
		}

		vars := ManageAdditionalInformationVars{AppVars: app, Form: form}

		err := client.UpdateAdditionalInformation(ctx, app.DeputyID(), form)

		if verr, ok := err.(sirius.ValidationError); ok {
			vars.Errors = util.RenameErrors(verr.Errors)
			vars.Form = form
			return tmpl.ExecuteTemplate(w, "page", vars)
		} else if err != nil {
			return err
		}

		return Redirect(fmt.Sprintf("/%d?success=additionalInformation", app.DeputyID()))
	}
}
