package sirius

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestDeputyClientsReturned(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := ` {
	  "clients": [
		{
			"id": 67,
			  "caseRecNumber": "67422477",
			  "email": "john.fearless@example.com",
			  "firstname": "John",
			  "surname": "Fearless",
			  "addressLine1": "94 Duckpit Lane",
			  "addressLine2": "Upper Oddington",
			  "addressLine3": "Canvey Island",
			  "town": "",
			  "county": "",
			  "postcode": "GL566WQ",
			  "country": "",
			  "phoneNumber": "07960209814",
			  "clientAccommodation": {
				"handle": "FAMILY MEMBER/FRIEND'S HOME",
				"label": "Family Member/Friend's Home (including spouse/civil partner)",
				"deprecated": false
			  },
			  "orders": [
				{
				  "id": 59,
				  "latestSupervisionLevel": {
					"appliesFrom": "01/12/2020",
					"supervisionLevel": {
					  "handle": "GENERAL",
					  "label": "General",
					  "deprecated": null
					}
				  },
				  "orderDate": "01/12/2020",
				  "orderStatus": {
					"handle": "ACTIVE",
					"label": "Active",
					"deprecated": false
				  }
				},
				{
				  "id": 60,
				  "latestSupervisionLevel": {
					 "appliesFrom": "01/12/2017",
					"supervisionLevel": {
					  "handle": "GENERAL",
					  "label": "General",
					  "deprecated": null
					}
				  },
				  "orderDate": "01/12/2017",
				  "orderStatus": {
					"handle": "ACTIVE",
					"label": "Active",
					"deprecated": false
				  }
				}
			  ],
			  "oldestNonLodgedAnnualReport": {
				"dueDate": "01/01/2016",
				"revisedDueDate": "01/05/2016",
				"status": {
				  "label": "Pending"
				}
			  },
			  "riskScore": 5,
				"hasActiveREMWarning": true
			}
		  ],
		  "pages": {
			"current": 1,
			"total": 1
		  },
		"metadata": {
			"totalActiveClients": 1
		  },
		  "total": 1
		} `

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := ClientList{
		Metadata:     Metadata{TotalActiveClients: 1},
		TotalClients: 1,
	}

	deputyClientDetails, err := client.GetDeputyClients(getContext(nil), ClientListParams{1})

	assert.Equal(t, 1, deputyClientDetails.Metadata.TotalActiveClients)
	assert.Equal(t, expectedResponse, deputyClientDetails)
	assert.Equal(t, nil, err)
}
