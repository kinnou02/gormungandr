package checker

import (
	"testing"

	"github.com/CanalTP/gonavitia"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func IsValidRouteSchedule(t *testing.T, schedule *gonavitia.RouteSchedule) {
	assert := assert.New(t)
	require := require.New(t)
	require.NotNil(schedule.DisplayInfo)
	assert.NotEmpty(schedule.DisplayInfo.Direction)
	assert.NotEmpty(schedule.DisplayInfo.Label)
	assert.NotEmpty(schedule.DisplayInfo.Network)
	assert.NotEmpty(schedule.DisplayInfo.Name)

	//TODO: check shape

	require.NotNil(schedule.Table)
	require.NotNil(schedule.Table.Headers)
	require.NotNil(schedule.Table.Rows)

	for _, h := range schedule.Table.Headers {
		IsValidRouteScheduleHeader(t, h)
	}

	for _, r := range schedule.Table.Rows {
		IsValidRouteScheduleRow(t, r)
	}

}

func IsValidRouteScheduleHeader(t *testing.T, header *gonavitia.Header) {
	assert := assert.New(t)
	require := require.New(t)
	assert.NotNil(header.DisplayInfo)
	require.NotNil(header.Links)

	links := make(map[string]string)
	for _, l := range header.Links {
		links[*l.Type] = *l.Id
	}
	assert.NotEmpty(links["vehicle_journey"])
	assert.NotEmpty(links["physical_mode"])
	//TODO: check optional note

}

func IsValidRouteScheduleRow(t *testing.T, row gonavitia.Row) {
	assert := assert.New(t)

	for _, d := range row.DateTimes {
		assert.NotNil(d.AdditionalInfo)
		assert.NotEmpty(d.Links)
		//unmarshalling was a success so the datetime is valid or empty
	}
	IsValidStopPoint(t, row.StopPoint, 1)

}

func IsValidStopPoint(t *testing.T, sp *gonavitia.StopPoint, depth int) {
	assert := assert.New(t)
	assert.NotEmpty(sp.Name)
	assert.NotEmpty(sp.Label)
	require.NotNil(t, sp.Coord)
	IsValidCoord(t, *sp.Coord)

	//TODO: check comments
	for _, m := range sp.PhysicalModes {
		IsValidPhysicalMode(t, m)
	}
	for _, m := range sp.CommercialModes {
		IsValidCommercialMode(t, m)
	}

	if depth > 0 {
		require.NotNil(t, sp.StopArea)
		IsValidStopArea(t, *sp.StopArea, depth-1)
	}

	if depth >= 3 {
		require.NotNil(t, sp.Address)
		IsValidAddress(t, *sp.Address)
	}
}

func IsValidCoord(t *testing.T, c gonavitia.Coord) {
	assert := assert.New(t)
	assert.Truef(c.Lon <= 180.0, "invalid longitude for coord")
	assert.Truef(c.Lon >= -180.0, "invalid longitude for coord")

	assert.Truef(c.Lat <= 90.0, "invalid latitude for coord")
	assert.Truef(c.Lat >= -90.0, "invalid latitude for coord")
}

func IsValidAddress(t *testing.T, a gonavitia.Address) {
	assert := assert.New(t)
	assert.NotEmpty(a.Id)
	assert.NotEmpty(a.HouseNumber)
	assert.NotEmpty(a.Name)
	require.NotNil(t, a.Coord)
	IsValidCoord(t, *a.Coord)
}

func IsValidStopArea(t *testing.T, sa gonavitia.StopArea, depth int) {
	assert := assert.New(t)
	assert.NotEmpty(sa.Name)
	assert.NotEmpty(sa.Label)
	require.NotNil(t, sa.Coord)
	IsValidCoord(t, *sa.Coord)

	//TODO: check comments
	for _, m := range sa.PhysicalModes {
		IsValidPhysicalMode(t, m)
	}
	for _, m := range sa.CommercialModes {
		IsValidCommercialMode(t, m)
	}
}

func IsValidPhysicalMode(t *testing.T, mode gonavitia.PhysicalMode) {
	assert := assert.New(t)
	assert.NotEmpty(mode.Name)
	assert.NotEmpty(mode.Id)
}

func IsValidCommercialMode(t *testing.T, mode gonavitia.CommercialMode) {
	assert := assert.New(t)
	assert.NotEmpty(mode.Name)
	assert.NotEmpty(mode.Id)
}
