package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Log es la instancia global del logger
	Log *zap.Logger
)

// InitLogger inicializa el logger de la aplicación
func InitLogger(debug bool) error {
	var config zap.Config

	if debug {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		config = zap.NewProductionConfig()
	}

	// Configuración común
	config.EncoderConfig.TimeKey = "time"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.MessageKey = "message"
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	// Construir el logger
	var err error
	Log, err = config.Build()
	if err != nil {
		return err
	}

	// Redirigir logs estándar
	zap.RedirectStdLog(Log)

	return nil
}

// Sync cierra el buffer de logs
func Sync() error {
	return Log.Sync()
}
