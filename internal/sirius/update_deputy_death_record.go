package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/shared"
)

type DeputyDeathRecord struct {
	DeputyID                      int         `json:"id"`
	DateOfDeath                   shared.Date `json:"dateOfDeath,omitempty"`
	ProofOfDeathReceived          string      `json:"proofOfDeathReceived,omitempty"`
	DateDeathCertificateReceieved shared.Date `json:"dateDeathCertificateReceived,omitempty"`
	DateNotified                  shared.Date `json:"dateNotified,omitempty"`
	NotifiedBy                    RefData     `json:"notifiedBy"`
	NotificationMethod            string      `json:"notificationMethod,omitempty"`
	Notes                         string      `json:"notes,omitempty"`
}

func (c *Client) UpdateDeputyDeathRecord(ctx Context, deathRecord DeputyDeathRecord) error {
	var body bytes.Buffer

	err := json.NewEncoder(&body).Encode(deathRecord)

	if err != nil {
		return err
	}
	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf(shared.SupervisionAPIPath+"/v1/deputies/lay/%d/death-notification", deathRecord.DeputyID), &body)

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
			return ValidationError{Errors: v.ValidationErrors}
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return newStatusError(resp)
	}

	return nil
}
