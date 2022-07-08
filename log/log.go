package log

// And just go global.
var defaultLogger Logger

// init creates the default   This can be changed
func init() {
	defaultLogger, _ = NewZapLogger(false)
}

// SetLogger sets the default logger to be used by this package
func SetLogger(logger Logger) {
	defaultLogger = logger
}

// Error logs an error message with the parameters
func Error(args ...interface{}) {
	defaultLogger.Error(args...)
}

// Errorf logs a formatted error message
func Errorf(format string, args ...interface{}) {
	defaultLogger.Errorf(format, args...)
}

// Infof logs a formatted info message
func Infof(format string, args ...interface{}) {
	defaultLogger.Infof(format, args...)
}

// Info logs an info message with the parameters
func Info(args ...interface{}) {
	defaultLogger.Info(args...)
}

// Debugf logs a formatted debug message
func Debugf(format string, args ...interface{}) {
	defaultLogger.Debugf(format, args...)
}

// Debug logs an debug message with the parameters
func Debug(args ...interface{}) {
	defaultLogger.Debug(args...)
}
