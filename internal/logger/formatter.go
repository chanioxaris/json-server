package logger

import (
	"fmt"

	"github.com/gookit/color"
	"github.com/sirupsen/logrus"
)

// CustomFormatter for logrus logger.
type CustomFormatter struct {
}

// Format renders a single custom log entry.
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Log method field.
	if method, ok := entry.Data["method"]; ok {
		fmt.Printf("%v", method)
	}

	// Log url field.
	if url, ok := entry.Data["url"]; ok {
		fmt.Printf(" %v", url)
	}

	// Log status field.
	if status, ok := entry.Data["status"]; ok {
		switch entry.Level {
		case logrus.InfoLevel:
			color.Green.Printf(" %v", status)
		case logrus.WarnLevel:
			color.Yellow.Printf(" %v", status)
		case logrus.ErrorLevel:
			color.Red.Printf(" %v", status)
		default:
			color.White.Printf(" %v", status)
		}
	}

	// Log duration field.
	if duration, ok := entry.Data["duration"]; ok {
		fmt.Printf(" - %v", duration)
	}

	// Log size field.
	if size, ok := entry.Data["size"]; ok {
		fmt.Printf(" - %v Bytes", size)
	}

	fmt.Println()

	return nil, nil
}
