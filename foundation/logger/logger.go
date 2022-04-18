package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//NewLogger creates a new instance of zap.Logger
//For now it's creating a production level logger
//Next step is to pass the env param and depending on him
//Construct the appropriate logger
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

	log.Info("Logs Construct")

	return log, nil
}
