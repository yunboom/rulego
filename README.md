# RuleGo

[![GoDoc](https://pkg.go.dev/badge/github.com/yunboom/rulego)](https://pkg.go.dev/github.com/yunboom/rulego) 
[![Go Report](https://goreportcard.com/badge/github.com/yunboom/rulego)](https://goreportcard.com/report/github.com/yunboom/rulego)
[![codecov](https://codecov.io/gh/rulego/rulego/graph/badge.svg?token=G6XCGY7KVN)](https://codecov.io/gh/rulego/rulego)
[![test](https://github.com/yunboom/rulego/workflows/test/badge.svg)](https://github.com/yunboom/rulego/actions/workflows/test.yml)
[![build](https://github.com/yunboom/rulego/workflows/build/badge.svg)](https://github.com/yunboom/rulego/actions/workflows/build.yml)
[![build](https://github.com/yunboom/rulego/workflows/build/badge.svg)](https://github.com/yunboom/rulego/actions/workflows/build.yml)
[![QQ-720103251](https://img.shields.io/badge/QQ-720103251-orange)](https://qm.qq.com/q/8RDaYcOry8)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go?tab=readme-ov-file#iot-internet-of-things)

English| [简体中文](README_ZH.md)

[Official Website](https://rulego.cc) | [Docs](https://rulego.cc/en/pages/0f6af2/) | [Contribution Guide](CONTRIBUTION.md)

<img src="doc/imgs/logo.png" alt="logo" width="100"/>   

`RuleGo` is a lightweight, high-performance, embedded, orchestrable component-based rule engine built on the Go language.

It can help you quickly build loosely coupled and flexible systems that can respond and adjust to changes in business requirements in real time.

`RuleGo` also provides a large number of reusable components that support the aggregation, filtering, distribution, transformation, enrichment, and execution of various actions on data, and can also interact and integrate with various protocols and systems.
It has a wide range of application potential in low-code, business code orchestration, data integration, workflows, large model intelligent agents, edge computing, automation, IoT, and other scenarios.

## Features

* **Lightweight:** No external middleware dependencies, efficient data processing and linkage on low-cost devices, suitable for IoT edge computing.
* **High Performance:** Thanks to Go's high-performance characteristics, RuleGo also employs technologies such as coroutine pools and object pools.
* **Dual Mode:** Embedded and Standalone Deployment modes. Supports embedding `RuleGo` into existing applications. It can also be deployed independently as middleware, providing rule engine and orchestration services.
* **Componentized:** All business logic is component-based, allowing flexible configuration and reuse.
* **Rule Chains:** Flexibly combine and reuse different components to achieve highly customized and scalable business processes.
* **Workflow Orchestration:** Supports dynamic orchestration of rule chain components, replacing or adding business logic without restarting the application.
* **Easy Extension:** Provides rich and flexible extension interfaces, making it easy to implement custom components or introduce third-party components.
* **Dynamic Loading:** Supports dynamic loading of components and extensions through Go plugins.
* **Nested Rule Chains:** Supports nesting of sub-rule chains to reuse processes.
* **Built-in Components:** Includes a large number of components such as `Message Type Switch`, `JavaScript Switch`, `JavaScript Filter`, `JavaScript Transformer`, `HTTP Push`, `MQTT Push`, `Send Email`, `Log Recording`, etc. Other components can be extended as needed.
* **Context Isolation Mechanism:** Reliable context isolation mechanism, no need to worry about data streaming in high concurrency situations.
* **AOP Mechanism:** Allows adding extra behavior to the execution of rule chains or directly replacing the original logic of rule chains or nodes without modifying their original logic.
* **Data Integration:** Allows dynamic configuration of Endpoints, such as `HTTP Endpoint`, `MQTT Endpoint`, `TCP/UDP Endpoint`, `UDP Endpoint`, `Kafka Endpoint`, `Schedule Endpoint`, etc.

## Use Cases

`RuleGo` is an orchestrable rule engine that excels at decoupling your systems.

- If your system's business is complex and the code is bloated
- If your business scenarios are highly customized or frequently changing
- If your system needs to interface with a large number of third-party applications or protocols
- Or if you need an end-to-end IoT solution
- Or if you need centralized processing of heterogeneous system data
- Or if you want to try hot deployment in the Go language...
  Then the RuleGo framework will be a very good solution.

#### Typical Use Cases

* **Edge Computing:** Deploy RuleGo on edge servers to preprocess data, filter, aggregate, or compute before reporting to the cloud. Data processing rules and distribution rules can be dynamically configured and modified through rule chains without restarting the system.
* **IoT:** Collect device data reports, make rule judgments through rule chains, and trigger one or more actions, such as sending emails, alarms, and linking with other devices or systems.
* **Data Distribution:** Distribute data to different systems using HTTP, MQTT, or gRPC based on different message types.
* **Application Integration:** Use RuleGo as glue to connect various systems or protocols, such as SSH, webhook, Kafka, message queues, databases, ChatGPT, third-party application systems.
* **Centralized Processing of Heterogeneous System Data:** Receive data from different sources (such as MQTT, HTTP, WS, TCP/UDP, etc.), then filter, format convert, and distribute to databases, business systems, or dashboards.
* **Highly Customized Business:** Decouple highly customized or frequently changing business and manage it with RuleGo rule chains. Business requirements change without needing to restart the main program.
* **Complex Business Orchestration:** Encapsulate business into custom components, orchestrate and drive these custom components through RuleGo, and support dynamic adjustment and replacement of business logic.
* **Microservice Orchestration:** Orchestrate and drive microservices through RuleGo, or dynamically call third-party services to process business and return results.
* **Decoupling of Business Code and Logic:** For example, user points calculation systems, risk control systems.
* **Automation:** For example, CI/CD systems, process automation systems, marketing automation systems.
* **Low Code:** For example, low-code platforms, iPaaS systems, ETL, LangFlow-like systems (interfacing with large models to extract user intent, then triggering rule chains to interact with other systems or process business).

## Architecture Diagram

<img src="doc/imgs/architecture.png" width="100%">
<p align="center">RuleGo Architecture Diagram</p>

## Rule Chain Running Example Diagram

  <img src="doc/imgs/rulechain/demo.png" style="height:40%;width:100%;"/>

[More Running Modes](http://8.134.32.225:9090/ui/)

## Installation

Install `RuleGo` using the `go get` command:

```bash
go get github.com/yunboom/rulego
# or
go get gitee.com/rulego/rulego
```

## Usage

`RuleGo` is extremely simple to use. Just follow these 3 steps:

1. Define rule chains using JSON:
   - [Rule Chain DSL Doc](https://rulego.cc/en/pages/10e1c0/) 
   - [example_chain.json](testdata/rule/chain_call_rest_api.json)

2. Import the RuleGo package and use the rule chain definition to create a rule engine instance:

```go
import "github.com/yunboom/rulego"
//Load the rule chain definition file.
ruleFile := fs.LoadFile("chain_call_rest_api.json")
// Create a rule engine instance using the rule chain definition
ruleEngine, err := rulego.New("rule01", ruleFile)
```

3. Hand over the message payload, message type, and message metadata to the rule engine instance for processing, and then the rule engine will process the message according to the rule chain's definition:

```go
// Define message metadata
metaData := types.NewMetadata()
metaData.PutValue("productType", "test01")
// Define message payload and message type
msg := types.NewMsg(0, "TELEMETRY_MSG", types.JSON, metaData, "{\"temperature\":35}")

// Hand over the message to the rule engine for processing
ruleEngine.OnMsg(msg)
```
> Real time update of rule chain logic without restarting the application
### Rule Engine Management API

- Dynamically update rule chains
```go
// Dynamically update rule chain logic
err := ruleEngine.ReloadSelf(ruleFile)
// Update a node under the rule chain
ruleEngine.ReloadChild("node01", nodeFile)
// Get the rule chain definition
ruleEngine.DSL()
```

- Rule Engine Instance Management:
```go
// Load all rule chain definitions in a folder into the rule engine pool
rulego.Load("/rules", rulego.WithConfig(config))
// Get an already created rule engine instance by ID
ruleEngine, ok := rulego.Get("rule01")
// Delete an already created rule engine instance
rulego.Del("rule01")
```

- Config:[Documentation](https://rulego.cc/en/pages/d59341/)
```go
// Create a default configuration
config := rulego.NewConfig()
// Debug node callback, the node configuration must be set to debugMode:true to trigger the call
// Both node entry and exit information will call this callback function
config.OnDebug = func (chainId,flowType string, nodeId string, msg types.RuleMsg, relationType string, err error) {
}
// Use the configuration
ruleEngine, err := rulego.New("rule01", []byte(ruleFile), rulego.WithConfig(config))
```

### Rule Chain Definition DSL
[Rule Chain Definition DSL](https://rulego.cc/en/pages/10e1c0/)

### Rule Chain Node Components
The core feature of `RuleGo` is its component-based architecture, where all business logic is encapsulated in components that can be flexibly configured and reused. Currently, 
`RuleGo` has built-in a vast array of commonly used components.
- [Standard Components](https://rulego.cc/en/pages/88fc3c/)
- [rulego-components](https://github.com/yunboom/rulego-components)  :[Documentation](https://rulego.cc/en/pages/d7fc43/)
- [rulego-components-ai](https://github.com/yunboom/rulego-components-ai)
- [rulego-components-ci](https://github.com/yunboom/rulego-components-ci)
- [rulego-components-iot](https://github.com/yunboom/rulego-components-iot)
- [rulego-components-etl](https://github.com/yunboom/rulego-components-etl)
- [rulego-marketplace](https://github.com/yunboom/rulego-marketplace) :Dynamic component and rule chain marketplace
- [Custom Node Component Example](examples/custom_component) :[Documentation](https://rulego.cc/en/pages/caed1b/)

## Data Integration
`RuleGo` provides the Endpoint module for unified data integration and processing of heterogeneous systems. For details, refer to: [Endpoint](endpoint/README.md)

### Input Endpoint Components
- [Endpoint Components](https://rulego.cc/en/pages/691dd3/)
- [Endpoint DSL](https://rulego.cc/en/pages/390ad7/)

### RuleGo-Server
We have provided a ready-to-use automation workflow platform based on `RuleGo`: [RuleGo-Server](https://github.com/yunboom/rulego-server) ,
and a front-end visual editor: [RuleGo-Editor](https://editor.rulego.cc/) .
## Performance

`RuleGo` completes most of its work during initialization, so running the rule chain almost doesn't add extra overhead to the system, and the resource consumption is extremely low, making it particularly suitable for running on edge servers.
Additionally, RuleGo uses Directed Acyclic Graph (DAG) to represent the rule chain, where each input message only needs to be processed along the path in the graph without matching all the rules, which greatly improves the efficiency of message processing and routing, and also saves system resources.

Performance test cases:
```
Machine: Raspberry Pi 2 (900MHz Cortex-A7*4,1GB LPDDR2)  
Data size: 260B   
Rule chain: JS script filtering->JS complex transformation->HTTP push   
Test results: 100 concurrent and 500 concurrent, memory consumption does not change much around 19M
```

[More performance test cases](https://rulego.cc/en/pages/f60381/)

## Ecosystem
- [RuleGo-Editor](https://editor.rulego.cc/) : Rule chain visual editor
- [RuleGo-Server](https://github.com/yunboom/rulego-server) : Automation workflow platform
- [rulego-components](https://gitee.com/rulego/rulego-components) : Extension component library
- [rulego-components-ai](https://github.com/yunboom/rulego-components-ai) : AI scenario component library
- [rulego-components-ci](https://github.com/yunboom/rulego-components-ci) : CI/CD scenario component library
- [rulego-components-iot](https://github.com/yunboom/rulego-components-iot) : IoT scenario component library
- [rulego-components-etl](https://github.com/yunboom/rulego-components-etl) : ETL scenario component library
- [rulego-marketplace](https://github.com/yunboom/rulego-marketplace) : Dynamic component and rule chain marketplace

## Contribution
Any form of contribution is welcome, including submitting issues, suggestions, documentation, tests or code. [Contribution Guide](CONTRIBUTION.md)

## License
`RuleGo` uses Apache 2.0 license, please refer to [LICENSE](LICENSE) file for details.

## Contact Us
Email: rulego@outlook.com