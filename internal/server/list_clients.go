package server

import (
	"net/http"
	"sort"
	"strings"

	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/shared"
	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/sirius"
)

type ListClientsVars struct {
	AppVars
	Clients []DeputyClient
}

type DeputyClient struct {
	ClientId     int
	Firstname    string
	Surname      string
	CourtRef     string
	Relationship string
	DeputyStatus string
	Orders       []Order
}

type Order struct {
	OrderDate       shared.Date
	OrderStatus     string
	OrderType       string
	AppointmentType string
}

type DeputyClientsClient interface {
	GetDeputyClients(sirius.Context, sirius.ClientListParams) (sirius.ClientList, error)
}

func renderTemplateForClientTab(client DeputyClientsClient, tmpl Template) Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		clients, err := client.GetDeputyClients(getContext(r), sirius.ClientListParams{
			DeputyId: app.DeputyId(),
		})
		if err != nil {
			return err
		}

		app.PageName = "Clients"

		var items []DeputyClient
		for _, client := range clients.Clients {
			items = append(items, parseClientData(client))
		}

		return tmpl.ExecuteTemplate(w, "page", ListClientsVars{
			AppVars: app,
			Clients: items,
		})
	}
}

func parseClientData(client sirius.DeputyClient) DeputyClient {
	item := DeputyClient{
		ClientId:     client.ClientId,
		Firstname:    client.Firstname,
		Surname:      client.Surname,
		CourtRef:     client.CourtRef,
		Relationship: client.Relationship,
		DeputyStatus: client.DeputyStatus,
	}
	for _, order := range client.Orders {
		item.Orders = append(item.Orders, Order{
			OrderDate:       order.OrderDate,
			OrderStatus:     order.OrderStatus.Label,
			OrderType:       strings.ToUpper(order.CaseSubType),
			AppointmentType: order.AppointmentType,
		})
	}
	item.Orders = mostRelevantOrders(item.Orders)
	return item
}

func mostRelevantOrders(orders []Order) []Order {
	byType := make(map[string]Order)
	// Select the most relevant order for each type, preferring active orders and then the latest order date.
	for _, order := range orders {
		if selected, exists := byType[order.OrderType]; !exists || orderIsMoreRelevant(order, selected) {
			byType[order.OrderType] = order
		}
	}

	var selected []Order
	hasActiveOrder := false
	for _, order := range byType {
		selected = append(selected, order)
		hasActiveOrder = hasActiveOrder || order.OrderStatus == "Active"
	}
	// PFA comes first so it wins when inactive order dates are equal.
	sort.Slice(selected, func(i, j int) bool {
		return selected[i].OrderType > selected[j].OrderType
	})
	// If any order is active, or there is only one order, return all the most relevant orders.
	if hasActiveOrder || len(selected) < 2 {
		return selected
	}

	// If there are multiple inactive orders, return the most relevant order.
	latest := selected[0]
	for _, order := range selected[1:] {
		if orderIsMoreRelevant(order, latest) {
			latest = order
		}
	}
	return []Order{latest}
}

func orderIsMoreRelevant(candidate, selected Order) bool {
	candidateActive := candidate.OrderStatus == "Active"
	selectedActive := selected.OrderStatus == "Active"
	if candidateActive != selectedActive {
		return candidateActive
	}
	return candidate.OrderDate.After(selected.OrderDate)
}
