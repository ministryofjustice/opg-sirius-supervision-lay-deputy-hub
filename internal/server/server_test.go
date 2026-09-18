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
	"github.com/stretchr/testify/require"
)

type mockClient struct{}

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
