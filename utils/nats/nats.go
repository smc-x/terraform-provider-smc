// Package nats provides a wrapper for easier use of nats.
package nats

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

var logger = logrus.WithField("n", "nats")

func logIf(err error) {
	if err != nil {
		logger.Error(err)
	}
}

// Wrapper wraps around the use of nats.
type Wrapper struct {
	nc *nats.Conn
}

// New initiates a wrapper.
func New(token, endpoint string, skipVerify bool) (*Wrapper, func(), error) {
	drained := make(chan bool, 1)
	nc, err := nats.Connect(
		fmt.Sprintf("tls://%s@%s", token, endpoint),
		nats.Secure(&tls.Config{
			InsecureSkipVerify: skipVerify,
		}),
		nats.ClosedHandler(func(*nats.Conn) {
			drained <- true
		}),
	)
	if err != nil {
		return nil, func() {}, err
	}
	return &Wrapper{nc}, func() {
		err := nc.Drain()
		if err == nil {
			<-drained
		} else {
			logIf(err)
		}
	}, nil
}

// Remote invokes a remote method specified by the concrete subj.
func (w *Wrapper) Remote(subj string, data []byte, timeout time.Duration) (
	reply []byte, err error,
) {
	msg, err := w.nc.Request(subj, data, timeout)
	if err != nil {
		return nil, err
	}
	return msg.Data, nil
}

// Serve serves a request-based pattern with the given callback.
func (w *Wrapper) Serve(pattern string, cb Callback) (clean func(), err error) {
	sub, err := w.nc.QueueSubscribe(
		pattern, "default", func(msg *nats.Msg) {
			defer recoverCbErr(msg.Subject)
			panicIfCbErr(msg.Respond(cb(msg.Subject, msg.Data)))
		},
	)
	if err != nil {
		return func() {}, err
	}
	return func() {
		logIf(sub.Drain())
	}, nil
}

type Callback func(subj string, data []byte) (reply []byte)

var loggerCb = logrus.WithField("n", "nats-cb")

func recoverCbErr(subj string) {
	if e := recover(); e != nil {
		loggerCb.WithField("subj", subj).Error(e)
	}
}

func panicIfCbErr(err error) {
	if err != nil {
		panic(err)
	}
}
