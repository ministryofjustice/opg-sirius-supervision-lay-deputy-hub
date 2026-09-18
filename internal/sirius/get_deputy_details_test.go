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
			"deputyNumber": 46382901,
			"deputyStatus": "Active",
			"deputyType": {
				"handle": "LAY",
				"label": "Lay Deputy"
			}
		}`))
	}))
	t.Cleanup(server.Close)

	client, err := NewClient(server.Client(), server.URL)
	require.NoError(t, err)

	deputyDetails, err := client.GetDeputyDetails(getContext(nil), 1)

	require.NoError(t, err)
	assert.Equal(t, DeputyDetails{
		ID:              1,
		Salutation:      "Mr",
		DeputyFirstName: "Test",
		DeputySurname:   "Deputy",
		DeputyNumber:    46382901,
		DeputyStatus:    "Active",
		DeputyType: DeputyType{
			Handle: "LAY",
			Label:  "Lay Deputy",
		},
	}, deputyDetails)
}
