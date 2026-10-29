package shared

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatDateTime(t *testing.T) {
	tests := []struct {
		name                              string
		currentFormat                     string
		dateString                        string
		displayFormat                     string
		expectedResponseForFormatDateTime string
	}{
		{
			name:                              "IsoDateTimeZone to SiriusDate",
			currentFormat:                     IsoDateTimeZone,
			dateString:                        "1985-05-03T13:55:23+01:00",
			displayFormat:                     SiriusDate,
			expectedResponseForFormatDateTime: "03/05/1985",
		},
		{
			name:                              "IsoDateTime to SiriusDate",
			currentFormat:                     IsoDateTime,
			dateString:                        "1985-05-03 13:55:23",
			displayFormat:                     SiriusDate,
			expectedResponseForFormatDateTime: "03/05/1985",
		},
		{
			name:                              "IsoDate to SiriusDate",
			currentFormat:                     IsoDate,
			dateString:                        "1985-05-03",
			displayFormat:                     SiriusDate,
			expectedResponseForFormatDateTime: "03/05/1985",
		},
		{
			name:                              "SiriusDateTime to SiriusDate",
			currentFormat:                     SiriusDateTime,
			dateString:                        "03/05/1985 13:55:23",
			displayFormat:                     SiriusDate,
			expectedResponseForFormatDateTime: "03/05/1985",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expectedResponseForFormatDateTime, FormatDateTime(test.currentFormat, test.dateString, test.displayFormat))
		})
	}
}
