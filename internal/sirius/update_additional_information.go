package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type AdditionalInformationForm struct {
	InterpreterRequired               string                                `json:"interpreterRequired"`
	SpecialCorrespondenceRequirements SpecialCorrespondenceRequirementsType `json:"specialCorrespondenceRequirements"`
	Notes                             string                                `json:"notes"`
	CorrespondenceByWelsh             bool                                  `json:"correspondenceByWelsh"`
	CorrespondenceByEmail             bool                                  `json:"correspondenceByEmail"`
	CorrespondenceByPhone             bool                                  `json:"correspondenceByPhone"`
	CorrespondenceByPost              bool                                  `json:"correspondenceByPost"`
}

func (c *Client) UpdateAdditionalInformation(ctx Context, deputyID int, additionalInformation AdditionalInformationForm) error {
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(additionalInformation)
	if err != nil {
		return err
	}

	requestURL := fmt.Sprintf(SupervisionAPIPath+"/v1/deputies/%d/additional-details", deputyID)

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
