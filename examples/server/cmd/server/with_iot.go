//go:build with_iot

package main

import (
	// 注册扩展组件库
	// 使用`go build -tags with_iot .`把扩展组件编译到运行文件
	_ "github.com/yunboom/rulego-components-iot/endpoint/opcua"
	_ "github.com/yunboom/rulego-components-iot/external/modbus"
	_ "github.com/yunboom/rulego-components-iot/external/opcua"
)
