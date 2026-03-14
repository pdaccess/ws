package servers

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/pdaccess/ws/internal/core/ports"
	"github.com/pdaccess/ws/internal/platform/handlers"
	"github.com/pdaccess/ws/internal/platform/handlers/custom"
	"github.com/pdaccess/ws/internal/platform/handlers/external"
	"github.com/rs/zerolog"
)

var logger = zerolog.New(nil).With().Timestamp().Logger()

func NewHttpServer(svc ports.Service) http.Handler {
	router := chi.NewRouter()

	router.Use(loggingMiddleware)

	healthController := custom.NewHealthController()
	healthController.Routes(router)

	apiCall := router.With()

	Iservice := handlers.NewHttpHandler(svc)
	handler := external.NewStrictHandler(Iservice, []external.StrictMiddlewareFunc{})
	external.HandlerWithOptions(handler, external.ChiServerOptions{
		BaseRouter:  apiCall,
		Middlewares: []external.MiddlewareFunc{},
	})

	return router
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(ww, r)

		duration := time.Since(start)

		logger.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", ww.statusCode).
			Dur("duration", duration).
			Str("remote_addr", r.RemoteAddr).
			Msg("request")
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	return rw.ResponseWriter.Write(b)
}
