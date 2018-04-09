package auth

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestGetTokenBasicAuth(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("Get", "/", nil)

	assert.Equal(t, "", getToken(c))

	c.Request.SetBasicAuth("mykey", "")
	assert.Equal(t, "mykey", getToken(c))

	c.Request.SetBasicAuth("mykey", "unpassword")
	assert.Equal(t, "mykey", getToken(c))

	c.Request.SetBasicAuth("mykeyé$€", "")
	assert.Equal(t, "mykeyé$€", getToken(c))

	c.Request.SetBasicAuth("115aa17b-63d3-4a31-acd6-edebebd4d415", "")
	assert.Equal(t, "115aa17b-63d3-4a31-acd6-edebebd4d415", getToken(c))
}

func TestGetTokenHeader(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("Get", "/", nil)

	assert.Equal(t, "", getToken(c))

	c.Request.Header.Set("Authorization", "mykey")
	assert.Equal(t, "mykey", getToken(c))

	c.Request.Header.Set("Authorization", "mykeyé$€")
	assert.Equal(t, "mykeyé$€", getToken(c))

	c.Request.Header.Set("Authorization", "115aa17b-63d3-4a31-acd6-edebebd4d415")
	assert.Equal(t, "115aa17b-63d3-4a31-acd6-edebebd4d415", getToken(c))
}

func TestGetTokenParams(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	c.Request = httptest.NewRequest("Get", "/?key=mykey", nil)
	assert.Equal(t, "mykey", getToken(c))

	c.Request = httptest.NewRequest("Get", "/?key=mykeyé$€", nil)
	assert.Equal(t, "mykeyé$€", getToken(c))
}

func TestMiddlewareNoToken(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("Get", "/coverage/fr-idf", nil)
	db, _ := newMock()
	defer db.Close()
	middleware(c, db)
	assert.True(t, c.IsAborted())
}

func TestMiddlewareAuthFail(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("Get", "/coverage/fr-idf", nil)
	c.Request.SetBasicAuth("mykey", "")
	db, mock := newMock()
	defer db.Close()
	mock = expectAuthNoResult(mock)
	middleware(c, db)
	assert.True(t, c.IsAborted())
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestMiddlewareNotAuthorized(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("Get", "/coverage/fr-idf", nil)
	c.Request.SetBasicAuth("mykey", "")
	db, mock := newMock()
	defer db.Close()
	mock = expectAuthSuccess(mock)
	mock = expectIsAuthorizeNoResult(mock)
	middleware(c, db)
	assert.True(t, c.IsAborted())
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestMiddlewareAuthorized(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("Get", "/coverage/fr-idf", nil)
	c.Request.SetBasicAuth("mykey", "")
	db, mock := newMock()
	defer db.Close()
	mock = expectAuthSuccess(mock)
	mock = expectIsAuthorizeSuccess(mock)
	middleware(c, db)
	assert.False(t, c.IsAborted())
	assert.Nil(t, mock.ExpectationsWereMet())
}
