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
	"fmt"
	"github.com/yunboom/rulego"
	"github.com/yunboom/rulego/api/types"
	"log"
	"strconv"
	"strings"
	"time"
)

// 数据经过js转换后，增加deviceId变量，然后往主题 topic: /device/msg/${deviceId}发送处理后msg数据
// 其中${deviceId}为元数据的变量
func main() {

	config := rulego.NewConfig()

	metaData := types.NewMetadata()
	metaData.PutValue("productType", "test01")

	//js处理后，并调用http推送
	ruleEngine, err := rulego.New("rule01", []byte(chainJsonFile), rulego.WithConfig(config))
	if err != nil {
		log.Fatal(err)
	}

	var i = 1
	for i <= 5 {
		go func(index int) {
			msg := types.NewMsg(0, "TEST_MSG_TYPE1", types.JSON, metaData, "{\"temperature\":"+strconv.Itoa(index)+"}")
			ruleEngine.OnMsg(msg, types.WithEndFunc(func(ctx types.RuleContext, msg types.RuleMsg, err error) {
				fmt.Println("msg处理结果=====")
				//得到规则链处理结果
				fmt.Println(msg, err)
			}))
		}(i)

		i++
	}

	time.Sleep(time.Second * 1)

	//更新规则链节点配置，mqtt连接错误
	updateChain := strings.Replace(chainJsonFile, "127.0.0.1:1883", "127.0.0.1:1885", -1)

	err = ruleEngine.ReloadSelf([]byte(updateChain), rulego.WithConfig(config))

	//更新失败
	if err != nil {
		fmt.Println(err)
	}

	//继续使用之前的规则链发送
	for i <= 10 {
		go func(index int) {
			msg := types.NewMsg(0, "TEST_MSG_TYPE1", types.JSON, metaData, "{\"temperature\":"+strconv.Itoa(index)+"}")
			ruleEngine.OnMsg(msg, types.WithEndFunc(func(ctx types.RuleContext, msg types.RuleMsg, err error) {
				fmt.Println("msg处理结果=====")
				//得到规则链处理结果
				fmt.Println(msg, err)
			}))
		}(i)

		i++
	}

	time.Sleep(time.Second * 2)
}

var chainJsonFile = `
{
  "ruleChain": {
	"id":"rule01",
    "name": "测试规则链",
    "root": true
  },
  "metadata": {
    "nodes": [
       {
        "id": "s1",
        "type": "jsTransform",
        "name": "转换",
        "configuration": {
          "jsScript": "metadata['name']='test02';\n metadata['deviceId']='id01';\n msg['addField']='addValue2'; return {'msg':msg,'metadata':metadata,'msgType':msgType};"
        }
      },
      {
        "id": "s2",
        "type": "mqttClient",
        "name": "往mqtt Broker 推送数据",
        "configuration": {
          "server": "127.0.0.1:1883",
          "topic": "/device/msg/${deviceId}"
        }
      }
    ],
    "connections": [
      {
        "fromId": "s1",
        "toId": "s2",
        "type": "Success"
      }
    ]
  }
}
`
