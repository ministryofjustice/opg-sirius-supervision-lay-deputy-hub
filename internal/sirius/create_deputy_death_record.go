package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/shared"
)

type DeputyDeathRecord struct {
	DeputyId                      int    `json:"id"`
	DateOfDeath                   string `json:"dateOfDeath,omitempty"`
	ProofOfDeathReceived          string `json:"proofOfDeathReceived,omitempty"`
	DateDeathCertificateReceieved string `json:"dateDeathCertificateReceived,omitempty"`
	DateNotified                  string `json:"dateNotified,omitempty"`
	NotifiedBy                    string `json:"notifiedBy,omitempty"`
	NotificationMethod            string `json:"notificationMethod,omitempty"`
	Notes                         string `json:"notes,omitempty"`
}

func (c *Client) CreateDeputyDeathRecord(ctx Context, deathRecord DeputyDeathRecord, logger *slog.Logger) error {
	var body bytes.Buffer

	err := json.NewEncoder(&body).Encode(deathRecord)

	if err != nil {
		return err
	}
	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf(shared.SupervisionAPIPath+"/v1/deputies/lay/%d/death-notification", deathRecord.DeputyId), &body)

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
		logger.Info("200 300 resp code")
		return newStatusError(resp)
	}

	return nil
}
