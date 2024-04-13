package logger

var defaultLogger *Logger

func Init(level Level, context string) {
	defaultLogger = &Logger{
		level:   level,
		context: context,
	}
}

func getDefaultLogger() *Logger {
	if defaultLogger == nil {
		Init(InfoLevel, "App")
	}

	return defaultLogger
}

func Info(message string, args ...interface{}) {
	getDefaultLogger().Info(message, args...)
}

func Warn(message string, args ...interface{}) {
	getDefaultLogger().Warn(message, args...)
}

func Error(message string, err error, args ...interface{}) {
	getDefaultLogger().Error(message, err, args...)
}

func Fatal(message string, err error, args ...interface{}) {
	getDefaultLogger().Fatal(message, err, args...)
}

func Debug(message string, args ...interface{}) {
	getDefaultLogger().Debug(message, args...)
}
