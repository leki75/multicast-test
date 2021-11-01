package log

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func init() {
	config := zap.NewProductionConfig()
	config.EncoderConfig.CallerKey = zapcore.OmitKey
	config.EncoderConfig.EncodeTime = zapcore.EpochNanosTimeEncoder
	config.Sampling = nil
	config.OutputPaths = []string{"stdout"}

	var err error
	Logger, err = config.Build()
	if err != nil {
		log.Fatalf("error initializing zap logger: %s", err)
	}
}
