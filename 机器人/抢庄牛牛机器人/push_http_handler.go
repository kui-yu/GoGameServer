package main

import (
	"encoding/json"
	"io/ioutil"
	"strconv"

	"net/http"

	"github.com/tidwall/gjson"
)

// 机器人数量缓存
var robotNumCache int = 0

//初始化http请求
func InitHttpHandle() {
	http.HandleFunc("/"+gameConfig.OpenWebInterface.GetRobotConfigUrl, Handler_GetRobotConfig)  // 获得机器人配置
	http.HandleFunc("/"+gameConfig.OpenWebInterface.PutRobotConfigUrl, Handler_PutRobotConfig)  // 设置机器人配置
	http.HandleFunc("/"+gameConfig.OpenWebInterface.CheckRobotUrl, Handler_CheckRobot)          // 检查机器人配置
	http.HandleFunc("/"+gameConfig.OpenWebInterface.OfflineRobotUrl, Handler_OfflineRobotUrl)   // 设置选择的机器人
	http.HandleFunc("/"+gameConfig.OpenWebInterface.Forceoffroboturl, Handler_Forceoffroboturl) // 设置选择的机器人
	http.HandleFunc("/"+gameConfig.OpenWebInterface.Forceonroboturl, Handler_Forceonroboturl)   // 设置选择的机器人
}

func Handler_GetRobotConfig(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-type", "application/json")
	//参数获取
	defer req.Body.Close()

	ReturnRsp(w, &struct {
		Code int         `json:"code"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
	}{
		Code: 200,
		Msg:  "ok",
		Data: OnGetRobotConfig(),
	})
}

func Handler_PutRobotConfig(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-type", "application/json")
	//参数获取
	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)

	DebugLog("接收到后台提交配置 ", string(body))

	if err != nil {
		ReturnRsp(w, &struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
		}{
			Code: 255,
			Msg:  "解析错误",
		})
		return
	}

	token := req.Header.Get("token")
	DebugLog("提交配置参数 token:", token)
	if token != GRobotServer.Info.BgRobotToken {
		ReturnRsp(w, &struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
		}{
			Code: 601,
			Msg:  "error token",
		})
		return
	}

	for i, item := range gameConfig.RobotConfigItems {
		if item.Type == "input" {
			array := gjson.ParseBytes(body).Array()
			for _, gitem := range array {
				if gitem.Get("name").String() == item.Name {
					// 向控制器发送配置改变事件
					oldValue := item.Value
					newValue := gitem.Get("value").String()
					if oldValue != newValue {
						item.Value = newValue
						gameConfig.RobotConfigItems[i] = item
					} else {
						break
					}

					controller.sendEvent(EVENT_CONT_CONFIG_CHANGE, struct {
						Key      string
						OldValue string
						NewValue string
					}{
						Key:      item.Name,
						OldValue: oldValue,
						NewValue: newValue,
					})
					break
				}
			}

		}
	}

	SaveConf()

	ReturnRsp(w, &struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}{
		Code: 200,
		Msg:  "ok",
	})
}

func Handler_OfflineRobotUrl(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-type", "application/json")
	//参数获取
	defer req.Body.Close()

	token := req.Header.Get("token")
	if token != GRobotServer.Info.BgRobotToken {
		ReturnRsp(w, &struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
		}{
			Code: 601,
			Msg:  "error token",
		})
		return
	}

	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		ReturnRsp(w, &struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
		}{
			Code: 255,
			Msg:  "解析错误",
		})
		return
	}

	noexistRobot := []string{}
	jsonItems := gjson.ParseBytes(body).Array()
	for _, v := range jsonItems {
		if controller.findRobotFromToken(v.String()) == nil {
			noexistRobot = append(noexistRobot, v.String())
		} else {
			controller.sendEvent(EVENT_CONT_OFFLINEROBOT, v.String())
		}
	}

	ReturnRsp(w, &struct {
		Code int      `json:"code"`
		Msg  string   `json:"msg"`
		Data []string `json:"data"`
	}{
		Code: 200,
		Msg:  "ok",
		Data: noexistRobot,
	})
}

// 查询机器人并返回结果
func Handler_CheckRobot(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-type", "application/json")
	//参数获取
	defer req.Body.Close()

	token := req.Header.Get("token")
	if token != GRobotServer.Info.BgRobotToken {
		ReturnRsp(w, &struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
		}{
			Code: 601,
			Msg:  "error token",
		})
		return
	}

	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		ReturnRsp(w, &struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
		}{
			Code: 255,
			Msg:  "解析错误",
		})
		return
	}

	ReturnRsp(w, &struct {
		Code int         `json:"code"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
	}{
		Code: 200,
		Msg:  "ok",
		Data: OnCheckRobotData(string(body)),
	})
}

