package registration

import (
	"github.com/golang/glog"
	"github.com/turbonomic/turbo-go-sdk/pkg/builder"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

const (
	TargetIdentifierField string = "targetIdentifier"
	//ProbeCategory string = "Cloud Management"
	//TargetType    string = "Knative"
	propertyId            string = "id"
)

// Implements the TurboRegistrationClient interface
type KnativeTurboRegistrationClient struct {
}

func (p *KnativeTurboRegistrationClient) GetSupplyChainDefinition() []*proto.TemplateDTO {
	glog.Infoln("Building a supply chain ..........")

	supplyChainFactory := &SupplyChainFactory{}
	templateDtos, err := supplyChainFactory.CreateSupplyChain()
	if err != nil {
		glog.Error("Error creating Supply chain for KnativeTurbo")
		return nil
	}
	glog.Infoln("Supply chain for KnativeTurbo is created.")
	return templateDtos
}

func (p *KnativeTurboRegistrationClient) GetIdentifyingFields() string {
	return TargetIdentifierField
}

func (p *KnativeTurboRegistrationClient) GetAccountDefinition() []*proto.AccountDefEntry {

	targetIDAcctDefEntry := builder.NewAccountDefEntryBuilder(TargetIdentifierField, "kubeconfig",
		"Path to the kubeconfig", ".*", true, false).Create()

	return []*proto.AccountDefEntry{
		targetIDAcctDefEntry,
	}
}


func (rclient *KnativeTurboRegistrationClient) GetEntityMetadata() []*proto.EntityIdentityMetadata {
	glog.V(3).Infof("Begin to build EntityIdentityMetadata")

	result := []*proto.EntityIdentityMetadata{}

	entities := []proto.EntityDTO_EntityType{
		proto.EntityDTO_VIRTUAL_APPLICATION,
	}

	for _, etype := range entities {
		meta := rclient.newIdMetaData(etype, []string{propertyId})
		result = append(result, meta)
	}

	glog.V(4).Infof("EntityIdentityMetaData: %++v", result)

	return result
}

func (rclient *KnativeTurboRegistrationClient) newIdMetaData(etype proto.EntityDTO_EntityType, names []string) *proto.EntityIdentityMetadata {
	data := []*proto.EntityIdentityMetadata_PropertyMetadata{}
	for _, name := range names {
		dat := &proto.EntityIdentityMetadata_PropertyMetadata{
			Name: &name,
		}
		data = append(data, dat)
	}

	result := &proto.EntityIdentityMetadata{
		EntityType:            &etype,
		NonVolatileProperties: data,
	}

	return result
}
