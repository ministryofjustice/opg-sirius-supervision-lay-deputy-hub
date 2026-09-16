package server

import (
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/ministryofjustice/opg-go-common/securityheaders"
	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/ministryofjustice/opg-sirius-lay-deputy-hub/internal/sirius"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Client interface {
	LayDeputyHubClient
}

type Template interface {
	ExecuteTemplate(io.Writer, string, interface{}) error
}

func New(logger *slog.Logger, client Client, templates map[string]*template.Template, envVars EnvironmentVars) http.Handler {
	mux := http.NewServeMux()

	wrap := wrapHandler(logger, client, templates["error.gotmpl"], envVars)

	static := http.FileServer(http.Dir(envVars.WebDir + "/static"))
	mux.Handle("/assets/", static)
	mux.Handle("/javascript/", static)
	mux.Handle("/stylesheets/", static)

	// Health check
	mux.Handle("/health-check", healthCheck())

	deputyDetails := wrap(renderTemplateForDeputyHub(templates["deputy-details.gotmpl"]))
	clients := wrap(renderTemplateForClientTab(templates["clients.gotmpl"]))
	timeline := wrap(renderTemplateForDeputyHubEvents(templates["timeline.gotmpl"]))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(pathParts) == 1 && pathParts[0] != "" {
			r.SetPathValue("id", pathParts[0])
			deputyDetails.ServeHTTP(w, r)
			return
		}

		if len(pathParts) == 2 && pathParts[0] != "" {
			r.SetPathValue("id", pathParts[0])
			switch pathParts[1] {
			case "clients":
				clients.ServeHTTP(w, r)
				return
			case "timeline":
				timeline.ServeHTTP(w, r)
				return
			}
		}

		w.WriteHeader(http.StatusNotFound)
		_ = templates["error.gotmpl"].ExecuteTemplate(w, "page", ErrorVars{
			Code:            http.StatusNotFound,
			Error:           "Page not found",
			EnvironmentVars: envVars,
		})
	})

	return otelhttp.NewHandler(http.StripPrefix(envVars.Prefix, telemetry.Middleware(logger)(securityheaders.Use(mux))), "supervision-lay-deputy-hub")
}

func getContext(r *http.Request) sirius.Context {
	token := ""

	if r.Method == http.MethodGet {
		if cookie, err := r.Cookie("XSRF-TOKEN"); err == nil {
			token, _ = url.QueryUnescape(cookie.Value)
		}
	} else {
		token = r.FormValue("xsrfToken")
	}

	return sirius.Context{
		Context:   r.Context(),
		Cookies:   r.Cookies(),
		XSRFToken: token,
	}
}
