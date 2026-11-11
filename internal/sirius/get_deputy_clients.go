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

type OrderStatus struct {
	Label string `json:"label"`
}

type Order struct {
	OrderDate       shared.Date `json:"orderDate"`
	OrderStatus     OrderStatus `json:"orderStatus"`
	CaseSubType     string      `json:"casesubtype"`
	AppointmentType string      `json:"appointmentType"`
}

type DeputyClient struct {
	ClientId        int     `json:"id"`
	Firstname       string  `json:"firstname"`
	Surname         string  `json:"surname"`
	CourtRef        string  `json:"caseRecNumber"`
	Relationship    string  `json:"relationshipToClient"`
	DeputyStatus    string  `json:"statusOnCase"`
	AppointmentType string  `json:"appointmentType"`
	Orders          []Order `json:"orders"`
}

type ClientList struct {
	Clients      []DeputyClient `json:"clients"`
	TotalClients int            `json:"total"`
	Metadata     Metadata       `json:"metadata"`
}

type ClientListParams struct {
	DeputyId int
}

func (c *Client) GetDeputyClients(ctx Context, params ClientListParams) (ClientList, error) {
	var clientList ClientList

	requestURL := fmt.Sprintf(shared.SupervisionAPIPath+"/v1/deputies/%d/clients", params.DeputyId)

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
