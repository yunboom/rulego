# server

[English](README.md)| 中文

`RuleGo-Server`一个独立运行的开箱即用规则引擎服务，该工程也是一个开发RuleGo应用的脚手架。你可以基于该工程进行二次开发，也可以直接下载可执行[二进制文件](https://github.com/yunboom/rulego/releases)。

可视化编辑器：[RuleGo-Editor](https://editor.rulego.cc/) ，配置该工程HTTP API，可以对规则链管理和调试。

- 体验地址1：[http://8.134.32.225:9090/editor/](http://8.134.32.225:9090/editor/)
- 体验地址2：[http://8.134.32.225:9090/ui/](http://8.134.32.225:9090/ui/)

## 特性
- 基于RuleGo 独立运行的开箱即用规则引擎服务
- 可应用于边缘计算、IoT、大模型编排、应用编排、数据处理网关、自动化等应用场景
- 自动扫描组件和组件表单，提供给编辑器使用
- 支持对规则链进行可视化管理，调试，部署和对外提供API执行规则链等
- 支持RuleGo-Editor可视化前端
- 部署简单、开箱即用、不需要数据库
- 轻量级，内存小，性能高
- 自动把所有组件和规则链注册成MCP工具，对外提供给AI助手调用。详情：[rulego-server-mcp](https://rulego.cc/pages/rulego-server-mcp/)

## HTTP API

[API 文档](https://apifox.com/apidoc/shared-d17a63fe-2201-4e37-89fb-f2e8c1cbaf40/234016936e0)

* 获取所有组件列表
    - GET /api/v1/components

* 执行规则链并得到执行结果API
    - POST /api/v1/rules/:chainId/execute/:msgType
    - chainId：处理数据的规则链ID
    - msgType：消息类型
    - body：消息体
  
* 往规则链上报数据API，不关注执行结果
  - POST /api/v1/rules/:chainId/notify/:msgType
  - chainId：处理数据的规则链ID
  - msgType：消息类型
  - body：消息体
  
* 查询规则链
    - GET /api/v1/rules/{chainId}
    - chainId：规则链ID

* 保存或更新规则链
    - POST /api/v1/rule/{chainId}
    - chainId：规则链ID
    - body：更新规则链DSL内容
  
* 保存规则链Configuration
    - POST /api/v1/rules/:chainId/config/:varType
    - chainId：规则链ID
    - varType: vars/secrets 变量/秘钥
    - body：配置内容

* 获取节点调试日志API
    - Get /api/v1/logs/debug?&chainId={chainId}&nodeId={nodeId}
    - chainId：规则链ID
    - nodeId：节点ID

  当节点debugMode打开后，会记录调试日志。目前该接口日志存放在内存，每个节点保存最新的40条，如果需要获取历史数据，请实现接口存储到数据库。

## 多租户/多用户
该工程支持多租户/用户，每个用户的规则链数据是隔离的，用户数据存在`data/workflows/{username}`目录下。

用户权限校验默认是关闭，所有操作都是基于默认用户操作。开启权限校验方法：

- 通关过用户名密码获取token，然后通过token访问其他接口。示例:
```ini
# api是否开启jwt认证，如果关闭，则以默认用户(admin)身份操作
require_auth = true
# jwt secret key
jwt_secret_key = r6G7qZ8xk9P0y1Q2w3E4r5T6y7U8i9O0pL7z8x9CvBnM3k2l1
# jwt expire time
jwt_expire_time = 43200000
# jwt issuer
jwt_issuer = rulego.cc
# 用户列表
# 配置用户和密码，格式 username=password[,apiKey]，apiKey可选。
# 如果配置apiKey 调用方可以不需要登录，直接通过apiKey访问其他接口。
[users]
admin = admin
user01 = user01
```
前端通过登录接口(`/api/v1/login`)，获取token，然后通过token访问其他接口。示例：
```shell
curl -H "Authorization: Bearer token" http://localhost:8080/api/resource
```
- 通过`api_key`方式访问其他接口。示例：
```ini
# api是否开启jwt认证，如果关闭，则以默认用户(admin)身份操作
require_auth = true
# 用户列表
# 配置用户和密码，格式 username=password[,apiKey]，apiKey可选。
# 如果配置apiKey 调用方可以不需要登录，直接通过apiKey访问其他接口。
[users]
admin = admin,2af255ea-5618-467d-914c-67a8beeca31d
user01 = user01
```

然后通过token访问其他接口。示例：
```shell
curl -H "Authorization: Bearer apiKey" http://localhost:8080/api/resource
```

## server编译

为了节省编译后文件大小，默认不引入扩展组件[rulego-components](https://github.com/yunboom/rulego-components) ，默认编译：

```shell
cd cmd/server
go build .
```

如果需要引入扩展组件[rulego-components](https://github.com/yunboom/rulego-components) ，使用`with_extend`tag进行编译：

```shell
cd cmd/server
go build -tags with_extend .
```
其他扩展组件库tags：
- 注册扩展组件[rulego-components](https://github.com/yunboom/rulego-components) ，使用`with_extend`tag进行编译：
- 注册AI扩展组件[rulego-components-ai](https://github.com/yunboom/rulego-components-ai) ，使用`with_ai`tag进行编译
- 注册CI/CD扩展组件[rulego-components-ci](https://github.com/yunboom/rulego-components-ci) ，使用`with_ci`tag进行编译
- 注册IoT扩展组件[rulego-components-iot](https://github.com/yunboom/rulego-components-iot) ，使用`with_iot`tag进行编译
- 注册ETL扩展组件[rulego-components-etl](https://github.com/yunboom/rulego-components-etl) ，使用`with_etl`tag进行编译
- 使用`fasthttp`代替标准`endpoint/http`和`restApiCall`组件 ，使用`use_fasthttp`tag进行编译

如果需要同时引入多个扩展组件库，可以使用`go build -tags "with_extend,with_ai,with_ci,with_iot,with_etl,use_fasthttp" .` tag进行编译。

## server启动

```shell
./server -c="./config.conf"
```

或者后台启动

```shell
nohup ./server -c="./config.conf" >> console.log &
```
## RuleGo-Editor
RuleGo-Editor 是 RuleGo-Server 的UI界面，可以对规则链进行可视化管理，调试，部署等。

使用步骤：
- 解压下载好的`editor.zip`到当前目录，打开浏览器访问`http://localhost:9090/` ，即可访问RuleGo-Editor。
- 可以通过`config.conf`的 resource_mapping 配置修改rulego-editor目录。
- 可以通过`editor/config/config.js`的 baseUrl 配置修改rulego-editor后端api地址。

> RuleGo-Editor仅用于学习，商用请向我们购买授权。Email：rulego@outlook.com

## RuleGo-Server-MCP
RuleGo-Server 支持 MCP（Model Context Protocol，模型上下文协议），开启后，系统会自动将所有注册的组件、规则链以及 API 注册为 MCP 工具。这使得 AI 助手（如 Windsurf、Cursor、Codeium 等）能够通过 MCP 协议直接调用这些工具，实现与应用系统的深度融合。
文档: [rulego-server-mcp](https://rulego.cc/pages/rulego-server-mcp/)

## 配置文件参数
```ini
# 数据目录
data_dir = ./data
# cmd组件命令白名单
cmd_white_list = cp,scp,mvn,npm,yarn,git,make,cmake,docker,kubectl,helm,ansible,puppet,pytest,python,python3,pip,go,java,dotnet,gcc,g++,ctest
# 是否加载lua第三方库
load_lua_libs = true
# http server
server = :9090
# 默认用户
default_username = admin
# 是否把节点执行日志打印到日志文件
debug = true
# 最大节点日志大小，默认40
max_node_log_size =40
# 资源映射，支持通配符，多个映射用逗号分隔，格式：/url/*filepath=/path/to/file
resource_mapping = /editor/*filepath=./editor,/images/*filepath=./editor/images
# 节点池文件，规则链json格式，示例：./node_pool.json
node_pool_file=./node_pool.json
# save run log to file
save_run_log = false
# script max execution time
script_max_execution_time = 5000
# api是否开启jwt认证
require_auth = false
# jwt secret key
jwt_secret_key = r6G7qZ8xk9P0y1Q2w3E4r5T6y7U8i9O0pL7z8x9CvBnM3k2l1
# jwt expire time，单位毫秒
jwt_expire_time = 43200000
# jwt issuer
jwt_issuer = rulego.cc
# Set the default HTTP server as a shared node
share_http_server = true
# mcp server config
[mcp]
# Whether to enable the MCP service
enable = true
# Whether to use the component as an MCP tool
load_components_as_tool = true
# Whether to use the rule chain as an MCP tool
load_chains_as_tool = true
# Whether to add a rule chain api tool
load_apis_as_tool = true
# Exclude component list
exclude_components = comment,iterator,delay,groupAction,ref,fork,join,*Filter
# Exclude rule chain list
exclude_chains =

# pprof配置
[pprof]
# 是否开启pprof
enable = false
# pprof地址
addr = 0.0.0.0:6060

# 全局自定义配置，组件可以通过${global.xxx}方式取值
[global]
# 例子
sqlDriver = mysql
sqlDsn = root:root@tcp(127.0.0.1:3306)/test

# 用户列表
# 配置用户和密码，格式 username=password[,apiKey]，apiKey可选。
# 如果配置apiKey 调用方可以不需要登录，直接通过apiKey访问其他接口。
[users]
admin = admin
user01 = user01
```