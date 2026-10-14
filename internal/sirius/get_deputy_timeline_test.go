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
			Id:              20,
			Timestamp:       "28/09/2026 17:13:12",
			TimestampAsDate: "28/09/2026",
			EventType:       "PersonDetailsChanged",
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
						FieldName:         "firstname",
						OldValue:          "",
						NewValue:          "test6aba9217958a1",
						Type:              "string",
						TimelineValueName: "First name",
					},
					{
						FieldName:         "surname",
						OldValue:          "",
						NewValue:          "Banks",
						Type:              "string",
						TimelineValueName: "Surname",
					},
					{
						FieldName:         "deputyType",
						OldValue:          "",
						NewValue:          "Lay",
						Type:              "string",
						TimelineValueName: "Type of deputy",
					},
					{
						FieldName:         "deputySubType",
						OldValue:          "",
						NewValue:          "Person",
						Type:              "string",
						TimelineValueName: "Deputy Subtype",
					},
				},
			},
		},
		{
			Id:              23,
			Timestamp:       "28/09/2026 17:13:12",
			TimestampAsDate: "28/09/2026",
			EventType:       "DeputyLinkedToOrder",
			User: User{
				Id:          61,
				PhoneNumber: "12345678",
				DisplayName: "case manager",
				Email:       "case.manager@opgtest.com",
			},
			Event: Event{
				OrderType:      "PFA",
				OrderUid:       "7000-0000-1334",
				OrderId:        "86",
				OrderCourtRef:  "31924130",
				CourtReference: "31924130",
				PersonType:     "Deputy",
				PersonId:       "49",
				PersonUid:      "7000-0000-1367",
				PersonName:     "Mrs test6aba9217958a1 Banks",
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
			Id:              20,
			Timestamp:       "28/09/2026 17:13:12",
			TimestampAsDate: "28/09/2026",
			EventType:       "PersonDetailsChanged",
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
						FieldName:         "firstname",
						OldValue:          "",
						NewValue:          "test6aba9217958a1",
						Type:              "string",
						TimelineValueName: "First name",
					},
					{
						FieldName:         "surname",
						OldValue:          "",
						NewValue:          "Banks",
						Type:              "string",
						TimelineValueName: "Surname",
					},
					{
						FieldName:         "deputyType",
						OldValue:          "",
						NewValue:          "Lay",
						Type:              "string",
						TimelineValueName: "Type of deputy",
					},
					{
						FieldName:         "deputySubType",
						OldValue:          "",
						NewValue:          "Person",
						Type:              "string",
						TimelineValueName: "Deputy Subtype",
					},
				},
			},
		},
		{
			Id:              23,
			Timestamp:       "28/09/2026 17:13:12",
			TimestampAsDate: "28/09/2026",
			EventType:       "DeputyLinkedToOrder",
			User: User{
				Id:          61,
				PhoneNumber: "12345678",
				DisplayName: "case manager",
				Email:       "case.manager@opgtest.com",
			},
			Event: Event{
				OrderType:      "PFA",
				OrderUid:       "7000-0000-1334",
				OrderId:        "86",
				OrderCourtRef:  "31924130",
				CourtReference: "31924130",
				PersonType:     "Deputy",
				PersonId:       "49",
				PersonUid:      "7000-0000-1367",
				PersonName:     "Mrs test6aba9217958a1 Banks",
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
				OrderType:      "pfa",
				OrderUid:       "7000-0000-1334",
				OrderId:        "86",
				OrderCourtRef:  "31924130",
				CourtReference: "31924130",
				PersonType:     "Deputy",
				PersonId:       "49",
				PersonUid:      "7000-0000-1367",
				PersonName:     "Mrs test6aba9217958a1 Banks",
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

func TestReformatEvent(t *testing.T) {
	tests := []struct {
		name          string
		myEvent       Event
		myEventType   string
		expectedEvent Event
	}{
		{
			"Death notification deleted event",
			Event{
				DateDeathCertificateReceived: "2026-10-28T16:13:12+00:00",
				DateNotified:                 "2026-09-28T16:13:12+00:00",
				DateOfDeath:                  "2026-09-29T16:13:12+00:00",
				Notes:                        "Notes about death notification being deleted",
				NotificationMethod:           "Email",
				NotifiedBy:                   "EXECUTOR",
				ParentEventType:              "PersonEvent",
				PersonCourtRef:               "",
				PersonId:                     "49",
				PersonName:                   "Mrs Betty Banks",
				PersonType:                   "Deputy",
				PersonUid:                    "7000-0000-1367",
				ProofOfDeathReceived:         false,
				AdditionalPersons:            []AdditionalPerson{},
			},
			"Opg\\Core\\Model\\Event\\Person\\DeathNotificationDeleted",
			Event{
				DateDeathCertificateReceived: "28/10/2026",
				DateNotified:                 "28/09/2026",
				DateOfDeath:                  "29/09/2026",
				Notes:                        "Notes about death notification being deleted",
				NotificationMethod:           "Email",
				NotifiedBy:                   "Executor",
				ParentEventType:              "PersonEvent",
				PersonCourtRef:               "",
				PersonId:                     "49",
				PersonName:                   "Mrs Betty Banks",
				PersonType:                   "Deputy",
				PersonUid:                    "7000-0000-1367",
				ProofOfDeathReceived:         false,
				AdditionalPersons:            []AdditionalPerson{},
			},
		},
		{
			"Person death notified event",
			Event{
				DateDeathCertificateReceived: "2026-02-20 00:00:00",
				DateNotified:                 "2026-01-20 00:00:00",
				DateOfDeath:                  "2026-01-19 00:00:00",
				Notes:                        "Notes about death notification",
				NotificationMethod:           "Phone",
				NotifiedBy:                   "OTHER",
				NotifiedByOther:              "Client notified us",
				ParentEventType:              "PersonEvent",
				PersonCourtRef:               "",
				PersonId:                     "49",
				PersonName:                   "Mrs Betty Banks",
				PersonType:                   "Deputy",
				PersonUid:                    "7000-0000-1367",
				ProofOfDeathReceived:         true,
				AdditionalPersons:            []AdditionalPerson{},
			},
			"Opg\\Core\\Model\\Event\\Person\\PersonDeathNotified",
			Event{
				DateDeathCertificateReceived: "20/02/2026",
				DateNotified:                 "20/01/2026",
				DateOfDeath:                  "19/01/2026",
				Notes:                        "Notes about death notification",
				NotificationMethod:           "Phone",
				NotifiedBy:                   "Other",
				NotifiedByOther:              "Client notified us",
				ParentEventType:              "PersonEvent",
				PersonCourtRef:               "",
				PersonId:                     "49",
				PersonName:                   "Mrs Betty Banks",
				PersonType:                   "Deputy",
				PersonUid:                    "7000-0000-1367",
				ProofOfDeathReceived:         true,
				AdditionalPersons:            []AdditionalPerson{},
			},
		},
		{
			"Address changed event",
			Event{
				PersonCourtRef: "",
				PersonId:       "49",
				PersonName:     "Mrs Betty Banks",
				PersonType:     "Deputy",
				PersonUid:      "7000-0000-1367",
				Changes: []Changeset{
					{
						FieldName: "addressLine1",
						OldValue:  "",
						NewValue:  "Flat 7 Bollerworth",
						Type:      "string",
					},
					{
						FieldName: "town",
						OldValue:  "",
						NewValue:  "Birmingham",
						Type:      "string",
					},
					{
						FieldName: "isAirmailRequired",
						OldValue:  "",
						NewValue:  "1",
						Type:      "string",
					},
				},
			},
			"Opg\\Core\\Model\\Event\\Person\\AddressChanged",
			Event{
				PersonCourtRef: "",
				PersonId:       "49",
				PersonName:     "Mrs Betty Banks",
				PersonType:     "Deputy",
				PersonUid:      "7000-0000-1367",
				Changes: []Changeset{
					{
						FieldName:         "addressLine1",
						OldValue:          "",
						NewValue:          "Flat 7 Bollerworth",
						Type:              "string",
						TimelineValueName: "Address line 1",
					},
					{
						FieldName:         "town",
						OldValue:          "",
						NewValue:          "Birmingham",
						Type:              "string",
						TimelineValueName: "Town",
					},
					{
						FieldName:         "isAirmailRequired",
						OldValue:          "",
						NewValue:          "1",
						Type:              "string",
						TimelineValueName: "Is airmail required",
					},
				},
			},
		},
		{
			"Unlink",
			Event{
				CourtReference: "31924130",
				OrderCourtRef:  "31924130",
				OrderId:        "86",
				OrderType:      "pfa",
				OrderUid:       "7000-0000-1334",
				PersonCourtRef: "",
				PersonId:       "49",
				PersonName:     "Mrs Betty Banks",
				PersonType:     "Deputy",
				PersonUid:      "7000-0000-1367",
				AdditionalPersons: []AdditionalPerson{
					{
						PersonCourtRef: "28498842",
						PersonId:       "47",
						PersonName:     "Mr Luke Crete",
						PersonType:     "Client",
						PersonUid:      "7000-0000-1243",
					},
				},
			},
			"Opg\\Core\\Model\\Event\\Order\\DeputyUnlinkedFromOrder",
			Event{
				CourtReference: "31924130",
				OrderCourtRef:  "31924130",
				OrderId:        "86",
				OrderType:      "PFA",
				OrderUid:       "7000-0000-1334",
				PersonCourtRef: "",
				PersonId:       "49",
				PersonName:     "Mrs Betty Banks",
				PersonType:     "Deputy",
				PersonUid:      "7000-0000-1367",
				AdditionalPersons: []AdditionalPerson{
					{
						PersonCourtRef: "28498842",
						PersonId:       "47",
						PersonName:     "Mr Luke Crete",
						PersonType:     "Client",
						PersonUid:      "7000-0000-1243",
					},
				},
			},
		},
		{
			"Order details changed",
			Event{
				CourtReference: "31924130",
				OrderId:        "86",
				OrderType:      "pfa",
				OrderUid:       "7000-0000-1334",
				PersonCourtRef: "",
				PersonId:       "49",
				PersonName:     "Mrs Betty Banks",
				PersonType:     "Deputy",
				PersonUid:      "7000-0000-1367",
				AdditionalPersons: []AdditionalPerson{
					{
						PersonCourtRef: "28498842",
						PersonId:       "47",
						PersonName:     "Mr Luke Crete",
						PersonType:     "Client",
						PersonUid:      "7000-0000-1243",
					},
				},
				Changes: []Changeset{
					{
						FieldName: "statusOnCase",
						OldValue:  "ERROR",
						NewValue:  "DECEASED",
						Type:      "string",
					},
					{
						FieldName: "statusChangeDate",
						OldValue:  "",
						NewValue:  "2026-10-12T00:00:00+00:00",
						Type:      "string",
					},
					{
						FieldName: "mainCorrespondent",
						OldValue:  "",
						NewValue:  "1",
						Type:      "string",
					},
					{
						FieldName: "relationshipToClient",
						OldValue:  "",
						NewValue:  "SPOUSE",
						Type:      "string",
					},
				},
			},
			"Opg\\Core\\Model\\Event\\Order\\DeputyOrderDetailsChanged",
			Event{
				CourtReference: "31924130",
				OrderId:        "86",
				OrderType:      "PFA",
				OrderUid:       "7000-0000-1334",
				PersonCourtRef: "",
				PersonId:       "49",
				PersonName:     "Mrs Betty Banks",
				PersonType:     "Deputy",
				PersonUid:      "7000-0000-1367",
				AdditionalPersons: []AdditionalPerson{
					{
						PersonCourtRef: "28498842",
						PersonId:       "47",
						PersonName:     "Mr Luke Crete",
						PersonType:     "Client",
						PersonUid:      "7000-0000-1243",
					},
				},
				Changes: []Changeset{
					{
						FieldName:         "statusOnCase",
						OldValue:          "Error",
						NewValue:          "Deceased",
						Type:              "string",
						TimelineValueName: "Status on case",
					},
					{
						FieldName:         "statusChangeDate",
						OldValue:          "",
						NewValue:          "12/10/2026",
						Type:              "string",
						TimelineValueName: "Status change date",
					},
					{
						FieldName:         "mainCorrespondent",
						OldValue:          "",
						NewValue:          "1",
						Type:              "string",
						TimelineValueName: "Correspondent",
					},
					{
						FieldName:         "relationshipToClient",
						OldValue:          "",
						NewValue:          "Spouse",
						Type:              "string",
						TimelineValueName: "Relationship to client",
					},
				},
			},
		},
		{
			"Person details changed",
			Event{
				PersonCourtRef:    "",
				PersonId:          "49",
				PersonName:        "Mrs Betty Banks",
				PersonType:        "Deputy",
				PersonUid:         "7000-0000-1367",
				AdditionalPersons: []AdditionalPerson{},
				Changes: []Changeset{
					{
						FieldName: "salutation",
						OldValue:  "",
						NewValue:  "Mrs",
						Type:      "string",
					},
					{
						FieldName: "firstname",
						OldValue:  "",
						NewValue:  "Betty",
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
			"Opg\\Core\\Model\\Event\\Person\\PersonDetailsChanged",
			Event{
				PersonCourtRef:    "",
				PersonId:          "49",
				PersonName:        "Mrs Betty Banks",
				PersonType:        "Deputy",
				PersonUid:         "7000-0000-1367",
				AdditionalPersons: []AdditionalPerson{},
				Changes: []Changeset{
					{
						FieldName:         "salutation",
						OldValue:          "",
						NewValue:          "Mrs",
						Type:              "string",
						TimelineValueName: "Title",
					},
					{
						FieldName:         "firstname",
						OldValue:          "",
						NewValue:          "Betty",
						Type:              "string",
						TimelineValueName: "First name",
					},
					{
						FieldName:         "deputyType",
						OldValue:          "",
						NewValue:          "Lay",
						Type:              "string",
						TimelineValueName: "Type of deputy",
					},
					{
						FieldName:         "deputySubType",
						OldValue:          "",
						NewValue:          "Person",
						Type:              "string",
						TimelineValueName: "Deputy Subtype",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := reformatEvent(tt.myEvent, tt.myEventType)
			assert.Equal(t, tt.expectedEvent, got)
		})
	}
}

//func TestReformatEventType(t *testing.T) {
//	tests := []struct {
//		name               string
//		assurances         []model.Assurance
//		want               bool
//		wantedErrorMessage string
//	}{
//		{
//			"No assurances",
//			[]model.Assurance{},
//			false,
//			"",
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got := reformatEventType(tt.assurances)
//			if got != tt.want {
//				t.Errorf("isAddVisitDisabled() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
