package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type DeputyDetailsForm struct {
	DeputyFirstName        string        `json:"firstname"`
	DeputyMiddleName       string        `json:"middleName"`
	DeputySurname          string        `json:"surname"`
	DeputyOtherNames       string        `json:"otherNames"`
	DeputySubType          DeputySubType `json:"deputySubType"`
	DeputyDateOfBirth      string        `json:"dob"`
	Email                  string        `json:"email"`
	MobileTelephoneNumber  string        `json:"mobileNumber"`
	DaytimeTelephoneNumber string        `json:"workPhoneNumber"`
	EveningTelephoneNumber string        `json:"homePhoneNumber"`
	AddressLine1           string        `json:"addressLine1"`
	AddressLine2           string        `json:"addressLine2"`
	AddressLine3           string        `json:"addressLine3"`
	Town                   string        `json:"town"`
	County                 string        `json:"county"`
	Country                string        `json:"country"`
	Postcode               string        `json:"postcode"`
	IsAirmailRequired      bool          `json:"isAirmailRequired"`
}

func (c *Client) UpdateDeputyDetails(ctx Context, deputyID int, deputyDetails DeputyDetailsForm) error {
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(deputyDetails)
	if err != nil {
		return err
	}

	requestURL := fmt.Sprintf(SupervisionAPIPath+"/v1/deputies/%d/contact-details", deputyID)

	req, err := c.newRequest(ctx, http.MethodPut, requestURL, &body)

	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}

	defer unchecked(resp.Body.Close)

	if resp.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}
	if resp.StatusCode == http.StatusBadRequest {
		var v struct {
			ValidationErrors ValidationErrors `json:"validation_errors"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil && len(v.ValidationErrors) > 0 {
			return ValidationError{
				Errors: v.ValidationErrors,
			}
		}
	}

	if resp.StatusCode != http.StatusOK {
		return newStatusError(resp)
	}

	return nil
}
