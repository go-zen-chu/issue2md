package log

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
)

var (
	once   sync.Once
	logger Logger
)

// Init initialize logger just once in application lifecycle
func Init(debug bool) error {
	var err error
	once.Do(func() {
		var l *zap.Logger
		if debug {
			l, err = zap.NewDevelopment()
		} else {
			l, err = zap.NewProduction()
		}
		logger = &zapLogger{
			lgr: l,
		}
	})
	if err != nil {
		return fmt.Errorf("initializing logger: %w", err)
	}
	return nil
}

// Debug logs the vals at Debug level.
func Debug(vals ...interface{}) {
	logger.Debug(vals...)
}

// Debugf logs the formatted message at Debug level.
func Debugf(format string, vals ...interface{}) {
	logger.Debugf(format, vals...)
}

// Info logs the vals at Info level.
func Info(vals ...interface{}) {
	logger.Info(vals...)
}

// Infof logs the formatted message at Info level.
func Infof(format string, vals ...interface{}) {
	logger.Infof(format, vals...)
}

// Warn logs the vals at Warn level.
func Warn(vals ...interface{}) {
	logger.Warn(vals...)
}

// Warnf logs the formatted message at Warn level.
func Warnf(format string, vals ...interface{}) {
	logger.Warnf(format, vals...)
}

// Error logs the vals at Error level.
func Error(vals ...interface{}) {
	logger.Error(vals...)
}

// Errorf logs the formatted message at Error level.
func Errorf(format string, vals ...interface{}) {
	logger.Errorf(format, vals...)
}

// Fatal logs the vals at Fatal level, then calls os.Exit(1).
func Fatal(vals ...interface{}) {
	logger.Fatal(vals...)
}

// Fatalf logs the formatted message at Fatal level, then calls os.Exit(1).
func Fatalf(format string, vals ...interface{}) {
	logger.Fatalf(format, vals...)
}

// Close calls finalizer of logger implementations.
func Close() error {
	return logger.Close()
}
