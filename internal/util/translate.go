package util

import "github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/sirius"

type pair struct {
	k string
	v string
}

var validationMappings = map[string]map[string]pair{
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
	"country": {
		"stringLengthTooLong": pair{"country", "The country must be 255 characters or fewer"},
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
	"middleName": {
		"stringLengthTooLong": pair{"middleName", "The deputy middle name must be 255 characters or fewer"},
	},
	"surname": {
		"stringLengthTooLong": pair{"surname", "The deputy surname must be 255 characters or fewer"},
		"isEmpty":             pair{"surname", "The deputy surname is required and can't be empty"},
	},
	"otherNames": {
		"stringLengthTooLong": pair{"otherNames", "The deputy other name must be 255 characters or fewer"},
	},
	"homePhoneNumber": {
		"stringLengthTooLong": pair{"homePhoneNumber", "The evening telephone number must be 255 characters or fewer"},
	},
	"mobileNumber": {
		"stringLengthTooLong": pair{"mobileNumber", "The mobile number must be 255 characters or fewer"},
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
