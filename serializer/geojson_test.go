package serializer

import (
	"testing"

	"github.com/CanalTP/gonavitia/pbnavitia"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestNewGeoJsonNil(t *testing.T) {
	assert.Nil(t, New().NewGeoJson(nil))
}

func TestNewGeoJson(t *testing.T) {
	pb_section := pbnavitia.Section{
		Shape: []*pbnavitia.GeographicalCoord{
			{Lat: proto.Float64(1), Lon: proto.Float64(2)},
		},
		Length: proto.Int32(19),
	}
	geo := New().NewGeoJson(&pb_section)
	assert.Equal(t, "LineString", geo.Type)
	assert.Equal(t, 1, len(geo.Properties))
	assert.Equal(t, int32(19), geo.Properties[0]["length"])
	assert.Equal(t, len(pb_section.Shape), len(geo.Coordinates))
}
