package shared

import (
	"time"
)

const IsoDateTimeZone string = "2006-01-02T15:04:05Z07:00"
const IsoDateTime string = "2006-01-02 15:04:05"
const IsoDate string = "2006-01-02"
const SiriusDate string = "02/01/2006"
const SiriusDateTime string = "02/01/2006 15:04:05"

var londonLocation = func() *time.Location {
	loc, err := time.LoadLocation("Europe/London")
	if err != nil {
		return time.UTC
	}
	return loc
}()

func FormatDateTime(currentFormat string, dateString string, displayFormat string) string {
	if dateString == "" {
		return dateString
	}
	stringToDateTime, err := time.Parse(currentFormat, dateString)
	if err != nil {
		return dateString
	}
	dateTime := stringToDateTime.In(londonLocation).Format(displayFormat)
	return dateTime
}
