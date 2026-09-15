package server

import (
	"github.com/ministryofjustice/opg-go-common/paginate"
	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/urlbuilder"
)

type ListPage struct {
	AppVars
	AppliedFilters []string
	Sort           urlbuilder.Sort
	Error          string
	Pagination     paginate.Pagination
	PerPage        int
	UrlBuilder     urlbuilder.UrlBuilder
}
