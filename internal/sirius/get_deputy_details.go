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

type DeputyDetails struct {
	ID              int           `json:"id"`
	Salutation      string        `json:"salutation"`
	DeputyFirstName string        `json:"firstname"`
	DeputySurname   string        `json:"surname"`
	DeputyCasrecId  int           `json:"deputyCasrecId"`
	DisplayName     string        `json:"displayName"`
	CanDelete       bool          `json:"canDelete"`
	DeputyNumber    int           `json:"deputyNumber"`
	DeputySubType   DeputySubType `json:"deputySubType"`
	DeputyStatus    string        `json:"deputyStatus"`
	DeputyType      DeputyType    `json:"deputyType"`
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
