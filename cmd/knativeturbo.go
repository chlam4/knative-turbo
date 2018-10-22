package main

import (
	"flag"
	"github.com/golang/glog"
	"github.com/pallavidn/knative-turbo/pkg"
	"github.com/pallavidn/knative-turbo/pkg/conf"
)

func main() {
	// The default is to log to both of stderr and file
	// These arguments can be overloaded from the command-line args
	flag.Set("logtostderr", "false")
	flag.Set("alsologtostderr", "true")
	//flag.Set("log_dir", "/var/log")
	defer glog.Flush()

	args := conf.NewKnativeTurboArgs(flag.CommandLine)
	flag.Parse()

	glog.Info("Starting KnativeTurbo...")
	s, err := pkg.NewKnativeTAPService(args)

	if err != nil {
		glog.Fatal("Failed creating KnativeTurbo: %v", err)
	}

	s.Start()
}
