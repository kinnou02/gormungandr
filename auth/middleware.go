package auth

import (
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func getToken(c *gin.Context) string {
	username, _, ok := c.Request.BasicAuth()
	if ok {
		return username
	} else if authorization := c.GetHeader("Authorization"); len(authorization) > 0 {
		return authorization
	} else {
		return c.Query("key")
	}
}

func AuthenticationMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		middleware(c, db)
	}
}

func middleware(c *gin.Context, db *sql.DB) {
	coverage := c.Param("coverage")
	token := getToken(c)
	logger := logrus.WithFields(logrus.Fields{
		"coverage": coverage,
		"token":    token,
	})
	logger.Debug("authentifying request")
	if token == "" {
		c.Header("WWW-Authenticate", "basic realm=\"Token Required\"")
		c.AbortWithStatusJSON(401, gin.H{"message": "no token"})
		return
	}
	user, err := Authenticate(token, time.Now(), db)
	if err == AuthenticationFailed {
		c.Header("WWW-Authenticate", "basic realm=\"Token Required\"")
		c.AbortWithStatusJSON(401, gin.H{"message": "authentication failed"})
		return
	} else if err != nil {
		panic(err)
	}
	ok, err := IsAuthorized(user, coverage, db)
	if err != nil {
		panic(err)
	}
	if !ok {
		c.Header("WWW-Authenticate", "basic realm=\"Token Required\"")
		c.AbortWithStatusJSON(403, gin.H{"message": "authentication failed"})
		return
	}
	c.Next()
}
