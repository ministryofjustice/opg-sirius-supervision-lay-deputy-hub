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

type mockApiClient struct {
	error          error
	Clients        sirius.ClientList
	DeputyDetails  sirius.DeputyDetails
	TimelineEvents []sirius.TimelineEvent
}

func (m mockApiClient) GetDeputyClients(_ sirius.Context, _ sirius.ClientListParams) (sirius.ClientList, error) {
	return sirius.ClientList{
		TotalClients: 3,
		Metadata: sirius.Metadata{
			TotalActiveClients: 3,
		},
	}, nil
}

func (m mockApiClient) GetDeputyDetails(_ sirius.Context, deputyID int) (sirius.DeputyDetails, error) {
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

func (m mockApiClient) GetDeputyTimeline(_ sirius.Context, deputyID int) ([]sirius.TimelineEvent, error) {
	return m.TimelineEvents, m.error
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
		mockApiClient{},
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
