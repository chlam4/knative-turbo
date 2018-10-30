package discovery

import (
	"fmt"
	"github.com/turbonomic/turbo-go-sdk/pkg/builder"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	"strings"
)

const (
	StitchingAttr        string = "VappIds"       //"vAppUuid"
	GatewayStitchingAttr string = "GatewayVappId" //"vAppUuid"
	KubernetesPropAttr   string = "KubernetesVappId"

	LocalNameAttr    string = "LocalName"
	AltNameAttr      string = "altName"
	ExternalNameAttr string = "externalnames"

	DefaultPropertyNamespace string = "DEFAULT"
	PropertyUsed                    = "used"
	PropertyCapacity                = "capacity"
)

type KnativeDTOBuilder struct {
}

func (dtoBuilder *KnativeDTOBuilder) buildFunctionDto(funcSvc *KnativeFunction) (*builder.EntityDTOBuilder, error) {
	if funcSvc == nil {
		return nil, fmt.Errorf("Null service for %++v", funcSvc)
	}

	// id.
	var vappId, localId, altId string
	//localId = funcSvc.HostName

	//if vappId == "" {
	//	return nil, fmt.Errorf("Cannot create function vapp without ID %++v", localId)
	//}

	// uuid and display name.
	vappId = fmt.Sprintf("%s/%s/%s", "knative", funcSvc.FunctionName, funcSvc.HostName)
	//funcName := fmt.Sprintf("%s/%s", "vApp", funcSvc.Namespace, funcSvc.Revision)
	fmt.Printf("**** vapp id : %s\n", vappId)

	localId = fmt.Sprintf("%s/%s-%s", funcSvc.Namespace, funcSvc.Revision, "service")
	altId =  funcSvc.HostName
	fmt.Printf("**** local name: %s\n", localId)
	fmt.Printf("**** alt name: %s\n", altId)

	commodities := []*proto.CommodityDTO{}
	commodity, _ := builder.NewCommodityDTOBuilder(proto.CommodityDTO_TRANSACTION).Key(localId).Create()
	commodities = append(commodities, commodity)
	commodity, _ = builder.NewCommodityDTOBuilder(proto.CommodityDTO_RESPONSE_TIME).Key(localId).Create()
	commodities = append(commodities, commodity)
	commodity, _ = builder.NewCommodityDTOBuilder(proto.CommodityDTO_APPLICATION).Key(localId).Create()
	commodities = append(commodities, commodity)

	entityDTOBuilder := builder.NewEntityDTOBuilder(proto.EntityDTO_VIRTUAL_APPLICATION, vappId).
		DisplayName(vappId).
		WithProperty(getEntityProperty(LocalNameAttr, localId)).
		WithProperty(getEntityProperty(AltNameAttr, altId)).
		//WithProperty(getEntityProperty(KubernetesPropAttr, vappId)).
		//WithProperty(getEntityProperty(GatewayStitchingAttr, funcSvc.HostName)).
		SellsCommodities(commodities).
		ReplacedBy(getReplacementMetaData(proto.EntityDTO_VIRTUAL_APPLICATION))

	fmt.Printf("Created function dto builder\n")
	return entityDTOBuilder, nil
}

func getReplacementMetaData(entityType proto.EntityDTO_EntityType,
) *proto.EntityDTO_ReplacementEntityMetaData {
	extAttr := ExternalNameAttr	//StitchingAttr
	intAttr := LocalNameAttr	//KubernetesPropAttr
	useTopoExt := true

	b := builder.NewReplacementEntityMetaDataBuilder().
		Matching(intAttr).
		MatchingExternal(&proto.ServerEntityPropDef{
			Entity:     &entityType,
			Attribute:  &extAttr,
			UseTopoExt: &useTopoExt,
		})
		//PatchBuyingWithProperty(proto.CommodityDTO_TRANSACTION, []string{PropertyUsed}).
		//PatchBuyingWithProperty(proto.CommodityDTO_RESPONSE_TIME, []string{PropertyUsed}).
		//PatchSellingWithProperty(proto.CommodityDTO_TRANSACTION, []string{PropertyUsed, PropertyCapacity}).
		//PatchSellingWithProperty(proto.CommodityDTO_RESPONSE_TIME, []string{PropertyUsed, PropertyCapacity})

	return b.Build()
}

func getEntityProperty(attr, value string) *proto.EntityDTO_EntityProperty {
	ns := DefaultPropertyNamespace

	return &proto.EntityDTO_EntityProperty{
		Namespace: &ns,
		Name:      &attr,
		Value:     &value,
	}
}

func (dtoBuilder *KnativeDTOBuilder) buildContainerDto(funcSvc *KnativeFunction) (*builder.EntityDTOBuilder, error) {
	if funcSvc == nil {
		return nil, fmt.Errorf("Null service for %++v", funcSvc)
	}

	// id.
	containerUuid := strings.Join([]string{"cg", funcSvc.FunctionName}, "-")
	fmt.Printf("**** vapp id : %s\n", containerUuid)

	// display name.
	displayName := containerUuid

	commodities := []*proto.CommodityDTO{}
	commodity, _ := builder.NewCommodityDTOBuilder(proto.CommodityDTO_APPLICATION).Key(containerUuid).Create()
	commodities = append(commodities, commodity)
	commodity, _ = builder.NewCommodityDTOBuilder(proto.CommodityDTO_VCPU).Create()
	commodities = append(commodities, commodity)
	commodity, _ = builder.NewCommodityDTOBuilder(proto.CommodityDTO_VMEM).Create()
	commodities = append(commodities, commodity)

	entityDTOBuilder := builder.NewEntityDTOBuilder(proto.EntityDTO_CONTAINER, containerUuid).
		DisplayName(displayName).
		SellsCommodities(commodities)

	return entityDTOBuilder, nil
}
