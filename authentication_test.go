package gormungandr

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
