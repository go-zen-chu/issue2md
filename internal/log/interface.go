package log

type Logger interface {
	// Debug logs the vals at Debug level.
	Debug(vals ...interface{})
	// Debugf logs the formatted message at Debug level.
	Debugf(format string, vals ...interface{})
	// Info logs the vals at Info level.
	Info(vals ...interface{})
	// Infof logs the formatted message at Info level.
	Infof(format string, vals ...interface{})
	// Warn logs the vals at Warn level.
	Warn(vals ...interface{})
	// Warnf logs the formatted message at Warn level.
	Warnf(format string, vals ...interface{})
	// Error logs the vals at Error level.
	Error(vals ...interface{})
	// Errorf logs the formatted message at Error level.
	Errorf(format string, vals ...interface{})
	// Fatal logs the vals at Fatal level, then calls os.Exit(1).
	Fatal(vals ...interface{})
	// Fatalf logs the formatted message at Fatal level, then calls os.Exit(1).
	Fatalf(format string, vals ...interface{})
	// Close calls finalizer of logger implementations.
	Close() error
}
