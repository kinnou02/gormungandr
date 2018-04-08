package serializer

import (
	"strings"
	"time"

	"github.com/canaltp/gonavitia"
	"github.com/canaltp/gonavitia/pbnavitia"
	"github.com/golang/protobuf/proto"
)

func NewRouteSchedulesResponse(pb *pbnavitia.Response) *gonavitia.RouteScheduleResponse {
	if pb == nil {
		return nil
	}
	response := gonavitia.RouteScheduleResponse{
		Error:          NewError(pb.Error),
		RouteSchedules: make([]*gonavitia.RouteSchedule, 0),
		Pagination:     NewPagination(pb.Pagination),
		FeedPublishers: make([]*gonavitia.FeedPublisher, 0, len(pb.FeedPublishers)),
	}
	for _, r := range pb.RouteSchedules {
		response.RouteSchedules = append(response.RouteSchedules, NewRouteSchedule(r))
	}
	for _, f := range pb.FeedPublishers {
		response.FeedPublishers = append(response.FeedPublishers, NewFeedPublisher(f))
	}
	return &response
}

func NewRouteSchedule(pb *pbnavitia.RouteSchedule) *gonavitia.RouteSchedule {
	if pb == nil {
		return nil
	}
	var additionalInfo *string
	info := pb.ResponseStatus
	if info != nil {
		tmp := strings.ToLower(info.Enum().String())
		additionalInfo = &tmp
	}
	return &gonavitia.RouteSchedule{
		DisplayInfo:    NewPtDisplayInfo(pb.PtDisplayInformations),
		Table:          NewTable(pb.Table),
		AdditionalInfo: additionalInfo,
		GeoJson:        NewGeoJsonMultistring(pb.Geojson),
		Links:          NewLinksFromUris(pb.PtDisplayInformations),
	}
}

func NewTable(pb *pbnavitia.Table) *gonavitia.Table {
	if pb == nil {
		return nil
	}
	t := gonavitia.Table{
		Headers: make([]*gonavitia.Header, 0, len(pb.Headers)),
		Rows:    make([]gonavitia.Row, 0, len(pb.Rows)),
	}
	for _, h := range pb.Headers {
		t.Headers = append(t.Headers, NewHeader(h))
	}
	for _, r := range pb.Rows {
		t.Rows = append(t.Rows, NewRow(r))
	}
	return &t
}

func NewHeader(pb *pbnavitia.Header) *gonavitia.Header {
	if pb == nil {
		return nil
	}
	return &gonavitia.Header{
		DisplayInfo: NewPtDisplayInfo(pb.PtDisplayInformations),
		Links:       NewLinksFromUris(pb.PtDisplayInformations),
		//AdditionalInfo []string       `json:"additional_informations"`
	}
}

func NewRow(pb *pbnavitia.RouteScheduleRow) gonavitia.Row {
	if pb == nil {
		return gonavitia.Row{}
	}
	r := gonavitia.Row{
		StopPoint: NewStopPoint(pb.StopPoint),
		DateTimes: make([]gonavitia.DateTime, 0, len(pb.DateTimes)),
	}
	for _, d := range pb.DateTimes {
		r.DateTimes = append(r.DateTimes, NewDatetime(d))
	}
	return r
}

func NewDatetime(pb *pbnavitia.ScheduleStopTime) gonavitia.DateTime {
	if pb == nil {
		return gonavitia.DateTime{}
	}
	rtLevel := strings.ToLower(pb.GetRealtimeLevel().Enum().String())
	return gonavitia.DateTime{
		DateTime:       gonavitia.NavitiaDatetime(time.Unix(int64(pb.GetDate()+pb.GetTime()), 0)),
		BaseDateTime:   gonavitia.NavitiaDatetime(time.Unix(int64(pb.GetBaseDateTime()), 0)),
		AdditionalInfo: make([]string, 0),
		DataFreshness:  rtLevel,
		Links:          NewLinksFromProperties(pb.Properties),
	}
}

func NewLinksFromProperties(pb *pbnavitia.Properties) []gonavitia.Link {
	result := make([]gonavitia.Link, 0, 1)
	if pb == nil {
		return result
	}
	if pb.VehicleJourneyId != nil {
		result = append(result, gonavitia.Link{
			Id:    pb.VehicleJourneyId,
			Value: pb.VehicleJourneyId,
			Type:  proto.String("vehicle_journey"),
			Rel:   proto.String("vehicle_journeys"),
		})
	}
	return result
}
