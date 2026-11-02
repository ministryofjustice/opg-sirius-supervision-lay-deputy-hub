package shared

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testJsonDateStruct struct {
	TestDate Date `json:"testDate"`
}

func TestNewDate(t *testing.T) {
	want := Date{Time: time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC)}
	assert.Equal(t, want, NewDate("31/12/2020"))

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("NewDate should panic with incorrect date format")
		}
	}()
	NewDate("12/31/2020") // wrong format should trigger a panic
}

func TestDate_Before_And_After(t *testing.T) {
	tests := []struct {
		name              string
		date1             Date
		date2             Date
		wantForBeforeTest bool
		wantForAfterTest  bool
	}{
		{
			name:              "Date1 is before Date2",
			date1:             NewDate("01/01/2020"),
			date2:             NewDate("02/01/2020"),
			wantForBeforeTest: true,
			wantForAfterTest:  false,
		},
		{
			name:              "Date1 is after Date2",
			date1:             NewDate("02/01/2020"),
			date2:             NewDate("01/01/2020"),
			wantForBeforeTest: false,
			wantForAfterTest:  true,
		},
		{
			name:              "Date1 is the same as Date2",
			date1:             NewDate("01/01/2020"),
			date2:             NewDate("01/01/2020"),
			wantForBeforeTest: false,
			wantForAfterTest:  false,
		},
		{
			name:              "Date1 is empty",
			date1:             Date{},
			date2:             NewDate("02/01/2020"),
			wantForBeforeTest: true,
			wantForAfterTest:  false,
		},
		{
			name:              "Date2 is empty",
			date1:             NewDate("01/01/2020"),
			date2:             Date{},
			wantForBeforeTest: false,
			wantForAfterTest:  true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.wantForBeforeTest, test.date1.Before(test.date2))
			assert.Equal(t, test.wantForAfterTest, test.date1.After(test.date2))
		})
	}
}

func TestDate_IsNull(t *testing.T) {
	tests := []struct {
		date Date
		want bool
	}{
		{
			date: NewDate("01/01/2020"),
			want: false,
		},
		{
			date: NewDate("01/01/0001"),
			want: true,
		},
		{
			date: Date{},
			want: true,
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
			assert.Equal(t, test.want, test.date.IsNull())
		})
	}
}

func TestDate_MarshalJSON(t *testing.T) {
	v := testJsonDateStruct{TestDate: NewDate("01/01/2020")}
	b, err := json.Marshal(v)
	assert.Nil(t, err)
	assert.Equal(t, `{"testDate":"01\/01\/2020"}`, string(b))
}

func TestDate_DateAndDateTimeString(t *testing.T) {
	tests := []struct {
		name               string
		inputDate          string
		wantDate           Date
		wantDateString     string
		wantDateTimeString string
	}{
		{
			name:               "returns correct format for slashers",
			inputDate:          "01/01/2020",
			wantDate:           Date{Time: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)},
			wantDateString:     "01/01/2020",
			wantDateTimeString: "01/01/2020 00:00:00",
		},
		{
			name:               "returns correct format for dashers",
			inputDate:          "2024-10-01",
			wantDate:           Date{Time: time.Date(2024, time.October, 1, 0, 0, 0, 0, time.UTC)},
			wantDateString:     "01/10/2024",
			wantDateTimeString: "01/10/2024 00:00:00",
		},
		{
			name:               "returns correct format for date time string",
			inputDate:          "2025-01-02T18:07:10+00:00",
			wantDate:           Date{Time: time.Date(2025, time.January, 2, 18, 7, 10, 0, time.UTC)},
			wantDateString:     "02/01/2025",
			wantDateTimeString: "02/01/2025 18:07:10",
		},
		{
			name:               "returns correct format for date time string",
			inputDate:          "02/01/2025 18:07:10",
			wantDate:           Date{Time: time.Date(2025, time.January, 2, 18, 7, 10, 0, time.UTC)},
			wantDateString:     "02/01/2025",
			wantDateTimeString: "02/01/2025 18:07:10",
		},
		{
			name:               "returns correct format for date time string",
			inputDate:          "2025-01-02T18:07:10+00:00",
			wantDate:           Date{Time: time.Date(2025, time.January, 2, 18, 7, 10, 0, time.UTC)},
			wantDateString:     "02/01/2025",
			wantDateTimeString: "02/01/2025 18:07:10",
		},
		{
			name:               "returns correct format for date time string",
			inputDate:          "2025-01-02T18:07:10+00:00",
			wantDate:           Date{Time: time.Date(2025, time.January, 2, 18, 7, 10, 0, time.UTC)},
			wantDateString:     "02/01/2025",
			wantDateTimeString: "02/01/2025 18:07:10",
		},
		{
			name:               "returns nothing for default date string",
			inputDate:          "01/01/0001",
			wantDate:           Date{Time: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)},
			wantDateString:     "",
			wantDateTimeString: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.wantDate, NewDate(tt.inputDate), "dateAndTimeTest: wantDate(%v)", tt.inputDate)
			assert.Equalf(t, tt.wantDateString, NewDate(tt.inputDate).DateString(), "dateAndTimeTest: wantDateString(%v)", tt.inputDate)
			assert.Equalf(t, tt.wantDateTimeString, NewDate(tt.inputDate).DateTimeString(), "dateAndTimeTest: wantDateTimeString(%v)", tt.inputDate)
		})
	}
}

func TestDate_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		json     string
		wantDate Date
		want     string
	}{
		{
			json:     `{"testDate":"01\/01\/2020"}`,
			wantDate: Date{Time: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)},
			want:     "01/01/2020",
		},
		{
			json:     `{"testDate":"01/01/2020"}`,
			wantDate: Date{Time: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)},
			want:     "01/01/2020",
		},
		{
			json:     `{"testDate":"2020-01-01T20:01:02+00:00"}`,
			wantDate: Date{Time: time.Date(2020, time.January, 1, 20, 1, 2, 0, time.UTC)},
			want:     "01/01/2020",
		},
		{
			json:     `{"testDate":"2020-01-01 20:01:02"}`,
			wantDate: Date{Time: time.Date(2020, time.January, 1, 20, 1, 2, 0, time.UTC)},
			want:     "01/01/2020",
		},
		{
			json:     `{"testDate":"01\/01\/2020 20:01:02"}`,
			wantDate: Date{Time: time.Date(2020, time.January, 1, 20, 1, 2, 0, time.UTC)},
			want:     "01/01/2020",
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
			var v *testJsonDateStruct
			err := json.Unmarshal([]byte(test.json), &v)
			assert.Nil(t, err)
			assert.Equal(t, test.wantDate, v.TestDate)
			assert.Equal(t, test.want, v.TestDate.DateString())
		})
	}
}
