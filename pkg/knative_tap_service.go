package pkg

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/pallavidn/knative-turbo/pkg/conf"
	"github.com/pallavidn/knative-turbo/pkg/discovery"
	"github.com/pallavidn/knative-turbo/pkg/registration"
	"github.com/turbonomic/turbo-go-sdk/pkg/probe"
	"github.com/turbonomic/turbo-go-sdk/pkg/service"
	"hash/fnv"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"os/signal"
	"syscall"
)

type disconnectFromTurboFunc func()

type KnativeTAPService struct {
	tapService *service.TAPService
}

func NewKnativeTAPService(args *conf.KnativeTurboArgs) (*KnativeTAPService, error) {
	tapService, err := createTAPService(args)

	if err != nil {
		glog.Errorf("Error while building turbo TAP service on target %v", err)
		return nil, err
	}

	return &KnativeTAPService{tapService}, nil
}

func (p *KnativeTAPService) Start() {
	glog.V(0).Infof("Starting knative TAP service...")

	// Disconnect from Turbo server when kongturbo is shutdown
	handleExit(func() { p.tapService.DisconnectFromTurbo() })

	// Connect to the Turbo server
	p.tapService.ConnectToTurbo()

	select {}
}

func createTAPService(args *conf.KnativeTurboArgs) (*service.TAPService, error) {
	confPath := args.TurboConf

	conf, err := conf.NewKnativeTurboServiceSpec(confPath)
	if err != nil {
		glog.Errorf("Error while parsing the service config file %s: %v", confPath, err)
		os.Exit(1)
	}

	glog.V(3).Infof("Read service configuration from %s: %++v", confPath, conf)

	communicator := conf.TurboCommunicationConfig
	targetConf := conf.KnativeTurboTargetConf

	//if conf.KnativeTurboTargetConf == nil {
	//	conf.KnativeTurboTargetConf = &conf.KnativeTurboTargetConf{}
	//}

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", args.KubeConf)
	if err != nil {
		glog.Errorf("Fatal error: failed to get kubeconfig:  %s", err)
		os.Exit(1)
	}

	targetConf.TargetAddress = targetConf.TargetType + "-" + kubeConfig.Host //conf.KnativeTurboTargetConf.Kubeconfig

	registrationClient := &registration.KnativeTurboRegistrationClient{}
	discoveryClient := discovery.NewDiscoveryClient(targetConf, kubeConfig)

	targetType := targetConf.TargetType + "-" + fmt.Sprint(hash(targetConf.TargetAddress))

	return service.NewTAPServiceBuilder().
		WithTurboCommunicator(communicator).
		WithTurboProbe(probe.NewProbeBuilder(targetType, targetConf.ProbeCategory).
			WithDiscoveryOptions(probe.FullRediscoveryIntervalSecondsOption(int32(*args.DiscoveryIntervalSec))).
			RegisteredBy(registrationClient).
			WithEntityMetadata(registrationClient).
			DiscoversTarget(targetConf.TargetAddress, discoveryClient)).
		Create()
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

// TODO: Move the handle to turbo-sdk-probe as it should be common logic for similar probes
// handleExit disconnects the tap service from Turbo service when kongturbo is terminated
func handleExit(disconnectFunc disconnectFromTurboFunc) {
	glog.V(4).Infof("*** Handling Knativeturbo Termination ***")
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGHUP)

	go func() {
		select {
		case sig := <-sigChan:
			// Close the mediation container including the endpoints. It avoids the
			// invalid endpoints remaining in the server side. See OM-28801.
			glog.V(2).Infof("Signal %s received. Disconnecting from Turbo server...\n", sig)
			disconnectFunc()
		}
	}()
}
