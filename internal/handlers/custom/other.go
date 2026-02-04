package custom

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	_ "embed"
)

//go:embed openapi.yaml
var swagger []byte

type HealthController struct {
}

func NewHealthController() *HealthController {
	external := &HealthController{}

	return external
}

// Routes returns all the api routes for the PingController
func (c *HealthController) Routes(r *chi.Mux) http.Handler {
	r.Get("/health", c.Health)
	r.Get("/ping", c.Ping)
	r.Get("/openapi", c.Swagger)

	return r
}

// Ping: probe endpoint
func (c *HealthController) Health(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
}

func (c *HealthController) Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content", "application/yaml")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"status": "success"}`)
}

func (c *HealthController) Swagger(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content", "application/yaml")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", string(swagger))
}
