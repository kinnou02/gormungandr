package auth

import (
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const userContextKey = "gormungandr.auth.User"
const coverageContextKey = "gormungandr.auth.Coverage"

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
	c.Set(userContextKey, user)
	c.Set(coverageContextKey, coverage)
	c.Next()
}

//return the user associated to the context if no user has been set it returns (User{}, false)
func GetUser(c *gin.Context) (user User, ok bool) {
	tmp, ok := c.Get(userContextKey)
	if !ok {
		return User{}, ok
	}
	return tmp.(User), ok
}

//return the coverage associated to the context if no coverage has been set it returns ("", false)
func GetCoverage(c *gin.Context) (coverage string, ok bool) {
	tmp, ok := c.Get(coverageContextKey)
	if !ok {
		return "", ok
	}
	return tmp.(string), ok
}
