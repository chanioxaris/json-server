package logger

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

// CustomFormatter for logrus logger.
type CustomFormatter struct {
}

// Format renders a single custom log entry.
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Output string.
	var sb strings.Builder

	// Log method field.
	if method, ok := entry.Data["method"]; ok {
		sb.WriteString(fmt.Sprintf("%v", method))
	}

	// Log ulr field.
	if url, ok := entry.Data["url"]; ok {
		sb.WriteString(fmt.Sprintf(" %v", url))
	}

	// Log status field.
	if status, ok := entry.Data["status"]; ok {
		var textColor string

		switch entry.Level {
		case logrus.InfoLevel:
			textColor = "38;5;76m"
		case logrus.WarnLevel:
			textColor = "38;5;11m"
		case logrus.ErrorLevel:
			textColor = "38;5;196m"
		default:
			textColor = "30m"
		}

		sb.WriteString(fmt.Sprintf("\u001B[%s %v \u001B[0m", textColor, status))
	}

	// Log duration field.
	if duration, ok := entry.Data["duration"]; ok {
		sb.WriteString(fmt.Sprintf("- %v ", duration))
	}

	// Log size field.
	if size, ok := entry.Data["size"]; ok {
		sb.WriteString(fmt.Sprintf("- %v Bytes", size))
	}

	sb.WriteString("\n")

	return []byte(sb.String()), nil
}
