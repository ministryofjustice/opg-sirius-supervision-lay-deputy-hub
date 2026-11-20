package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

func TestGetDeputyDeathRecord(t *testing.T) {
	assert := assert.New(t)
	client := &mockApiClient{
		DeathRecord: sirius.DeputyDeathRecord{
			DeputyID:                      333,
			DateOfDeath:                   "2020-01-01",
			ProofOfDeathReceived:          "true",
			DateDeathCertificateReceieved: "2020-02-02",
			DateNotified:                  "2020-02-01",
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

	handler := renderTemplateForDeputyDeath(client, template)
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
			AppVars: AppVars{
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
			},
			DeputyDeathRecord: sirius.DeputyDeathRecord{
				DeputyID:                      333,
				DateOfDeath:                   "2020-01-01",
				ProofOfDeathReceived:          "true",
				DateDeathCertificateReceieved: "2020-02-02",
				DateNotified:                  "2020-02-01",
				NotificationMethod:            "LETTER",
				NotifiedBy: sirius.RefData{
					Label:  "EXECUTOR",
					Handle: "Executor",
				},
				Notes: "Some notes",
			},
		}, template.lastVars)
}

//func TestUpdateDeputyDeathRecord(t *testing.T) {
//	assert := assert.New(t)
//	client := &mockApiClient{
//		DeathRecord: sirius.DeputyDeathRecord{
//			DeputyID:                      333,
//			DateOfDeath:                   "2020-01-01",
//			ProofOfDeathReceived:          "true",
//			DateDeathCertificateReceieved: "2020-02-02",
//			DateNotified:                  "2020-02-01",
//			NotificationMethod:            "LETTER",
//			NotifiedBy: sirius.RefData{
//				Label:  "EXECUTOR",
//				Handle: "Executor",
//			},
//			Notes: "Some notes",
//		},
//	}
//	template := &mockTemplates{}
//
//	w := httptest.NewRecorder()
//	r, _ := http.NewRequest("POST", "/death-record", nil)
//
//	handler := renderTemplateForDeputyDeath(client, template)
//	err := handler(AppVars{
//		Path: "/deathRecord",
//		DeputyDetails: sirius.DeputyDetails{
//			ID:              333,
//			DeputyFirstName: "Test",
//			DeputySurname:   "Deputy",
//			DeputyType: sirius.DeputyType{
//				Handle: "LAY",
//				Label:  "Lay Deputy",
//			},
//		},
//	}, w, r)
//
//	assert.Nil(err)
//
//	resp := w.Result()
//	assert.Equal(http.StatusOK, resp.StatusCode)
//	assert.Equal(
//		deputyDeathRecordVars{
//			AppVars: AppVars{
//				Path:     "/deathRecord",
//				PageName: "Record death of Test Deputy",
//				DeputyDetails: sirius.DeputyDetails{
//					ID:              333,
//					DeputyFirstName: "Test",
//					DeputySurname:   "Deputy",
//					DeputyType: sirius.DeputyType{
//						Handle: "LAY",
//						Label:  "Lay Deputy",
//					},
//				},
//			},
//			DeputyDeathRecord: sirius.DeputyDeathRecord{
//				DeputyID:                      333,
//				DateOfDeath:                   "2020-01-01",
//				ProofOfDeathReceived:          "true",
//				DateDeathCertificateReceieved: "2020-02-02",
//				DateNotified:                  "2020-02-01",
//				NotificationMethod:            "LETTER",
//				NotifiedBy: sirius.RefData{
//					Label:  "EXECUTOR",
//					Handle: "Executor",
//				},
//				Notes: "Some notes",
//			},
//		}, template.lastVars)
//}

//
//func TestPostManageDeputyDetails(t *testing.T) {
//	assert := assert.New(t)
//
//	client := &mockManageDeputyDetails{}
//	template := &mockTemplates{}
//
//	w := httptest.NewRecorder()
//	r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))
//
//	returnedTemplate := renderTemplateForManageDeputyDetailsPost(client, template)(manageDeputyAppVars, w, r)
//
//	assert.Equal(Redirect("/123?success=deputyDetails"), returnedTemplate)
//}
//
//func TestPostManageDeputyDetailsValidationError(t *testing.T) {
//	assert := assert.New(t)
//
//	client := &mockManageDeputyDetails{}
//
//	validationErrors := sirius.ValidationErrors{
//		"firstname": {
//			"stringLengthTooLong": "What sirius gives us",
//		},
//	}
//
//	client.updateErr = sirius.ValidationError{
//		Errors: validationErrors,
//	}
//
//	template := &mockTemplates{}
//
//	w := httptest.NewRecorder()
//	r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))
//
//	returnedTemplate := renderTemplateForManageDeputyDetailsPost(client, template)(manageDeputyAppVars, w, r)
//	assert.Nil(returnedTemplate)
//	assert.Equal(1, client.count)
//	assert.Equal(1, template.count)
//	assert.Equal("page", template.lastName)
//	assert.Equal(ManageDeputyDetailsVars{
//		AppVars: AppVars{
//			DeputyDetails: manageDeputyAppVars.DeputyDetails,
//			Errors:        util.RenameErrors(validationErrors),
//			PageName:      "Manage deputy contact details",
//		},
//		Form: sirius.DeputyDetailsForm{
//			DeputySubType: sirius.DeputySubType{SubType: "PERSON"},
//		},
//	}, template.lastVars)
//}
//
//func TestManageDeputyDetailsHandlesErrors(t *testing.T) {
//	assert := assert.New(t)
//
//	client := &mockManageDeputyDetails{
//		updateErr: sirius.StatusError{Code: 500},
//	}
//
//	template := &mockTemplates{}
//
//	w := httptest.NewRecorder()
//	r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))
//
//	returnedTemplate := renderTemplateForManageDeputyDetailsPost(client, template)(AppVars{}, w, r)
//
//	assert.Equal(client.updateErr, returnedTemplate)
//}
