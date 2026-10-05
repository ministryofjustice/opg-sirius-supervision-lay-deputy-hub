package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type TimelineEvent struct {
	Id        int       `json:"id"`
	Hash      string    `json:"hash"`
	Timestamp time.Time `json:"timestamp"`
	EventType string    `json:"eventType"`
	User      struct {
		Id          int    `json:"id"`
		PhoneNumber string `json:"phoneNumber"`
		DisplayName string `json:"displayName"`
		Email       string `json:"email"`
	} `json:"user"`
	Event struct {
		PersonType     string      `json:"personType"`
		PersonId       string      `json:"personId"`
		PersonUid      string      `json:"personUid"`
		PersonName     string      `json:"personName"`
		PersonCourtRef interface{} `json:"personCourtRef"`
		Changes        []struct {
			FieldName string  `json:"fieldName"`
			OldValue  *string `json:"oldValue"`
			NewValue  string  `json:"newValue"`
			Type      string  `json:"type"`
		} `json:"changes,omitempty"`
		OrderType            string `json:"orderType,omitempty"`
		OrderUid             string `json:"orderUid,omitempty"`
		OrderId              string `json:"orderId,omitempty"`
		OrderCourtRef        string `json:"orderCourtRef,omitempty"`
		CourtReferenceNumber string `json:"courtReferenceNumber,omitempty"`
		CourtReference       string `json:"courtReference,omitempty"`
		AdditionalPersons    []struct {
			PersonType     string `json:"personType"`
			PersonId       string `json:"personId"`
			PersonUid      string `json:"personUid"`
			PersonName     string `json:"personName"`
			PersonCourtRef string `json:"personCourtRef"`
		} `json:"additionalPersons,omitempty"`
		PreviousStatus *string `json:"previousStatus,omitempty"`
		Status         string  `json:"status,omitempty"`
	} `json:"event"`
}

func (c *Client) GetDeputyTimeline(ctx Context, deputyID int) ([]TimelineEvent, error) {
	var v []TimelineEvent

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf(SupervisionAPIPath+"/v1/timeline/%d", deputyID), nil)
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
