package main

import (
	"github.com/tidwall/gjson"
)

// 获取游戏配置信息，自己可以附加一些统计数据作后台显示
func OnGetRobotConfig() interface{} {
	return gameConfig.RobotConfigItems
}

// 检查机器人是否在游戏中
func OnCheckRobotData(jsonStr string) interface{} {
	jsonItems := gjson.Parse(jsonStr).Array()
	itemLen := len(jsonItems)
	 
	// 查询数据并删除已存在的机器人
	var itemArray []CheckRobotItemInfo = []CheckRobotItemInfo{}
	for i := 0; i < itemLen; i++ {
		token := jsonItems[i].Get("token").String()
		if controller.findRobotFromToken(token) == nil {
			iteminfo := CheckRobotItemInfo{}

			err := gjson.Unmarshal([]byte(jsonItems[i].String()), &iteminfo)
			if err != nil {
				ErrorLog("查询机器人出错", err.Error())
			}
			itemArray = append(itemArray, iteminfo)
		}
	}

	return itemArray
}
