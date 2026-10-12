package sirius

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDeputyTimeline(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/supervision-api/v1/timeline/1/lay-deputy", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
    {
        "id": 20,
        "hash": "K",
        "timestamp": "2026-09-28T16:13:12+00:00",
        "eventType": "Opg\\Core\\Model\\Event\\Person\\PersonDetailsChanged",
        "user": {
            "id": 61,
            "phoneNumber": "12345678",
            "displayName": "case manager",
            "email": "case.manager@opgtest.com"
        },
        "event": {
            "personType": "Deputy",
            "personId": "49",
            "personUid": "7000-0000-1367",
            "personName": "Mrs test6aba9217958a1 Banks",
            "personCourtRef": null,
            "changes": [
                {
                    "fieldName": "firstname",
                    "oldValue": null,
                    "newValue": "test6aba9217958a1",
                    "type": "string"
                },
                {
                    "fieldName": "surname",
                    "oldValue": null,
                    "newValue": "Banks",
                    "type": "string"
                },
                {
                    "fieldName": "deputyType",
                    "oldValue": "",
                    "newValue": "LAY",
                    "type": "string"
                },
                {
                    "fieldName": "deputySubType",
                    "oldValue": "",
                    "newValue": "PERSON",
                    "type": "string"
                }
            ]
        }
    },
    {
        "id": 23,
        "hash": "N",
        "timestamp": "2026-09-28T16:13:12+00:00",
        "eventType": "Opg\\Core\\Model\\Event\\Order\\DeputyLinkedToOrder",
        "user": {
            "id": 61,
            "phoneNumber": "12345678",
            "displayName": "case manager",
            "email": "case.manager@opgtest.com"
        },
        "event": {
            "orderType": "pfa",
            "orderUid": "7000-0000-1334",
            "orderId": "86",
            "orderCourtRef": "31924130",
            "courtReferenceNumber": "31924130",
            "courtReference": "31924130",
            "personType": "Deputy",
            "personId": "49",
            "personUid": "7000-0000-1367",
            "personName": "Mrs test6aba9217958a1 Banks",
            "personCourtRef": null,
            "additionalPersons": [
                {
                    "personType": "Client",
                    "personId": "47",
                    "personUid": "7000-0000-1243",
                    "personName": "Mr Luke Crete",
                    "personCourtRef": "28498842"
                }
            ]
        }
    }]`))
	}))
	t.Cleanup(server.Close)

	client, err := NewClient(server.Client(), server.URL)
	require.NoError(t, err)

	deputyDetails, err := client.GetDeputyTimeline(getContext(nil), 1)

	require.NoError(t, err)
	assert.Equal(t, []TimelineEvent{
		{
			Id:        20,
			Timestamp: "28/09/2026 17:13:12",
			EventType: "PersonDetailsChanged",
			User: User{
				Id:          61,
				PhoneNumber: "12345678",
				DisplayName: "case manager",
				Email:       "case.manager@opgtest.com",
			},
			Event: Event{
				PersonType: "Deputy",
				PersonId:   "49",
				PersonUid:  "7000-0000-1367",
				PersonName: "Mrs test6aba9217958a1 Banks",
				Changes: []Changeset{
					{
						FieldName: "firstname",
						OldValue:  "",
						NewValue:  "test6aba9217958a1",
						Type:      "string",
					},
					{
						FieldName: "surname",
						OldValue:  "",
						NewValue:  "Banks",
						Type:      "string",
					},
					{
						FieldName: "deputyType",
						OldValue:  "",
						NewValue:  "LAY",
						Type:      "string",
					},
					{
						FieldName: "deputySubType",
						OldValue:  "",
						NewValue:  "PERSON",
						Type:      "string",
					},
				},
			},
		},
		{
			Id:        23,
			Timestamp: "28/09/2026 17:13:12",
			EventType: "DeputyLinkedToOrder",
			User: User{
				Id:          61,
				PhoneNumber: "12345678",
				DisplayName: "case manager",
				Email:       "case.manager@opgtest.com",
			},
			Event: Event{
				OrderType:            "pfa",
				OrderUid:             "7000-0000-1334",
				OrderId:              "86",
				OrderCourtRef:        "31924130",
				CourtReferenceNumber: "31924130",
				CourtReference:       "31924130",
				PersonType:           "Deputy",
				PersonId:             "49",
				PersonUid:            "7000-0000-1367",
				PersonName:           "Mrs test6aba9217958a1 Banks",
				AdditionalPersons: []AdditionalPerson{
					{
						PersonType:     "Client",
						PersonId:       "47",
						PersonUid:      "7000-0000-1243",
						PersonName:     "Mr Luke Crete",
						PersonCourtRef: "28498842",
					},
				},
			},
		},
	}, deputyDetails)
}

func TestEditDeputyEvents(t *testing.T) {
	expectedEvents := []TimelineEvent{
		{
			Id:        20,
			Timestamp: "28/09/2026 17:13:12",
			EventType: "PersonDetailsChanged",
			User: User{
				Id:          61,
				PhoneNumber: "12345678",
				DisplayName: "case manager",
				Email:       "case.manager@opgtest.com",
			},
			Event: Event{
				PersonType: "Deputy",
				PersonId:   "49",
				PersonUid:  "7000-0000-1367",
				PersonName: "Mrs test6aba9217958a1 Banks",
				Changes: []Changeset{
					{
						FieldName: "firstname",
						OldValue:  "",
						NewValue:  "test6aba9217958a1",
						Type:      "string",
					},
					{
						FieldName: "surname",
						OldValue:  "",
						NewValue:  "Banks",
						Type:      "string",
					},
					{
						FieldName: "deputyType",
						OldValue:  "",
						NewValue:  "LAY",
						Type:      "string",
					},
					{
						FieldName: "deputySubType",
						OldValue:  "",
						NewValue:  "PERSON",
						Type:      "string",
					},
				},
			},
		},
		{
			Id:        23,
			Timestamp: "28/09/2026 17:13:12",
			EventType: "DeputyLinkedToOrder",
			User: User{
				Id:          61,
				PhoneNumber: "12345678",
				DisplayName: "case manager",
				Email:       "case.manager@opgtest.com",
			},
			Event: Event{
				OrderType:            "pfa",
				OrderUid:             "7000-0000-1334",
				OrderId:              "86",
				OrderCourtRef:        "31924130",
				CourtReferenceNumber: "31924130",
				CourtReference:       "31924130",
				PersonType:           "Deputy",
				PersonId:             "49",
				PersonUid:            "7000-0000-1367",
				PersonName:           "Mrs test6aba9217958a1 Banks",
				AdditionalPersons: []AdditionalPerson{
					{
						PersonType:     "Client",
						PersonId:       "47",
						PersonUid:      "7000-0000-1243",
						PersonName:     "Mr Luke Crete",
						PersonCourtRef: "28498842",
					},
				},
			},
		},
	}

	actualEvents := editDeputyEvents([]TimelineEvent{
		{
			Id:        20,
			Timestamp: "2026-09-28T16:13:12+00:00",
			EventType: "Opg\\\\Core\\\\Model\\\\Event\\\\Person\\\\PersonDetailsChanged",
			User: User{
				Id:          61,
				PhoneNumber: "12345678",
				DisplayName: "case manager",
				Email:       "case.manager@opgtest.com",
			},
			Event: Event{
				PersonType: "Deputy",
				PersonId:   "49",
				PersonUid:  "7000-0000-1367",
				PersonName: "Mrs test6aba9217958a1 Banks",
				Changes: []Changeset{
					{
						FieldName: "firstname",
						OldValue:  "",
						NewValue:  "test6aba9217958a1",
						Type:      "string",
					},
					{
						FieldName: "surname",
						OldValue:  "",
						NewValue:  "Banks",
						Type:      "string",
					},
					{
						FieldName: "deputyType",
						OldValue:  "",
						NewValue:  "LAY",
						Type:      "string",
					},
					{
						FieldName: "deputySubType",
						OldValue:  "",
						NewValue:  "PERSON",
						Type:      "string",
					},
				},
			},
		},
		{
			Id:        23,
			Timestamp: "2026-09-28T16:13:12+00:00",
			EventType: "Opg\\\\Core\\\\Model\\\\Event\\\\Order\\\\DeputyLinkedToOrder",
			User: User{
				Id:          61,
				PhoneNumber: "12345678",
				DisplayName: "case manager",
				Email:       "case.manager@opgtest.com",
			},
			Event: Event{
				OrderType:            "pfa",
				OrderUid:             "7000-0000-1334",
				OrderId:              "86",
				OrderCourtRef:        "31924130",
				CourtReferenceNumber: "31924130",
				CourtReference:       "31924130",
				PersonType:           "Deputy",
				PersonId:             "49",
				PersonUid:            "7000-0000-1367",
				PersonName:           "Mrs test6aba9217958a1 Banks",
				AdditionalPersons: []AdditionalPerson{
					{
						PersonType:     "Client",
						PersonId:       "47",
						PersonUid:      "7000-0000-1243",
						PersonName:     "Mr Luke Crete",
						PersonCourtRef: "28498842",
					},
				},
			},
		},
	})

	assert.Equal(t, expectedEvents, actualEvents)
}
