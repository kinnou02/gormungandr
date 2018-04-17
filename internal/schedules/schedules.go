package schedules

import (
	"net/http"
	"strings"
	"time"

	"github.com/CanalTP/gonavitia/pbnavitia"
	"github.com/CanalTP/gormungandr"
	"github.com/CanalTP/gormungandr/serializer"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type RouteScheduleRequest struct {
	FromDatetime     time.Time `form:"from_datetime" time_format:"20060102T150405" binding:"required"`
	DisableGeojson   bool      `form:"disable_geojson"`
	StartPage        int32     `form:"start_page"`
	Count            int32     `form:"count"`
	Duration         int32     `form:"duration"`
	ForbiddenUris    []string  //mapping with Binding doesn't work
	Depth            int32     `form:"depth"`
	CurrentDatetime  time.Time `form:"_current_datetime"`
	ItemsPerSchedule int32     `form:"items_per_schedule"`
	DataFreshness    string    `form:"data_freshness"`
	Filters          []string
	User             gormungandr.User
	Coverage         string //requested coverage
	ID               uuid.UUID
}

func NewRouteScheduleRequest() RouteScheduleRequest {
	return RouteScheduleRequest{
		StartPage:        0,
		Count:            10,
		Duration:         86400,
		CurrentDatetime:  time.Now(),
		Depth:            2,
		ItemsPerSchedule: 10000,
		DataFreshness:    "base_schedudle",
	}
}

func RouteSchedule(c *gin.Context, kraken *gormungandr.Kraken, request *RouteScheduleRequest, publisher Publisher, logger *logrus.Entry) {
	pbReq := BuildRequestRouteSchedule(*request)
	resp, err := kraken.Call(pbReq)
	logger.Debug("calling kraken")
	if err != nil {
		logger.Errorf("Error while calling kraken: %+v\n", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err})
		return
	}
	logger.Debug("building response")
	r := serializer.NewRouteSchedulesResponse(resp)
	c.JSON(http.StatusOK, r)
	logger.Debug("handling stats")

	go func() {
		err = publisher.PublishRouteSchedule(*request, *r, *c.Copy())
		if err != nil {
			logger.Errorf("stat not sent %+v", err)
		} else {
			logger.Debug("stat sent")
		}
	}()
}

func BuildRequestRouteSchedule(req RouteScheduleRequest) *pbnavitia.Request {
	departureFilter := strings.Join(req.Filters, "and ")
	//TODO handle Realtime level from request
	pbReq := &pbnavitia.Request{
		RequestedApi: pbnavitia.API_ROUTE_SCHEDULES.Enum(),
		NextStopTimes: &pbnavitia.NextStopTimeRequest{
			DepartureFilter:  proto.String(departureFilter),
			ArrivalFilter:    proto.String(""),
			FromDatetime:     proto.Uint64(uint64(req.FromDatetime.Unix())),
			Duration:         proto.Int32(req.Duration),
			Depth:            proto.Int32(req.Depth),
			NbStoptimes:      proto.Int32(req.Count),
			Count:            proto.Int32(req.Count),
			StartPage:        proto.Int32(req.StartPage),
			DisableGeojson:   proto.Bool(req.DisableGeojson),
			ItemsPerSchedule: proto.Int32(req.ItemsPerSchedule),
			RealtimeLevel:    pbnavitia.RTLevel_BASE_SCHEDULE.Enum(),
		},
		XCurrentDatetime: proto.Uint64(uint64(req.CurrentDatetime.Unix())),
	}
	pbReq.NextStopTimes.ForbiddenUri = append(pbReq.NextStopTimes.ForbiddenUri, req.ForbiddenUris...)

	return pbReq
}
