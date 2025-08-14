package router

import (
	"errors"
	"examples/server/config"
	"examples/server/config/logger"
	"examples/server/internal/controller"
	"github.com/yunboom/rulego"
	"github.com/yunboom/rulego/api/types"
	endpointApi "github.com/yunboom/rulego/api/types/endpoint"
	"github.com/yunboom/rulego/endpoint"
	"github.com/yunboom/rulego/endpoint/rest"
	"github.com/yunboom/rulego/node_pool"
	"net/http"
	"strings"
)

const (
	// base HTTP paths.
	apiVersion  = "v1"
	apiBasePath = "/api/" + apiVersion
	moduleFlows = "rules"
	// moduleDcs 动态组件
	moduleDynamicComponents = "dynamic-components"
	// moduleSharedNodes 共享组件
	moduleSharedNodes = "shared-nodes"
	moduleLocales     = "locales"
	moduleLogs        = "logs"
	moduleMarketplace = "marketplace"
	ContentTypeKey    = "Content-Type"
	JsonContextType   = "application/json"
)

// SystemRulegoConfig 系统rulego配置
var SystemRulegoConfig types.Config

// SystemNodePool 系统内部节点池
var SystemNodePool *node_pool.NodePool

func InitRulegoConfig() {
	SystemRulegoConfig = rulego.NewConfig(types.WithDefaultPool(), types.WithLogger(logger.Logger))
	SystemNodePool = node_pool.NewNodePool(SystemRulegoConfig)
	SystemRulegoConfig.NodePool = SystemNodePool
}

