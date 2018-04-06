package gormungandr

import (
	"database/sql"
	"time"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type User struct {
	Id       int
	Username string
	AppName  string
	Type     string
}

var (
	AuthenticationFailed = errors.New("Authentication failed")
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
}

func Authenticate(token string, now time.Time, db *sql.DB) (User, error) {
	//return AuthenticationFailed if the the authentication fail
	var user User
	row := db.QueryRow(authenticationQuery, token, now)
	err := row.Scan(&user.Id, &user.Username, &user.AppName, &user.Type)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, AuthenticationFailed
		} else {
			return user, errors.Wrap(err, "error while authentication")
		}
	}
	return user, nil
}

func IsAuthorized(user User, coverage string, db *sql.DB) (bool, error) {
	var result bool
	if user.Type == "super_user" {
		return true, nil
	}
	row := db.QueryRow(authorizationQuery, coverage, user.Id)
	err := row.Scan(&result)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		} else {
			return false, errors.Wrap(err, "error while IsAuthorized")
		}
	}
	return result, nil
}

const authenticationQuery = `
	SELECT
		u.id,
		u.login,
		k.app_name,
		u.type
	FROM "user" u
	JOIN "key" k on u.id = k.user_id
	WHERE k.token = $1
	AND (k.valid_until > $2 or k.valid_until is null)
`

const authorizationQuery = `
	SELECT true
	FROM "instance" i,
		 "authorization" a,
		 "user" u
	WHERE i.name = $1
	AND u.id = $2
	AND (
		(u.type = 'with_free_instances' and i.is_free)
		OR (i.id=a.instance_id and u.id=a.user_id)
	)
`
