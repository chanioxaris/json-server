package middleware

import (
	"bufio"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// Logger is operating as middleware to log http requests info.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rww := newResponseWriterWrapper(w)
		next.ServeHTTP(rww, r)

		duration := time.Since(start)

		logrus.
			WithField("method", r.Method).
			WithField("url", r.URL.Path).
			WithField("status", rww.statusCode).
			WithField("duration", duration).
			WithField("size", rww.size).
			Log(rww.logLevel)
	})
}

// responseWriterWrapper that implements ResponseWriter interface to retrieve response status code and content size.
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
	size       int
	logLevel   logrus.Level
}

func newResponseWriterWrapper(w http.ResponseWriter) *responseWriterWrapper {
	return &responseWriterWrapper{
		ResponseWriter: w,
		statusCode:     200,
		logLevel:       logrus.InfoLevel,
	}
}

func (c *responseWriterWrapper) WriteHeader(statusCode int) {
	if statusCode < 200 {
		c.logLevel = logrus.TraceLevel
	} else if statusCode >= 200 && statusCode < 300 {
		c.logLevel = logrus.InfoLevel
	} else if statusCode >= 300 && statusCode < 400 {
		c.logLevel = logrus.WarnLevel
	} else {
		c.logLevel = logrus.ErrorLevel
	}

	c.statusCode = statusCode
	c.ResponseWriter.WriteHeader(statusCode)
}

func (c *responseWriterWrapper) Write(b []byte) (int, error) {
	size, err := c.ResponseWriter.Write(b)
	if err != nil {
		return 0, err
	}

	c.size += size

	return size, nil
}

func (c *responseWriterWrapper) Flush() {
	if f, ok := c.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (c *responseWriterWrapper) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := c.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("ResponseWriter does not implement the Hijacker interface")
	}

	return hijacker.Hijack()
}
