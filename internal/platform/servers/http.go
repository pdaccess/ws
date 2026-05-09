package servers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	commons_domain "github.com/pdaccess/commons/pkg/domain"
	"github.com/pdaccess/commons/pkg/middleware"
	pdws_domain "github.com/pdaccess/ws/internal/core/domain"
	"github.com/pdaccess/ws/internal/core/ports"
	"github.com/pdaccess/ws/internal/platform/handlers"
	"github.com/pdaccess/ws/internal/platform/handlers/custom"
	"github.com/pdaccess/ws/internal/platform/handlers/external"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var doubleSlashRegex = regexp.MustCompile(`//+`)

func init() {
	logLevel := os.Getenv("LOG_LEVEL")
	level := zerolog.InfoLevel
	switch strings.ToLower(logLevel) {
	case "debug":
		level = zerolog.DebugLevel
	case "warn":
		level = zerolog.WarnLevel
	case "error":
		level = zerolog.ErrorLevel
	}
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger().Level(level)
}

func NewHttpServer(svc ports.Service) http.Handler {
	router := chi.NewRouter()

	router.Use(recoveryMiddleware)

	healthController := custom.NewHealthController()
	healthController.Routes(router)

	apiCall := router.With()
	apiCall.Use(loggingMiddleware)
	apiCall.Use(middleware.RequestIdAdapter)
	apiCall.Use(authzAdapter)

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

func authzAdapter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		client := middleware.ClientFromCtx(r.Context())

		claims := &commons_domain.PdaccessClaims{}

		xUserId := r.Header.Get("X-User-Id")
		if xUserId != "" {
			claims.UserId = xUserId
			claims.Urk = r.Header.Get("X-User-Urk")
			claims.Role = r.Header.Get("X-User-Role")
			claims.Realm = r.Header.Get("X-User-Realm")
			claims.AuthId = r.Header.Get("X-User-Auth-Id")
			client.User = *claims
			log.Debug().Str("user_id", claims.UserId).Str("source", "header").Msg("authenticated")
			next.ServeHTTP(w, r)
			return
		}

		authStr := r.Header.Get("Authorization")
		tokenStr, found := strings.CutPrefix(authStr, "Bearer ")

		if !found || tokenStr == "" {
			cookie, err := r.Cookie("pdaccess_session")
			if err == nil && cookie != nil && cookie.Value != "" {
				tokenStr = cookie.Value
				found = true
			}
		}

		if !found || tokenStr == "" {
			log.Warn().
				Str("reason", "missing_auth").
				Str("has_auth_header", strconv.FormatBool(authStr != "")).
				Interface("cookies", r.Cookies()).
				Msg("unauthorized")
			response := commons_domain.ApiResponse{
				Code:      commons_domain.ErrNotAuth.Code(),
				Message:   "Not Authorized",
				RequestId: client.RequestId,
			}
			buf, _ := json.Marshal(response)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(buf)
			return
		}

		decodedToken, _, err := jwt.NewParser().ParseUnverified(tokenStr, &commons_domain.PdaccessClaims{})
		if err != nil {
			log.Warn().Err(err).Str("reason", "token_parse_error").Msg("unauthorized")
			response := commons_domain.ApiResponse{
				Code:      commons_domain.ErrNotAuth.Code(),
				Message:   "Not Authorized",
				RequestId: client.RequestId,
			}
			buf, _ := json.Marshal(response)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(buf)
			return
		}

		claims, ok := decodedToken.Claims.(*commons_domain.PdaccessClaims)
		if !ok {
			log.Warn().Str("reason", "invalid_claims_type").Msg("unauthorized")
			response := commons_domain.ApiResponse{
				Code:      commons_domain.ErrNotAuth.Code(),
				Message:   "Not Authorized",
				RequestId: client.RequestId,
			}
			buf, _ := json.Marshal(response)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(buf)
			return
		}

		client.User = *claims
		log.Debug().Str("user_id", claims.UserId).Str("request_id", client.RequestId).Msg("authenticated")

		next.ServeHTTP(w, r)
	})
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
					err = pdws_domain.ErrValidation
				}
				responseErrorHandler(w, r, err)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

type apiError struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(apiError{
		Error: struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}{
			Code:    code,
			Message: message,
		},
	})
}

func requestErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	status, code, message := getErrorInfo(err)
	writeError(w, status, code, message)
}

func responseErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	status, code, message := getErrorInfo(err)
	if status == http.StatusInternalServerError {
		log.Error().Err(err).Stack().Interface("error_type", fmt.Sprintf("%T", err)).Msg("internal server error")
	}
	writeError(w, status, code, message)
}

func getErrorInfo(err error) (int, string, string) {
	var validationErr *json.UnmarshalTypeError
	var syntaxErr *json.SyntaxError

	switch {
	case errors.As(err, &validationErr):
		return http.StatusBadRequest, pdws_domain.ErrCodeValidation, "invalid request format"
	case errors.As(err, &syntaxErr):
		return http.StatusBadRequest, pdws_domain.ErrCodeValidation, "invalid JSON syntax"
	case errors.Is(err, io.EOF):
		return http.StatusBadRequest, pdws_domain.ErrCodeValidation, "empty request body"
	case errors.Is(err, pdws_domain.ErrValidation):
		return http.StatusBadRequest, pdws_domain.ErrCodeValidation, "validation failed"
	case errors.Is(err, pdws_domain.ErrInvalidID):
		return http.StatusBadRequest, pdws_domain.ErrCodeInvalidID, "invalid ID format"
	case errors.Is(err, pdws_domain.ErrNotFound):
		return http.StatusNotFound, pdws_domain.ErrCodeNotFound, "resource not found"
	case errors.Is(err, pdws_domain.ErrUnauthorized):
		return http.StatusUnauthorized, pdws_domain.ErrCodeUnauthorized, "unauthorized"
	}

	var ve pdws_domain.ValidationError
	if errors.As(err, &ve) {
		return http.StatusBadRequest, ve.Code, ve.Message
	}
	var ie pdws_domain.InvalidIDError
	if errors.As(err, &ie) {
		return http.StatusBadRequest, ie.Code, ie.Message
	}
	var nfe pdws_domain.NotFoundError
	if errors.As(err, &nfe) {
		return http.StatusNotFound, nfe.Code, nfe.Error()
	}
	var ine pdws_domain.InternalError
	if errors.As(err, &ine) {
		return http.StatusInternalServerError, ine.Code, ine.Message
	}

	return http.StatusInternalServerError, pdws_domain.ErrCodeInternal, "internal server error"
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		authHeader := r.Header.Get("Authorization")
		authPrefix := ""
		if len(authHeader) >= 7 {
			authPrefix = authHeader[:7]
		}
		log.Debug().
			Str("raw_path", r.URL.Path).
			Str("method", r.Method).
			Str("auth_header", authHeader).
			Str("auth_prefix", authPrefix).
			Msg("request_start")

		ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(ww, r)

		duration := time.Since(start)

		event := log.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", ww.statusCode).
			Dur("duration", duration).
			Str("remote_addr", r.RemoteAddr)

		client := middleware.ClientFromCtx(r.Context())
		if client != nil {
			event.Str("user_id", client.User.UserId)
		}

		if reqID := r.Context().Value(middleware.ClientRequestType(middleware.ClientRequestStr)); reqID != nil {
			if cr, ok := reqID.(*middleware.ClientRequest); ok {
				event.Str("request_id", cr.RequestId)
			}
		}

		event.Msg("request")
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
