package serializer

import (
	"strings"
	"time"

	"github.com/CanalTP/gonavitia"
	"github.com/CanalTP/gonavitia/pbnavitia"
)

func (s *Serializer) NewPagination(pb *pbnavitia.Pagination) *gonavitia.Pagination {
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

func (s *Serializer) NewCode(pb *pbnavitia.Code) *gonavitia.Code {
	if pb == nil {
		return nil
	}
	return &gonavitia.Code{
		Type:  pb.Type,
		Value: pb.Value,
	}
}

func (s *Serializer) NewPlace(pb *pbnavitia.PtObject) *gonavitia.Place {
	if pb == nil {
		return nil
	}
	t := strings.ToLower(pb.EmbeddedType.String())
	place := gonavitia.Place{
		Id:           pb.Uri,
		Name:         pb.Name,
		EmbeddedType: &t,
		Quality:      pb.Quality,
		StopPoint:    s.NewStopPoint(pb.StopPoint),
		StopArea:     s.NewStopArea(pb.StopArea),
		Admin:        s.NewAdmin(pb.AdministrativeRegion),
		Address:      s.NewAddress(pb.Address),
	}
	return &place
}

func (s *Serializer) NewAdmin(pb *pbnavitia.AdministrativeRegion) *gonavitia.Admin {
	if pb == nil {
		return nil
	}
	admin := gonavitia.Admin{
		Id:      pb.Uri,
		Name:    pb.Name,
		Label:   pb.Label,
		Coord:   s.NewCoord(pb.Coord),
		Insee:   pb.Insee,
		ZipCode: pb.ZipCode,
		Level:   pb.GetLevel(),
	}
	return &admin
}

func (s *Serializer) NewCoord(pb *pbnavitia.GeographicalCoord) *gonavitia.Coord {
	if pb == nil {
		//this is what jormun does...
		return &gonavitia.Coord{
			Lat: 0,
			Lon: 0,
		}
	}
	coord := gonavitia.Coord{
		Lat: pb.GetLat(),
		Lon: pb.GetLon(),
	}
	return &coord
}

func (s *Serializer) NewContext(request *pbnavitia.Request, pb *pbnavitia.Response) *gonavitia.Context {
	if pb == nil || request == nil || pb.Metadatas == nil {
		return nil
	}
	return &gonavitia.Context{
		CurrentDatetime: s.NewNavitiaDatetime(int64(request.GetXCurrentDatetime())),
		Timezone:        pb.Metadatas.GetTimezone(),
	}
}

func (s *Serializer) NewStopPoint(pb *pbnavitia.StopPoint) *gonavitia.StopPoint {
	if pb == nil {
		return nil
	}
	sp := gonavitia.StopPoint{
		Id:              pb.Uri,
		Name:            pb.Name,
		Label:           pb.Label,
		Coord:           s.NewCoord(pb.Coord),
		Admins:          make([]*gonavitia.Admin, 0, len(pb.AdministrativeRegions)),
		StopArea:        s.NewStopArea(pb.StopArea),
		Codes:           make([]*gonavitia.Code, 0, len(pb.Codes)),
		Equipments:      s.NewEquipments(pb.HasEquipments),
		Links:           make([]*gonavitia.Link, 0),
		PhysicalModes:   s.NewPhysicalModes(pb.PhysicalModes),
		CommercialModes: s.NewCommercialModes(pb.CommercialModes),
		Address:         s.NewAddress(pb.Address),
	}
	for _, pb_admin := range pb.AdministrativeRegions {
		sp.Admins = append(sp.Admins, s.NewAdmin(pb_admin))
	}
	for _, code := range pb.Codes {
		sp.Codes = append(sp.Codes, s.NewCode(code))
	}
	return &sp
}

func (s *Serializer) NewStopArea(pb *pbnavitia.StopArea) *gonavitia.StopArea {
	if pb == nil {
		return nil
	}
	sa := gonavitia.StopArea{
		Id:              pb.Uri,
		Name:            pb.Name,
		Label:           pb.Label,
		Timezone:        pb.Timezone,
		Coord:           s.NewCoord(pb.Coord),
		Admins:          make([]*gonavitia.Admin, 0, len(pb.AdministrativeRegions)),
		Codes:           make([]*gonavitia.Code, 0, len(pb.Codes)),
		Links:           make([]*gonavitia.Link, 0),
		PhysicalModes:   s.NewPhysicalModes(pb.PhysicalModes),
		CommercialModes: s.NewCommercialModes(pb.CommercialModes),
	}
	for _, pb_admin := range pb.AdministrativeRegions {
		sa.Admins = append(sa.Admins, s.NewAdmin(pb_admin))
	}
	for _, code := range pb.Codes {
		sa.Codes = append(sa.Codes, s.NewCode(code))
	}
	for _, sp := range pb.StopPoints {
		sa.StopPoints = append(sa.StopPoints, s.NewStopPoint(sp))
	}
	return &sa
}

func (s *Serializer) NewAddress(pb *pbnavitia.Address) *gonavitia.Address {
	if pb == nil {
		return nil
	}
	address := gonavitia.Address{
		Id:          pb.Uri,
		Name:        pb.Name,
		Label:       pb.Label,
		Coord:       s.NewCoord(pb.Coord),
		HouseNumber: pb.HouseNumber,
		Admins:      make([]*gonavitia.Admin, 0),
	}
	for _, pb_admin := range pb.AdministrativeRegions {
		address.Admins = append(address.Admins, s.NewAdmin(pb_admin))
	}
	return &address
}

func (s *Serializer) NewFeedPublisher(pb *pbnavitia.FeedPublisher) *gonavitia.FeedPublisher {
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

func (s *Serializer) NewNavitiaDatetime(timestamp int64) gonavitia.NavitiaDatetime {
	return gonavitia.NavitiaDatetime(
		time.Unix(timestamp, 0).
			In(s.Location))
}
