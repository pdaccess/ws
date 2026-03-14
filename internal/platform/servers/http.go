package servers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/pdaccess/ws/internal/core/domain"
	"github.com/pdaccess/ws/internal/core/ports"
	"github.com/pdaccess/ws/internal/platform/handlers"
	"github.com/pdaccess/ws/internal/platform/handlers/custom"
	"github.com/pdaccess/ws/internal/platform/handlers/external"
	"github.com/rs/zerolog"
)

var logger = zerolog.New(nil).With().Timestamp().Logger()

func NewHttpServer(svc ports.Service) http.Handler {
	router := chi.NewRouter()

	router.Use(recoveryMiddleware)
	router.Use(loggingMiddleware)

	healthController := custom.NewHealthController()
	healthController.Routes(router)

	apiCall := router.With()

	Iservice := handlers.NewHttpHandler(svc)
	handler := external.NewStrictHandlerWithOptions(Iservice, []external.StrictMiddlewareFunc{}, external.StrictHTTPServerOptions{
		RequestErrorHandlerFunc:  requestErrorHandler,
		ResponseErrorHandlerFunc: responseErrorHandler,
	})
	external.HandlerWithOptions(handler, external.ChiServerOptions{
		BaseRouter:  apiCall,
		Middlewares: []external.MiddlewareFunc{},
	})

	return router
}

func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				var err error
				switch e := rec.(type) {
				case error:
					err = e
				default:
					err = domain.ErrValidation
				}
				responseErrorHandler(w, r, err)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func requestErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	status, code, message := getErrorInfo(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(external.GenericError{
		Error: struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}{
			Code:    code,
			Message: message,
		},
	})
}

func responseErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	status, code, message := getErrorInfo(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(external.GenericError{
		Error: struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}{
			Code:    code,
			Message: message,
		},
	})
}

func getErrorInfo(err error) (int, string, string) {
	var validationErr *json.UnmarshalTypeError
	var syntaxErr *json.SyntaxError

	switch {
	case errors.As(err, &validationErr):
		return http.StatusBadRequest, domain.ErrCodeValidation, domain.ErrCodeValidation
	case errors.As(err, &syntaxErr):
		return http.StatusBadRequest, domain.ErrCodeValidation, domain.ErrCodeValidation
	case errors.Is(err, io.EOF):
		return http.StatusBadRequest, domain.ErrCodeValidation, domain.ErrCodeValidation
	case errors.Is(err, domain.ErrValidation):
		return http.StatusBadRequest, domain.ErrCodeValidation, domain.ErrCodeValidation
	case errors.Is(err, domain.ErrInvalidID):
		return http.StatusBadRequest, domain.ErrCodeInvalidID, domain.ErrCodeInvalidID
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound, domain.ErrCodeNotFound, domain.ErrCodeNotFound
	case errors.Is(err, domain.ErrUnauthorized):
		return http.StatusUnauthorized, domain.ErrCodeUnauthorized, domain.ErrCodeUnauthorized
	}

	var ve domain.ValidationError
	if errors.As(err, &ve) {
		return http.StatusBadRequest, ve.Code, ve.Code
	}
	var ie domain.InvalidIDError
	if errors.As(err, &ie) {
		return http.StatusBadRequest, ie.Code, ie.Code
	}
	var nfe domain.NotFoundError
	if errors.As(err, &nfe) {
		return http.StatusNotFound, nfe.Code, nfe.Code
	}
	var ine domain.InternalError
	if errors.As(err, &ine) {
		return http.StatusInternalServerError, ine.Code, ine.Code
	}

	return http.StatusInternalServerError, domain.ErrCodeInternal, domain.ErrCodeInternal
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
