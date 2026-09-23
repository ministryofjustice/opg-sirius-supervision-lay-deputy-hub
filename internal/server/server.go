package server

import (
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"net/url"

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

	static := staticFileHandler(envVars.WebDir)
	mux.Handle("/static/assets/", static)
	mux.Handle("/static/javascript/", static)
	mux.Handle("/static/stylesheets/", static)

	mux.Handle("GET /{id}", wrap(renderTemplateForDeputyHub(templates["deputy-details.gotmpl"])))
	mux.Handle("GET /{id}/clients", wrap(renderTemplateForClientTab(templates["clients.gotmpl"])))
	mux.Handle("GET /{id}/timeline", wrap(renderTemplateForDeputyHubEvents(templates["timeline.gotmpl"])))

	mux.Handle("GET /health-check", healthCheck())

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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

func staticFileHandler(webDir string) http.Handler {
	h := http.FileServer(http.Dir(webDir))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "must-revalidate")
		h.ServeHTTP(w, r)
	})
}
