package auth

import (
	"database/sql"
	"testing"
	"time"

	"github.com/CanalTP/gormungandr"
	cache "github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func expectAuthSuccess(mock sqlmock.Sqlmock) sqlmock.Sqlmock {
	rows := sqlmock.NewRows([]string{"id", "login", "app_name", "type", "end_point_id", "end_point_name", "token"}).
		AddRow(42, "mylogin", "myapp", "with_free_instances", 1, "navio", "key")
	mock.ExpectQuery("SELECT u.id").WillReturnRows(rows)
	return mock
}

func expectAuthNoResult(mock sqlmock.Sqlmock) sqlmock.Sqlmock {
	rows := sqlmock.NewRows([]string{"id", "login", "app_name", "type", "end_point_id", "end_point_name", "token"})
	mock.ExpectQuery("SELECT u.id").WillReturnRows(rows)
	return mock
}

func expectAuthError(mock sqlmock.Sqlmock) sqlmock.Sqlmock {
	mock.ExpectQuery("SELECT u.id").WillReturnError(sql.ErrTxDone)
	return mock
}

func expectIsAuthorizedSuccess(mock sqlmock.Sqlmock) sqlmock.Sqlmock {
	rows := sqlmock.NewRows([]string{"bool"}).
		AddRow(true)
	mock.ExpectQuery("SELECT true").WillReturnRows(rows)
	return mock
}

func expectIsAuthorizedNoResult(mock sqlmock.Sqlmock) sqlmock.Sqlmock {
	rows := sqlmock.NewRows([]string{"bool"})
	mock.ExpectQuery("SELECT true").WillReturnRows(rows)
	return mock
}

func expectIsAuthorizedError(mock sqlmock.Sqlmock) sqlmock.Sqlmock {
	mock.ExpectQuery("SELECT true").WillReturnError(sql.ErrTxDone)
	return mock
}

func newMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	return db, mock
}

func TestAuthenticate(t *testing.T) {
	t.Parallel()
	db, mock := newMock()
	defer db.Close()
	mock = expectAuthSuccess(mock)
	user, err := authenticate("mytoken", time.Now(), db)
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Equal(t, "mylogin", user.Username)
	assert.Equal(t, 42, user.Id)
	assert.Equal(t, "myapp", user.AppName)
	assert.Equal(t, "with_free_instances", user.Type)
	assert.Equal(t, 1, user.EndPointId)
	assert.Equal(t, "navio", user.EndPointName)
	assert.Equal(t, "key", user.Token)
}

func TestAuthenticateFail(t *testing.T) {
	t.Parallel()
	db, mock := newMock()
	defer db.Close()
	mock = expectAuthNoResult(mock)
	_, err := authenticate("mytoken", time.Now(), db)
	assert.NotNil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Equal(t, ErrAuthenticationFailed, err)
}

