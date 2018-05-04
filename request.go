package gormungandr

import (
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type Request struct {
	User     User          `form:"-"`
	Coverage string        `form:"-"` //requested coverage
	ID       uuid.UUID     `form:"-"`
	Logger   *logrus.Entry `form:"-"` //logger associated to this request
}

func NewRequest() Request {
	id := uuid.NewV4()
	return Request{
		ID:     id,
		Logger: logrus.WithField("request_id", id),
	}
}
