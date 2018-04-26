package auth

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/CanalTP/gormungandr"
	cache "github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
)

var (
	ErrAuthenticationFailed = errors.New("Authentication failed")
)

type authResult struct {
	user gormungandr.User
	err  error
}

func getAuthKey(token string) string {
	return fmt.Sprintf("auth.CachedAuthenticate#%v", token)
}

// return AuthenticationFailed if the authentication fail
// triggers cache only if cache structure is provided
func CachedAuthenticate(token string, now time.Time, db *sql.DB, authCache *cache.Cache) (user gormungandr.User, err error) {
	if authCache == nil {
		return authenticate(token, now, db)
	}

	var k = getAuthKey(token)
	authRes, found := authCache.Get(k)
	if found {
		return authRes.(*authResult).user, authRes.(*authResult).err
	}

	user, err = authenticate(token, now, db)

	if err == nil || err == ErrAuthenticationFailed {
		authCache.SetDefault(k, &authResult{user, err})
	}
	return user, err
}

func authenticate(token string, now time.Time, db *sql.DB) (user gormungandr.User, err error) {
	row := db.QueryRow(authenticationQuery, token, now)
	err = row.Scan(&user.Id, &user.Username, &user.AppName, &user.Type, &user.EndPointId, &user.EndPointName, &user.Token)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, ErrAuthenticationFailed
		}
		return user, errors.Wrap(err, "error while authentication")
	}
	return user, nil
}

type isAuthorizedResult struct {
	isAuthorized bool
	err          error
}

func getIsAuthorizedKey(coverage string, userId int) string {
	return fmt.Sprintf("auth.CachedIsAuthorized#%v#%v", coverage, userId)
}

// triggers cache only if cache structure is provided
func CachedIsAuthorized(user gormungandr.User, coverage string, db *sql.DB, authCache *cache.Cache) (result bool, err error) {
	if authCache == nil {
		return isAuthorized(user, coverage, db)
	}

	var k = getIsAuthorizedKey(coverage, user.Id)
	isAuthorizedRes, found := authCache.Get(k)
	if found {
		return isAuthorizedRes.(*isAuthorizedResult).isAuthorized, isAuthorizedRes.(*isAuthorizedResult).err
	}

	result, err = isAuthorized(user, coverage, db)

	if err == nil {
		authCache.SetDefault(k, &isAuthorizedResult{result, err})
	}
	return result, err
}

func isAuthorized(user gormungandr.User, coverage string, db *sql.DB) (result bool, err error) {
	if user.Type == "super_user" {
		return true, nil
	}
	row := db.QueryRow(authorizationQuery, coverage, user.Id)
	err = row.Scan(&result)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, errors.Wrap(err, "error while IsAuthorized")
	}
	return result, nil
}

const authenticationQuery = `
	SELECT
		u.id,
		u.login,
		coalesce(k.app_name, ''),
		u.type,
		e.id,
		e.name,
		k.token
	FROM "user" u
	JOIN "key" k on u.id = k.user_id
	JOIN "end_point" e on u.end_point_id=e.id
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
