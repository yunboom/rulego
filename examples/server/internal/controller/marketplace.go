package controller

import (
	endpointApi "github.com/yunboom/rulego/api/types/endpoint"
	"github.com/yunboom/rulego/endpoint"
	"strconv"
)

// MarketplaceComponents 获取组件市场动态组件
func (c *node) MarketplaceComponents(url string) endpointApi.Router {
	return endpoint.NewRouter().From(url).Process(AuthProcess).Process(func(router endpointApi.Router, exchange *endpointApi.Exchange) bool {
		checkMyStr := exchange.In.GetParam("checkMy") //是否检查自己的组件
		var checkMy bool
		if i, err := strconv.ParseBool(checkMyStr); err == nil {
			checkMy = i
		}
		c.getCustomNodeList(true, checkMy, exchange)
		return true
	}).End()
}

func (c *rule) MarketplaceChains(url string) endpointApi.Router {
	return endpoint.NewRouter().From(url).Process(AuthProcess).Process(func(router endpointApi.Router, exchange *endpointApi.Exchange) bool {
		c.list(true, exchange)
		return true
	}).End()
}
