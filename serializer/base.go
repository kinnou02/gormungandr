package serializer

import (
	"strings"
	"time"

	"github.com/CanalTP/gonavitia"
	"github.com/CanalTP/gonavitia/pbnavitia"
)

func NewError(pb *pbnavitia.Error) *gonavitia.Error {
	if pb == nil {
		return nil
	}
	id := pb.Id.Enum().String()
	return &gonavitia.Error{
		Id:      &id,
		Message: pb.Message,
	}
}

func NewPagination(pb *pbnavitia.Pagination) *gonavitia.Pagination {
	if pb == nil {
		return nil
	}
	return &gonavitia.Pagination{
		StartPage:    pb.GetStartPage(),
		ItemsOnPage:  pb.GetItemsOnPage(),
		ItemsPerPage: pb.GetItemsPerPage(),
		TotalResult:  pb.GetTotalResult(),
	}
}

func NewCode(pb *pbnavitia.Code) *gonavitia.Code {
	if pb == nil {
		return nil
	}
	return &gonavitia.Code{
		Type:  pb.Type,
		Value: pb.Value,
	}
}

func NewPlace(pb *pbnavitia.PtObject) *gonavitia.Place {
	if pb == nil {
		return nil
	}
	t := strings.ToLower(pb.EmbeddedType.String())
	place := gonavitia.Place{
		Id:           pb.Uri,
		Name:         pb.Name,
		EmbeddedType: &t,
		Quality:      pb.Quality,
		StopPoint:    NewStopPoint(pb.StopPoint),
		StopArea:     NewStopArea(pb.StopArea),
		Admin:        NewAdmin(pb.AdministrativeRegion),
		Address:      NewAddress(pb.Address),
	}
	return &place
}

func NewAdmin(pb *pbnavitia.AdministrativeRegion) *gonavitia.Admin {
	if pb == nil {
		return nil
	}
	admin := gonavitia.Admin{
		Id:      pb.Uri,
		Name:    pb.Name,
		Label:   pb.Label,
		Coord:   NewCoord(pb.Coord),
		Insee:   pb.Insee,
		ZipCode: pb.ZipCode,
		Level:   pb.GetLevel(),
	}
	return &admin
}

func NewCoord(pb *pbnavitia.GeographicalCoord) *gonavitia.Coord {
	if pb == nil {
		return nil
	}
	coord := gonavitia.Coord{
		Lat: pb.GetLat(),
		Lon: pb.GetLon(),
	}
	return &coord
}

func NewContext(request *pbnavitia.Request, pb *pbnavitia.Response) *gonavitia.Context {
	if pb == nil || request == nil || pb.Metadatas == nil {
		return nil
	}
	return &gonavitia.Context{
		CurrentDatetime: gonavitia.NavitiaDatetime(time.Unix(int64(request.GetXCurrentDatetime()), 0)),
		Timezone:        pb.Metadatas.GetTimezone(),
	}
}

func NewStopPoint(pb *pbnavitia.StopPoint) *gonavitia.StopPoint {
	if pb == nil {
		return nil
	}
	sp := gonavitia.StopPoint{
		Id:              pb.Uri,
		Name:            pb.Name,
		Label:           pb.Label,
		Coord:           NewCoord(pb.Coord),
		Admins:          make([]*gonavitia.Admin, 0, len(pb.AdministrativeRegions)),
		StopArea:        NewStopArea(pb.StopArea),
		Codes:           make([]*gonavitia.Code, 0, len(pb.Codes)),
		Equipments:      NewEquipments(pb.HasEquipments),
		Links:           make([]*gonavitia.Link, 0),
		PhysicalModes:   NewPhysicalModes(pb.PhysicalModes),
		CommercialModes: NewCommercialModes(pb.CommercialModes),
		Address:         NewAddress(pb.Address),
	}
	for _, pb_admin := range pb.AdministrativeRegions {
		sp.Admins = append(sp.Admins, NewAdmin(pb_admin))
	}
	for _, code := range pb.Codes {
		sp.Codes = append(sp.Codes, NewCode(code))
	}
	return &sp
}

func NewStopArea(pb *pbnavitia.StopArea) *gonavitia.StopArea {
	if pb == nil {
		return nil
	}
	sa := gonavitia.StopArea{
		Id:              pb.Uri,
		Name:            pb.Name,
		Label:           pb.Label,
		Timezone:        pb.Timezone,
		Coord:           NewCoord(pb.Coord),
		Admins:          make([]*gonavitia.Admin, 0, len(pb.AdministrativeRegions)),
		Codes:           make([]*gonavitia.Code, 0, len(pb.Codes)),
		Links:           make([]*gonavitia.Link, 0),
		PhysicalModes:   NewPhysicalModes(pb.PhysicalModes),
		CommercialModes: NewCommercialModes(pb.CommercialModes),
	}
	for _, pb_admin := range pb.AdministrativeRegions {
		sa.Admins = append(sa.Admins, NewAdmin(pb_admin))
	}
	for _, code := range pb.Codes {
		sa.Codes = append(sa.Codes, NewCode(code))
	}
	for _, sp := range pb.StopPoints {
		sa.StopPoints = append(sa.StopPoints, NewStopPoint(sp))
	}
	return &sa
}

func NewAddress(pb *pbnavitia.Address) *gonavitia.Address {
	if pb == nil {
		return nil
	}
	address := gonavitia.Address{
		Id:          pb.Uri,
		Name:        pb.Name,
		Label:       pb.Label,
		Coord:       NewCoord(pb.Coord),
		HouseNumber: pb.HouseNumber,
		Admins:      make([]*gonavitia.Admin, 0),
	}
	for _, pb_admin := range pb.AdministrativeRegions {
		address.Admins = append(address.Admins, NewAdmin(pb_admin))
	}
	return &address
}

func NewFeedPublisher(pb *pbnavitia.FeedPublisher) *gonavitia.FeedPublisher {
	if pb == nil {
		return nil
	}
	return &gonavitia.FeedPublisher{
		Name:    pb.Name,
		Url:     pb.Url,
		Id:      pb.Id,
		License: pb.License,
	}
}
