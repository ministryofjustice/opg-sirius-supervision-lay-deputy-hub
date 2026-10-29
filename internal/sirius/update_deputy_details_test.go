package sirius

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestUpdateDeputyDetails(t *testing.T) {
	mockClient := &mocks.MockClient{}

	client, _ := NewClient(mockClient, "http://localhost:3000")
	json := `{
    "Email":                            "email",
    "PhoneNumber":                      "telephone",
    "AddressLine1":                     "address-line-1",
    "AddressLine2":                     "address-line-2",
    "AddressLine3":                     "address-line-3",
    "Town":                             "town",
    "County":                           "county",
    "Postcode":                         "postcode",
	}`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

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
		URL:    fmt.Sprintf("%v/v1/deputies/%d/contact-details", svr.URL+SupervisionAPIPath, 1),
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
