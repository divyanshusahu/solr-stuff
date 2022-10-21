package log

import (
	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

func InitLogger() {
	l, _ := zap.NewProduction()
	Logger = l.Sugar()
}
