package conf

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/turbonomic/turbo-go-sdk/pkg/service"
	"io/ioutil"
)

const (
// LocalDebugConfPath = "configs/knative-prometurbo-config.json"
// DefaultConfPath    = "/etc/kongturbo/turbo.config"
)

const (
	defaultProbeCategory = "Cloud Native"
	defaultTargetType    = "Knative"
)

type KnativeTurboServiceSpec struct {
	*service.TurboCommunicationConfig 	`json:"communicationConfig,omitempty"`
	*KnativeTurboTargetConf              	//`json:"knativeturboTargetConfig,omitempty"`
}

type KnativeTurboTargetConf struct {
	ProbeCategory string `json:"probeCategory,omitempty"`
	TargetType string `json:"probeCategory,omitempty"`
	TargetAddress string
	Kubeconfig       string `json:"kubeconfig,omitempty"`
}

func NewKnativeTurboServiceSpec(configFilePath string) (*KnativeTurboServiceSpec, error) {

	glog.Infof("Read configuration from %s", configFilePath)
	tapSpec, err := readConfig(configFilePath)

	if err != nil {
		return nil, err
	}

	if tapSpec.TurboCommunicationConfig == nil {
		return nil, fmt.Errorf("Unable to read the turbo communication config from %s", configFilePath)
	}

	tapSpec.KnativeTurboTargetConf = &KnativeTurboTargetConf{
		ProbeCategory:defaultProbeCategory,
		TargetType: defaultTargetType,
	}

	return tapSpec, nil
}

func readConfig(path string) (*KnativeTurboServiceSpec, error) {
	// path = "/Users/pallavinayak/GO/src/github.com/pallavidn/knative-turbo/probe/local.config/config"
	file, err := ioutil.ReadFile(path)
	if err != nil {
		glog.Errorf("File error: %v\n", err)
		return nil, err
	}
	glog.Infoln(string(file))

	var spec KnativeTurboServiceSpec
	err = json.Unmarshal(file, &spec)

	if err != nil {
		glog.Errorf("Unmarshall error :%v\n", err)
		return nil, err
	}
	glog.Infof("Results: %+v\n", spec)

	return &spec, nil
}
