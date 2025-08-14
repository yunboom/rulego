/*
 * Copyright 2024 The RuleGo Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package main demonstrates how to create a NET endpoint server
// that processes binary and JSON data using JavaScript transform node.
//
// This example shows:
// - Setting up a NET endpoint server with DSL configuration
// - Processing both binary and JSON data types
// - Using JavaScript transform node for data processing
// - Responding to client messages
package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yunboom/rulego"

	"github.com/yunboom/rulego/api/types"
	"github.com/yunboom/rulego/engine"
)

const (
	// Chain ID for the data processor
	// 数据处理器的链ID
	chainId = "dataProcessor"
)

func main() {
	// Load rule chain DSL configuration with embedded endpoint
	// 加载包含嵌入式端点的规则链DSL配置
	chainDSL, err := os.ReadFile("examples/net_endpoint_example/server/chain_dsl.json")
	if err != nil {
		// Try current directory if server/ doesn't exist
		chainDSL, err = os.ReadFile("chain_dsl.json")
		if err != nil {
			log.Fatalf("Failed to read chain DSL: %v", err)
		}
	}

	// Create rule engine configuration with debug callback
	// 创建带调试回调的规则引擎配置
	config := rulego.NewConfig(
		types.WithDefaultPool(),
		types.WithOnDebug(func(chainId, flowType string, nodeId string, msg types.RuleMsg, relationType string, err error) {
			errStr := ""
			if err != nil {
				errStr = fmt.Sprintf(", 错误: %v", err)
			}
			metadataStr := ""
			if msg.Metadata != nil && len(msg.Metadata.Values()) > 0 {
				metadataStr = fmt.Sprintf(", 元数据: %+v", msg.Metadata.Values())
			}
			if msg.GetDataType() == types.BINARY {
				fmt.Printf("[节点处理调试] 链: %s, 类型: %s, 节点: %s, 关系: %s, 消息类型: %s, 数据类型: %s, 数据: %s%s%s\n",
					chainId, flowType, nodeId, relationType, msg.Type, msg.DataType, hex.EncodeToString(msg.GetBytes()), errStr, metadataStr)
			} else {
				fmt.Printf("[节点处理调试] 链: %s, 类型: %s, 节点: %s, 关系: %s, 消息类型: %s, 数据类型: %s, 数据: %s%s%s\n",
					chainId, flowType, nodeId, relationType, msg.Type, msg.DataType, msg.GetData(), errStr, metadataStr)
			}

		}),
	)

	// Create rule chain with embedded endpoints from DSL
	// 从DSL创建包含嵌入式端点的规则链
	ruleEngine, err := rulego.New(chainId, chainDSL, engine.WithConfig(config))
	if err != nil {
		log.Fatalf("Failed to create rule chain with endpoints: %v", err)
	}

	fmt.Println("NET endpoint server started with embedded DSL configuration")
	fmt.Println("调试模式已启用，将显示详细的节点处理信息")
	fmt.Println("Server listening on :8090")
	fmt.Println("Supported data types:")
	fmt.Println("  - JSON data (messages starting with '{')")
	fmt.Println("  - Binary data (all other messages)")
	fmt.Println("Processing features:")
	fmt.Println("  - Temperature alerts for values > 30°C")
	fmt.Println("  - Humidity monitoring")
	fmt.Println("  - Battery level warnings")
	fmt.Println("  - Binary command parsing")
	fmt.Println("Press Ctrl+C to stop the server")

	// Wait for interrupt signal to gracefully shutdown
	// 等待中断信号以优雅关闭
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("\nShutting down server...")
	ruleEngine.Stop(context.Background())
	fmt.Println("Server stopped")
}
