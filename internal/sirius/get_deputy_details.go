package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func (s SpecialCorrespondenceRequirementsType) ActiveSpecialCorrespondenceRequirements() []string {
	requirementsToList := []string{}

	if s.AudioTape {
		requirementsToList = append(requirementsToList, "Audio")
	}
	if s.LargePrint {
		requirementsToList = append(requirementsToList, "Large Print")
	}
	if s.HearingImpaired {
		requirementsToList = append(requirementsToList, "Hearing impaired")
	}
	if s.SpellingOfNameRequiresCare {
		requirementsToList = append(requirementsToList, "Spelling of name requires care")
	}

	return requirementsToList
}

type DeputyDetails struct {
	ID                                int                                   `json:"id"`
	Salutation                        string                                `json:"salutation"`
	DeputyFirstName                   string                                `json:"firstName"`
	DeputySurname                     string                                `json:"surname"`
	DeputyOtherNames                  string                                `json:"otherNames"`
	DeputyNumber                      int                                   `json:"deputyNumber"`
	DeputyStatus                      string                                `json:"deputyStatus"`
	DeputyType                        DeputyType                            `json:"deputyType"`
	DeputyDateOfBirth                 string                                `json:"dob"`
	Email                             string                                `json:"email"`
	MobileTelephoneNumber             string                                `json:"mobileNumber"`
	DaytimeTelephoneNumber            string                                `json:"daytimeNumber"`
	EveningTelephoneNumber            string                                `json:"eveningNumber"`
	AddressLine1                      string                                `json:"addressLine1"`
	AddressLine2                      string                                `json:"addressLine2"`
	AddressLine3                      string                                `json:"addressLine3"`
	Town                              string                                `json:"town"`
	County                            string                                `json:"county"`
	Postcode                          string                                `json:"postcode"`
	IsAirmailRequired                 bool                                  `json:"isAirmailRequired"`
	SpecialCorrespondenceRequirements SpecialCorrespondenceRequirementsType `json:"specialCorrespondenceRequirements"`
	InterpreterRequired               string                                `json:"interpreterRequired"`
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

	return v, err
}
