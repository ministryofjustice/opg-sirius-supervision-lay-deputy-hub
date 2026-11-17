package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/shared"
)

func (c *Client) GetDeputyDeathRecord(ctx Context, DeputyId int) (DeputyDeathRecord, error) {
	var v DeputyDeathRecord

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf(shared.SupervisionAPIPath+"/v1/deputies/%d/death-notification", DeputyId), nil)
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
