package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Entry

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	Log = logrus.WithFields(
		logrus.Fields{
			"AppName": "JS_RUNNER",
		})
	logrus.SetOutput(os.Stderr)
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.DebugLevel)
}
