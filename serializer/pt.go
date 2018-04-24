package serializer

import (
	"strings"

	"github.com/CanalTP/gonavitia"
	"github.com/CanalTP/gonavitia/pbnavitia"
	"github.com/golang/protobuf/proto"
)

func NewPtDisplayInfoForRoute(pb *pbnavitia.PtDisplayInfo) *gonavitia.PtDisplayInfo {
	if pb == nil {
		return nil
	}
	var label *string
	if code := pb.GetCode(); len(code) > 0 {
		label = &code
	} else {
		label = pb.Name
	}
	info := gonavitia.PtDisplayInfo{
		Direction:      pb.GetDirection(),
		Code:           pb.GetCode(),
		Network:        pb.GetNetwork(),
		Color:          pb.GetColor(),
		Name:           pb.GetName(),
		TextColor:      pb.GetTextColor(),
		CommercialMode: pb.GetCommercialMode(),
		Label:          label,
		Links:          make([]gonavitia.Link, 0),
	}
	return &info

}

func NewPtDisplayInfoForVJ(pb *pbnavitia.PtDisplayInfo) *gonavitia.PtDisplayInfo {
	if pb == nil {
		return nil
	}
	info := NewPtDisplayInfoForRoute(pb)
	info.Description = proto.String(pb.GetDescription())
	info.PhysicalMode = proto.String(pb.GetPhysicalMode())
	info.Headsign = proto.String(pb.GetHeadsign())
	if len(pb.Headsigns) > 0 {
		info.Headsigns = make([]string, 0, len(pb.Headsigns))
		info.Headsigns = append(info.Headsigns, pb.Headsigns...)
	}
	info.Equipments = NewEquipments(pb.HasEquipments)
	return info
}

func NewAdditionalInformations(pb []pbnavitia.SectionAdditionalInformationType) []string {
	infos := make([]string, 0, len(pb))
	for _, v := range pb {
		additionalInfo := strings.ToLower(v.Enum().String())
		infos = append(infos, additionalInfo)
	}
	return infos
}

func NewEquipments(pb *pbnavitia.HasEquipments) []string {
	if pb == nil {
		return make([]string, 0)
	}
	equipments := make([]string, 0, len(pb.HasEquipments))
	for _, v := range pb.HasEquipments {
		e := strings.ToLower(v.Enum().String())
		equipments = append(equipments, e)
	}
	return equipments
}

func NewPhysicalModes(pb []*pbnavitia.PhysicalMode) []gonavitia.PhysicalMode {
	slice := make([]gonavitia.PhysicalMode, 0, len(pb))
	for _, v := range pb {
		if v != nil {
			slice = append(slice, NewPhysicalMode(*v))
		}
	}
	return slice
}

func NewPhysicalMode(pb pbnavitia.PhysicalMode) gonavitia.PhysicalMode {
	return gonavitia.PhysicalMode{
		Id:   pb.GetUri(),
		Name: pb.GetName(),
	}
}

func NewCommercialModes(pb []*pbnavitia.CommercialMode) []gonavitia.CommercialMode {
	slice := make([]gonavitia.CommercialMode, 0, len(pb))
	for _, v := range pb {
		if v != nil {
			slice = append(slice, NewCommercialMode(*v))
		}
	}
	return slice
}

func NewCommercialMode(pb pbnavitia.CommercialMode) gonavitia.CommercialMode {
	return gonavitia.CommercialMode{
		Id:   pb.GetUri(),
		Name: pb.GetName(),
	}
}
