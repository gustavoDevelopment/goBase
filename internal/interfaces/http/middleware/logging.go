package middleware

import (
	"log"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// LoggingMiddleware registra información sobre cada petición HTTP
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Crear un response writer personalizado para capturar el status code
		lrw := newLoggingResponseWriter(w)

		// Continuar con el siguiente handler
		next.ServeHTTP(lrw, r)

		// Calcular la duración
		duration := time.Since(start)

		// Obtener el request ID del contexto
		requestID := GetRequestID(r.Context())

		// Registrar la información de la petición
		logger := zap.L().With(
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

// loggingResponseWriter envuelve http.ResponseWriter para capturar el status code
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK} // Status por defecto 200
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// Helper para inicializar el logger
func getLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("No se pudo inicializar el logger: %v", err)
	}
	return logger
}
