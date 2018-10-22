package registration

import (
	"github.com/golang/glog"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	"github.com/turbonomic/turbo-go-sdk/pkg/supplychain"
)

var (
	transactionType proto.CommodityDTO_CommodityType = proto.CommodityDTO_TRANSACTION
	respTimeType    proto.CommodityDTO_CommodityType = proto.CommodityDTO_RESPONSE_TIME
	appCommType     proto.CommodityDTO_CommodityType = proto.CommodityDTO_APPLICATION
	fakeKey                                          = "key-placeholder"

	transactionTemplateCommWithKey *proto.TemplateCommodity = &proto.TemplateCommodity{Key: &fakeKey, CommodityType: &transactionType}
	respTimeTemplateCommWithKey    *proto.TemplateCommodity = &proto.TemplateCommodity{Key: &fakeKey, CommodityType: &respTimeType}
	applicationTemplateCommWithKey *proto.TemplateCommodity = &proto.TemplateCommodity{Key: &fakeKey, CommodityType: &appCommType}
)

type SupplyChainFactory struct{}

func (f *SupplyChainFactory) CreateSupplyChain() ([]*proto.TemplateDTO, error) {
	// Virtual application supply chain template
	vAppSupplyChainNode, err := f.buildVirtualApplicationSupplyBuilder()
	if err != nil {
		return nil, err
	}
	glog.V(4).Infof("supply chain node : %++v", vAppSupplyChainNode)

	supplyChainBuilder := supplychain.NewSupplyChainBuilder()
	supplyChainBuilder.Top(vAppSupplyChainNode)

	return supplyChainBuilder.Create()
}

func (f *SupplyChainFactory) buildVirtualApplicationSupplyBuilder() (*proto.TemplateDTO, error) {
	vAppSupplyChainNodeBuilder := supplychain.NewSupplyChainNodeBuilder(proto.EntityDTO_VIRTUAL_APPLICATION)
	vAppSupplyChainNodeBuilder = vAppSupplyChainNodeBuilder.
		Sells(applicationTemplateCommWithKey).
		Sells(transactionTemplateCommWithKey).
		Sells(respTimeTemplateCommWithKey)

	vAppSupplyChainNodeBuilder.SetPriority(-100)
	//vAppSupplyChainNodeBuilder.SetTemplateType(proto.TemplateDTO_BASE)

	return vAppSupplyChainNodeBuilder.Create()
}
