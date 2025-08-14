//go:build with_extend

package main

import (
	// 注册扩展组件库
	// 使用`go build -tags with_extend .`把扩展组件编译到运行文件
	_ "github.com/yunboom/rulego-components/endpoint/beanstalkd"
	_ "github.com/yunboom/rulego-components/endpoint/grpc_stream"
	_ "github.com/yunboom/rulego-components/endpoint/kafka"
	_ "github.com/yunboom/rulego-components/endpoint/nats"
	_ "github.com/yunboom/rulego-components/endpoint/nsq"
	_ "github.com/yunboom/rulego-components/endpoint/pulsar"
	_ "github.com/yunboom/rulego-components/endpoint/rabbitmq"
	_ "github.com/yunboom/rulego-components/endpoint/redis"
	_ "github.com/yunboom/rulego-components/endpoint/redis_stream"
	_ "github.com/yunboom/rulego-components/endpoint/wukongim"
	_ "github.com/yunboom/rulego-components/external/beanstalkd"
	_ "github.com/yunboom/rulego-components/external/grpc" //编译后文件大约增加7M
	_ "github.com/yunboom/rulego-components/external/kafka"
	_ "github.com/yunboom/rulego-components/external/mongodb"
	_ "github.com/yunboom/rulego-components/external/nats"
	_ "github.com/yunboom/rulego-components/external/nsq"
	_ "github.com/yunboom/rulego-components/external/opengemini"
	_ "github.com/yunboom/rulego-components/external/otel"
	_ "github.com/yunboom/rulego-components/external/pulsar"
	_ "github.com/yunboom/rulego-components/external/rabbitmq"
	_ "github.com/yunboom/rulego-components/external/redis"
	_ "github.com/yunboom/rulego-components/external/wukongim"
	_ "github.com/yunboom/rulego-components/filter"
	_ "github.com/yunboom/rulego-components/stats/streamsql"
	_ "github.com/yunboom/rulego-components/transform"
)
