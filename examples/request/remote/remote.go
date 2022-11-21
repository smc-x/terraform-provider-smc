package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/smc-x/terraform-provider-smc/utils/nats"
)

var (
	token    string
	endpoint string
)

func init() {
	pflag.StringVar(&token, "token", "", "")
	pflag.StringVar(&endpoint, "endpoint", "", "")
	pflag.Parse()
}

func main() {
	mid, cleanMid, err := nats.New(token, endpoint, false)
	if err != nil {
		panic(err)
	}
	defer cleanMid()

	cleanW, err := mid.Serve("workers.*", func(subj string, data []byte) (reply []byte) {
		logrus.Infof("%s: %s", subj, data)
		return []byte("{\"worker\":\"worker.default\"}")
	})
	if err != nil {
		panic(err)
	}
	defer cleanW()

	cleanD, err := mid.Serve("worker.default.*", func(subj string, data []byte) (reply []byte) {
		logrus.Infof("%s: %s", subj, data)
		return []byte(fmt.Sprintf("{\"msg\":\"Hello, %s\"}", subj[22:]))
	})
	if err != nil {
		panic(err)
	}
	defer cleanD()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT)
	<-interrupt
}
