package kraken

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

type KrakenZMQ struct {
	Name    string
	Addr    string
	Timeout time.Duration

	cb       *gobreaker.CircuitBreaker
	observer krakenObserver
}

func NewKrakenZMQ(name, addr string, timeout time.Duration) Kraken {
	kraken := &KrakenZMQ{
		Name:    name,
		Timeout: timeout,
		Addr:    addr,
	}
	var st gobreaker.Settings
	st.Name = "Kraken"
	kraken.cb = gobreaker.NewCircuitBreaker(st)
	return kraken
}

func (k *KrakenZMQ) call(request []byte) ([]byte, error) {
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
	if _, err = requester.SendBytes(request, 0); err != nil {
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
	rawResp, err := p[0].Socket.RecvBytes(0)
	if err != nil {
		return nil, errors.Wrap(err, "error while receiving response")
	}
	return rawResp, nil
}

func (k *KrakenZMQ) Call(request *pbnavitia.Request) (*pbnavitia.Response, error) {
	data, err := proto.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(err, "error while marshalling")
	}
	rawResponse, err := k.cb.Execute(func() (interface{}, error) {
		o := k.observer.StartRequest(request.RequestedApi.String())
		defer o.Finish()
		return k.call(data)
	})
	if err != nil {
		k.observer.OnError(request.RequestedApi.String(), err)
		return nil, err
	}

	response := &pbnavitia.Response{}
	if err = proto.Unmarshal(rawResponse.([]byte), response); err != nil {
		return nil, errors.Wrap(err, "error while unmarshalling response")
	}
	return response, nil
}
