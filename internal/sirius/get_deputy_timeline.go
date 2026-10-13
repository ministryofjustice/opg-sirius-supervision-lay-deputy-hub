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
	AdditionalPersons               []AdditionalPersons `json:"additionalPersons,omitempty"`
	Changes                         []Changeset         `json:"changes,omitempty"`
	CourtReferenceNumber            string              `json:"courtReferenceNumber,omitempty"`
	CourtReference                  string              `json:"courtReference,omitempty"`
	DateOfDeath                     string              `json:"dateOfDeath,omitempty"`
	DateDeathCertificateReceived    string              `json:"dateDeathCertificateReceived,omitempty"`
	DateNotified                    string              `json:"dateNotified,omitempty"`
	DocumentId                      string              `json:"documentId,omitempty"`
	DocumentFilename                string              `json:"filename,omitempty"`
	DocumentDescription             string              `json:"description,omitempty"`
	DocumentDirection               string              `json:"direction,omitempty"`
	DocumentPublishedDate           string              `json:"publishedDate,omitempty"`
	DocumentType                    string              `json:"type,omitempty"`
	DocumentSource                  string              `json:"source,omitempty"`
	DocumentSystemType              string              `json:"systemType,omitempty"`
	DocumentMimeType                string              `json:"mimeType,omitempty"`
	DocumentPageCount               string              `json:"pageCount,omitempty"`
	DocumentPublishedById           string              `json:"publishedById,omitempty"`
	DocumentClientId                string              `json:"clientId,omitempty"`
	DocumentClientName              string              `json:"clientName,omitempty"`
	DocumentIsAnnualInvoiceDocument bool                `json:"isAnnualInvoiceDocument,omitempty"`
	Notes                           string              `json:"notes,omitempty"`
	NotificationMethod              string              `json:"notificationMethod,omitempty"`
	NotifiedBy                      string              `json:"notifiedBy,omitempty"`
	NotifiedByOther                 string              `json:"notifiedByOther,omitempty"`
	ParentEventType                 string              `json:"parentEventType"`
	PreviousStatus                  string              `json:"previousStatus,omitempty"`
	PersonType                      string              `json:"personType,omitempty"`
	PersonId                        string              `json:"personId,omitempty"`
	PersonUid                       string              `json:"personUid,omitempty"`
	PersonName                      string              `json:"personName,omitempty"`
	PersonCourtRef                  string              `json:"personCourtRef,omitempty"`
	ProofOfDeathReceived            bool                `json:"proofOfDeathReceived,omitempty"`
	MultiPersons                    []AdditionalPersons `json:"multiPersons,omitempty"`
	OrderType                       string              `json:"orderType,omitempty"`
	OrderUid                        string              `json:"orderUid,omitempty"`
	OrderId                         string              `json:"orderId,omitempty"`
	OrderCourtRef                   string              `json:"orderCourtRef,omitempty"`
	Status                          string              `json:"status,omitempty"`
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
			Event:           reformatEvent(e.Event, e.EventType),
		}

		list = append(list, event)
	}
	return list
}

func reformatEventType(s string) string {
	eventTypeArray := strings.Split(s, "\\")
	return eventTypeArray[len(eventTypeArray)-1]
}

func reformatEvent(myEvent Event, eventType string) Event {
	eventToReturn := myEvent

	if strings.Contains(eventType, "DeathNotificationDeleted") || strings.Contains(eventType, "PersonDeathNotified") {
		myTimeZone := IsoDateTime
		if strings.Contains(eventType, "DeathNotificationDeleted") {
			myTimeZone = IsoDateTimeZone
		}

		if myEvent.DateNotified != "" {
			eventToReturn.DateNotified = FormatDateTimeAsDate(myTimeZone, myEvent.DateNotified, SiriusDate)
		}
		if myEvent.DateOfDeath != "" {
			eventToReturn.DateOfDeath = FormatDateTimeAsDate(myTimeZone, myEvent.DateOfDeath, SiriusDate)
		}
		if myEvent.DateDeathCertificateReceived != "" {
			eventToReturn.DateDeathCertificateReceived = FormatDateTimeAsDate(myTimeZone, myEvent.DateDeathCertificateReceived, SiriusDate)
		}
	}

	if strings.Contains(eventType, "DeputyOrderDetailsChanged") || strings.Contains(eventType, "PersonDetailsChanged") {
		var updatedChanges []Changeset
		for _, e := range myEvent.Changes {
			changes := Changeset{}
			if e.FieldName == "statusChangeDate" {
				changes = Changeset{
					FieldName: e.FieldName,
					OldValue:  FormatDateTimeAsDate(IsoDateTimeZone, e.OldValue, SiriusDate),
					NewValue:  FormatDateTimeAsDate(IsoDateTimeZone, e.NewValue, SiriusDate),
					Type:      e.FieldName,
				}
			} else {
				changes = Changeset{
					FieldName: e.FieldName,
					OldValue:  strings.Title(strings.ToLower(e.OldValue)),
					NewValue:  strings.Title(strings.ToLower(e.NewValue)),
					Type:      e.FieldName,
				}
			}
			updatedChanges = append(updatedChanges, changes)

		}
		eventToReturn.Changes = updatedChanges
	}

	if strings.Contains(eventType, "PersonDeathNotified") || strings.Contains(eventType, "DeathNotificationDeleted") {
		eventToReturn.NotifiedBy = strings.Title(strings.ToLower(myEvent.NotifiedBy))
	}

	if myEvent.OrderType != "" {
		eventToReturn.OrderType = strings.ToUpper(myEvent.OrderType)
	}

	return eventToReturn
}
