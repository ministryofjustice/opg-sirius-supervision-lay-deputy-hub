package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/shared"
)

type DeputyDeathRecord struct {
	DeputyId                      int    `json:"personId"`
	DateOfDeath                   string `json:"dateOfDeath"`
	ProofOfDeathReceived          bool   `json:"proofOfDeathReceived"`
	DateDeathCertificateReceieved string `json:"dateDeathCertificateReceived"`
	DateNotified                  string `json:"dateNotified"`
	NotifiedBy                    string `json:"notifiedBy"`
	NotifiedByOther               string `json:"notifiedByOther"`
	NotificationMethod            string `json:"notificationMethod"`
	Notes                         string `json:"notes"`
}

func (c *Client) CreateDeputyDeathRecord(ctx Context, deathRecord DeputyDeathRecord) error {
	var body bytes.Buffer

	err := json.NewEncoder(&body).Encode(deathRecord)

	if err != nil {
		return err
	}
	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf(shared.SupervisionAPIPath+"/v1/deputies/%d/death-notification", deathRecord.DeputyId), &body)

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

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var v struct {
			ValidationErrors ValidationErrors `json:"validation_errors"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil && len(v.ValidationErrors) > 0 {
			return ValidationError{Errors: v.ValidationErrors}
		}

		return newStatusError(resp)
	}

	return nil
}
