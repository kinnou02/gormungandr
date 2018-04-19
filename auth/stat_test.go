package auth

import (
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/CanalTP/gonavitia"
	"github.com/CanalTP/gonavitia/pbnavitia"
	"github.com/CanalTP/gormungandr"
	"github.com/CanalTP/gormungandr/internal/schedules"
	"github.com/rafaeljesus/rabbus"
)

type mockRabbus struct {
	mock.Mock
}

func (r *mockRabbus) EmitAsync() chan<- rabbus.Message {
	args := r.MethodCalled("EmitAsync")
	return args.Get(0).(chan rabbus.Message)
}

func (r *mockRabbus) EmitErr() <-chan error {
	args := r.MethodCalled("EmitErr")
	return args.Get(0).(chan error)
}

func (r *mockRabbus) EmitOk() <-chan struct{} {
	args := r.MethodCalled("EmitOk")
	return args.Get(0).(chan struct{})
}

func init() {
	gin.SetMode(gin.TestMode)
}

func newMockRabbus(sizeAsync, sizeErr, sizeOk int, expectAsync, expectOk, expectErr bool) (mock *mockRabbus, emitAsync chan rabbus.Message, emitErr chan error, emitOK chan struct{}) {
	mock = new(mockRabbus)
	emitAsync = make(chan rabbus.Message, sizeAsync)
	emitErr = make(chan error, sizeErr)
	emitOK = make(chan struct{}, sizeOk)
	if expectOk {
		mock.On("EmitOk").Return(emitOK)
	}
	if expectErr {
		mock.On("EmitErr").Return(emitErr)
	}
	if expectAsync {
		mock.On("EmitAsync").Return(emitAsync)
	}
	return
}

func TestPublishForRouteSchedulesNil(t *testing.T) {
	t.Parallel()
	var (
		statPublisher *StatPublisher
		response      gonavitia.RouteScheduleResponse
		request       schedules.RouteScheduleRequest
	)
	ginContext, _ := gin.CreateTestContext(httptest.NewRecorder())
	ginContext.Request = httptest.NewRequest("Get", "/", nil)
	assert.NotPanics(t, func() { statPublisher.PublishRouteSchedule(request, response, *ginContext) })
}

func TestPublishForRouteSchedulesOk(t *testing.T) {
	t.Parallel()
	var (
		response gonavitia.RouteScheduleResponse
		request  schedules.RouteScheduleRequest
	)
	ginContext, _ := gin.CreateTestContext(httptest.NewRecorder())
	ginContext.Request = httptest.NewRequest("Get", "/", nil)
	mock, emitAsync, emitErr, emitOK := newMockRabbus(1, 1, 1, true, true, true)
	statPublisher := NewStatPublisher(mock, "test", 10*time.Millisecond)

	emitOK <- struct{}{}
	go statPublisher.PublishRouteSchedule(request, response, *ginContext)
	select {
	case m := <-emitAsync:
		assert.NotNil(t, m)
	case <-time.After(10 * time.Millisecond):
		require.Fail(t, "timeout on emitAsync")
	}
	//we wait a little for the goroutine to do it's job
	<-time.After(10 * time.Millisecond)
	assert.Empty(t, emitErr)
	assert.Empty(t, emitOK)
}

func TestPublishForRouteSchedulesErr(t *testing.T) {
	t.Parallel()
	var (
		response gonavitia.RouteScheduleResponse
		request  schedules.RouteScheduleRequest
	)
	ginContext, _ := gin.CreateTestContext(httptest.NewRecorder())
	ginContext.Request = httptest.NewRequest("Get", "/", nil)
	mock, emitAsync, emitErr, emitOK := newMockRabbus(1, 1, 1, true, true, true)
	statPublisher := NewStatPublisher(mock, "test", 10*time.Millisecond)

	emitErr <- fmt.Errorf("error")
	go statPublisher.PublishRouteSchedule(request, response, *ginContext)
	select {
	case m := <-emitAsync:
		assert.NotNil(t, m)
	case <-time.After(10 * time.Millisecond):
		require.Fail(t, "timeout on emitAsync")
	}
	//we wait a little for the goroutine to do it's job
	<-time.After(10 * time.Millisecond)
	assert.Empty(t, emitErr)
	assert.Empty(t, emitOK)
}

