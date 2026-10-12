package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type User struct {
	Id          int    `json:"id"`
	PhoneNumber string `json:"phoneNumber"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
}

type Changeset struct {
	FieldName string `json:"fieldName"`
	OldValue  string `json:"oldValue"`
	NewValue  string `json:"newValue"`
	Type      string `json:"type"`
}

type AdditionalPersons struct {
	PersonType     string `json:"personType,omitempty"`
	PersonId       string `json:"personId,omitempty"`
	PersonUid      string `json:"personUid,omitempty"`
	PersonName     string `json:"personName,omitempty"`
	PersonCourtRef string `json:"personCourtRef,omitempty"`
}

type Event struct {
	NotifiedBy                   string              `json:"notifiedBy,omitempty"`
	ProofOfDeathReceived         bool                `json:"proofOfDeathReceived,omitempty"`
	DateOfDeath                  string              `json:"dateOfDeath,omitempty"`
	DateDeathCertificateReceived string              `json:"dateDeathCertificateReceived,omitempty"`
	DateNotified                 string              `json:"dateNotified,omitempty"`
	NotificationMethod           string              `json:"notificationMethod,omitempty"`
	Notes                        string              `json:"notes,omitempty"`
	PersonType                   string              `json:"personType,omitempty"`
	PersonId                     string              `json:"personId,omitempty"`
	PersonUid                    string              `json:"personUid,omitempty"`
	PersonName                   string              `json:"personName,omitempty"`
	PersonCourtRef               string              `json:"personCourtRef,omitempty"`
	Changes                      []Changeset         `json:"changes,omitempty"`
	OrderType                    string              `json:"orderType,omitempty"`
	OrderUid                     string              `json:"orderUid,omitempty"`
	OrderId                      string              `json:"orderId,omitempty"`
	OrderCourtRef                string              `json:"orderCourtRef,omitempty"`
	CourtReferenceNumber         string              `json:"courtReferenceNumber,omitempty"`
	CourtReference               string              `json:"courtReference,omitempty"`
	AdditionalPersons            []AdditionalPersons `json:"additionalPersons,omitempty"`
	MultiPersons                 []AdditionalPersons `json:"multiPersons,omitempty"`
	ParentEventType              string              `json:"parentEventType"`
	PreviousStatus               string              `json:"previousStatus,omitempty"`
	Status                       string              `json:"status,omitempty"`
}
type TimelineEvent struct {
	Id              int    `json:"id"`
	Timestamp       string `json:"timestamp"`
	TimestampAsDate string `json:"timestampAsDate"`
	EventType       string `json:"eventType"`
	User            User   `json:"user"`
	Event           Event  `json:"event"`
}

func (c *Client) GetDeputyTimeline(ctx Context, deputyID int) ([]TimelineEvent, error) {
	var v []TimelineEvent

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf(SupervisionAPIPath+"/v1/timeline/%d/lay-deputy", deputyID), nil)
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

	return editDeputyEvents(v), nil
}

func editDeputyEvents(events []TimelineEvent) []TimelineEvent {
	var list []TimelineEvent
	for _, e := range events {
		event := TimelineEvent{
			Timestamp:       FormatDateTime(IsoDateTimeZone, e.Timestamp, SiriusDateTime),
			TimestampAsDate: FormatDateTimeAsDate(IsoDateTimeZone, e.Timestamp, SiriusDate),
			EventType:       reformatEventType(e.EventType),
			Id:              e.Id,
			User:            e.User,
			Event:           e.Event,
		}

		list = append(list, event)
	}
	return list
}

func reformatEventType(s string) string {
	eventTypeArray := strings.Split(s, "\\")
	return eventTypeArray[len(eventTypeArray)-1]
}
