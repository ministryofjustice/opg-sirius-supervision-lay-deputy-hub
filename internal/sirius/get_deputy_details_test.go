package sirius

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDeputyDetailsReturnsDeputy(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/supervision-api/v1/deputies/1", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": 1,
			"salutation": "Mr",
			"firstName": "Test",
			"surname": "Deputy",
			"otherNames": null,
			"deputyNumber": 46382901,
			"deputyStatus": "Active",
			"deputyType": {
				"handle": "LAY",
				"label": "Lay Deputy"
			},
			"dob": "01/01/1980",
			"email": "testemail@hotmail.co.uk",
			"mobileNumber": "0771 2345678",
			"phoneNumber": "0115 876 5574",
			"eveningNumber": "0115 2767825",
			"addressLine1": "Seax House",
			"addressLine2": "19 Market Rd",
			"addressLine3": null,
			"town": "Chelmsford",
			"county": "Essex",
			"postcode": "CM1 1GG",
			"country": null,
			"isAirmailRequired": false,
			"InterpreterRequired": null
		}`))
	}))
	t.Cleanup(server.Close)

	client, err := NewClient(server.Client(), server.URL)
	require.NoError(t, err)

	deputyDetails, err := client.GetDeputyDetails(getContext(nil), 1)

	require.NoError(t, err)
	assert.Equal(t, DeputyDetails{
		ID:               1,
		Salutation:       "Mr",
		DeputyFirstName:  "Test",
		DeputySurname:    "Deputy",
		DeputyOtherNames: "",
		DeputyNumber:     46382901,
		DeputyStatus:     "Active",
		DeputyType: DeputyType{
			Handle: "LAY",
			Label:  "Lay Deputy",
		},
		DeputyDateOfBirth:      "01/01/1980",
		Email:                  "testemail@hotmail.co.uk",
		MobileTelephoneNumber:  "0771 2345678",
		DaytimeTelephoneNumber: "0115 876 5574",
		EveningTelephoneNumber: "0115 2767825",
		AddressLine1:           "Seax House",
		AddressLine2:           "19 Market Rd",
		AddressLine3:           "",
		Town:                   "Chelmsford",
		County:                 "Essex",
		Postcode:               "CM1 1GG",
		IsAirmailRequired:      false,
		InterpreterRequired:    "",
	}, deputyDetails)
}

func TestGetDeputyDetailsReturnsActiveSpecialCorrespondenceRequirements(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/supervision-api/v1/deputies/1", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": 1,
			"specialCorrespondenceRequirements": {
				"audioTape": true,
				"largePrint": false,
				"hearingImpaired": true,
				"spellingOfNameRequiresCare": false
			}
		}`))
	}))
	t.Cleanup(server.Close)

	client, err := NewClient(server.Client(), server.URL)
	require.NoError(t, err)

	deputyDetails, err := client.GetDeputyDetails(getContext(nil), 1)

	require.NoError(t, err)
	assert.Equal(t, []string{"Audio", "Hearing impaired"}, deputyDetails.SpecialCorrespondenceRequirements.ActiveSpecialCorrespondenceRequirements())
}
