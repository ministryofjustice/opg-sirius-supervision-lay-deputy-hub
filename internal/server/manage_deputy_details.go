package server

import (
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/sirius"
)

type DeputyDetailsInformation interface {
	UpdateDeputyDetails(sirius.Context, int, sirius.DeputyDetailsForm) error
}

type ManageDeputyDetailsVars struct {
	AppVars
}

type pair struct {
	k string
	v string
}

var validationMappings = map[string]map[string]pair{
	// address
	"addressLine1": {
		"stringLengthTooLong": pair{"addressLine1", "The building or street must be 255 characters or fewer"},
		"isEmpty":             pair{"addressLine1", "The building or street is required and can't be empty"},
	},
	"addressLine2": {
		"stringLengthTooLong": pair{"addressLine2", "Address line 2 must be 255 characters or fewer"},
	},
	"addressLine3": {
		"stringLengthTooLong": pair{"addressLine3", "Address line 3 must be 255 characters or fewer"},
	},
	"town": {
		"stringLengthTooLong": pair{"town", "The town or city must be 255 characters or fewer"},
		"isEmpty":             pair{"town", "The town or city is required and can't be empty"},
	},
	"county": {
		"stringLengthTooLong": pair{"county", "The county must be 255 characters or fewer"},
	},
	"postcode": {
		"stringLengthTooLong": pair{"postcode", "The postcode must be 255 characters or fewer"},
		"isEmpty":             pair{"postcode", "The postcode is required and can't be empty"},
	},
	"phoneNumber": {
		"stringLengthTooLong": pair{"phoneNumber", "The telephone number must be 255 characters or fewer"},
		"isEmpty":             pair{"phoneNumber", "Enter a telephone number"},
	},
	"workPhoneNumber": {
		"stringLengthTooLong": pair{"workPhoneNumber", "The telephone number must be 255 characters or fewer"},
	},
	"email": {
		"stringLengthTooLong":         pair{"email", "The email must be 255 characters or fewer"},
		"isEmpty":                     pair{"email", "Enter an email address"},
		"emailAddressInvalidFormat":   pair{"email", "Enter an email address in the correct format, like name@example.com"},
		"emailAddressInvalidHostname": pair{"email", "The email address is invalid"},
		"hostnameInvalidHostname":     pair{"email", "The email address is invalid"},
		"hostnameLocalNameNotAllowed": pair{"email", "The email address is invalid"},
	},
	"firstname": {
		"stringLengthTooLong": pair{"firstname", "The deputy first name must be 255 characters or fewer"},
		"isEmpty":             pair{"firstname", "The deputy first name is required and can't be empty"},
	},
	"surname": {
		"stringLengthTooLong": pair{"surname", "The deputy surname must be 255 characters or fewer"},
		"isEmpty":             pair{"surname", "The deputy surname is required and can't be empty"},
	},
}

func RenameErrors(siriusError sirius.ValidationErrors) sirius.ValidationErrors {
	mappedErrors := sirius.ValidationErrors{}
	for fieldName, value := range siriusError {
		for errorType, errorMessage := range value {
			err := make(map[string]string)
			if mapping, ok := validationMappings[fieldName][errorType]; ok {
				err[errorType] = mapping.v
				mappedErrors[mapping.k] = err
			} else {
				err[errorType] = errorMessage
				mappedErrors[fieldName] = err
			}
		}
	}
	return mappedErrors
}

func renderTemplateForManageDeputyDetailsForm(client DeputyDetailsInformation, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)
		app.PageName = "Manage deputy contact details"

		vars := ManageDeputyDetailsVars{AppVars: app}

		switch r.Method {
		case http.MethodGet:
			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			form := sirius.DeputyDetailsForm{
				DeputyFirstName:        r.PostFormValue("deputy-first-name"),
				DeputyMiddleName:       r.PostFormValue("deputy-middle-name"),
				DeputySurname:          r.PostFormValue("deputy-last-name"),
				DeputyOtherNames:       r.PostFormValue("deputy-other-names"),
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
				DeputySubType:          sirius.DeputySubType{SubType: "PERSON"},
			}

			err := client.UpdateDeputyDetails(ctx, app.DeputyId(), form)

			if verr, ok := err.(sirius.ValidationError); ok {
				vars.Errors = RenameErrors(verr.Errors)
				vars.DeputyDetails = sirius.DeputyDetails{
					ID:                     app.DeputyId(),
					DeputyFirstName:        form.DeputyFirstName,
					DeputyMiddleName:       form.DeputyMiddleName,
					DeputySurname:          form.DeputySurname,
					DeputyOtherNames:       form.DeputyOtherNames,
					DeputyDateOfBirth:      form.DeputyDateOfBirth,
					MobileTelephoneNumber:  form.MobileTelephoneNumber,
					DaytimeTelephoneNumber: form.DaytimeTelephoneNumber,
					EveningTelephoneNumber: form.EveningTelephoneNumber,
					Email:                  form.Email,
					AddressLine1:           form.AddressLine1,
					AddressLine2:           form.AddressLine2,
					AddressLine3:           form.AddressLine3,
					Town:                   form.Town,
					County:                 form.County,
					Postcode:               form.Postcode,
					Country:                form.Country,
					IsAirmailRequired:      form.IsAirmailRequired,
				}
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
