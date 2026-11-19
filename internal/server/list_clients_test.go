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
					ClientId:  67,
					Firstname: "Robert",
					Surname:   "Robson",
					CourtRef:  "12323984",
					Orders: []sirius.Order{
						{
							OrderDate:       shared.NewDate("01/01/2024"),
							CaseSubType:     "pfa",
							OrderStatus:     sirius.RefData{Label: "Active"},
							AppointmentType: sirius.RefData{Label: "Joint and Severally"},
							Deputies: []sirius.OrderDeputy{
								{
									Deputy: struct {
										ID int `json:"id"`
									}{ID: 333},
									Relationship: sirius.RefData{Label: "Sibling"},
									DeputyStatus: sirius.RefData{Label: "Active"},
								},
							},
						},
					},
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
			ClientId:  67,
			Firstname: "Robert",
			Surname:   "Robson",
			CourtRef:  "12323984",
			Orders: []Order{{
				OrderDate:       shared.NewDate("01/01/2024"),
				OrderType:       "PFA",
				OrderStatus:     "Active",
				AppointmentType: "Joint and Severally",
				Relationship:    "Sibling",
				StatusOnCase:    "Active",
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

func TestParseClientData(t *testing.T) {
	deputy := func(id int, relationship, status string) sirius.OrderDeputy {
		var result sirius.OrderDeputy
		result.Deputy.ID = id
		result.Relationship = sirius.RefData{Label: relationship}
		result.DeputyStatus = sirius.RefData{Label: status}
		return result
	}
	pfa := sirius.Order{
		OrderDate:       shared.NewDate("01/01/2019"),
		CaseSubType:     "pfa",
		OrderStatus:     sirius.RefData{Label: "Active"},
		AppointmentType: sirius.RefData{Label: "Sole"},
		Deputies: []sirius.OrderDeputy{
			deputy(999, "Spouse", "Discharged"),
			deputy(333, "Sibling", "Active"),
			deputy(444, "Parent", "Error"),
		},
	}
	hw := sirius.Order{
		OrderDate:       shared.NewDate("01/01/2020"),
		CaseSubType:     "hw",
		OrderStatus:     sirius.RefData{Label: "Active"},
		AppointmentType: sirius.RefData{Label: "Joint and Severally"},
		Deputies:        []sirius.OrderDeputy{deputy(333, "Sibling", "Suspended")},
	}
	expectedPFA := Order{
		OrderDate:       pfa.OrderDate,
		OrderType:       "PFA",
		OrderStatus:     "Active",
		AppointmentType: "Sole",
		Relationship:    "Sibling",
		StatusOnCase:    "Active",
	}
	expectedHW := Order{
		OrderDate:       hw.OrderDate,
		OrderType:       "HW",
		OrderStatus:     "Active",
		AppointmentType: "Joint and Severally",
		Relationship:    "Sibling",
		StatusOnCase:    "Suspended",
	}
	unrelated := pfa
	unrelated.OrderDate = shared.NewDate("01/01/2025")
	unrelated.Deputies = []sirius.OrderDeputy{deputy(999, "Spouse", "Active")}
	older := pfa
	older.OrderDate = shared.NewDate("01/01/2018")
	noDeputies := pfa
	noDeputies.Deputies = nil
	noType := pfa
	noType.CaseSubType = ""
	missingDetails := hw
	missingDetails.AppointmentType = sirius.RefData{}
	missingDetails.Deputies = []sirius.OrderDeputy{deputy(333, "", "")}
	duplicateDeputy := pfa
	duplicateDeputy.Deputies = []sirius.OrderDeputy{
		deputy(333, "Sibling", "Active"),
		deputy(333, "Spouse", "Discharged"),
	}

	for _, test := range []struct {
		name     string
		deputyID int
		orders   []sirius.Order
		expected []Order
	}{
		{name: "no orders", deputyID: 333},
		{
			name:     "maps details from the matching deputy rather than other deputies",
			deputyID: 333, orders: []sirius.Order{pfa}, expected: []Order{expectedPFA},
		},
		{
			name:     "keeps separate deputy details for each order with PFA first",
			deputyID: 333, orders: []sirius.Order{hw, pfa}, expected: []Order{expectedPFA, expectedHW},
		},
		{
			name:     "filters unrelated newer orders before selecting the latest per type",
			deputyID: 333, orders: []sirius.Order{unrelated, older, pfa}, expected: []Order{expectedPFA},
		},
		{
			name:     "no orders for the requested deputy",
			deputyID: 123, orders: []sirius.Order{pfa, hw},
		},
		{
			name:     "ignores orders without deputies",
			deputyID: 333, orders: []sirius.Order{noDeputies},
		},
		{
			name:     "preserves missing appointment and deputy details",
			deputyID: 333, orders: []sirius.Order{missingDetails},
			expected: []Order{{OrderDate: hw.OrderDate, OrderType: "HW", OrderStatus: "Active"}},
		},
		{
			name:     "uses the first matching deputy only",
			deputyID: 333, orders: []sirius.Order{duplicateDeputy}, expected: []Order{expectedPFA},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			client := sirius.DeputyClient{
				ClientId: 67, Firstname: "Robert", Surname: "Robson", CourtRef: "12323984",
				Orders: test.orders,
			}
			original := client
			if client.Orders != nil {
				original.Orders = append([]sirius.Order(nil), client.Orders...)
				for i, order := range client.Orders {
					original.Orders[i].Deputies = append([]sirius.OrderDeputy(nil), order.Deputies...)
				}
			}

			item := parseClientData(client, test.deputyID)

			assert.Equal(t, DeputyClient{
				ClientId: 67, Firstname: "Robert", Surname: "Robson", CourtRef: "12323984",
				Orders: test.expected,
			}, item)
			assert.Equal(t, original, client)
		})
	}
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
