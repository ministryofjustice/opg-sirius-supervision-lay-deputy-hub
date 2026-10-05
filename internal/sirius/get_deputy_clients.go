package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Metadata struct {
	TotalActiveClients int `json:"totalActiveClients"`
}

type ClientList struct {
	TotalClients int
	Metadata     Metadata
}

type ClientListParams struct {
	DeputyId int
}

func (c *Client) GetDeputyClients(ctx Context, params ClientListParams) (ClientList, error) {
	var clientList ClientList

	url := fmt.Sprintf(SupervisionAPIPath+"/v1/deputies/%d/clients", params.DeputyId)

	req, err := c.newRequest(ctx, http.MethodGet, url, nil)

	if err != nil {
		return clientList, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return clientList, err
	}

	defer unchecked(resp.Body.Close)

	if resp.StatusCode == http.StatusUnauthorized {
		return clientList, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return clientList, newStatusError(resp)
	}

	if err = json.NewDecoder(resp.Body).Decode(&clientList); err != nil {
		return clientList, err
	}

	clientList.TotalClients = clientList.Metadata.TotalActiveClients

	return clientList, err
}
