package log

import "fmt"

type DefaultLogger struct {
}

func (DefaultLogger) Info(msg string) {
	fmt.Printf("[INFO] %s\n", msg)
}

func (DefaultLogger) Debug(msg string) {
	fmt.Printf("[DEBUG] %s\n", msg)
}

func (DefaultLogger) Error(msg string) {
	fmt.Printf("[ERROR] %s\n", msg)
}

func (DefaultLogger) Warn(msg string) {
	fmt.Printf("[WARNNING] %s\n", msg)
}
