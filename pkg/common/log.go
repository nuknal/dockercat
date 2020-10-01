package common

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

// NewLogger new logrus logger
func NewLogger() *logrus.Entry {
	var log *logrus.Logger
	if os.Getenv("DEBUG") == "TRUE" {
		log = newDevelopmentLogger()
	} else {
		log = newProductionLogger()
	}

	log.Formatter = &logrus.JSONFormatter{}

	return log.WithFields(logrus.Fields{
		"version": "0.1.0",
	})
}

func newDevelopmentLogger() *logrus.Logger {
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	file, err := os.OpenFile("development.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("unable to log to file")
		os.Exit(1)
	}
	log.SetOutput(file)
	return log
}

func newProductionLogger() *logrus.Logger {
	log := logrus.New()
	log.Out = ioutil.Discard
	log.SetLevel(logrus.ErrorLevel)
	return log
}
