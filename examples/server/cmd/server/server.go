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

package main

import (
	"examples/server/config"
	"examples/server/config/logger"
	"examples/server/internal/router"
	"examples/server/internal/service"
	"flag"
	"fmt"
	endpointApi "github.com/yunboom/rulego/api/types/endpoint"
	"github.com/yunboom/rulego/node_pool"
	"github.com/yunboom/rulego/utils/str"
	"gopkg.in/ini.v1"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
)

const (
	version = "1.0.0"
)

var (
	//是否是查询版本
	ver bool
	//配置文件
	configFile string
)

func init() {
	flag.StringVar(&configFile, "c", "", "配置文件")
	flag.BoolVar(&ver, "v", false, "打印版本")
}

func main() {
	flag.Parse()

	if ver {
		fmt.Printf("RuleGo-Ci Server v%s", version)
		os.Exit(0)
	}

	var c config.Config
	if configFile == "" {
		c = config.DefaultConfig
	} else if cfg, err := ini.Load(configFile); err != nil {
		log.Fatal("error:", err)
	} else {
		if err := cfg.MapTo(&c); err != nil {
			log.Fatal("error:", err)
		}
		if section, err := cfg.GetSection("global"); err == nil {
			c.Global = section.KeysHash()
		}
		if section, err := cfg.GetSection("users"); err == nil {
			c.Users = section.KeysHash()
		}
	}
	config.Set(c)
	logger.Set(initLogger(c))

	//pprof
	if c.Pprof.Enable {
		addr := c.Pprof.Addr
		if addr == "" {
			addr = "0.0.0.0:6060"
		}
		log.Printf("pprof enabled, addr=%s \n", addr)
		go http.ListenAndServe(addr, nil)
	}

	log.Printf("Get Converter Info: %s \n", str.GetConverterInfo())

	//初始化用户名、密码、apiKey之间的映射
	c.InitUserMap()
	log.Printf("use config file=%s \n", configFile)

	if err := loadNodePool(c); err != nil {
		log.Fatal("loadNodePool error:", err)
	} else {
		log.Printf("loadNodePool file=%s \n", c.NodePoolFile)
	}
	//初始化rulego配置
	router.InitRulegoConfig()
	//创建http服务
	ep, err := router.NewRestServe(c)
	if err != nil {
		log.Fatal("error:", err)
	}
	//启动http服务
	if err := ep.Start(); err != nil {
		log.Fatal("error:", err)
	}
	//创建websocket服务
	if restEp, ok := ep.(endpointApi.HttpEndpoint); ok {
		wsEp, err := router.NewWebsocketServe(c, restEp)
		if err != nil {
			log.Fatal("websocket endpoint creation error:", err)
		}
		if err := wsEp.Start(); err != nil {
			log.Fatal("websocket start error:", err)
		}
	}
	//初始化服务
	if err := service.Setup(c); err != nil {
		log.Fatal("setup service error:", err)
	}

	sigs := make(chan os.Signal, 1)
	// 监听系统信号，包括中断信号和终止信号
	signal.Notify(sigs, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigs:
		if ep != nil {
			ep.Destroy()
		}
		log.Println("stopped server")
		os.Exit(0)
	}
}

// 初始化日志记录器
func initLogger(c config.Config) *log.Logger {
	if c.LogFile == "" {
		return log.New(os.Stdout, "", log.LstdFlags)
	} else {
		f, err := os.OpenFile(c.LogFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal(err)
		}
		return log.New(f, "", log.LstdFlags)
	}
}

func loadNodePool(c config.Config) error {
	file := c.NodePoolFile
	if file != "" {
		if buf, err := os.ReadFile(file); err != nil {
			return err
		} else {
			_, err = node_pool.DefaultNodePool.Load(buf)
			return err
		}
	}
	return nil
}
