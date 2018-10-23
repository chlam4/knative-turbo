package discovery

import (
	"fmt"
	"github.com/turbonomic/turbo-go-sdk/pkg/builder"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	"strings"
)

const (
	StitchingAttr            string = "VappIds"	//"vAppUuid"
	DefaultPropertyNamespace string = "DEFAULT"
	extPropAttr string = "DisplayName"
	PropertyUsed        = "used"
	PropertyCapacity    = "capacity"
)

type KnativeDTOBuilder struct {
}

func (dtoBuilder *KnativeDTOBuilder) buildFunctionDto(funcSvc *KnativeFunction) (*builder.EntityDTOBuilder, error) {
	if funcSvc == nil {
		return nil, fmt.Errorf("Null service for %++v", funcSvc)
	}

	// id.
	var vappId string
	vappId = funcSvc.HostName

	if vappId == "" {
		return nil, fmt.Errorf("Cannot create function vapp without ID %++v", vappId)
	}
	vappId = fmt.Sprintf("%s/%s-%s", funcSvc.Namespace, funcSvc.Revision, "service")
	//vappId = fmt.Sprintf("%s/%s-%s", "vApp", funcSvc.Namespace, funcSvc.Revision, "service")

	// display name.
	funcName := fmt.Sprintf("%s/%s", funcSvc.Namespace, funcSvc.Revision)
	//funcName := fmt.Sprintf("%s/%s", "vApp", funcSvc.Namespace, funcSvc.Revision)
	fmt.Printf("**** vapp id : %s\n", funcName)

	commodities := []*proto.CommodityDTO{}
	commodity, _ := builder.NewCommodityDTOBuilder(proto.CommodityDTO_TRANSACTION).Key(vappId).Create()
	commodities = append(commodities, commodity)
	commodity, _ = builder.NewCommodityDTOBuilder(proto.CommodityDTO_RESPONSE_TIME).Key(vappId).Create()
	commodities = append(commodities, commodity)
	commodity, _ = builder.NewCommodityDTOBuilder(proto.CommodityDTO_APPLICATION).Key(funcSvc.FunctionName).Create()
	commodities = append(commodities, commodity)

	entityDTOBuilder := builder.NewEntityDTOBuilder(proto.EntityDTO_VIRTUAL_APPLICATION, funcName).
		DisplayName(funcName).
		WithProperty(getEntityProperty(StitchingAttr, vappId)).		// + "," + funcSvc.HostName)).
			//WithProperty(getEntityProperty(extPropAttr, vappId)).
		SellsCommodities(commodities).
		ReplacedBy(getReplacementMetaData(proto.EntityDTO_VIRTUAL_APPLICATION))
		//Provider(provider).BuysCommodities(boughtCommodities)

	fmt.Printf("Created function dto builder\n")
	return entityDTOBuilder, nil
}


func getReplacementMetaData(entityType proto.EntityDTO_EntityType,
) *proto.EntityDTO_ReplacementEntityMetaData {
	attr := StitchingAttr
	useTopoExt := true

	b := builder.NewReplacementEntityMetaDataBuilder().
		Matching(attr).
		//MatchingExternalProperty(extPropAttr)
		MatchingExternal(&proto.ServerEntityPropDef{
			Entity:    &entityType,
			Attribute: &attr,
			UseTopoExt: &useTopoExt,
		}).
		PatchBuyingWithProperty(proto.CommodityDTO_TRANSACTION, []string{PropertyUsed}).
		PatchBuyingWithProperty(proto.CommodityDTO_RESPONSE_TIME, []string{PropertyUsed}).
		PatchSellingWithProperty(proto.CommodityDTO_TRANSACTION, []string{PropertyUsed, PropertyCapacity}).
		PatchSellingWithProperty(proto.CommodityDTO_RESPONSE_TIME, []string{PropertyUsed, PropertyCapacity})

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
	containerUuid := strings.Join( []string{"cg",funcSvc.FunctionName}, "-")
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
