package logger

import (
	"gorooster/helpers"

	"go.uber.org/zap"
)

type LoggerImpl struct {
	Log  *zap.SugaredLogger
	data []interface{}
}

func GetLogger() *LoggerImpl {
	logConfig := helpers.GetZapLoggerSetup()
	loggerInstence := zap.Must(logConfig.Build())
	defer loggerInstence.Sync()
	sugar := loggerInstence.Sugar()

	logger := &LoggerImpl{
		Log: sugar,
	}
	return logger
}

// Add data log args
func (li *LoggerImpl) AddData(key string, value any) {
	li.data = append(li.data, key)
	li.data = append(li.data, value)
}

// Get data log args
func (li *LoggerImpl) Data() []any {
	return li.data
}

// Clear data log args
func (li *LoggerImpl) ClearData() {
	li.data = []any{}
}
