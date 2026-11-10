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

type mockManageAdditionalInformation struct {
	count     int
	lastCtx   sirius.Context
	updateErr error
}

func (m *mockManageAdditionalInformation) UpdateAdditionalInformation(ctx sirius.Context, _ int, _ sirius.AdditionalInformationForm) error {
	m.count += 1
	m.lastCtx = ctx

	return m.updateErr
}

var manageAdditionalInformationAppVars = AppVars{
	DeputyDetails: sirius.DeputyDetails{
		ID: 123,
	},
}

func TestGetManageAdditionalInformation(t *testing.T) {
	assert := assert.New(t)

	client := &mockManageAdditionalInformation{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	handler := renderTemplateForManageAdditionalInformationForm(client, template)
	err := handler(manageAdditionalInformationAppVars, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
}

func TestPostManageAdditionalInformation(t *testing.T) {
	assert := assert.New(t)

	client := &mockManageAdditionalInformation{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))

	returnedTemplate := renderTemplateForManageAdditionalInformationForm(client, template)(manageAdditionalInformationAppVars, w, r)

	assert.Equal(Redirect("/123?success=additionalInformation"), returnedTemplate)
}

func TestPostManageAdditionalInformationValidationError(t *testing.T) {
	assert := assert.New(t)

	client := &mockManageAdditionalInformation{}

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

	returnedTemplate := renderTemplateForManageAdditionalInformationForm(client, template)(manageAdditionalInformationAppVars, w, r)

	assert.Nil(returnedTemplate)
	assert.Equal(1, client.count)
	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(ManageAdditionalInformationVars{
		AppVars: AppVars{
			DeputyDetails: manageAdditionalInformationAppVars.DeputyDetails,
			Errors:        util.RenameErrors(validationErrors),
			PageName:      "Manage additional information",
		},
	}, template.lastVars)
}

func TestManageAdditionalInformationHandlesErrors(t *testing.T) {
	assert := assert.New(t)

	client := &mockManageAdditionalInformation{
		updateErr: sirius.StatusError{Code: 500},
	}

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))

	returnedTemplate := renderTemplateForManageAdditionalInformationForm(client, template)(AppVars{}, w, r)

	assert.Equal(client.updateErr, returnedTemplate)
}
