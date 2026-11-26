package sirius

import (
	"encoding/json"
	"net/http"
	"slices"

	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/shared"
)

type UserDetails struct {
	ID       int      `json:"id"`
	Roles    []string `json:"roles"`
	Username string   `json:"displayName"`
}

func (u User) IsManager() bool {
	return slices.Contains(u.Roles, "Manager")
}

func (c *Client) GetUserDetails(ctx Context) (User, error) {
	var v User

	req, err := c.newRequest(ctx, http.MethodGet, shared.SupervisionAPIPath+"/v1/users/current", nil)
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
