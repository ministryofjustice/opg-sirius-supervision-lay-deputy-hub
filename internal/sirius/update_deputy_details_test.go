package sirius

import (
	"bytes"
	stdjson "encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/mocks"
	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/shared"
	"github.com/stretchr/testify/assert"
)

func TestUpdateDeputyDetails(t *testing.T) {
	mockClient := &mocks.MockClient{}

	client, _ := NewClient(mockClient, "http://localhost:3000")
	r := io.NopCloser(bytes.NewReader([]byte(`{}`)))

	mocks.GetDoFunc = func(req *http.Request) (*http.Response, error) {
		body, err := io.ReadAll(req.Body)
		assert.NoError(t, err)
		var payload map[string]interface{}
		assert.NoError(t, stdjson.Unmarshal(body, &payload))
		assert.Equal(t, map[string]interface{}{
			"salutation": "Dr", "firstname": "Alex", "middleNames": "Jordan",
			"surname": "Smith", "otherNames": "Alexander", "dob": "03/02/1985",
			"email": "email", "mobileNumber": "mobile", "phoneNumber": "telephone",
			"eveningNumber": "evening", "addressLine1": "address-line-1",
			"addressLine2": "address-line-2", "addressLine3": "address-line-3",
			"town": "town", "county": "county", "country": "country",
			"postcode": "postcode", "isAirmailRequired": true,
		}, payload)
		assert.Equal(t, http.MethodPut, req.Method)
		assert.Equal(t, "/supervision-api/v1/deputies/1/contact-details", req.URL.Path)

		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	formData := DeputyDetailsForm{
		Title:                  "Dr",
		FirstName:              "Alex",
		MiddleName:             "Jordan",
		Surname:                "Smith",
		OtherNames:             "Alexander",
		DateOfBirth:            shared.NewDate("03/02/1985"),
		Email:                  "email",
		MobileTelephoneNumber:  "mobile",
		DaytimeTelephoneNumber: "telephone",
		EveningTelephoneNumber: "evening",
		AddressLine1:           "address-line-1",
		AddressLine2:           "address-line-2",
		AddressLine3:           "address-line-3",
		Town:                   "town",
		County:                 "county",
		Country:                "country",
		Postcode:               "postcode",
		IsAirmailRequired:      true,
	}

	err := client.UpdateDeputyDetails(getContext(nil), 1, formData)
	assert.Nil(t, err)
}

func TestUpdateDeputyContactDetailsReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("{}"))
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	formData := DeputyDetailsForm{
		Email:                  "email",
		DaytimeTelephoneNumber: "telephone",
		AddressLine1:           "address-line-1",
		AddressLine2:           "address-line-2",
		AddressLine3:           "address-line-3",
		Town:                   "town",
		County:                 "county",
		Postcode:               "postcode",
	}

	err := client.UpdateDeputyDetails(getContext(nil), 1, formData)

	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    fmt.Sprintf("%v/v1/deputies/%d/contact-details", svr.URL+shared.SupervisionAPIPath, 1),
		Method: http.MethodPut,
	}, err)
}

func TestUpdateDeputyContactDetailsReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	formData := DeputyDetailsForm{
		Email:                  "email",
		DaytimeTelephoneNumber: "telephone",
		AddressLine1:           "address-line-1",
		AddressLine2:           "address-line-2",
		AddressLine3:           "address-line-3",
		Town:                   "town",
		County:                 "county",
		Postcode:               "postcode",
	}

	err := client.UpdateDeputyDetails(getContext(nil), 1, formData)

	assert.Equal(t, ErrUnauthorized, err)
}
