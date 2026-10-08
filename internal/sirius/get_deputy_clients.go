package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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

	query := url.Values{}
	filter := "order-status:ACTIVE"
	query.Set("filter", filter)
	requestURL := fmt.Sprintf(SupervisionAPIPath+"/v1/deputies/%d/clients?%s", params.DeputyId, query.Encode())

	req, err := c.newRequest(ctx, http.MethodGet, requestURL, nil)

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

//func (p ClientListParams) CreateFilter() string {
//	var filter string
//	for _, s := range p.OrderStatuses {
//		filter += "order-status:" + s + ","
//	}
//	return strings.TrimRight(filter, ",")
//}
