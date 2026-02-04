package servers

import (
	"net/http"

	"git.h2hsecure.com/core/ws/internal/handlers"
	"git.h2hsecure.com/core/ws/internal/handlers/custom"
	"git.h2hsecure.com/core/ws/internal/handlers/external"
	"github.com/go-chi/chi/v5"
)

func NewHttpServer() http.Handler {
	router := chi.NewRouter()

	healthController := custom.NewHealthController()
	healthController.Routes(router)

	apiCall := router.With()

	Iservice := external.NewStrictHandler(handlers.NewHttpHandler(), nil)
	external.HandlerWithOptions(Iservice, external.ChiServerOptions{
		BaseRouter:  apiCall,
		Middlewares: []external.MiddlewareFunc{},
	})

	return router
}
