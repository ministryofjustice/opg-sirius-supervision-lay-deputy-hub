package sirius

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/mocks"
	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/shared"
	"github.com/stretchr/testify/assert"
)

func TestUpdateDeputyDeathRecord(t *testing.T) {
	mockClient := &mocks.MockClient{}

	client, _ := NewClient(mockClient, "http://localhost:3000")
	json := `{
    "1":                            "id",
    "2020-02-01":                      "dateOfDeath",
    "true":                     "proofOfDeathReceived",
    "2020-02-03":                     "dateDeathCertificateReceived",
    "2020-02-02":                     "dateNotified",
    "EXECUTOR":                             "notifiedBy",
    "LETTER":                           "notificationMethod",
    "Some notes":                         "notes",
	}`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	formData := DeputyDeathRecord{
		DeputyID:                      1,
		DateOfDeath:                   "2020-02-01",
		ProofOfDeathReceived:          "true",
		DateDeathCertificateReceieved: "2020-02-03",
		DateNotified:                  "2020-02-02",
		NotifiedBy:                    RefData{Handle: "EXECUTOR"},
		NotificationMethod:            "LETTER",
		Notes:                         "notes",
	}

	err := client.UpdateDeputyDeathRecord(getContext(nil), formData)
	assert.Nil(t, err)
}

func TestUpdateDeputyDeathRecordReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("{}"))
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	formData := DeputyDeathRecord{
		DeputyID:                      1,
		DateOfDeath:                   "2020-02-01",
		ProofOfDeathReceived:          "true",
		DateDeathCertificateReceieved: "2020-02-03",
		DateNotified:                  "2020-02-02",
		NotifiedBy:                    RefData{Handle: "EXECUTOR"},
		NotificationMethod:            "LETTER",
		Notes:                         "notes",
	}

	err := client.UpdateDeputyDeathRecord(getContext(nil), formData)

	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    fmt.Sprintf("%v/v1/deputies/lay/%d/death-notification", svr.URL+shared.SupervisionAPIPath, 1),
		Method: http.MethodPost,
	}, err)
}

func TestUpdateDeputyDeathRecordReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	formData := DeputyDeathRecord{
		DeputyID:                      1,
		DateOfDeath:                   "2020-02-01",
		ProofOfDeathReceived:          "true",
		DateDeathCertificateReceieved: "2020-02-03",
		DateNotified:                  "2020-02-02",
		NotifiedBy:                    RefData{Handle: "EXECUTOR"},
		NotificationMethod:            "LETTER",
		Notes:                         "notes",
	}

	err := client.UpdateDeputyDeathRecord(getContext(nil), formData)

	assert.Equal(t, ErrUnauthorized, err)
}
