package goshopify

import "fmt"

// Logger is an interface the caller should implement when wanting to override
// the default logging.
type Logger interface {
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
}

// defaultLogger is a very naive logger that just prints to standard output.
type defaultLogger struct{}

func (l defaultLogger) Info(format string, args ...interface{}) {
	fmt.Printf("[INFO] "+format+"\n", args...)
}

func (l defaultLogger) Warn(format string, args ...interface{}) {
	fmt.Printf("[WARN] "+format+"\n", args...)
}
