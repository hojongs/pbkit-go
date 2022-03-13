package util

import "go.uber.org/zap"

var Sugar = getSugar()

func getSugar() *zap.SugaredLogger {
	logger, _ := zap.NewDevelopment()
	return logger.Sugar()
}
