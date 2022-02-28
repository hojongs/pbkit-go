package log

import "go.uber.org/zap"

var Sugar = getSugar()

func Infow(msg string, keysAndValues ...interface{}) {
	Sugar.Infow(msg, keysAndValues)
}

func Fatalw(msg string, cause interface{}, keysAndValues ...interface{}) {
	Sugar.Fatalw(msg, append(keysAndValues, "cause", cause)...)
}

func getSugar() *zap.SugaredLogger {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // flushes buffer, if any
	return logger.Sugar()
}
