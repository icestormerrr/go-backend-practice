package logger

import "go.uber.org/zap"

func New() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stderr"}
	cfg.DisableCaller = true
	cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)

	return cfg.Build()
}
