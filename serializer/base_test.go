package serializer

import (
	"testing"
	"time"

	"github.com/CanalTP/gonavitia"
	"github.com/CanalTP/gonavitia/pbnavitia"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestNewPlaceNil(t *testing.T) {
	assert.Nil(t, New().NewPlace(nil))
}

func TestNewPlace(t *testing.T) {
	pb := pbnavitia.PtObject{
		Uri:          proto.String("foo"),
		Name:         proto.String("bar"),
		EmbeddedType: pbnavitia.NavitiaType_STOP_AREA.Enum(),
	}
	place := New().NewPlace(&pb)
	assert.Equal(t, *place.Id, "foo")
	assert.Equal(t, *place.Name, "bar")
	assert.Equal(t, *place.EmbeddedType, "stop_area")
}

func TestNewPaginationNil(t *testing.T) {
	assert.Nil(t, New().NewPagination(nil))
}

func TestNewPagination(t *testing.T) {
	pb := pbnavitia.Pagination{
		ItemsOnPage:  proto.Int32(1),
		ItemsPerPage: proto.Int32(2),
		StartPage:    proto.Int32(3),
		TotalResult:  proto.Int32(4),
	}
	pagination := New().NewPagination(&pb)
	assert.Equal(t, pagination.ItemsOnPage, int32(1))
	assert.Equal(t, pagination.ItemsPerPage, int32(2))
	assert.Equal(t, pagination.StartPage, int32(3))
	assert.Equal(t, pagination.TotalResult, int32(4))
}

func TestNewFeedPublisherNil(t *testing.T) {
	assert.Nil(t, New().NewFeedPublisher(nil))
}

func TestNewFeedPublisher(t *testing.T) {
	pb := pbnavitia.FeedPublisher{
		Id:      proto.String("id"),
		Name:    proto.String("name"),
		Url:     proto.String("url"),
		License: proto.String("license"),
	}
	fp := New().NewFeedPublisher(&pb)
	assert.NotNil(t, fp)
	assert.Equal(t, "id", *fp.Id)
	assert.Equal(t, "name", *fp.Name)
	assert.Equal(t, "url", *fp.Url)
	assert.Equal(t, "license", *fp.License)
}

func TestNewNavitiaDatetime(t *testing.T) {
	serializer := New()

	assert.Equal(t, gonavitia.NavitiaDatetime(time.Unix(1525348246, 0).In(time.UTC)), serializer.NewNavitiaDatetime(1525348246))
	location, err := time.LoadLocation("Europe/Paris")
	assert.NoError(t, err)
	serializer.Location = location
	assert.Equal(t, gonavitia.NavitiaDatetime(time.Unix(1525348246, 0).In(location)), serializer.NewNavitiaDatetime(1525348246))
}
