package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//NewLogger create an instance of the logger use throughout
//the application providing human-readable log messages
func NewLogger(service string) (*zap.Logger, error) {
	config := zap.NewProductionConfig()

	config.OutputPaths = []string{"stdout"}
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.DisableStacktrace = true

	config.InitialFields = map[string]interface{}{
		"service": service,
	}

	log, err := config.Build()

	if err != nil {
		return nil, err
	}

	return log, nil
}
