package gormungandr

import (
	"time"

	"github.com/CanalTP/gonavitia/pbnavitia"
	"github.com/golang/protobuf/proto"
	zmq "github.com/pebbe/zmq4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/sony/gobreaker"
)

var (
	KrakenTimeout = errors.New("kraken timeout")
)

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
		requester, _ := zmq.NewSocket(zmq.REQ)
		err := requester.Connect(k.Addr)
		if err != nil {
			return nil, errors.Wrap(err, "error while connecting")
		}
		defer func() {
			if err = requester.Close(); err != nil {
				logrus.Warnf("error while closing the socket %s", err)
			}
		}()
		data, _ := proto.Marshal(request)
		_, err = requester.Send(string(data), 0)
		if err != nil {
			return nil, errors.Wrap(err, "error while sending")
		}
		poller := zmq.NewPoller()
		poller.Add(requester, zmq.POLLIN)
		p, err := poller.Poll(k.Timeout)
		if err != nil {
			return nil, errors.Wrap(err, "error during polling")
		}
		if len(p) < 1 {
			return nil, errors.Errorf("kraken %s timeout", k.Name)
		}
		rawResp, _ := p[0].Socket.Recv(0)
		resp := &pbnavitia.Response{}
		_ = proto.Unmarshal([]byte(rawResp), resp)
		return resp, nil
	})
	if err != nil {
		return nil, err
	}
	return rep.(*pbnavitia.Response), nil
}
