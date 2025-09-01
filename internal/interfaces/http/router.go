package http

import (
	"api-ptf-core-business-orchestrator-go-ms/internal/interfaces/http/middleware"
	"api-ptf-core-business-orchestrator-go-ms/internal/interfaces/routes"
	"api-ptf-core-business-orchestrator-go-ms/internal/models"
	"api-ptf-core-business-orchestrator-go-ms/internal/pkg/logger"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// loggingResponseWriter envuelve http.ResponseWriter para capturar el status code
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// NewRouter creates a new HTTP router with all the routes
func NewRouter(a *models.Application) *mux.Router {
	r := mux.NewRouter()

	// Use the base path from config
	api := r.PathPrefix(a.Configs().HTTP.BasePath).Subrouter()

	routes.SetupRoutes(api, a)

	// Add middleware
	r.Use(middleware.RequestIDMiddleware)
	r.Use(loggingMiddleware)
	r.Use(mux.CORSMethodMiddleware(r))

	return r
}

// loggingMiddleware registra información detallada de cada petición HTTP
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Crear un response writer personalizado para capturar el status code
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Continuar con el siguiente handler
		next.ServeHTTP(lrw, r)

		// Calcular la duración
		duration := time.Since(start)

		// Obtener el request ID del contexto
		requestID := middleware.GetRequestID(r.Context())

		// Registrar la información de la petición
		logger := logger.Logger().With(
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Duration("duration", duration),
			zap.Int("status", lrw.statusCode),
			zap.String("request_id", requestID),
		)

		// Usar nivel de log apropiado basado en el status code
		if lrw.statusCode >= 400 {
			logger.Error("Request completed with error")
		} else {
			logger.Info("Request completed")
		}
	})
}
