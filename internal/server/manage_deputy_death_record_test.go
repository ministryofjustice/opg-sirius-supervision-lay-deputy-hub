package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/shared"
	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestGetDeputyDeathRecord(t *testing.T) {
	assert := assert.New(t)
	client := &mockApiClient{
		DeathRecord: sirius.DeputyDeathRecord{
			DeputyID:                      333,
			DateOfDeath:                   shared.NewDate("2020-01-01"),
			ProofOfDeathReceived:          "true",
			DateDeathCertificateReceieved: shared.NewDate("2020-02-02"),
			DateNotified:                  shared.NewDate("2020-02-01"),
			NotificationMethod:            "LETTER",
			NotifiedBy: sirius.RefData{
				Label:  "EXECUTOR",
				Handle: "Executor",
			},
			Notes: "Some notes",
		},
	}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/death-record", nil)

	handler := renderTemplateForDeputyDeathGet(client, template)
	err := handler(AppVars{
		Path: "/deathRecord",
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
	assert.Equal(
		deputyDeathRecordVars{
			Path:     "/deathRecord",
			PageName: "Record death of Test Deputy",
			DeputyDetails: sirius.DeputyDetails{
				ID:              333,
				DeputyFirstName: "Test",
				DeputySurname:   "Deputy",
				DeputyType: sirius.DeputyType{
					Handle: "LAY",
					Label:  "Lay Deputy",
				},
			},
			DeputyDeathRecord: sirius.DeputyDeathRecord{
				DeputyID:                      333,
				DateOfDeath:                   shared.NewDate("2020-01-01"),
				ProofOfDeathReceived:          "true",
				DateDeathCertificateReceieved: shared.NewDate("2020-02-02"),
				DateNotified:                  shared.NewDate("2020-02-01"),
				NotificationMethod:            "LETTER",
				NotifiedBy: sirius.RefData{
					Label:  "EXECUTOR",
					Handle: "Executor",
				},
				Notes: "Some notes",
			},
		}, template.lastVars)
}

func TestUpdateDeputyDeathRecord(t *testing.T) {
	assert := assert.New(t)
	client := &mockApiClient{
		DeathRecord: sirius.DeputyDeathRecord{},
	}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(
		"POST",
		"/death-record",
		strings.NewReader("{date-of-death:'2020-01-01', proof-of-death:'true', date-proof-of-death-received:'2020-02-02', date-notified:'2020-02-01', notified-by:'DEPUTY', notification-method:'PHONE', notes:'some thoughts here'}"))

	handler := renderTemplateForDeputyDeathPost(client, template)
	err := handler(AppVars{
		Path: "/deathRecord",
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

	assert.Equal(Redirect("/333?success=recordDeath"), err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
}

func TestUpdateDeputyDeathRecordValidationError(t *testing.T) {
	assert := assert.New(t)
	client := &mockApiClient{
		DeathRecord: sirius.DeputyDeathRecord{},
	}
	template := &mockTemplates{}
	validationErrors := sirius.ValidationErrors{
		"notes": {
			"stringLengthTooLong": "What sirius gives us",
		},
	}

	client.error = sirius.ValidationError{
		Errors: validationErrors,
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(
		"POST",
		"/death-record",
		strings.NewReader("{date-of-death:'2020-01-01', proof-of-death:'true', date-proof-of-death-received:'2020-02-02', date-notified:'2020-02-01', notified-by:'DEPUTY', notification-method:'PHONE', notes:'some thoughts here'}"))

	handler := renderTemplateForDeputyDeathPost(client, template)
	err := handler(AppVars{
		Path: "/deathRecord",
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

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(deputyDeathRecordVars{
		Path: "/deathRecord",
		DeputyDetails: sirius.DeputyDetails{
			ID:              333,
			DeputyFirstName: "Test",
			DeputySurname:   "Deputy",
			DeputyType: sirius.DeputyType{
				Handle: "LAY",
				Label:  "Lay Deputy",
			},
		},
		Errors:   util.RenameErrors(validationErrors),
		PageName: "Record death of Test Deputy",
		DeputyDeathRecord: sirius.DeputyDeathRecord{
			DeputyID: 333,
		},
	}, template.lastVars)

	assert.Nil(err)
}

func TestUpdateDeputyDeathRecordHandlesError(t *testing.T) {
	assert := assert.New(t)
	client := &mockApiClient{
		DeathRecord: sirius.DeputyDeathRecord{},
	}
	template := &mockTemplates{}
	client.error = sirius.StatusError{Code: 500}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(
		"POST",
		"/death-record",
		strings.NewReader("{date-of-death:'2020-01-01', proof-of-death:'true', date-proof-of-death-received:'2020-02-02', date-notified:'2020-02-01', notified-by:'DEPUTY', notification-method:'PHONE', notes:'some thoughts here'}"))

	handler := renderTemplateForDeputyDeathPost(client, template)
	err := handler(AppVars{
		Path: "/deathRecord",
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

	assert.Equal(client.error, err)
}