// NewRestServe rest服务 接收端点
func NewRestServe(config config.Config) (endpointApi.HttpEndpoint, error) {
	//初始化日志
	addr := config.Server
	if strings.HasPrefix(addr, ":") {
		logger.Logger.Println("RuleGo-Server now running at http://127.0.0.1" + addr)
	} else {
		logger.Logger.Println("RuleGo-Server now running at http://" + addr)
	}

	ep, err := endpoint.Registry.New(
		rest.Type,
		SystemRulegoConfig,
		rest.Config{
			Server:    addr,
			AllowCors: true,
		},
	)
	if err != nil {
		return nil, err
	}
	var restEndpoint endpointApi.HttpEndpoint
	if ep, ok := ep.(endpointApi.HttpEndpoint); !ok {
		return nil, errors.New("is not HttpEndpoint type error")
	} else {
		restEndpoint = ep
	}
	//添加全局拦截器
	restEndpoint.AddInterceptors(func(router endpointApi.Router, exchange *endpointApi.Exchange) bool {
		if out, ok := exchange.Out.(endpointApi.HeaderModifier); ok {
			out.AddHeader(ContentTypeKey, JsonContextType)
		} else {
			exchange.Out.Headers().Set(ContentTypeKey, JsonContextType)
		}
		return true
	})
	//重定向UI界面
	restEndpoint.GET(endpoint.NewRouter().From("/").Process(func(router endpointApi.Router, exchange *endpointApi.Exchange) bool {
		r, ok1 := exchange.In.(*rest.RequestMessage)
		w, ok2 := exchange.Out.(*rest.ResponseMessage)
		if ok1 && ok2 {
			http.Redirect(w.Response(), r.Request(), "/editor/", http.StatusFound)
		}
		return false
	}).End())
	//创建获取所有规则引擎组件列表路由
	restEndpoint.GET(controller.Node.Components(apiBasePath + "/components"))

	//获取所有共享组件
	restEndpoint.GET(controller.Node.ListNodePool(apiBasePath + "/" + moduleSharedNodes))

	//获取组件市场组件列表
	restEndpoint.GET(controller.Node.MarketplaceComponents(apiBasePath + "/" + moduleMarketplace + "/components"))
	//获取组件市场规则链列表
	restEndpoint.GET(controller.Rule.MarketplaceChains(apiBasePath + "/" + moduleMarketplace + "/chains"))

	//获取用户所有自定义动态组件列表
	restEndpoint.GET(controller.Node.CustomNodeList(apiBasePath + "/" + moduleDynamicComponents))
	//获取自定义动态组件DSL
	restEndpoint.GET(controller.Node.CustomNodeDSL(apiBasePath + "/" + moduleDynamicComponents + "/:id"))
	//安装/升级自定义动态组件
	restEndpoint.POST(controller.Node.CustomNodeUpgrade(apiBasePath + "/" + moduleDynamicComponents + "/:id"))
	//卸装自定义动态组件
	restEndpoint.DELETE(controller.Node.CustomNodeUninstall(apiBasePath + "/" + moduleDynamicComponents + "/:id"))

	//获取所有规则链列表
	restEndpoint.GET(controller.Rule.List(apiBasePath + "/" + moduleFlows))
	//获取最新修改的规则链DSL 实际是：/api/v1/rules/get/latest
	restEndpoint.GET(controller.Rule.GetLatest(apiBasePath + "/" + moduleFlows + "/:id/latest"))
	//获取规则链DSL
	restEndpoint.GET(controller.Rule.Get(apiBasePath + "/" + moduleFlows + "/:id"))
	//新增/修改规则链DSL
	restEndpoint.POST(controller.Rule.Save(apiBasePath + "/" + moduleFlows + "/:id"))
	//删除规则链
	restEndpoint.DELETE(controller.Rule.Delete(apiBasePath + "/" + moduleFlows + "/:id"))
	//保存规则链附加信息
	restEndpoint.POST(controller.Rule.SaveBaseInfo(apiBasePath + "/" + moduleFlows + "/:id/base"))
	//保存规则链配置信息
	restEndpoint.POST(controller.Rule.SaveConfiguration(apiBasePath + "/" + moduleFlows + "/:id/config/:varType"))
	//执行规则链,并得到规则链处理结果
	restEndpoint.POST(controller.Rule.Execute(apiBasePath + "/" + moduleFlows + "/:id/execute/:msgType"))
	//处理数据上报请求，并转发到规则引擎，不等待规则引擎处理结果
	restEndpoint.POST(controller.Rule.PostMsg(apiBasePath + "/" + moduleFlows + "/:id/notify/:msgType"))
	//部署或者下线规则链
	restEndpoint.POST(controller.Rule.Operate(apiBasePath + "/" + moduleFlows + "/:id/operate/:type"))

	//获取节点调试日志列表
	restEndpoint.GET(controller.Log.GetDebugLogs(apiBasePath + "/" + moduleLogs + "/debug"))
	//获取规则链运行日志列表
	restEndpoint.GET(controller.Log.List(apiBasePath + "/" + moduleLogs + "/runs"))
	//获取规则链运行日志详情
	restEndpoint.DELETE(controller.Log.Delete(apiBasePath + "/" + moduleLogs + "/runs"))

	restEndpoint.GET(controller.Locale.Locales(apiBasePath + "/" + moduleLocales))
	restEndpoint.POST(controller.Locale.Save(apiBasePath + "/" + moduleLocales))
	//创建用户登录路由
	restEndpoint.POST(controller.Base.Login(apiBasePath + "/login"))

	if config.MCP.Enable {
		restEndpoint.GET(controller.MCP.Handler(apiBasePath + "/mcp/:apiKey/sse"))
		restEndpoint.POST(controller.MCP.Handler(apiBasePath + "/mcp/:apiKey/message"))
		logger.Logger.Println("RuleGo-Server mcp server running at http://127.0.0.1" + addr + apiBasePath + "/mcp/" +
			config.GetApiKeyByUsername(config.DefaultUsername) + "/sse")
	}
	// 加载静态文件映射
	restEndpoint.RegisterStaticFiles(config.ResourceMapping)

	//把默认HTTP服务设置成共享节点
	if config.ShareHttpServer {
		_, _ = node_pool.DefaultNodePool.AddNode(restEndpoint)
	}
	//把默认HTTP服务添加到系统节点池
	_, _ = SystemNodePool.AddNode(restEndpoint)
	return restEndpoint, nil
}

// LoadServeFiles 加载静态文件映射
func LoadServeFiles(c config.Config, restEndpoint endpointApi.HttpEndpoint) {
	if c.ResourceMapping != "" {
		restEndpoint.RegisterStaticFiles(c.ResourceMapping)
	}
}
