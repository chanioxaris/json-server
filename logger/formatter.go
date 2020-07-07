package logger

import (
	"bytes"
	"fmt"

	"github.com/sirupsen/logrus"
)

type CustomFormatter struct {
}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Output buffer.
	b := &bytes.Buffer{}

	// Log method field.
	method, ok := entry.Data["method"]
	if ok {
		fmt.Fprintf(b, "%v", method)
	}

	// Log ulr field.
	url, ok := entry.Data["url"]
	if ok {
		fmt.Fprintf(b, " %v", url)
	}

	// Log status field.
	status, ok := entry.Data["status"]
	if ok {
		switch entry.Level {
		case logrus.InfoLevel:
			fmt.Fprintf(b, "\u001B[38;5;76m %v \u001B[0m", status)
		case logrus.WarnLevel:
			fmt.Fprintf(b, "\u001B[38;5;11m %v \u001B[0m", status)
		case logrus.ErrorLevel:
			fmt.Fprintf(b, "\u001B[38;5;196m %v \u001B[0m", status)
		default:
			fmt.Fprintf(b, "\u001B[30m %v \u001B[0m", status)
		}
	}

	// Log duration field.
	duration, ok := entry.Data["duration"]
	if ok {
		fmt.Fprintf(b, "- %v ", duration)
	}

	// Log size field.
	size, ok := entry.Data["size"]
	if ok {
		fmt.Fprintf(b, "- %v Bytes", size)
	}

	b.WriteByte('\n')

	return b.Bytes(), nil
}
