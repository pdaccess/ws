package servers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pdaccess/ws/internal/platform/handlers"
	"github.com/pdaccess/ws/internal/platform/handlers/custom"
	"github.com/pdaccess/ws/internal/platform/handlers/external"
)

func NewHttpServer() http.Handler {
	router := chi.NewRouter()

	healthController := custom.NewHealthController()
	healthController.Routes(router)

	apiCall := router.With()

	Iservice := external.NewStrictHandler(handlers.NewHttpHandlerWithDefault(), nil)
	external.HandlerWithOptions(Iservice, external.ChiServerOptions{
		BaseRouter:  apiCall,
		Middlewares: []external.MiddlewareFunc{},
	})

	return router
}
