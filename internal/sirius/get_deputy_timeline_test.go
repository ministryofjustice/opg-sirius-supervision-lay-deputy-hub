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
            "personName": "Mrs Betty Banks",
            "personCourtRef": null,
            "changes": [
                {
                    "fieldName": "firstname",
                    "oldValue": null,
                    "newValue": "Betty",
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
            "personName": "Mrs Betty Banks",
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
    },
    {
        "id": 367,
        "hash": "A7",
        "timestamp": "2026-10-12T12:35:35+00:00",
        "eventType": "Opg\\Core\\Model\\Event\\Client\\ClientFeePayerRemoved",
        "user": {
            "id": 101,
            "phoneNumber": "12345678",
            "displayName": "PROTeam1 User1",
            "email": "pro1@opgtest.com"
        },
        "event": {
            "multiPersons": [
                {
                    "personType": "Deputy",
                    "personId": "49",
                    "personUid": "7000-0000-1367",
                    "personName": "Mrs Betty Banks",
                    "personCourtRef": null
                },
                {
                    "personType": "Client",
                    "personId": "47",
                    "personUid": "7000-0000-1243",
                    "personName": "Mr Luke Crete",
                    "personCourtRef": "28498842"
                }
            ],
            "parentEventType": "MultiPersonEvent"
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
				PersonName: "Mrs Betty Banks",
				Changes: []Changeset{
					{
						FieldName:         "firstname",
						OldValue:          "",
						NewValue:          "Betty",
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
				PersonName:     "Mrs Betty Banks",
				AdditionalPersons: []AdditionalPerson{
					{
						PersonType: "Client",
						PersonId:   "47",
						PersonUid:  "7000-0000-1243",
						PersonName: "Mr Luke Crete",
					},
				},
			},
		},
		{
			Id:              367,
			Timestamp:       "12/10/2026 13:35:35",
			TimestampAsDate: "12/10/2026",
			EventType:       "ClientFeePayerRemoved",
			User: User{
				Id:          101,
				PhoneNumber: "12345678",
				DisplayName: "PROTeam1 User1",
				Email:       "pro1@opgtest.com",
			},
			Event: Event{
				MultiPersons: []AdditionalPerson{
					{
						PersonType: "Deputy",
						PersonId:   "49",
						PersonUid:  "7000-0000-1367",
						PersonName: "Mrs Betty Banks",
					},
					{
						PersonType: "Client",
						PersonId:   "47",
						PersonUid:  "7000-0000-1243",
						PersonName: "Mr Luke Crete",
					},
				},
				ParentEventType: "MultiPersonEvent",
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
				PersonName: "Mrs Betty Banks",
				Changes: []Changeset{
					{
						FieldName:         "firstname",
						OldValue:          "",
						NewValue:          "Betty",
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
			Id:              20,
			Timestamp:       "27/09/2026 17:13:12",
			TimestampAsDate: "27/09/2026 17:13:12",
			EventType:       "ClientFeePayerRemoved",
			User: User{
				Id:          61,
				PhoneNumber: "12345678",
				DisplayName: "case manager",
				Email:       "case.manager@opgtest.com",
			},
			Event: Event{
				PersonId:   "49",
				PersonUid:  "7000-0000-1367",
				PersonType: "Deputy",
				PersonName: "Mrs Betty Banks",
				MultiPersons: []AdditionalPerson{
					{
						PersonType: "Client",
						PersonId:   "47",
						PersonUid:  "7000-0000-1243",
						PersonName: "Mr Luke Crete",
					},
					{
						PersonType: "Deputy",
						PersonId:   "49",
						PersonUid:  "7000-0000-1367",
						PersonName: "Mrs Betty Banks",
					},
				},
				ParentEventType: "MultiPersonEvent",
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
				PersonName:     "Mrs Betty Banks",
				AdditionalPersons: []AdditionalPerson{
					{
						PersonType: "Client",
						PersonId:   "47",
						PersonUid:  "7000-0000-1243",
						PersonName: "Mr Luke Crete",
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
				PersonName: "Mrs Betty Banks",
				Changes: []Changeset{
					{
						FieldName: "firstname",
						OldValue:  "",
						NewValue:  "Betty",
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
			Id:        20,
			Timestamp: "27/09/2026 17:13:12",
			EventType: "Opg\\Core\\Model\\Event\\Client\\ClientFeePayerRemoved",
			User: User{
				Id:          61,
				PhoneNumber: "12345678",
				DisplayName: "case manager",
				Email:       "case.manager@opgtest.com",
			},
			Event: Event{
				PersonId:   "49",
				PersonUid:  "7000-0000-1367",
				PersonType: "Deputy",
				PersonName: "Mrs Betty Banks",
				MultiPersons: []AdditionalPerson{
					{
						PersonType: "Client",
						PersonId:   "47",
						PersonUid:  "7000-0000-1243",
						PersonName: "Mr Luke Crete",
					},
					{
						PersonType: "Deputy",
						PersonId:   "49",
						PersonUid:  "7000-0000-1367",
						PersonName: "Mrs Betty Banks",
					},
				},
				ParentEventType: "MultiPersonEvent",
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
				PersonName:     "Mrs Betty Banks",
				AdditionalPersons: []AdditionalPerson{
					{
						PersonType: "Client",
						PersonId:   "47",
						PersonUid:  "7000-0000-1243",
						PersonName: "Mr Luke Crete",
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
				PersonId:                     "49",
				PersonName:                   "Mrs Betty Banks",
				PersonType:                   "Deputy",
				PersonUid:                    "7000-0000-1367",
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
				PersonId:                     "49",
				PersonName:                   "Mrs Betty Banks",
				PersonType:                   "Deputy",
				PersonUid:                    "7000-0000-1367",
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
				PersonId:                     "49",
				PersonName:                   "Mrs Betty Banks",
				PersonType:                   "Deputy",
				PersonUid:                    "7000-0000-1367",
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
				PersonId:                     "49",
				PersonName:                   "Mrs Betty Banks",
				PersonType:                   "Deputy",
				PersonUid:                    "7000-0000-1367",
				AdditionalPersons:            []AdditionalPerson{},
			},
		},
		{
			"Address changed event",
			Event{
				PersonId:   "49",
				PersonName: "Mrs Betty Banks",
				PersonType: "Deputy",
				PersonUid:  "7000-0000-1367",
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
				PersonId:   "49",
				PersonName: "Mrs Betty Banks",
				PersonType: "Deputy",
				PersonUid:  "7000-0000-1367",
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
				PersonId:       "49",
				PersonName:     "Mrs Betty Banks",
				PersonType:     "Deputy",
				PersonUid:      "7000-0000-1367",
				AdditionalPersons: []AdditionalPerson{
					{
						PersonId:   "47",
						PersonName: "Mr Luke Crete",
						PersonType: "Client",
						PersonUid:  "7000-0000-1243",
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
				PersonId:       "49",
				PersonName:     "Mrs Betty Banks",
				PersonType:     "Deputy",
				PersonUid:      "7000-0000-1367",
				AdditionalPersons: []AdditionalPerson{
					{
						PersonId:   "47",
						PersonName: "Mr Luke Crete",
						PersonType: "Client",
						PersonUid:  "7000-0000-1243",
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
				PersonId:       "49",
				PersonName:     "Mrs Betty Banks",
				PersonType:     "Deputy",
				PersonUid:      "7000-0000-1367",
				AdditionalPersons: []AdditionalPerson{
					{
						PersonId:   "47",
						PersonName: "Mr Luke Crete",
						PersonType: "Client",
						PersonUid:  "7000-0000-1243",
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
				PersonId:       "49",
				PersonName:     "Mrs Betty Banks",
				PersonType:     "Deputy",
				PersonUid:      "7000-0000-1367",
				AdditionalPersons: []AdditionalPerson{
					{
						PersonId:   "47",
						PersonName: "Mr Luke Crete",
						PersonType: "Client",
						PersonUid:  "7000-0000-1243",
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

func TestUpdateChangeset(t *testing.T) {
	tests := []struct {
		name              string
		myChangeset       Changeset
		expectedChangeset Changeset
	}{
		{
			name: "Address changed event",
			myChangeset: Changeset{
				FieldName: "addressLine1",
				OldValue:  "",
				NewValue:  "Flat 7 Bollerworth",
				Type:      "string",
			},
			expectedChangeset: Changeset{
				FieldName:         "addressLine1",
				OldValue:          "",
				NewValue:          "Flat 7 Bollerworth",
				Type:              "string",
				TimelineValueName: "Address line 1",
			},
		},
		{
			name: "Status changed date",
			myChangeset: Changeset{
				FieldName: "statusChangeDate",
				OldValue:  "",
				NewValue:  "2026-10-12T00:00:00+00:00",
				Type:      "string",
			},
			expectedChangeset: Changeset{
				FieldName:         "statusChangeDate",
				OldValue:          "",
				NewValue:          "12/10/2026",
				Type:              "string",
				TimelineValueName: "Status change date",
			},
		},
		{
			name: "Status on case",
			myChangeset: Changeset{
				FieldName: "statusOnCase",
				OldValue:  "DECEASED",
				NewValue:  "ERROR",
				Type:      "string",
			},
			expectedChangeset: Changeset{
				FieldName:         "statusOnCase",
				OldValue:          "Deceased",
				NewValue:          "Error",
				Type:              "string",
				TimelineValueName: "Status on case",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := updateChangeset(tt.myChangeset)
			assert.Equal(t, tt.expectedChangeset, got)
		})
	}
}

func TestGetTimelineValueForField(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"unknownEntity", "UnknownEntity"},
		{"addressLine1", "Address line 1"},
		{"addressLine2", "Address line 2"},
		{"addressLine3", "Address line 3"},
		{"audioTape", "Audio"},
		{"correspondenceByEmail", "Email"},
		{"correspondenceByPhone", "Phone"},
		{"correspondenceByPost", "Post"},
		{"correspondenceByWelsh", "Welsh"},
		{"county", "County"},
		{"country", "Country"},
		{"deputyType", "Type of deputy"},
		{"deputySubType", "Deputy Subtype"},
		{"dob", "Date of birth"},
		{"email", "Email address"},
		{"firstname", "First name"},
		{"hearingImpaired", "Hearing impaired"},
		{"homePhoneNumber", "Evening phone number"},
		{"interpreterRequired", "Interpreter required"},
		{"isAirmailRequired", "Is airmail required"},
		{"largePrint", "Large print"},
		{"mainCorrespondent", "Correspondent"},
		{"middleNames", "Middle names"},
		{"mobileNumber", "Mobile telephone number"},
		{"occupation", "Occupation"},
		{"otherNames", "Other names"},
		{"postcode", "Postcode"},
		{"relationshipToClient", "Relationship to client"},
		{"spellingOfNameRequiresCare", "Spelling of name requires care"},
		{"salutation", "Title"},
		{"statusChangeDate", "Status change date"},
		{"statusOnCase", "Status on case"},
		{"surname", "Surname"},
		{"town", "Town"},
		{"workPhoneNumber", "Daytime telephone number"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetTimelineValueForField(tt.name)
			if got != tt.want {
				t.Errorf("getTimelineValueForField() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReformatCapitalvalues(t *testing.T) {
	tests := []struct {
		fieldName string
		want      bool
	}{
		{"statusChangeDate", true},
		{"statusOnCase", true},
		{"deputyType", true},
		{"deputySubType", true},
		{"relationshipToClient", true},
		{"middleNames", false},
		{"email", false},
		{"town", false},
		{"audioTape", false},
		{"dob", false},
	}
	for _, tt := range tests {
		t.Run(tt.fieldName, func(t *testing.T) {
			got := reformatCapitalValues(tt.fieldName)
			if got != tt.want {
				t.Errorf("reformatCapitalValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReformatEventType(t *testing.T) {
	tests := []struct {
		eventType string
		want      string
	}{
		{"Opg\\\\Core\\\\Model\\\\Event\\\\Person\\\\PersonDetailsChanged", "PersonDetailsChanged"},
		{"Opg\\Core\\Model\\Event\\Person\\PersonContactDetailsChanged", "PersonContactDetailsChanged"},
		{"Opg\\Core\\Model\\Event\\Person\\AddressChanged", "AddressChanged"},
		{"Opg\\Core\\Model\\Event\\Order\\DeputyLinkedToOrder", "DeputyLinkedToOrder"},
		{"Opg\\Core\\Model\\Event\\Person\\PersonStatusChanged", "PersonStatusChanged"},
	}
	for _, tt := range tests {
		t.Run(tt.eventType, func(t *testing.T) {
			got := reformatEventType(tt.eventType)
			if got != tt.want {
				t.Errorf("reformat event type() = %v, want %v", got, tt.want)
			}
		})
	}
}
