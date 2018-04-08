package serializer

import "github.com/canaltp/gonavitia/pbnavitia"
import "testing"
import "github.com/stretchr/testify/assert"
import "github.com/golang/protobuf/proto"

func TestNewPlaceNil(t *testing.T) {
	assert.Nil(t, NewPlace(nil))
}

func TestNewPlace(t *testing.T) {
	pb := pbnavitia.PtObject{
		Uri:          proto.String("foo"),
		Name:         proto.String("bar"),
		EmbeddedType: pbnavitia.NavitiaType_STOP_AREA.Enum(),
	}
	place := NewPlace(&pb)
	assert.Equal(t, *place.Id, "foo")
	assert.Equal(t, *place.Name, "bar")
	assert.Equal(t, *place.EmbeddedType, "stop_area")
}

func TestNewPaginationNil(t *testing.T) {
	assert.Nil(t, NewPagination(nil))
}

func TestNewPagination(t *testing.T) {
	pb := pbnavitia.Pagination{
		ItemsOnPage:  proto.Int32(1),
		ItemsPerPage: proto.Int32(2),
		StartPage:    proto.Int32(3),
		TotalResult:  proto.Int32(4),
	}
	pagination := NewPagination(&pb)
	assert.Equal(t, pagination.ItemsOnPage, int32(1))
	assert.Equal(t, pagination.ItemsPerPage, int32(2))
	assert.Equal(t, pagination.StartPage, int32(3))
	assert.Equal(t, pagination.TotalResult, int32(4))
}
