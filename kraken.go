package gormungandr

import (
	"fmt"
	"time"

	"github.com/CanalTP/gonavitia/pbnavitia"
	"github.com/golang/protobuf/proto"
	zmq "github.com/pebbe/zmq4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/sony/gobreaker"
)

type KrakenTimeout struct {
	message string
}

func NewKrakenTimeout(message string) *KrakenTimeout {
	return &KrakenTimeout{
		message: message,
	}
}

func (e *KrakenTimeout) Error() string {
	return e.message
}

type Kraken struct {
	Name    string
	Addr    string
	Timeout time.Duration

	cb *gobreaker.CircuitBreaker
}

func NewKraken(name, addr string, timeout time.Duration) *Kraken {
	kraken := &Kraken{
		Name:    name,
		Timeout: timeout,
		Addr:    addr,
	}
	var st gobreaker.Settings
	st.Name = "Kraken"
	kraken.cb = gobreaker.NewCircuitBreaker(st)
	return kraken

}

func (k *Kraken) Call(request *pbnavitia.Request) (*pbnavitia.Response, error) {
	rep, err := k.cb.Execute(func() (interface{}, error) {
		requester, err := zmq.NewSocket(zmq.REQ)
		if err != nil {
			return nil, errors.Wrap(err, "error while creating ZMQ socket")
		}
		if err = requester.Connect(k.Addr); err != nil {
			return nil, errors.Wrap(err, "error while connecting")
		}
		defer func() {
			if err = requester.Close(); err != nil {
				logrus.Warnf("error while closing the socket %s", err)
			}
		}()
		data, err := proto.Marshal(request)
		if err != nil {
			return nil, errors.Wrap(err, "error while marshalling")
		}
		if _, err = requester.Send(string(data), 0); err != nil {
			return nil, errors.Wrap(err, "error while sending")
		}
		poller := zmq.NewPoller()
		poller.Add(requester, zmq.POLLIN)
		p, err := poller.Poll(k.Timeout)
		if err != nil {
			return nil, errors.Wrap(err, "error during polling")
		}
		if len(p) < 1 {
			return nil, NewKrakenTimeout(fmt.Sprintf("kraken %s timeout", k.Name))
		}
		rawResp, err := p[0].Socket.Recv(0)
		if err != nil {
			return nil, errors.Wrap(err, "error while receiving response")
		}
		resp := &pbnavitia.Response{}
		if err = proto.Unmarshal([]byte(rawResp), resp); err != nil {
			return nil, errors.Wrap(err, "error while unmarshalling response")
		}

		return resp, nil
	})
	if err != nil {
		return nil, err
	}
	return rep.(*pbnavitia.Response), nil
}