func TestPublishOK(t *testing.T) {
	t.Parallel()
	var pb pbnavitia.StatRequest
	mock, emitAsync, emitErr, emitOK := newMockRabbus(1, 1, 1, true, true, true)
	statPublisher := NewStatPublisher(mock, "test", time.Millisecond)

	emitOK <- struct{}{}
	assert.NoError(t, statPublisher.publish(pb))
	select {
	case m := <-emitAsync:
		assert.NotNil(t, m)
	case <-time.After(1 * time.Millisecond):
		assert.Fail(t, "timeout on emitAsync")
	}
	assert.Empty(t, emitErr)
	assert.Empty(t, emitOK)
}

func TestPublishErr(t *testing.T) {
	t.Parallel()
	var pb pbnavitia.StatRequest
	mock, emitAsync, emitErr, emitOK := newMockRabbus(1, 1, 1, true, true, true)
	statPublisher := NewStatPublisher(mock, "test", time.Millisecond)

	emitErr <- fmt.Errorf("error")
	err := statPublisher.publish(pb)
	assert.Error(t, err)
	select {
	case m := <-emitAsync:
		assert.NotNil(t, m)
	case <-time.After(1 * time.Millisecond):
		assert.Fail(t, "timeout on emitAsync")
	}
	assert.Empty(t, emitErr)
	//next message is handled correctly
	emitOK <- struct{}{}
	assert.NoError(t, statPublisher.publish(pb))
	select {
	case m := <-emitAsync:
		assert.NotNil(t, m)
	case <-time.After(1 * time.Millisecond):
		assert.Fail(t, "timeout on emitAsync")
	}
	assert.Empty(t, emitErr)
	assert.Empty(t, emitOK)
}

func TestTimeoutSend(t *testing.T) {
	t.Parallel()
	var pb pbnavitia.StatRequest
	mock, emitAsync, _, _ := newMockRabbus(0, 0, 0, true, false, false)
	statPublisher := NewStatPublisher(mock, "test", time.Millisecond)
	assert.Error(t, statPublisher.publish(pb))
	select {
	case <-emitAsync:
		assert.Fail(t, "message should have been canceled")
	case <-time.After(1 * time.Millisecond):
	}

}

func TestBuildRouteSchedules(t *testing.T) {
	t.Parallel()
	var (
		response gonavitia.RouteScheduleResponse
		request  schedules.RouteScheduleRequest
	)
	ginContext, _ := gin.CreateTestContext(httptest.NewRecorder())
	ginContext.Request = httptest.NewRequest("Get", "/", nil)
	response.Pagination = &gonavitia.Pagination{ItemsOnPage: 42}
	request.User = gormungandr.User{
		Username:     "bob",
		Id:           7,
		AppName:      "Bricolage&co",
		Token:        "key",
		EndPointId:   3,
		EndPointName: "navitia.io",
	}
	request.Coverage = "fr-idf"

	pb := buildStatForRouteSchedule(request, response, *ginContext)
	assert.Equal(t, "bob", pb.GetUserName())
	assert.Equal(t, "Bricolage&co", pb.GetApplicationName())
	assert.Equal(t, int32(7), pb.GetUserId())
	assert.Equal(t, "key", pb.GetToken())
	assert.Equal(t, "navitia.io", pb.GetEndPointName())
	assert.Equal(t, int32(3), pb.GetEndPointId())
	require.Len(t, pb.GetCoverages(), 1)
	assert.Equal(t, "fr-idf", pb.Coverages[0].GetRegionId())
	require.NotNil(t, pb.InfoResponse)
	assert.Equal(t, int32(42), pb.InfoResponse.GetObjectCount())
}

func TestBuildStatInfoResponse(t *testing.T) {
	t.Parallel()
	assert.Nil(t, buildStatInfoResponse(nil))

	pagination := &gonavitia.Pagination{ItemsOnPage: 42}
	pb := buildStatInfoResponse(pagination)

	require.NotNil(t, pb)
	assert.Equal(t, int32(42), pb.GetObjectCount())

	pagination = &gonavitia.Pagination{
		ItemsOnPage:  42,
		ItemsPerPage: 50,
		StartPage:    3,
		TotalResult:  500,
	}
	pb = buildStatInfoResponse(pagination)

	require.NotNil(t, pb)
	assert.Equal(t, int32(42), pb.GetObjectCount())
}
