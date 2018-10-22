package conf

import (
	"flag"
)

const (
	defaultDiscoveryIntervalSec = 600
	DefaultConfPath    = "/etc/knativeturbo/turbo.config"
)

type KnativeTurboArgs struct {
	DiscoveryIntervalSec *int
	TurboConf string
	KubeConf string
}

func NewKnativeTurboArgs(fs *flag.FlagSet) *KnativeTurboArgs {
	p := &KnativeTurboArgs{}

	p.DiscoveryIntervalSec = fs.Int("discovery-interval-sec", defaultDiscoveryIntervalSec, "The discovery interval in seconds")
	fs.StringVar(&p.TurboConf,  "turboconfig", p.TurboConf,  "Path to the turbo config file.")
	fs.StringVar(&p.KubeConf, "kubeconfig", p.KubeConf, "Path to kubeconfig file with authorization and master location information.")

	return p
}
