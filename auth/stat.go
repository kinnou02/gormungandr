package auth

import (
	"time"

	"github.com/CanalTP/gonavitia"
	"github.com/CanalTP/gonavitia/pbnavitia"
	"github.com/CanalTP/gormungandr/internal/schedules"
	"github.com/golang/protobuf/proto"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rafaeljesus/rabbus"
)

var (
	statErrorCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "gormungandr",
		Subsystem: "stats",
		Name:      "errors_count",
		Help:      "stats request errors count",
	},
	)
)

func init() {
	prometheus.MustRegister(statErrorCounter)
}

type sender interface {
	EmitAsync() chan<- rabbus.Message
	EmitErr() <-chan error
	EmitOk() <-chan struct{}
}

type StatPublisher struct {
	Rabbus      sender
	Exchange    string
	SendTimeout time.Duration
}

func NewStatPublisher(rabbus sender, exchange string, sendTimeout time.Duration) *StatPublisher {
	return &StatPublisher{
		Rabbus:      rabbus,
		Exchange:    exchange,
		SendTimeout: sendTimeout,
	}
}

func (s *StatPublisher) publish(stat pbnavitia.StatRequest) error {
	data, err := proto.Marshal(&stat)
	if err != nil {
		return errors.Wrap(err, "error while marshaling stats")
	}
	msg := rabbus.Message{
		Exchange:        s.Exchange,
		Kind:            "topic",
		Key:             stat.GetApi(),
		Payload:         data,
		ContentType:     "application/data",
		ContentEncoding: "binary",
	}

	//TODO if we are disconncted from rabbit we loose all messages even if rabbit is back
	select {
	case s.Rabbus.EmitAsync() <- msg:
		select {
		case <-s.Rabbus.EmitOk():
			return nil
		case err := <-s.Rabbus.EmitErr():
			return errors.Wrap(err, "failed to sent message")
		}
	case <-time.After(s.SendTimeout):
		return errors.Errorf("timeout while sending message")
	}
}

func buildStatForRouteSchedule(request schedules.RouteScheduleRequest, response gonavitia.RouteScheduleResponse, c echo.Context) pbnavitia.StatRequest {
	return pbnavitia.StatRequest{
		RequestDate:     proto.Uint64(uint64(time.Now().Unix())),
		Api:             proto.String("v1.route_schedules"),
		UserId:          proto.Int32(int32(request.User.Id)),
		UserName:        &request.User.Username,
		ApplicationName: &request.User.AppName,
		Token:           &request.User.Token, //we should'nt do this...
		EndPointId:      proto.Int32(int32(request.User.EndPointId)),
		EndPointName:    &request.User.EndPointName,
		ApplicationId:   proto.Int32(-1), //same as jormun
		RequestDuration: proto.Int32(0),  //TODO
		Path:            proto.String(c.Request().URL.Path),
		Host:            proto.String(c.Request().URL.Host),
		Client:          proto.String(c.RealIP()),
		ResponseSize:    proto.Int32(0), //always been wrong in jormun...
		InfoResponse:    buildStatInfoResponse(response.Pagination),
		Coverages:       []*pbnavitia.StatCoverage{{RegionId: &request.Coverage}},
	}
}

func buildStatInfoResponse(pagination *gonavitia.Pagination) *pbnavitia.StatInfoResponse {
	if pagination == nil {
		return nil
	}
	return &pbnavitia.StatInfoResponse{
		ObjectCount: proto.Int32(pagination.ItemsOnPage),
	}
}

func (s *StatPublisher) PublishRouteSchedule(request schedules.RouteScheduleRequest,
	response gonavitia.RouteScheduleResponse, c echo.Context) error {

	if s == nil { //if nil is passed by interface we don't want to panic
		return nil
	}

	pb := buildStatForRouteSchedule(request, response, c)
	err := s.publish(pb)
	if err != nil {
		statErrorCounter.Inc()
	}

	return err
}
