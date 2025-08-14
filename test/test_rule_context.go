/*
 * Copyright 2023 The RuleGo Authors.
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

package test

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/yunboom/rulego/api/types"
	"github.com/yunboom/rulego/utils/cache"
)

var _ types.RuleContext = (*NodeTestRuleContext)(nil)

// NodeTestRuleContext
// 只为测试单节点，临时创建的上下文
// 无法把多个节点组成链式
// callback 回调处理结果
type NodeTestRuleContext struct {
	context  context.Context
	config   types.Config
	callback func(msg types.RuleMsg, relationType string, err error)
	self     types.Node
	selfId   string
	//所有子节点处理完成事件，只执行一次
	onAllNodeCompleted func()
	onEndFunc          types.OnEndFunc
	childrenNodes      sync.Map
	out                types.RuleMsg
	globalCache        types.Cache
	chainCache         types.Cache
	mutex              sync.RWMutex // Add mutex for thread safety
}

func (ctx *NodeTestRuleContext) GlobalCache() types.Cache {
	return ctx.globalCache
}

func (ctx *NodeTestRuleContext) ChainCache() types.Cache {
	return ctx.chainCache
}

func NewRuleContext(config types.Config, callback func(msg types.RuleMsg, relationType string, err error)) types.RuleContext {
	globalCache := cache.NewMemoryCache(time.Minute * 5)
	return &NodeTestRuleContext{
		context:     context.TODO(),
		config:      config,
		callback:    callback,
		globalCache: globalCache,
		chainCache:  cache.NewNamespaceCache(globalCache, "test"),
	}
}

func NewRuleContextFull(config types.Config, self types.Node, childrenNodes map[string]types.Node, callback func(msg types.RuleMsg, relationType string, err error)) types.RuleContext {
	ctx := &NodeTestRuleContext{
		config:      config,
		self:        self,
		callback:    callback,
		context:     context.TODO(),
		globalCache: config.Cache,
		chainCache:  cache.NewNamespaceCache(config.Cache, "test"),
	}
	for k, v := range childrenNodes {
		ctx.childrenNodes.Store(k, v)
	}
	return ctx
}

func (ctx *NodeTestRuleContext) TellSuccess(msg types.RuleMsg) {
	ctx.mutex.RLock()
	callback := ctx.callback
	onEndFunc := ctx.onEndFunc
	ctx.mutex.RUnlock()

	if callback != nil {
		callback(msg, types.Success, nil)
	}
	if onEndFunc != nil {
		onEndFunc(ctx, msg, nil, types.Success)
	}
}

func (ctx *NodeTestRuleContext) TellFailure(msg types.RuleMsg, err error) {
	ctx.mutex.RLock()
	callback := ctx.callback
	onEndFunc := ctx.onEndFunc
	ctx.mutex.RUnlock()

	if callback != nil {
		callback(msg, types.Failure, err)
	}
	if onEndFunc != nil {
		onEndFunc(ctx, msg, err, types.Failure)
	}
}

func (ctx *NodeTestRuleContext) TellNext(msg types.RuleMsg, relationTypes ...string) {
	ctx.mutex.RLock()
	callback := ctx.callback
	onEndFunc := ctx.onEndFunc
	ctx.mutex.RUnlock()

	for _, relationType := range relationTypes {
		if callback != nil {
			callback(msg, relationType, nil)
		}
		if onEndFunc != nil {
			onEndFunc(ctx, msg, nil, relationType)
		}
	}
}

func (ctx *NodeTestRuleContext) TellSelf(msg types.RuleMsg, delayMs int64) {
	time.AfterFunc(time.Millisecond*time.Duration(delayMs), func() {
		if ctx.self != nil {
			ctx.self.OnMsg(ctx, msg)
		}
	})
}
func (ctx *NodeTestRuleContext) TellNextOrElse(msg types.RuleMsg, defaultRelationType string, relationTypes ...string) {
	ctx.TellNext(msg, relationTypes...)
}
func (ctx *NodeTestRuleContext) NewMsg(msgType string, metaData *types.Metadata, data string) types.RuleMsg {
	return types.NewMsg(0, msgType, types.JSON, metaData, data)
}
func (ctx *NodeTestRuleContext) GetSelfId() string {
	ctx.mutex.RLock()
	defer ctx.mutex.RUnlock()
	return ctx.selfId
}
func (ctx *NodeTestRuleContext) Self() types.NodeCtx {
	return nil
}

func (ctx *NodeTestRuleContext) From() types.NodeCtx {
	return nil
}
func (ctx *NodeTestRuleContext) RuleChain() types.NodeCtx {
	return nil
}
func (ctx *NodeTestRuleContext) Config() types.Config {
	return ctx.config
}
func (ctx *NodeTestRuleContext) SubmitTack(task func()) {
	ctx.SubmitTask(task)
}
func (ctx *NodeTestRuleContext) SubmitTask(task func()) {
	go task()
}

func (ctx *NodeTestRuleContext) SetEndFunc(onEndFunc types.OnEndFunc) types.RuleContext {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	ctx.onEndFunc = onEndFunc
	return ctx
}

func (ctx *NodeTestRuleContext) GetEndFunc() types.OnEndFunc {
	ctx.mutex.RLock()
	defer ctx.mutex.RUnlock()
	return ctx.onEndFunc
}

func (ctx *NodeTestRuleContext) SetContext(c context.Context) types.RuleContext {
	ctx.context = c
	return ctx
}

func (ctx *NodeTestRuleContext) GetContext() context.Context {
	return ctx.context
}

func (ctx *NodeTestRuleContext) TellFlow(chainId string, msg types.RuleMsg, opts ...types.RuleContextOption) {
	for _, opt := range opts {
		opt(ctx)
	}
	if chainId == "" {
		if ctx.onEndFunc != nil {
			ctx.onEndFunc(ctx, msg, errors.New("chainId can not nil"), types.Failure)
		}

	} else if chainId == "notfound" {
		if ctx.onEndFunc != nil {
			ctx.onEndFunc(ctx, msg, fmt.Errorf("ruleChain id=%s not found", chainId), types.Failure)
		}
		if ctx.onAllNodeCompleted != nil {
			ctx.onAllNodeCompleted()
		}
	} else if chainId == "toTrue" {
		if ctx.onEndFunc != nil {
			ctx.onEndFunc(ctx, msg, nil, types.True)
		}
		if ctx.onAllNodeCompleted != nil {
			ctx.onAllNodeCompleted()
		}
	} else {
		if ctx.onEndFunc != nil {
			ctx.onEndFunc(ctx, msg, nil, types.Success)
		}
		if ctx.onAllNodeCompleted != nil {
			ctx.onAllNodeCompleted()
		}
	}
}

// TellNode 独立执行某个节点，通过callback获取节点执行情况，用于节点分组类节点控制执行某个节点
func (ctx *NodeTestRuleContext) TellNode(context context.Context, nodeId string, msg types.RuleMsg, skipTellNext bool, callback types.OnEndFunc, onAllNodeCompleted func()) {
	if v, ok := ctx.childrenNodes.Load(nodeId); ok {
		// 线程安全地设置 selfId
		ctx.mutex.Lock()
		ctx.selfId = nodeId
		ctx.mutex.Unlock()

		subCtx := NewRuleContext(ctx.config, func(msg types.RuleMsg, relationType string, err error) {
			if callback != nil {
				callback(ctx, msg, err, relationType)
			}

			if onAllNodeCompleted != nil {
				onAllNodeCompleted()
			}
		})

		v.(types.Node).OnMsg(subCtx, msg)
	} else {
		if callback != nil {
			callback(ctx, msg, fmt.Errorf("node id=%s not found", nodeId), types.Failure)
		}
		if onAllNodeCompleted != nil {
			onAllNodeCompleted()
		}
	}
}

// TellChainNode 独立执行某个节点，通过callback获取节点执行情况，用于节点分组类节点控制执行某个节点
func (ctx *NodeTestRuleContext) TellChainNode(context context.Context, chainId string, nodeId string, msg types.RuleMsg, skipTellNext bool, callback types.OnEndFunc, onAllNodeCompleted func()) {
	ctx.TellNode(context, nodeId, msg, skipTellNext, callback, onAllNodeCompleted)
}

// SetOnAllNodeCompleted 设置所有节点执行完回调
func (ctx *NodeTestRuleContext) SetOnAllNodeCompleted(onAllNodeCompleted func()) {
	ctx.onAllNodeCompleted = onAllNodeCompleted
}

func (ctx *NodeTestRuleContext) DoOnEnd(msg types.RuleMsg, err error, relationType string) {

}

// SetCallbackFunc 设置回调函数
func (ctx *NodeTestRuleContext) SetCallbackFunc(functionName string, f interface{}) {

}

// GetCallbackFunc 获取回调函数
func (ctx *NodeTestRuleContext) GetCallbackFunc(functionName string) interface{} {
	return nil
}

// OnDebug 调用配置的OnDebug回调函数
func (ctx *NodeTestRuleContext) OnDebug(ruleChainId string, flowType string, nodeId string, msg types.RuleMsg, relationType string, err error) {
}

func (ctx *NodeTestRuleContext) SetExecuteNode(nodeId string, relationTypes ...string) {

}
func (ctx *NodeTestRuleContext) TellCollect(msg types.RuleMsg, callback func(msgList []types.WrapperMsg)) bool {
	callback(nil)
	return true
}

func (ctx *NodeTestRuleContext) GetOut() types.RuleMsg {
	ctx.mutex.RLock()
	defer ctx.mutex.RUnlock()
	return ctx.out
}

// setOut safely sets the out field
func (ctx *NodeTestRuleContext) setOut(msg types.RuleMsg) {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	ctx.out = msg
}

func (ctx *NodeTestRuleContext) GetErr() error {
	return nil
}

func (ctx *NodeTestRuleContext) TellStream(msg types.RuleMsg) {
	ctx.TellNext(msg, types.Stream)
}

// GetEnv 获取环境变量和元数据
func (ctx *NodeTestRuleContext) GetEnv(msg types.RuleMsg, useMetadata bool) map[string]interface{} {
	// 创建环境变量map
	envVars := make(map[string]interface{})

	// 设置基础环境变量
	envVars["id"] = msg.GetId()
	envVars["ts"] = msg.GetTs()
	envVars["data"] = msg.GetData()
	envVars["msgType"] = msg.GetType()
	envVars["type"] = msg.GetType()
	envVars["dataType"] = string(msg.GetDataType())
	// 使用 GetJsonData() 避免重复JSON解析
	if msg.DataType == types.JSON {
		if jsonData, err := msg.GetJsonData(); err == nil {
			envVars[types.MsgKey] = jsonData
		} else {
			// 解析失败，使用原始数据
			envVars[types.MsgKey] = msg.GetData()
		}
	} else {
		// 如果不是 JSON 类型，直接使用原始数据
		envVars[types.MsgKey] = msg.GetData()
	}
	// 优化 metadata 处理
	if msg.Metadata != nil {
		if useMetadata {
			// 遍历metadata，将键值对添加到环境变量中 - use zero-copy ForEach
			msg.Metadata.ForEach(func(k, v string) bool {
				envVars[k] = v
				return true // continue iteration
			})
		}
		envVars[types.MetadataKey] = msg.Metadata.Values()
	}

	return envVars
}
