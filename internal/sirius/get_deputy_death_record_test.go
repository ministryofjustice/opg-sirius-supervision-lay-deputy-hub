package sirius

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDeputyDeathRecord(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/supervision-api/v1/deputies/2/death-notification", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": 2,
			"proofOfDeathReceived": true,
			"dateDeathCertificateReceived": "01/06/2025",
			"dateNotified": "01/02/2025",
			"notifiedBy": {
				"handle": "DEPUTY",
				"label": "Deputy"
			},
			"notificationMethod": "Phone",
			"notes": "Some notes",
			"person": {
				"dateOfDeath": "01/01/2020"
			}
		}`))
	}))
	t.Cleanup(server.Close)

	client, err := NewClient(server.Client(), server.URL)
	require.NoError(t, err)

	deathRecord, err := client.GetDeputyDeathRecord(getContext(nil), 2)

	require.NoError(t, err)
	assert.Equal(t, DeputyDeathRecord{
		DeputyID:                      2,
		DateOfDeath:                   "2020-01-01",
		ProofOfDeathReceived:          "true",
		DateDeathCertificateReceieved: "2025-06-01",
		DateNotified:                  "2025-02-01",
		NotifiedBy:                    RefData{Handle: "DEPUTY", Label: "Deputy"},
		NotificationMethod:            "Phone",
		Notes:                         "Some notes",
	}, deathRecord)
}

func TestGetDeputyDeathRecordReturnsNilWhenReceives404(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("{}"))
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	deathRecord, err := client.GetDeputyDeathRecord(getContext(nil), 2)

	assert.Equal(t, deathRecord, DeputyDeathRecord{})
	require.NoError(t, err)
}

func TestGetDeputyDeathRecordReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	deathRecord, err := client.GetDeputyDeathRecord(getContext(nil), 2)

	assert.Equal(t, deathRecord, DeputyDeathRecord{})
	assert.Equal(t, ErrUnauthorized, err)
}
