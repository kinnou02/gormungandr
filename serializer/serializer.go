package serializer

import (
	"time"

	"github.com/CanalTP/gonavitia/pbnavitia"
)

type Serializer struct {
	Location *time.Location
}

func New() *Serializer {
	return &Serializer{
		Location: time.UTC,
	}
}

// initialize the serializer by setting the timezone from the response
//Panic if timezone cannot be determined as this should never happen
func (s *Serializer) Init(pb *pbnavitia.Metadatas) {
	if pb != nil {
		location, err := time.LoadLocation(pb.GetTimezone())
		if err != nil {
			panic(err)
		}
		s.Location = location
	} else {
		panic("no metadatas in response")
	}
}
