package auth

import (
	"database/sql"
	"time"

	"github.com/pkg/errors"
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

//return AuthenticationFailed if the the authentication fail
func Authenticate(token string, now time.Time, db *sql.DB) (User, error) {
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
		 "user" u
	WHERE i.name = $1
	AND u.id = $2
	AND u.type = 'with_free_instances' and i.is_free
	UNION ALL
	SELECT true
	FROM "instance" i
	JOIN "authorization" a on a.instance_id=i.id
	JOIN "user" u ON u.id=a.user_id
	WHERE i.name = $1
	AND u.id = $2
`
