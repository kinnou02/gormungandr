package auth

import (
	"database/sql"
	"testing"
	"time"

	"github.com/CanalTP/gormungandr"
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

func expectIsAuthorizeSuccess(mock sqlmock.Sqlmock) sqlmock.Sqlmock {
	rows := sqlmock.NewRows([]string{"bool"}).
		AddRow(true)
	mock.ExpectQuery("SELECT true").WillReturnRows(rows)
	return mock
}

func expectIsAuthorizeNoResult(mock sqlmock.Sqlmock) sqlmock.Sqlmock {
	rows := sqlmock.NewRows([]string{"bool"})
	mock.ExpectQuery("SELECT true").WillReturnRows(rows)
	return mock
}

func expectIsAuthorizeError(mock sqlmock.Sqlmock) sqlmock.Sqlmock {
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
	user, err := Authenticate("mytoken", time.Now(), db)
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
	_, err := Authenticate("mytoken", time.Now(), db)
	assert.NotNil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Equal(t, AuthenticationFailed, err)
}

func TestAuthenticateError(t *testing.T) {
	t.Parallel()
	db, mock := newMock()
	defer db.Close()
	mock = expectAuthError(mock)
	_, err := Authenticate("mytoken", time.Now(), db)
	assert.NotNil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestIsAuthorized(t *testing.T) {
	t.Parallel()
	db, mock := newMock()
	defer db.Close()
	mock = expectIsAuthorizeSuccess(mock)

	user := gormungandr.User{
		Id:   42,
		Type: "with_free_instances",
	}

	result, err := IsAuthorized(user, "fr-idf", db)
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

	result, err := IsAuthorized(user, "fr-idf", db)
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Equal(t, true, result)
}

func TestIsAuthorizedFailed(t *testing.T) {
	t.Parallel()
	db, mock := newMock()
	defer db.Close()
	mock = expectIsAuthorizeNoResult(mock)

	user := gormungandr.User{
		Id:   42,
		Type: "with_free_instances",
	}

	result, err := IsAuthorized(user, "fr-idf", db)
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Equal(t, false, result)
}

func TestIsAuthorizedError(t *testing.T) {
	t.Parallel()
	db, mock := newMock()
	defer db.Close()
	mock = expectIsAuthorizeError(mock)

	user := gormungandr.User{
		Id:   42,
		Type: "with_free_instances",
	}

	_, err := IsAuthorized(user, "fr-idf", db)
	assert.NotNil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}