func TestAuthenticateError(t *testing.T) {
	t.Parallel()
	db, mock := newMock()
	defer db.Close()
	mock = expectAuthError(mock)
	_, err := authenticate("mytoken", time.Now(), db)
	assert.NotNil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestIsAuthorized(t *testing.T) {
	t.Parallel()
	db, mock := newMock()
	defer db.Close()
	mock = expectIsAuthorizedSuccess(mock)

	user := gormungandr.User{
		Id:   42,
		Type: "with_free_instances",
	}

	result, err := isAuthorized(user, "fr-idf", db)
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Equal(t, true, result)
}

func TestIsAuthorizedSuperuser(t *testing.T) {
	t.Parallel()
	db, mock := newMock()
	defer db.Close()

	user := gormungandr.User{
		Id:   42,
		Type: "super_user",
	}

	result, err := isAuthorized(user, "fr-idf", db)
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Equal(t, true, result)
}

func TestIsAuthorizedFailed(t *testing.T) {
	t.Parallel()
	db, mock := newMock()
	defer db.Close()
	mock = expectIsAuthorizedNoResult(mock)

	user := gormungandr.User{
		Id:   42,
		Type: "with_free_instances",
	}

	result, err := isAuthorized(user, "fr-idf", db)
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Equal(t, false, result)
}

func TestIsAuthorizedError(t *testing.T) {
	t.Parallel()
	db, mock := newMock()
	defer db.Close()
	mock = expectIsAuthorizedError(mock)

	user := gormungandr.User{
		Id:   42,
		Type: "with_free_instances",
	}

	_, err := isAuthorized(user, "fr-idf", db)
	assert.NotNil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestCachedAuthenticateSuccess(t *testing.T) {
	t.Parallel()
	db, mock := newMock()
	defer db.Close()
	mock = expectAuthSuccess(mock)
	authCache := cache.New(300*time.Second, 600*time.Second)

	user, err := CachedAuthenticate("mytoken", time.Now(), db, authCache)
	assert.Nil(t, err)
	assert.Equal(t, "mylogin", user.Username) // not testing all here
	user, err = CachedAuthenticate("mytoken", time.Now(), db, authCache)
	assert.Nil(t, err)
	assert.Equal(t, "mylogin", user.Username)
	assert.Equal(t, 42, user.Id)
	assert.Equal(t, "myapp", user.AppName)
	assert.Equal(t, "with_free_instances", user.Type)
	assert.Equal(t, 1, user.EndPointId)
	assert.Equal(t, "navio", user.EndPointName)
	assert.Equal(t, "key", user.Token)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestCachedAuthenticateFail(t *testing.T) {
	t.Parallel()
	db, mock := newMock()
	defer db.Close()
	mock = expectAuthNoResult(mock)
	authCache := cache.New(300*time.Second, 600*time.Second)

	_, err := CachedAuthenticate("mytoken", time.Now(), db, authCache)
	assert.NotNil(t, err)
	assert.Equal(t, ErrAuthenticationFailed, err)
	_, err = CachedAuthenticate("mytoken", time.Now(), db, authCache)
	assert.NotNil(t, err)
	assert.Equal(t, ErrAuthenticationFailed, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestCachedAuthenticateFailOutdatedThenSuccess(t *testing.T) {
	t.Parallel()
	db, mock := newMock()
	defer db.Close()
	mock = expectAuthNoResult(mock)
	mock = expectAuthSuccess(mock)
	authCache := cache.New(1*time.Microsecond, 1*time.Microsecond)

	_, err := CachedAuthenticate("mytoken", time.Now(), db, authCache)
	assert.NotNil(t, err)
	assert.Equal(t, ErrAuthenticationFailed, err)
	time.Sleep(2 * time.Microsecond)
	user, err := CachedAuthenticate("mytoken", time.Now(), db, authCache)
	assert.Nil(t, err)
	assert.Equal(t, "mylogin", user.Username)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestCachedAuthenticateSuccessNocache(t *testing.T) {
	t.Parallel()
	db, mock := newMock()
	defer db.Close()
	mock = expectAuthSuccess(mock)
	mock = expectAuthSuccess(mock)

	user, err := CachedAuthenticate("mytoken", time.Now(), db, nil)
	assert.Nil(t, err)
	assert.Equal(t, "mylogin", user.Username)
	user, err = CachedAuthenticate("mytoken", time.Now(), db, nil)
	assert.Nil(t, err)
	assert.Equal(t, "mylogin", user.Username)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestCachedAuthenticateErrorThenSuccess(t *testing.T) {
	t.Parallel()
	db, mock := newMock()
	defer db.Close()
	mock = expectAuthError(mock)
	mock = expectAuthSuccess(mock)
	authCache := cache.New(300*time.Second, 600*time.Second)

	_, err := CachedAuthenticate("mytoken", time.Now(), db, authCache)
	assert.NotNil(t, err)
	user, err := CachedAuthenticate("mytoken", time.Now(), db, authCache)
	assert.Nil(t, err)
	assert.Equal(t, "mylogin", user.Username)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestCachedIsAuthorizedSuccess(t *testing.T) {
	t.Parallel()
	db, mock := newMock()
	defer db.Close()
	mock = expectIsAuthorizedSuccess(mock)
	authCache := cache.New(300*time.Second, 600*time.Second)

	user := gormungandr.User{
		Id:   42,
		Type: "with_free_instances",
	}

	result, err := CachedIsAuthorized(user, "fr-idf", db, authCache)
	assert.Nil(t, err)
	assert.Equal(t, true, result)
	result, err = CachedIsAuthorized(user, "fr-idf", db, authCache)
	assert.Nil(t, err)
	assert.Equal(t, true, result)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestCachedIsAuthorizedFailed(t *testing.T) {
	t.Parallel()
	db, mock := newMock()
	defer db.Close()
	mock = expectIsAuthorizedNoResult(mock)
	authCache := cache.New(300*time.Second, 600*time.Second)

	user := gormungandr.User{
		Id:   42,
		Type: "with_free_instances",
	}

	result, err := CachedIsAuthorized(user, "fr-idf", db, authCache)
	assert.Nil(t, err)
	assert.Equal(t, false, result)
	result, err = CachedIsAuthorized(user, "fr-idf", db, authCache)
	assert.Nil(t, err)
	assert.Equal(t, false, result)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestCachedIsAuthorizedFailOutdatedThenSuccess(t *testing.T) {
	t.Parallel()
	db, mock := newMock()
	defer db.Close()
	mock = expectIsAuthorizedNoResult(mock)
	mock = expectIsAuthorizedSuccess(mock)
	authCache := cache.New(1*time.Microsecond, 1*time.Microsecond)

	user := gormungandr.User{
		Id:   42,
		Type: "with_free_instances",
	}

	result, err := CachedIsAuthorized(user, "fr-idf", db, authCache)
	assert.Nil(t, err)
	assert.Equal(t, false, result)
	time.Sleep(2 * time.Microsecond)
	result, err = CachedIsAuthorized(user, "fr-idf", db, authCache)
	assert.Nil(t, err)
	assert.Equal(t, true, result)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestCachedIsAuthorizedSuccessNoCache(t *testing.T) {
	t.Parallel()
	db, mock := newMock()
	defer db.Close()
	mock = expectIsAuthorizedSuccess(mock)
	mock = expectIsAuthorizedSuccess(mock)

	user := gormungandr.User{
		Id:   42,
		Type: "with_free_instances",
	}

	result, err := CachedIsAuthorized(user, "fr-idf", db, nil)
	assert.Nil(t, err)
	assert.Equal(t, true, result)
	result, err = CachedIsAuthorized(user, "fr-idf", db, nil)
	assert.Nil(t, err)
	assert.Equal(t, true, result)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestCachedIsAuthorizedErrorThenSuccess(t *testing.T) {
	t.Parallel()
	db, mock := newMock()
	defer db.Close()
	mock = expectIsAuthorizedError(mock)
	mock = expectIsAuthorizedSuccess(mock)
	authCache := cache.New(300*time.Second, 600*time.Second)

	user := gormungandr.User{
		Id:   42,
		Type: "with_free_instances",
	}

	_, err := CachedIsAuthorized(user, "fr-idf", db, authCache)
	assert.NotNil(t, err)
	result, err := CachedIsAuthorized(user, "fr-idf", db, authCache)
	assert.Nil(t, err)
	assert.Equal(t, true, result)
	assert.Nil(t, mock.ExpectationsWereMet())
}
