package sirius

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/mocks"
	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/shared"
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
			  "relationshipToClient": "Sibling",
			  "statusOnCase": "Active",
			  "appointmentType": "Joint and Severally",
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
				  "casesubtype": "pfa",
				  "appointmentType": "Joint",
				  "orderDate": "01/12/2020",
				  "orderStatus": {
					"label": "Active"
				  }
				},
				{
				  "id": 60,
				  "casesubtype": "hw",
				  "appointmentType": "Sole",
				  "orderDate": "01/12/2017",
				  "orderStatus": {
					"label": "Closed"
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

	mocks.GetDoFunc = func(req *http.Request) (*http.Response, error) {
		assert.Empty(t, req.URL.RawQuery)
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := ClientList{
		Clients: []DeputyClient{{
			ClientId:        67,
			Firstname:       "John",
			Surname:         "Fearless",
			CourtRef:        "67422477",
			Relationship:    "Sibling",
			DeputyStatus:    "Active",
			AppointmentType: "Joint and Severally",
			Orders: []Order{
				{OrderDate: shared.NewDate("01/12/2020"), OrderStatus: OrderStatus{Label: "Active"}, CaseSubType: "pfa", AppointmentType: "Joint"},
				{OrderDate: shared.NewDate("01/12/2017"), OrderStatus: OrderStatus{Label: "Closed"}, CaseSubType: "hw", AppointmentType: "Sole"},
			},
		}},
		Metadata:     Metadata{TotalActiveClients: 1},
		TotalClients: 1,
	}

	deputyClientDetails, err := client.GetDeputyClients(getContext(nil), ClientListParams{DeputyId: 1})

	assert.Equal(t, 1, deputyClientDetails.Metadata.TotalActiveClients)
	assert.Equal(t, expectedResponse, deputyClientDetails)
	assert.Equal(t, nil, err)
}

func TestOrderDateDeserialization(t *testing.T) {
	for _, test := range []struct {
		name     string
		value    string
		expected string
	}{
		{name: "Sirius date", value: `"01/12/2020"`, expected: "01/12/2020"},
		{name: "ISO date", value: `"2020-12-01"`, expected: "01/12/2020"},
		{name: "timestamp", value: `"2020-12-01T12:30:00+00:00"`, expected: "01/12/2020"},
		{name: "empty", value: `""`},
		{name: "null", value: `null`},
	} {
		t.Run(test.name, func(t *testing.T) {
			var order Order
			err := json.Unmarshal([]byte(`{"orderDate":`+test.value+`}`), &order)

			assert.NoError(t, err)
			assert.Equal(t, test.expected, order.OrderDate.DateString())
			assert.Equal(t, test.expected == "", order.OrderDate.IsNull())
		})
	}

	t.Run("invalid date", func(t *testing.T) {
		var order Order
		assert.Error(t, json.Unmarshal([]byte(`{"orderDate":"not-a-date"}`), &order))
	})
}