// 机器人强制下线
func Handler_Forceoffroboturl(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-type", "application/json")

	//参数获取
	defer req.Body.Close()

	token := req.Header.Get("token")
	if token != GRobotServer.Info.BgRobotToken {
		ReturnRsp(w, &struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
		}{
			Code: 601,
			Msg:  "error token",
		})
		return
	}

	num := gameConfig.getGameConfigInt("num")
	if num == 0 {
		ReturnRsp(w, &struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
		}{
			Code: 200,
			Msg:  "ok",
		})
		return
	}

	robotNumCache = num

	// 通知机器人管理，当前数量变化
	GRobotServer.AddMsgNative(MSG_GAME_ROBOT_NUMCHANGE, struct {
		Id         int
		RobotCount int
	}{
		Id:         MSG_GAME_ROBOT_NUMCHANGE,
		RobotCount: 0,
	})

	for i, item := range gameConfig.RobotConfigItems {
		if item.Name == "num" {
			item.Value = "0"
			gameConfig.RobotConfigItems[i] = item
			SaveConf()

			robots := []*ExtRobotClient{}
			robots = append(robots, controller.getRobotClients()...)

			for _, v := range robots {
				controller.sendEvent(EVENT_CONT_ROBOTSHIFT, v)
			}
			break
		}
	}

	ReturnRsp(w, &struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}{
		Code: 200,
		Msg:  "ok",
	})
}

// 恢复强制下线的机器人
func Handler_Forceonroboturl(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-type", "application/json")

	//参数获取
	defer req.Body.Close()

	token := req.Header.Get("token")
	if token != GRobotServer.Info.BgRobotToken {
		ReturnRsp(w, &struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
		}{
			Code: 601,
			Msg:  "error token",
		})
		return
	}

	num := gameConfig.getGameConfigInt("num")
	if num != 0 {
		ReturnRsp(w, &struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
		}{
			Code: 601,
			Msg:  "机器人已在线",
		})
		return
	}

	if robotNumCache == 0 {
		ReturnRsp(w, &struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
		}{
			Code: 602,
			Msg:  "没有设置过机器人数量",
		})
		return
	}

	// 通知机器人管理，当前数量变化
	GRobotServer.AddMsgNative(MSG_GAME_ROBOT_NUMCHANGE, struct {
		Id         int
		RobotCount int
	}{
		Id:         MSG_GAME_ROBOT_NUMCHANGE,
		RobotCount: robotNumCache,
	})

	for i, item := range gameConfig.RobotConfigItems {
		if item.Name == "num" {
			item.Value = strconv.Itoa(robotNumCache)
			gameConfig.RobotConfigItems[i] = item
			SaveConf()

			robots := []*ExtRobotClient{}
			robots = append(robots, controller.getRobotClients()...)

			controller.sendEvent(EVENT_CONT_ADDROBOT, robotNumCache)
			break
		}
	}

	robotNumCache = 0

	ReturnRsp(w, &struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}{
		Code: 200,
		Msg:  "ok",
	})
}

//公共接口，返回http请求数据
func ReturnRsp(w http.ResponseWriter, rsp interface{}) {
	re, err := json.Marshal(rsp)
	if err != nil {
		ErrorLog("ReturnRsp error:", err)
		w.Write([]byte(err.Error()))
	} else {
		w.Write(re)
	}
}
