package main

import (
	"github.com/michaelquigley/pfxlog"
	"github.com/openziti/ziti/ziti-probe/subcmd"
	"github.com/openziti/foundation/util/debugz"
	"github.com/sirupsen/logrus"
)

func init() {
	pfxlog.Global(logrus.InfoLevel)
	pfxlog.SetPrefix("github.com/netfoundry/")
}

func main() {
	debugz.AddStackDumpHandler()
	subcmd.Execute()
}
