package server

import (
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockClient struct{}

type mockTemplates struct {
	mock.Mock
	count    int
	lastName string
	lastVars interface{}
	error    error
}

func (m *mockTemplates) ExecuteTemplate(i io.Writer, name string, vars interface{}) error {
	m.count += 1
	m.lastName = name
	m.lastVars = vars

	return nil
}

func (m mockClient) GetDeputyClients(_ sirius.Context, _ sirius.ClientListParams) (sirius.ClientList, error) {
	return sirius.ClientList{
		TotalClients: 3,
		Metadata: sirius.Metadata{
			TotalActiveClients: 3,
		},
	}, nil
}

func (m mockClient) GetDeputyDetails(_ sirius.Context, deputyID int) (sirius.DeputyDetails, error) {
	return sirius.DeputyDetails{
		ID:              deputyID,
		DeputyFirstName: "Test",
		DeputySurname:   "Deputy",
		DeputyType: sirius.DeputyType{
			Handle: "LAY",
			Label:  "Lay Deputy",
		},
	}, nil
}

func (m mockClient) GetDeputyTimeline(_ sirius.Context, deputyID int) ([]sirius.TimelineEvent, error) {
	return []sirius.TimelineEvent{
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
	}, nil
}

func TestPagesLoadSuccessfully(t *testing.T) {
	templates := map[string]*template.Template{
		"error.gotmpl":          testTemplate(t, "error.gotmpl"),
		"deputy-details.gotmpl": testTemplate(t, "deputy-details.gotmpl"),
		"clients.gotmpl":        testTemplate(t, "clients.gotmpl"),
		"timeline.gotmpl":       testTemplate(t, "timeline.gotmpl"),
	}

	handler := New(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		mockClient{},
		templates,
		EnvironmentVars{},
	)

	for _, test := range []struct {
		name string
		path string
	}{
		{name: "deputy details", path: "/123"},
		{name: "clients", path: "/123/clients"},
		{name: "timeline", path: "/123/timeline"},
	} {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, test.path, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}

func testTemplate(t *testing.T, name string) *template.Template {
	t.Helper()

	tmpl, err := template.New(name).Parse(`{{ define "page" }}OK{{ end }}`)
	require.NoError(t, err)

	return tmpl
}
