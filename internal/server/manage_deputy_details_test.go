package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	handler := renderTemplateForManageDeputyDetailsGet(template)
	err := handler(app, w, r)

	assert.Nil(err)
	assert.Equal("Dr", template.lastVars.(ManageDeputyDetailsVars).Form.Title)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
}

func TestPostManageDeputyDetails(t *testing.T) {
	assert := assert.New(t)

	client := &mockManageDeputyDetails{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", strings.NewReader("title=Dr"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	returnedTemplate := renderTemplateForManageDeputyDetailsPost(client, template)(manageDeputyAppVars, w, r)

	assert.Equal(Redirect("/123?success=deputyDetails"), returnedTemplate)
	assert.Equal("Dr", client.lastForm.Title)
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
	r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))

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
			DeputySubType: sirius.DeputySubType{SubType: "PERSON"},
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
