package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/shared"
)

func (c *Client) GetDeputyDeathRecord(ctx Context, DeputyID int) (DeputyDeathRecord, error) {
	var r struct {
		ID                           int         `json:"id"`
		ProofOfDeathReceived         bool        `json:"proofOfDeathReceived"`
		DateDeathCertificateReceived shared.Date `json:"dateDeathCertificateReceived"`
		DateNotified                 shared.Date `json:"dateNotified"`
		NotifiedBy                   RefData     `json:"notifiedBy"`
		NotifiedByOther              string      `json:"notifiedByOther"`
		NotificationMethod           string      `json:"notificationMethod"`
		Notes                        string      `json:"notes"`
		Person                       struct {
			DateOfDeath shared.Date `json:"dateOfDeath"`
		} `json:"person"`
	}

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf(shared.SupervisionAPIPath+"/v1/deputies/%d/death-notification", DeputyID), nil)
	if err != nil {
		return DeputyDeathRecord{}, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return DeputyDeathRecord{}, err
	}

	defer unchecked(resp.Body.Close)

	if resp.StatusCode == http.StatusUnauthorized {
		return DeputyDeathRecord{}, ErrUnauthorized
	}

	//Api throws a 404 exception if there's no record
	if resp.StatusCode == http.StatusNotFound {
		return DeputyDeathRecord{}, nil
	}

	if resp.StatusCode != http.StatusOK {
		return DeputyDeathRecord{}, newStatusError(resp)
	}

	err = json.NewDecoder(resp.Body).Decode(&r)

	notes := r.NotifiedByOther
	if r.NotifiedByOther == "" {
		notes = r.Notes
	}

	v := DeputyDeathRecord{
		DeputyID:                      DeputyID,
		ProofOfDeathReceived:          fmt.Sprintf("%v", r.ProofOfDeathReceived),
		DateDeathCertificateReceieved: r.DateDeathCertificateReceived,
		DateNotified:                  r.DateNotified,
		NotifiedBy:                    r.NotifiedBy,
		NotificationMethod:            strings.ToUpper(r.NotificationMethod),
		Notes:                         notes,
		DateOfDeath:                   r.Person.DateOfDeath,
	}
	return v, err
}
