package serializer

import (
	"strings"
	"time"

	"github.com/canaltp/gonavitia"
	"github.com/canaltp/gonavitia/pbnavitia"
)

func NewDisruption(pb *pbnavitia.Impact) *gonavitia.Disruption {
	if pb == nil {
		return nil
	}
	status := strings.ToLower(pb.Status.String())
	d := gonavitia.Disruption{
		ID:            pb.Uri,
		Uri:           pb.Uri,
		DisruptionUri: pb.DisruptionUri,
		ImpactID:      pb.Uri,
		DisruptionID:  pb.Uri,
		Cause:         pb.Cause,
		Contributor:   pb.Contributor,
		Category:      pb.Category,
		UpdatedAt:     gonavitia.NavitiaDatetime(time.Unix(int64(pb.GetUpdatedAt()), 0)),
		Status:        &status,
		Severity:      NewSeverity(pb.Severity),
	}
	for _, message := range pb.Messages {
		d.Messages = append(d.Messages, NewMessage(message))
	}
	for _, period := range pb.ApplicationPeriods {
		d.ApplicationPeriods = append(d.ApplicationPeriods, NewPeriod(period))
	}
	//TODO add properties
	return &d
}

func NewPeriod(pb *pbnavitia.Period) *gonavitia.Period {
	if pb == nil {
		return nil
	}
	p := gonavitia.Period{
		Begin: gonavitia.NavitiaDatetime(time.Unix(int64(pb.GetBegin()), 0)),
		End:   gonavitia.NavitiaDatetime(time.Unix(int64(pb.GetEnd()), 0)),
	}
	return &p
}

func NewSeverity(pb *pbnavitia.Severity) *gonavitia.Severity {
	if pb == nil {
		return nil
	}
	s := gonavitia.Severity{
		Name:     pb.Name,
		Priority: pb.GetPriority(),
		Color:    pb.Color,
		Effect:   pb.Effect,
	}
	return &s
}

func NewMessage(pb *pbnavitia.MessageContent) *gonavitia.Message {
	if pb == nil {
		return nil
	}
	message := gonavitia.Message{
		Text:    pb.Text,
		Channel: NewChannel(pb.Channel),
	}
	return &message
}

func NewChannel(pb *pbnavitia.Channel) *gonavitia.Channel {
	if pb == nil {
		return nil
	}
	channel := gonavitia.Channel{
		ID:          pb.Id,
		Name:        pb.Name,
		ContentType: pb.ContentType,
	}
	for _, types := range pb.ChannelTypes {
		channel.Types = append(channel.Types, types.String())
	}
	return &channel
}
