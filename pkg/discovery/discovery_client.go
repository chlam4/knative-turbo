package discovery

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/turbonomic/turbo-go-sdk/pkg/probe"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	"github.com/pallavidn/knative-turbo/pkg/registration"
	"k8s.io/client-go/rest"
	"github.com/pallavidn/knative-turbo/pkg/conf"
)

// Implements the TurboDiscoveryClient interface
type KnativeDiscoveryClient struct {
	//kubeconfig string
	kubeConfig *rest.Config
	targetConfig         *conf.KnativeTurboTargetConf
}

func NewDiscoveryClient(targetConfig *conf.KnativeTurboTargetConf, kubeConfig *rest.Config) *KnativeDiscoveryClient {
	glog.V(2).Infof("New Discovery client for kubernetes host: %s", kubeConfig.Host)
	return &KnativeDiscoveryClient{
		//kubeconfig: kubeconfig,
		targetConfig: targetConfig,
		kubeConfig: kubeConfig,
	}
}

// Get the Account Values to create VMTTarget in the turbo server corresponding to this client
func (d *KnativeDiscoveryClient) GetAccountValues() *probe.TurboTargetInfo {
	targetId := registration.TargetIdentifierField
	targetConf := d.targetConfig
	targetIdVal := &proto.AccountValue{
		Key:         &targetId,
		StringValue: &targetConf.TargetAddress,
	}

	accountValues := []*proto.AccountValue{
		targetIdVal,
	}

	targetInfo := probe.NewTurboTargetInfoBuilder(targetConf.ProbeCategory, targetConf.TargetType,
		registration.TargetIdentifierField, accountValues).Create()

	return targetInfo
}

// Validate the Target
func (d *KnativeDiscoveryClient) Validate(accountValues []*proto.AccountValue) (*proto.ValidationResponse, error) {
	glog.V(2).Infof("Validating Knative target %s", accountValues)
	fmt.Printf("Validating Knative target: %s\n", accountValues)
	// TODO: Add logic for validation
	validationResponse := &proto.ValidationResponse{}

	// Validation fails if no exporter responses
	return validationResponse, nil
}

// Discover the Target Topology
func (d *KnativeDiscoveryClient) Discover(accountValues []*proto.AccountValue) (*proto.DiscoveryResponse, error) {
	glog.V(2).Infof("Discovering Knative target %s", accountValues)
	fmt.Printf("Discovering Knative target %s\n", accountValues)
	var entities []*proto.EntityDTO

	var discoveryResponse *proto.DiscoveryResponse
	entities, err := DiscoverKnative(d.kubeConfig)
	if err != nil {
		fmt.Printf("Discovery failure %++v\n", err)
		return d.failDiscovery(), nil
	}

	discoveryResponse = &proto.DiscoveryResponse {
		EntityDTO: entities,
	}
	fmt.Printf("DONE Discovering Knative target: %s\n", d.kubeConfig.Host)
	return discoveryResponse, nil
}

func (d *KnativeDiscoveryClient) failDiscovery() *proto.DiscoveryResponse {
	description := fmt.Sprintf("KnativeTurbo probe discovery failed")
	glog.Errorf(description)
	severity := proto.ErrorDTO_CRITICAL
	errorDTO := &proto.ErrorDTO{
		Severity:    &severity,
		Description: &description,
	}
	discoveryResponse := &proto.DiscoveryResponse{
		ErrorDTO: []*proto.ErrorDTO{errorDTO},
	}
	return discoveryResponse
}

func (d *KnativeDiscoveryClient) failValidation() *proto.ValidationResponse {
	description := fmt.Sprintf("KnativeTurbo probe validation failed")
	glog.Errorf(description)
	severity := proto.ErrorDTO_CRITICAL
	errorDto := &proto.ErrorDTO{
		Severity:    &severity,
		Description: &description,
	}

	validationResponse := &proto.ValidationResponse{
		ErrorDTO: []*proto.ErrorDTO{errorDto},
	}
	return validationResponse
}
