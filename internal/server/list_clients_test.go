package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/shared"
	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

func TestNavigateToClients(t *testing.T) {
	assert := assert.New(t)

	client := &mockApiClient{
		Clients: sirius.ClientList{
			Clients: []sirius.DeputyClient{
				{
					ClientId:        67,
					Firstname:       "Robert",
					Surname:         "Robson",
					CourtRef:        "12323984",
					Relationship:    "Sibling",
					AppointmentType: "Joint and Severally",
					Orders: []sirius.Order{
						{
							OrderDate:       shared.NewDate("01/01/2024"),
							CaseSubType:     "pfa",
							OrderStatus:     sirius.OrderStatus{Label: "Active"},
							AppointmentType: "Joint and Severally",
						},
					},
					DeputyStatus: "Active",
				},
			},
			TotalClients: 1,
			Metadata:     sirius.Metadata{TotalActiveClients: 1},
		},
	}
	template := &mockTemplates{}
	app := AppVars{
		Path: "/333/clients",
		DeputyDetails: sirius.DeputyDetails{
			ID:              333,
			DeputyFirstName: "Test",
			DeputySurname:   "Deputy",
		},
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/333/clients", nil)

	err := renderTemplateForClientTab(client, template)(app, w, r)

	assert.NoError(err)
	assert.Equal(http.StatusOK, w.Result().StatusCode)
	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	app.PageName = "Clients"
	assert.Equal(ListClientsVars{
		AppVars: app,
		Clients: []DeputyClient{{
			ClientId:     67,
			Firstname:    "Robert",
			Surname:      "Robson",
			CourtRef:     "12323984",
			Relationship: "Sibling",
			DeputyStatus: "Active",
			Orders: []Order{{
				OrderDate:       shared.NewDate("01/01/2024"),
				OrderType:       "PFA",
				OrderStatus:     "Active",
				AppointmentType: "Joint and Severally",
			}},
		}},
	}, template.lastVars)
}

func TestClientTabWithNoClients(t *testing.T) {
	assert := assert.New(t)

	client := &mockApiClient{}
	template := &mockTemplates{}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/333/clients", nil)

	err := renderTemplateForClientTab(client, template)(AppVars{}, w, r)

	assert.NoError(err)
	assert.Equal(ListClientsVars{
		AppVars: AppVars{PageName: "Clients"},
	}, template.lastVars)
}

func TestHandlesErrorReturnedByGetDeputyClients(t *testing.T) {
	assert := assert.New(t)

	client := &mockApiClient{
		error: sirius.StatusError{Code: http.StatusInternalServerError},
	}
	template := &mockTemplates{}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/333/clients", nil)

	returnedError := renderTemplateForClientTab(client, template)(AppVars{}, w, r)

	assert.Equal(sirius.StatusError{Code: http.StatusInternalServerError}, returnedError)
	assert.Zero(template.count)
}

func TestParseClientDataDoesNotModifyAPIDataOrDefaultAppointments(t *testing.T) {
	client := sirius.DeputyClient{
		AppointmentType: "Joint",
		Orders: []sirius.Order{
			{OrderDate: shared.NewDate("01/01/2020"), CaseSubType: "hw", OrderStatus: sirius.OrderStatus{Label: "Active"}},
			{OrderDate: shared.NewDate("01/01/2019"), CaseSubType: "pFa", OrderStatus: sirius.OrderStatus{Label: "Active"}, AppointmentType: "Sole"},
			{OrderDate: shared.NewDate("01/01/2018"), CaseSubType: "pfa", OrderStatus: sirius.OrderStatus{Label: "Active"}, AppointmentType: "Joint and Severally"},
		},
	}
	originalOrders := append([]sirius.Order(nil), client.Orders...)

	item := parseClientData(client)

	assert.Equal(t, []Order{
		{OrderDate: shared.NewDate("01/01/2019"), OrderType: "PFA", OrderStatus: "Active", AppointmentType: "Sole"},
		{OrderDate: shared.NewDate("01/01/2020"), OrderType: "HW", OrderStatus: "Active"},
	}, item.Orders)
	assert.Equal(t, originalOrders, client.Orders)
}

func TestMostRelevantOrders(t *testing.T) {
	activePFA := Order{OrderDate: shared.NewDate("01/12/2020"), OrderStatus: "Active", OrderType: "PFA"}
	newerActivePFA := Order{OrderDate: shared.NewDate("01/01/2021"), OrderStatus: "Active", OrderType: "PFA"}
	closedPFA := Order{OrderDate: shared.NewDate("01/12/2024"), OrderStatus: "Closed", OrderType: "PFA"}
	activeHW := Order{OrderDate: shared.NewDate("01/12/2020"), OrderStatus: "Active", OrderType: "HW"}
	newerActiveHW := Order{OrderDate: shared.NewDate("01/01/2021"), OrderStatus: "Active", OrderType: "HW"}
	closedHW := Order{OrderDate: shared.NewDate("01/12/2023"), OrderStatus: "Closed", OrderType: "HW"}
	tiedClosedHW := Order{OrderDate: closedPFA.OrderDate, OrderStatus: "Closed", OrderType: "HW"}

	for _, test := range []struct {
		name     string
		orders   []Order
		expected []Order
	}{
		{name: "no orders"},
		{
			name:     "one PFA order",
			orders:   []Order{activePFA},
			expected: []Order{activePFA},
		},
		{
			name:     "dated order takes precedence over a missing date",
			orders:   []Order{activePFA, {OrderStatus: "Active", OrderType: "PFA"}},
			expected: []Order{activePFA},
		},
		{
			name:     "dated order replaces an order with a missing date",
			orders:   []Order{{OrderStatus: "Active", OrderType: "PFA"}, activePFA},
			expected: []Order{activePFA},
		},
		{
			name:     "latest active PFA despite a newer inactive order",
			orders:   []Order{activePFA, closedPFA, newerActivePFA},
			expected: []Order{newerActivePFA},
		},
		{
			name:     "latest active HW despite a newer inactive order",
			orders:   []Order{newerActiveHW, closedHW, activeHW},
			expected: []Order{newerActiveHW},
		},
		{
			name:     "latest inactive order when one type has no active orders",
			orders:   []Order{closedPFA, {OrderDate: shared.NewDate("01/01/2020"), OrderStatus: "Closed", OrderType: "PFA"}},
			expected: []Order{closedPFA},
		},
		{
			name:     "latest active order of each type",
			orders:   []Order{activeHW, closedPFA, newerActivePFA, activePFA, newerActiveHW, closedHW},
			expected: []Order{newerActivePFA, newerActiveHW},
		},
		{
			name:     "both active types with the same date",
			orders:   []Order{activeHW, activePFA},
			expected: []Order{activePFA, activeHW},
		},
		{
			name:     "active PFA and inactive HW",
			orders:   []Order{closedHW, activePFA},
			expected: []Order{activePFA, closedHW},
		},
		{
			name:     "inactive PFA and active HW",
			orders:   []Order{activeHW, closedPFA},
			expected: []Order{closedPFA, activeHW},
		},
		{
			name:     "both inactive with newer PFA",
			orders:   []Order{closedHW, closedPFA},
			expected: []Order{closedPFA},
		},
		{
			name:     "both inactive with newer HW",
			orders:   []Order{closedHW, {OrderDate: shared.NewDate("01/01/2020"), OrderStatus: "Duplicate", OrderType: "PFA"}},
			expected: []Order{closedHW},
		},
		{
			name:     "inactive dates tied favour PFA",
			orders:   []Order{tiedClosedHW, closedPFA},
			expected: []Order{closedPFA},
		},
		{
			name:     "inactive dates tied favour PFA in either input order",
			orders:   []Order{closedPFA, tiedClosedHW},
			expected: []Order{closedPFA},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			original := append([]Order(nil), test.orders...)
			assert.Equal(t, test.expected, mostRelevantOrders(test.orders))
			assert.Equal(t, original, test.orders)
		})
	}
}
