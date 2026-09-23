package server

import (
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/sirius"
	"golang.org/x/sync/errgroup"
)

type AppVars struct {
	Path           string
	XSRFToken      string
	DeputyDetails  sirius.DeputyDetails
	SuccessMessage string
	PageName       string
	Error          string
	Errors         sirius.ValidationErrors
	EnvironmentVars
}

type AppVarsClient interface {
	GetDeputyDetails(sirius.Context, int) (sirius.DeputyDetails, error)
}

func NewAppVars(client AppVarsClient, r *http.Request, envVars EnvironmentVars) (*AppVars, error) {
	ctx := getContext(r)
	_, groupCtx := errgroup.WithContext(ctx.Context)
	deputyId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return nil, StatusError(http.StatusBadRequest)
	}

	vars := AppVars{
		Path:            r.URL.Path,
		XSRFToken:       ctx.XSRFToken,
		EnvironmentVars: envVars,
	}

	deputy, err := client.GetDeputyDetails(ctx.With(groupCtx), deputyId)
	if err != nil {
		return nil, StatusError(http.StatusBadRequest)
	}
	vars.DeputyDetails = deputy

	return &vars, nil
}
