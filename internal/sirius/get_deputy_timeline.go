package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/shared"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type User struct {
	Id          int    `json:"id"`
	PhoneNumber string `json:"phoneNumber"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
}

type Changeset struct {
	FieldName         string `json:"fieldName"`
	OldValue          string `json:"oldValue"`
	NewValue          string `json:"newValue"`
	Type              string `json:"type"`
	TimelineValueName string
}

type AdditionalPerson struct {
	PersonType string `json:"personType"`
	PersonId   string `json:"personId"`
	PersonUid  string `json:"personUid"`
	PersonName string `json:"personName"`
}

type Event struct {
	AdditionalPersons            []AdditionalPerson `json:"additionalPersons"`
	Changes                      []Changeset        `json:"changes"`
	CourtReference               string             `json:"courtReference"`
	DateOfDeath                  shared.Date        `json:"dateOfDeath"`
	DateDeathCertificateReceived shared.Date        `json:"dateDeathCertificateReceived"`
	DateNotified                 shared.Date        `json:"dateNotified"`
	DocumentId                   string             `json:"documentId"`
	DocumentClientId             string             `json:"clientId"`
	DocumentSystemType           string             `json:"systemType"`
	Notes                        string             `json:"notes"`
	NotificationMethod           string             `json:"notificationMethod"`
	NotifiedBy                   string             `json:"notifiedBy"`
	NotifiedByOther              string             `json:"notifiedByOther"`
	ParentEventType              string             `json:"parentEventType"`
	PreviousStatus               string             `json:"previousStatus"`
	PersonType                   string             `json:"personType"`
	PersonId                     string             `json:"personId"`
	PersonUid                    string             `json:"personUid"`
	PersonName                   string             `json:"personName"`
	MultiPersons                 []AdditionalPerson `json:"multiPersons"`
	OrderType                    string             `json:"orderType"`
	OrderUid                     string             `json:"orderUid"`
	OrderId                      string             `json:"orderId"`
	OrderCourtRef                string             `json:"orderCourtRef"`
	Status                       string             `json:"status"`
}
type TimelineEvent struct {
	Id              int    `json:"id"`
	Timestamp       string `json:"timestamp"`
	TimestampAsDate string
	EventType       string `json:"eventType"`
	User            User   `json:"user"`
	Event           Event  `json:"event"`
}

func (c *Client) GetDeputyTimeline(ctx Context, deputyID int) ([]TimelineEvent, error) {
	var v []TimelineEvent

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf(shared.SupervisionAPIPath+"/v1/timeline/%d/lay-deputy", deputyID), nil)
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

	if err = json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return v, err
	}
	return editDeputyEvents(v), nil
}

func editDeputyEvents(events []TimelineEvent) []TimelineEvent {
	var list []TimelineEvent
	for _, e := range events {
		event := TimelineEvent{
			Timestamp:       shared.FormatDateTime(shared.IsoDateTimeZone, e.Timestamp, shared.SiriusDateTime),
			TimestampAsDate: shared.NewDate(e.Timestamp).String(),
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
	caser := cases.Title(language.English)

	eventToReturn := myEvent

	var updatedChanges []Changeset
	for _, e := range myEvent.Changes {
		updatedChangeset := updateChangeset(e)
		updatedChanges = append(updatedChanges, updatedChangeset)
	}
	eventToReturn.Changes = updatedChanges

	if strings.Contains(eventType, "PersonDeathNotified") || strings.Contains(eventType, "DeathNotificationDeleted") {
		eventToReturn.NotifiedBy = caser.String(myEvent.NotifiedBy)
	}

	if myEvent.OrderType != "" {
		eventToReturn.OrderType = strings.ToUpper(myEvent.OrderType)
	}

	return eventToReturn
}

func updateChangeset(c Changeset) Changeset {
	caser := cases.Title(language.English)

	if c.FieldName == "statusChangeDate" || c.FieldName == "dob" {
		return Changeset{
			FieldName:         c.FieldName,
			OldValue:          shared.FormatDateTime(shared.IsoDateTimeZone, c.OldValue, shared.SiriusDate),
			NewValue:          shared.FormatDateTime(shared.IsoDateTimeZone, c.NewValue, shared.SiriusDate),
			Type:              c.Type,
			TimelineValueName: GetTimelineValueForField(c.FieldName),
		}
	}

	if reformatAllCapitalValues(c.FieldName) {
		c.OldValue = caser.String(c.OldValue)
		c.NewValue = caser.String(c.NewValue)
	}
	return Changeset{
		FieldName:         c.FieldName,
		OldValue:          c.OldValue,
		NewValue:          c.NewValue,
		Type:              c.Type,
		TimelineValueName: GetTimelineValueForField(c.FieldName),
	}

}

func GetTimelineValueForField(myFieldName string) string {
	caser := cases.Title(language.English)

	switch myFieldName {
	case "addressLine1":
		return "Address line 1"
	case "addressLine2":
		return "Address line 2"
	case "addressLine3":
		return "Address line 3"
	case "audioTape":
		return "Audio"
	case "correspondenceByEmail":
		return "Email"
	case "correspondenceByPhone":
		return "Phone"
	case "correspondenceByPost":
		return "Post"
	case "correspondenceByWelsh":
		return "Welsh"
	case "deputyType":
		return "Type of deputy"
	case "deputySubType":
		return "Deputy Subtype"
	case "dob":
		return "Date of birth"
	case "email":
		return "Email address"
	case "firstname":
		return "First name"
	case "hearingImpaired":
		return "Hearing impaired"
	case "homePhoneNumber":
		return "Evening phone number"
	case "interpreterRequired":
		return "Interpreter required"
	case "isAirmailRequired":
		return "Is airmail required"
	case "largePrint":
		return "Large print"
	case "mainCorrespondent":
		return "Correspondent"
	case "middleNames":
		return "Middle names"
	case "mobileNumber":
		return "Mobile telephone number"
	case "occupation":
		return "Occupation"
	case "otherNames":
		return "Other names"
	case "relationshipToClient":
		return "Relationship to client"
	case "spellingOfNameRequiresCare":
		return "Spelling of name requires care"
	case "salutation":
		return "Title"
	case "statusChangeDate":
		return "Status change date"
	case "statusOnCase":
		return "Status on case"
	case "workPhoneNumber":
		return "Daytime telephone number"

	default:
		return caser.String(myFieldName)
	}
}

func reformatAllCapitalValues(myFieldName string) bool {
	switch myFieldName {
	case "deputyType", "deputySubType", "statusChangeDate", "statusOnCase", "relationshipToClient":
		return true
	default:
		return false
	}
}
