package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/shared"
)

type Metadata struct {
	TotalActiveClients int `json:"totalActiveClients"`
}

type RefData struct {
	Handle string `json:"handle"`
	Label  string `json:"label"`
}

type OrderDeputy struct {
	Deputy struct {
		ID int `json:"id"`
	} `json:"deputy"`
	Relationship RefData `json:"relationshipToClient"`
	DeputyStatus RefData `json:"statusOnCase"`
}

type Order struct {
	OrderDate       shared.Date   `json:"orderDate"`
	OrderStatus     RefData       `json:"orderStatus"`
	CaseSubType     string        `json:"casesubtype"`
	AppointmentType RefData       `json:"howDeputyAppointed"`
	Deputies        []OrderDeputy `json:"deputies"`
}

type DeputyClient struct {
	ClientId  int     `json:"id"`
	Firstname string  `json:"firstname"`
	Surname   string  `json:"surname"`
	CourtRef  string  `json:"caseRecNumber"`
	Orders    []Order `json:"orders"`
}

type ClientList struct {
	Clients      []DeputyClient `json:"clients"`
	TotalClients int            `json:"total"`
	Metadata     Metadata       `json:"metadata"`
}

type ClientListParams struct {
	DeputyID int
}

func (c *Client) GetDeputyClients(ctx Context, params ClientListParams) (ClientList, error) {
	var clientList ClientList

	requestURL := fmt.Sprintf(shared.SupervisionAPIPath+"/v1/deputies/%d/clients", params.DeputyID)

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

	return clientList, err
}
