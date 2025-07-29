package log

import (
	"log"
	"os"
)

func Info(format string, args ...any) {
	log.Printf("[INFO] "+format, args...)
}

func Warn(format string, args ...any) {
	log.Printf("[WARN] "+format, args...)
}

func Debug(format string, args ...any) {
	log.Printf("[DEBUG] "+format, args...)
}

func Error(format string, args ...any) {
	log.Printf("[ERROR] "+format, args...)
	os.Exit(1)
}
