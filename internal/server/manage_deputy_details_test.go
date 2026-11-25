package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/shared"
	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/util"
	"github.com/stretchr/testify/assert"
)

type mockManageDeputyDetails struct {
	count     int
	lastCtx   sirius.Context
	lastForm  sirius.DeputyDetailsForm
	updateErr error
}

func (m *mockManageDeputyDetails) UpdateDeputyDetails(ctx sirius.Context, _ int, form sirius.DeputyDetailsForm) error {
	m.count += 1
	m.lastCtx = ctx
	m.lastForm = form

	return m.updateErr
}

var manageDeputyAppVars = AppVars{
	DeputyDetails: sirius.DeputyDetails{
		ID: 123,
	},
}

func TestGetManageDeputyDetails(t *testing.T) {
	assert := assert.New(t)
	template := &mockTemplates{}
	app := manageDeputyAppVars
	app.DeputyDetails.Title = "Dr"
	app.DeputyDetails.DeputyFirstName = "Alex"
	app.DeputyDetails.DeputyMiddleName = "Jordan"
	app.DeputyDetails.DeputySurname = "Smith"
	app.DeputyDetails.DeputyOtherNames = "Alexander"
	app.DeputyDetails.DeputyDateOfBirth = shared.NewDate("03/02/1985")

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	handler := renderTemplateForManageDeputyDetailsGet(template)
	err := handler(app, w, r)

	assert.Nil(err)
	assert.Equal(sirius.DeputyDetailsForm{
		Title: "Dr", FirstName: "Alex", MiddleName: "Jordan", Surname: "Smith",
		OtherNames: "Alexander", DateOfBirth: shared.NewDate("03/02/1985"),
	}, template.lastVars.(ManageDeputyDetailsVars).Form)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
}

func TestPostManageDeputyDetails(t *testing.T) {
	assert := assert.New(t)

	client := &mockManageDeputyDetails{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	values := url.Values{
		"title": {"Dr"}, "deputy-first-name": {"Alex"}, "deputy-middle-name": {"Jordan"},
		"deputy-last-name": {"Smith"}, "deputy-other-names": {"Alexander"},
		"deputy-date-of-birth": {"1985-02-03"}, "deputy-phone-number": {"020 7946 0123"},
		"deputy-evening-number": {"020 7946 0456"},
	}
	r, _ := http.NewRequest("POST", "/123", strings.NewReader(values.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	returnedTemplate := renderTemplateForManageDeputyDetailsPost(client, template)(manageDeputyAppVars, w, r)

	assert.Equal(Redirect("/123?success=deputyDetails"), returnedTemplate)
	assert.Equal(sirius.DeputyDetailsForm{
		Title: "Dr", FirstName: "Alex", MiddleName: "Jordan", Surname: "Smith",
		OtherNames: "Alexander", DateOfBirth: shared.NewDate("03/02/1985"),
		DaytimeTelephoneNumber: "020 7946 0123", EveningTelephoneNumber: "020 7946 0456",
	}, client.lastForm)
}

func TestPostManageDeputyDetailsValidationError(t *testing.T) {
	assert := assert.New(t)

	client := &mockManageDeputyDetails{}

	validationErrors := sirius.ValidationErrors{
		"firstname": {
			"stringLengthTooLong": "What sirius gives us",
		},
	}

	client.updateErr = sirius.ValidationError{
		Errors: validationErrors,
	}

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", strings.NewReader("deputy-first-name=Alex&deputy-middle-name=Jordan&deputy-date-of-birth=1985-02-03"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	returnedTemplate := renderTemplateForManageDeputyDetailsPost(client, template)(manageDeputyAppVars, w, r)
	assert.Nil(returnedTemplate)
	assert.Equal(1, client.count)
	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(ManageDeputyDetailsVars{
		AppVars: AppVars{
			DeputyDetails: manageDeputyAppVars.DeputyDetails,
			Errors:        util.RenameErrors(validationErrors),
			PageName:      "Manage deputy contact details",
		},
		Form: sirius.DeputyDetailsForm{
			FirstName: "Alex", MiddleName: "Jordan", DateOfBirth: shared.NewDate("03/02/1985"),
		},
	}, template.lastVars)
}

func TestManageDeputyDetailsHandlesErrors(t *testing.T) {
	assert := assert.New(t)

	client := &mockManageDeputyDetails{
		updateErr: sirius.StatusError{Code: 500},
	}

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))

	returnedTemplate := renderTemplateForManageDeputyDetailsPost(client, template)(AppVars{}, w, r)

	assert.Equal(client.updateErr, returnedTemplate)
}
