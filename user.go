package gormungandr

import (
	"github.com/gin-gonic/gin"
)

const userContextKey = "gormungandr.User"
const coverageContextKey = "gormungandr.Coverage"

type User struct {
	Id       int
	Username string
	AppName  string
	Type     string
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

//store the user in the context
func SetUser(c *gin.Context, user User) {
	c.Set(userContextKey, user)
}

//store the coverage in the context
func SetCoverage(c *gin.Context, coverage string) {
	c.Set(coverageContextKey, coverage)
}
