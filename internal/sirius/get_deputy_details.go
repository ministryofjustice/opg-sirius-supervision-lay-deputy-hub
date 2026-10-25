package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/utility"
)

type DeputyType struct {
	Handle string `json:"handle"`
	Label  string `json:"label"`
}

type DeputySubType struct {
	SubType string `json:"handle"`
}

type SpecialCorrespondenceRequirementsType struct {
	AudioTape                  bool `json:"audioTape"`
	LargePrint                 bool `json:"largePrint"`
	HearingImpaired            bool `json:"hearingImpaired"`
	SpellingOfNameRequiresCare bool `json:"spellingOfNameRequiresCare"`
}

type DeputyDetails struct {
	ID                                int                                   `json:"id"`
	Salutation                        string                                `json:"salutation"`
	DeputyFirstName                   string                                `json:"firstName"`
	DeputyMiddleName                  string                                `json:"middleName"`
	DeputySurname                     string                                `json:"surname"`
	DeputyOtherNames                  string                                `json:"otherNames"`
	DeputyNumber                      int                                   `json:"deputyNumber"`
	DeputyStatus                      string                                `json:"deputyStatus"`
	DeputyType                        DeputyType                            `json:"deputyType"`
	DeputyDateOfBirth                 string                                `json:"dob"`
	Email                             string                                `json:"email"`
	MobileTelephoneNumber             string                                `json:"mobileNumber"`
	DaytimeTelephoneNumber            string                                `json:"phoneNumber"`
	EveningTelephoneNumber            string                                `json:"eveningNumber"`
	AddressLine1                      string                                `json:"addressLine1"`
	AddressLine2                      string                                `json:"addressLine2"`
	AddressLine3                      string                                `json:"addressLine3"`
	Town                              string                                `json:"town"`
	County                            string                                `json:"county"`
	Country                           string                                `json:"country"`
	Postcode                          string                                `json:"postcode"`
	IsAirmailRequired                 bool                                  `json:"isAirmailRequired"`
	SpecialCorrespondenceRequirements SpecialCorrespondenceRequirementsType `json:"specialCorrespondenceRequirements"`
	InterpreterRequired               string                                `json:"interpreterRequired"`
	CorrespondenceByWelsh             bool                                  `json:"correspondenceByWelsh"`
}

func (d DeputyDetails) ActiveSpecialCorrespondenceRequirements() []string {
	requirementsToList := []string{}

	if d.SpecialCorrespondenceRequirements.AudioTape {
		requirementsToList = append(requirementsToList, "Audio")
	}
	if d.SpecialCorrespondenceRequirements.LargePrint {
		requirementsToList = append(requirementsToList, "Large Print")
	}
	if d.SpecialCorrespondenceRequirements.HearingImpaired {
		requirementsToList = append(requirementsToList, "Hearing impaired")
	}
	if d.SpecialCorrespondenceRequirements.SpellingOfNameRequiresCare {
		requirementsToList = append(requirementsToList, "Spelling of name requires care")
	}
	if d.CorrespondenceByWelsh {
		requirementsToList = append(requirementsToList, "Welsh")
	}

	return requirementsToList
}

func (c *Client) GetDeputyDetails(ctx Context, deputyID int) (DeputyDetails, error) {
	var v DeputyDetails

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf(SupervisionAPIPath+"/v1/deputies/%d", deputyID), nil)
	if err != nil {
		return v, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return v, err
	}

	defer unchecked(resp.Body.Close)

	if resp.StatusCode == http.StatusUnauthorized {
		return v, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return v, newStatusError(resp)
	}

	err = json.NewDecoder(resp.Body).Decode(&v)

	return formatDeputyDetails(v), err
}

func formatDeputyDetails(v DeputyDetails) DeputyDetails {
	v.DeputyDateOfBirth = utility.FormatDateTime(utility.SiriusDate, v.DeputyDateOfBirth, utility.IsoDate)

	return v
}
