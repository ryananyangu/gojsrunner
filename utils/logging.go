package utils

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
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
	logrus.SetLevel(logrus.ErrorLevel)
}

//Log to file
func LogrusLogger() gin.HandlerFunc {

	// instantiation
	logger := logrus.New()

	//Set output
	logger.Out = os.Stderr

	//Set log level
	logger.SetLevel(logrus.ErrorLevel)

	//Format log
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	// logger.SetFormatter(&logrus.TextFormatter{})

	logger.SetReportCaller(true)

	return func(c *gin.Context) {
		//Start time
		startTime := time.Now()

		//Process request
		c.Next()

		//End time
		endTime := time.Now()

		//Execution time
		latencyTime := endTime.Sub(startTime)

		//Request method
		reqMethod := c.Request.Method

		//Request routing
		reqUri := c.Request.RequestURI

		// status code
		statusCode := c.Writer.Status()

		// request IP
		clientIP := c.ClientIP()

		//Log format
		logger.WithFields(logrus.Fields{
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"client_ip":    clientIP,
			"req_method":   reqMethod,
			"req_uri":      reqUri,
		}).Info()

	}
}
