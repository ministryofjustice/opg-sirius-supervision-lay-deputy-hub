package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/shared"
	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

func TestNavigateToTimeline(t *testing.T) {
	assert := assert.New(t)

	client := &mockApiClient{
		TimelineEvents: []sirius.TimelineEvent{
			{
				Id:              20,
				Timestamp:       "28/09/2026 17:13:12",
				TimestampAsDate: shared.NewDate("28/09/2026 17:13:12").String(),
				EventType:       "PersonDetailsChanged",
				User: sirius.User{
					Id:          61,
					PhoneNumber: "12345678",
					DisplayName: "case manager",
					Email:       "case.manager@opgtest.com",
				},
				Event: sirius.Event{
					PersonType: "Deputy",
					PersonId:   "49",
					PersonUid:  "7000-0000-1367",
					PersonName: "Mrs Betty Banks",
					Changes: []sirius.Changeset{
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
				TimestampAsDate: shared.NewDate("28/09/2026 17:13:12").String(),
				EventType:       "DeputyLinkedToOrder",
				User: sirius.User{
					Id:          61,
					PhoneNumber: "12345678",
					DisplayName: "case manager",
					Email:       "case.manager@opgtest.com",
				},
				Event: sirius.Event{
					OrderType:      "PFA",
					OrderUid:       "7000-0000-1334",
					OrderId:        "86",
					OrderCourtRef:  "31924130",
					CourtReference: "31924130",
					PersonType:     "Deputy",
					PersonId:       "49",
					PersonUid:      "7000-0000-1367",
					PersonName:     "Mrs Betty Banks",
					AdditionalPersons: []sirius.AdditionalPerson{
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
				TimestampAsDate: shared.NewDate("12/10/2026 13:35:35").String(),
				EventType:       "ClientFeePayerRemoved",
				User: sirius.User{
					Id:          101,
					PhoneNumber: "12345678",
					DisplayName: "PROTeam1 User1",
					Email:       "pro1@opgtest.com",
				},
				Event: sirius.Event{
					MultiPersons: []sirius.AdditionalPerson{
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
		},
	}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/timeline", nil)

	handler := renderTemplateForTimeline(client, template)
	err := handler(AppVars{
		Path: "/timeline",
		DeputyDetails: sirius.DeputyDetails{
			ID:              333,
			DeputyFirstName: "Test",
			DeputySurname:   "Deputy",
			DeputyType: sirius.DeputyType{
				Handle: "LAY",
				Label:  "Lay Deputy",
			},
		},
	}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(timelineVars{
		AppVars: AppVars{
			Path:      "/timeline",
			XSRFToken: "",
			DeputyDetails: sirius.DeputyDetails{
				ID:              333,
				DeputyFirstName: "Test",
				DeputySurname:   "Deputy",
				DeputyType: sirius.DeputyType{
					Handle: "LAY",
					Label:  "Lay Deputy",
				},
			},
			SuccessMessage: "",
			PageName:       "Timeline",
		},
		Timeline: []sirius.TimelineEvent{
			{
				Id:              20,
				Timestamp:       "28/09/2026 17:13:12",
				TimestampAsDate: "28/09/2026",
				EventType:       "PersonDetailsChanged",
				User: sirius.User{
					Id:          61,
					PhoneNumber: "12345678",
					DisplayName: "case manager",
					Email:       "case.manager@opgtest.com",
				},
				Event: sirius.Event{
					PersonType: "Deputy",
					PersonId:   "49",
					PersonUid:  "7000-0000-1367",
					PersonName: "Mrs Betty Banks",
					Changes: []sirius.Changeset{
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
				User: sirius.User{
					Id:          61,
					PhoneNumber: "12345678",
					DisplayName: "case manager",
					Email:       "case.manager@opgtest.com",
				},
				Event: sirius.Event{
					OrderType:      "PFA",
					OrderUid:       "7000-0000-1334",
					OrderId:        "86",
					OrderCourtRef:  "31924130",
					CourtReference: "31924130",
					PersonType:     "Deputy",
					PersonId:       "49",
					PersonUid:      "7000-0000-1367",
					PersonName:     "Mrs Betty Banks",
					AdditionalPersons: []sirius.AdditionalPerson{
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
				User: sirius.User{
					Id:          101,
					PhoneNumber: "12345678",
					DisplayName: "PROTeam1 User1",
					Email:       "pro1@opgtest.com",
				},
				Event: sirius.Event{
					MultiPersons: []sirius.AdditionalPerson{
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
		},
	}, template.lastVars)
}

func TestHandlesErrorReturnedByGetDeputyTimeline(t *testing.T) {
	assert := assert.New(t)

	client := &mockApiClient{
		error: sirius.StatusError{Code: 500},
		TimelineEvents: []sirius.TimelineEvent{
			{
				Id:              20,
				Timestamp:       "28/09/2026 17:13:12",
				TimestampAsDate: "28/09/2026",
				EventType:       "PersonDetailsChanged",
				User: sirius.User{
					Id:          61,
					PhoneNumber: "12345678",
					DisplayName: "case manager",
					Email:       "case.manager@opgtest.com",
				},
				Event: sirius.Event{
					PersonType: "Deputy",
					PersonId:   "49",
					PersonUid:  "7000-0000-1367",
					PersonName: "Mrs Betty Banks",
					Changes: []sirius.Changeset{
						{
							FieldName:         "firstname",
							OldValue:          "",
							NewValue:          "Betty",
							Type:              "string",
							TimelineValueName: "First name",
						},
					},
				},
			},
		},
	}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", strings.NewReader(""))

	returnedError := renderTemplateForTimeline(client, template)(AppVars{}, w, r)

	assert.Equal(sirius.StatusError{Code: 500}, returnedError)
}
