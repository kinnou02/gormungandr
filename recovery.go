package gormungandr

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http/httputil"
	"runtime/debug"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				httprequest, _ := httputil.DumpRequest(c.Request, false)
				logrus.WithFields(logrus.Fields{
					"request": string(httprequest),
					"err":     err,
					"stack":   string(debug.Stack()),
				}).Errorf("panic recovered")
				c.AbortWithStatusJSON(500, gin.H{"message": "Internal Server Error"})
			}
		}()
		c.Next()
	}
}
