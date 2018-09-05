package serializer

import (
	"github.com/CanalTP/gonavitia"
	"github.com/CanalTP/gonavitia/pbnavitia"
)

func (s *Serializer) NewError(pb *pbnavitia.Error) *gonavitia.Error {
	if pb == nil {
		return nil
	}
	id := pb.Id.Enum().String()
	return &gonavitia.Error{
		Id:      &id,
		Message: pb.Message,
		Code:    s.NewErrorCode(pb),
	}
}

func (s *Serializer) NewErrorCode(pb *pbnavitia.Error) gonavitia.ErrorCode {
	if pb == nil {
		return gonavitia.ErrorOk
	}
	switch pb.GetId() {
	case pbnavitia.Error_service_unavailable:
		return gonavitia.ErrorServiceUnavailable
	case pbnavitia.Error_internal_error:
		return gonavitia.ErrorInternalError
	case pbnavitia.Error_date_out_of_bounds:
		return gonavitia.ErrorDateOutOfBounds
	case pbnavitia.Error_no_origin:
		return gonavitia.ErrorNoOrigin
	case pbnavitia.Error_no_destination:
		return gonavitia.ErrorNoDestination
	case pbnavitia.Error_no_origin_nor_destination:
		return gonavitia.ErrorNOriginNorDestination
	case pbnavitia.Error_unknown_object:
		return gonavitia.ErrorUnknownObject
	case pbnavitia.Error_unable_to_parse:
		return gonavitia.ErrorUnableToParse
	case pbnavitia.Error_bad_filter:
		return gonavitia.ErrorBadFilter
	case pbnavitia.Error_unknown_api:
		return gonavitia.ErrorUnkownApi
	case pbnavitia.Error_bad_format:
		return gonavitia.ErrorBadFormat
	case pbnavitia.Error_no_solution:
		return gonavitia.ErrorNoSolution
	default:
		return gonavitia.ErrorOk
	}

}
