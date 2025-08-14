package controller

import (
	"examples/server/internal/constants"
	"examples/server/internal/service"
	"github.com/yunboom/rulego/api/types"
	endpointApi "github.com/yunboom/rulego/api/types/endpoint"
	"github.com/yunboom/rulego/endpoint"
	"github.com/yunboom/rulego/utils/json"
	"net/http"
	"strconv"
)

var Log = &log{}

type log struct {
}

// GetDebugLogs 创建获取节点调试数据路由
func (c *log) GetDebugLogs(url string) endpointApi.Router {
	return endpoint.NewRouter().From(url).Process(AuthProcess).Process(func(router endpointApi.Router, exchange *endpointApi.Exchange) bool {
		msg := exchange.In.GetMsg()
		chainId := msg.Metadata.GetValue(constants.KeyChainId)
		nodeId := msg.Metadata.GetValue(constants.KeyNodeId)
		username := msg.Metadata.GetValue(constants.KeyUsername)
		var current = 1
		var pageSize = 20
		currentStr := msg.Metadata.GetValue(constants.KeyPage)
		if i, err := strconv.Atoi(currentStr); err == nil {
			current = i
		}
		pageSizeStr := msg.Metadata.GetValue(constants.KeySize)
		if i, err := strconv.Atoi(pageSizeStr); err == nil {
			pageSize = i
		}
		if s, ok := service.UserRuleEngineServiceImpl.Get(username); ok {
			page := s.DebugData().GetToPage(chainId, nodeId, pageSize, current)
			if v, err := json.Marshal(page); err != nil {
				exchange.Out.SetStatusCode(http.StatusInternalServerError)
				exchange.Out.SetBody([]byte(err.Error()))
			} else {
				exchange.Out.SetBody(v)
			}
		} else {
			return userNotFound(username, exchange)
		}
		return true
	}).End()
}

func (c *log) List(url string) endpointApi.Router {
	return endpoint.NewRouter().From(url).Process(AuthProcess).Process(func(router endpointApi.Router, exchange *endpointApi.Exchange) bool {
		msg := exchange.In.GetMsg()
		chainId := msg.Metadata.GetValue(constants.KeyChainId)
		id := msg.Metadata.GetValue(constants.KeyId)
		username := msg.Metadata.GetValue(constants.KeyUsername)
		var result interface{}
		if id == "" {
			var current = 1
			var pageSize = 20
			currentStr := msg.Metadata.GetValue(constants.KeyPage)
			if i, err := strconv.Atoi(currentStr); err == nil {
				current = i
			}
			pageSizeStr := msg.Metadata.GetValue(constants.KeySize)
			if i, err := strconv.Atoi(pageSizeStr); err == nil {
				pageSize = i
			}
			if v, total, err := service.EventServiceImpl.List(username, chainId, current, pageSize); err != nil {
				exchange.Out.SetStatusCode(http.StatusNotFound)
				exchange.Out.SetBody([]byte(err.Error()))
				return false
			} else {
				result = map[string]interface{}{
					"page":  current,
					"size":  pageSize,
					"total": total,
					"items": v,
				}
			}
		} else {
			if v, err := service.EventServiceImpl.Get(username, chainId, id); err != nil {
				exchange.Out.SetStatusCode(http.StatusNotFound)
				exchange.Out.SetBody([]byte(err.Error()))
				return false
			} else {
				result = v
			}
		}

		if v, err := json.Marshal(result); err != nil {
			exchange.Out.SetStatusCode(http.StatusInternalServerError)
			exchange.Out.SetBody([]byte(err.Error()))
		} else {
			exchange.Out.SetBody(v)
		}
		return true
	}).End()
}

func (c *log) Delete(url string) endpointApi.Router {
	return endpoint.NewRouter().From(url).Process(AuthProcess).Process(func(router endpointApi.Router, exchange *endpointApi.Exchange) bool {
		msg := exchange.In.GetMsg()
		chainId := msg.Metadata.GetValue(constants.KeyChainId)
		id := msg.Metadata.GetValue(constants.KeyId)
		username := msg.Metadata.GetValue(constants.KeyUsername)

		if err := service.EventServiceImpl.Delete(username, chainId, id); err != nil {
			exchange.Out.SetStatusCode(http.StatusInternalServerError)
			exchange.Out.SetBody([]byte(err.Error()))
		}
		return true
	}).End()
}

func (c *log) WsNodeLogRouter(url string) endpointApi.Router {
	return endpoint.NewRouter().From(url).Process(AuthProcess).Process(func(router endpointApi.Router, exchange *endpointApi.Exchange) bool {
		msg := exchange.In.GetMsg()
		username := msg.Metadata.GetValue(constants.KeyUsername)
		if s, ok := service.UserRuleEngineServiceImpl.Get(username); ok {
			chainId := exchange.In.GetParam(constants.KeyChainId)
			clientId := exchange.In.GetParam(constants.KeyClientId)
			s.AddOnDebugObserver(chainId, clientId, func(chainId, flowType string, nodeId string, msg types.RuleMsg, relationType string, err error) {
				exchange.Out.SetBody([]byte(msg.GetData()))
			})
		}
		return true
	}).End()
}
