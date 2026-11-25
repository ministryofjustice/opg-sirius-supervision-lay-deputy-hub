package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/shared"
)

func (c *Client) DeleteDeputyDeathRecord(ctx Context, deputyId int) error {
	req, err := c.newRequest(ctx, http.MethodDelete, fmt.Sprintf(shared.SupervisionAPIPath+"/v1/deputies/%d/death-notification", deputyId), nil)

	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}

	defer unchecked(resp.Body.Close)

	if resp.StatusCode == http.StatusNotFound {
		return nil
	}

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
