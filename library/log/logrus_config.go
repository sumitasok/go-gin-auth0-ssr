package log

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"time"
)

// https://programmer.help/blogs/gin-framework-logging-with-logrus.html
func WithLogrus() gin.HandlerFunc {
	logger := log.New()
	logger.WithField("app", "availability-service")
	logger.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	return func(c *gin.Context) {
		startTime := time.Now()               // start time
		c.Next()                              // Processing request
		endTime := time.Now()                 // End time
		latencyTime := endTime.Sub(startTime) // execution time
		reqMethod := c.Request.Method         // Request mode
		reqUri := c.Request.RequestURI        // Request routing
		statusCode := c.Writer.Status()       // Status code
		clientIP := c.ClientIP()              // Request IP

		logger.Infof("| %3d | %13v | %15s | %s | %s |",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		) // Log format
	}
}
