package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/sirius"
	"github.com/ministryofjustice/opg-go-common/paginate"
	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/urlbuilder"

	"github.com/ministryofjustice/opg-sirius-supervision-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mocktimelineVars struct {
	count                int
	lastCtx              sirius.Context
	GetDeputyTimelineErr error
}

func (m *mocktimelineVars) GetDeputyTimeline() ([]sirius.TimelineEvent, error) {
	m.count += 1
	m.lastCtx = ctx

	var myTimelineEvents = []sirius.TimelineEvent{
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
	}

	return myTimelineEvents, m.GetDeputyTimelineErr
}

func TestNavigateToTimeline(t *testing.T) {
	assert := assert.New(t)

	client := &mockDeputyHubTimelineInformation{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForDeputyHubEvents(client, template, EnvironmentVars{})
	err := handler(AppVars{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(deputyHubEventVars{
		ListPage: ListPage{
			AppVars: AppVars{
				PageName: "Timeline",
			},
			Pagination: paginate.Pagination{
				CurrentPage:     1,
				TotalPages:      1,
				TotalElements:   2,
				ElementsPerPage: 25,
				ElementName:     "timeline event(s)",
				PerPageOptions:  []int{25, 50, 100},
				UrlBuilder: urlbuilder.UrlBuilder{
					OriginalPath:    "/path",
					SelectedPerPage: 25,
				},
			},
			PerPage: 25,
		},
		DeputyEvents: sirius.DeputyEvents{
			{
				ID: 1,
			},
			{
				ID: 2,
			},
		},
	}, template.lastVars)
}

//func TestDeputyEventsReturnsErrors(t *testing.T) {
//	assert := assert.New(t)
//	client := &mockDeputyHubTimelineInformation{
//		GetDeputyEventsErr: sirius.StatusError{Code: 500},
//	}
//
//	template := &mockTemplates{}
//
//	w := httptest.NewRecorder()
//	r, _ := http.NewRequest("GET", "/path", strings.NewReader(""))
//
//	returnedError := renderTemplateForDeputyHubEvents(client, template, EnvironmentVars{})(AppVars{}, w, r)
//
//	assert.Equal(client.GetDeputyEventsErr, returnedError)
//}
//
//func TestCreateUrlBuilder(t *testing.T) {
//	myUrlBuilder := CreateUrlBuilder("67/timeline", 25, "/supervision/deputies/")
//
//	assert.Equal(t, urlbuilder.UrlBuilder{
//		OriginalPath:    "/supervision/deputies/67/timeline",
//		SelectedPerPage: 25,
//		SelectedFilters: nil,
//		SelectedSort:    urlbuilder.Sort{},
//	}, myUrlBuilder)
//
//	myUrlBuilder = CreateUrlBuilder("99/url", 100, "/test/path/")
//
//	assert.Equal(t, urlbuilder.UrlBuilder{
//		OriginalPath:    "/test/path/99/url",
//		SelectedPerPage: 100,
//		SelectedFilters: nil,
//		SelectedSort:    urlbuilder.Sort{},
//	}, myUrlBuilder)
//}
