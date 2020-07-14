// Package logger contains custom settings regarding the log functionality.
package logger

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

// Setup logger options.
func Setup(show bool) {
	logrus.SetFormatter(&CustomFormatter{})

	if !show {
		logrus.SetOutput(ioutil.Discard)
	}
}
