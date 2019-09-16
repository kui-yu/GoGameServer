package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"logs"
	"net/http"

	"github.com/tidwall/gjson"
)

//初始化http请求
func InitHttpHandle() {
	fmt.Println("初始化gamelist接口", GCONFIG.GetGameListUrlLastInterface)
	http.HandleFunc(GCONFIG.GetGameListUrlLastInterface, Handler_InterfaceGameList) // 获得机器人信息
}

// web端口更新, 通知机器人端
func Handler_InterfaceGameList(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-type", "application/json")
	//参数获取
	defer req.Body.Close()

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

	isFilterGameId, filterGameId := false, 0
	isFilterGroupId, filterGroupId := false, 0

	if gjson.GetBytes(body, "gameId").Exists() {
		filterGameId = int(gjson.GetBytes(body, "gameId").Int())
		isFilterGameId = true
	}

	if gjson.GetBytes(body, "roomId").Exists() {
		filterGroupId = int(gjson.GetBytes(body, "roomId").Int())
		isFilterGroupId = true
	}

	sendRobots := []*RobotClient{}
	if isFilterGameId {
		for _, robot := range robotClients {
			if robot.GameId == uint32(filterGameId) {
				sendRobots = append(sendRobots, robot)
			}
		}
	} else if isFilterGroupId {
		for _, robot := range robotClients {
			if robot.GroupId == uint32(filterGroupId) {
				sendRobots = append(sendRobots, robot)
			}
		}
	} else {
		sendRobots = append(sendRobots, robotClients...)
	}

	count := len(sendRobots)

	// sendRobots 排序
	for i := 0; i < count; i++ {
		for j := i + 1; j < count; j++ {
			a, b := sendRobots[i], sendRobots[j]
			isShift := false
			if a.GroupId > b.GroupId {
				isShift = true
			} else {
				if a.GameId > b.GameId {
					isShift = true
				} else {
					if a.GradeId > b.GradeId {
						isShift = true
					}
				}
			}

			if isShift {
				sendRobots[i], sendRobots[j] = b, a
			}
		}
	}

	start := 0
	totalCount := count

	infoList := struct {
		Code     uint32     `json:"code"`
		Msg      string     `json:"msg"`
		GameList []GameInfo `json:"gamelist"`
		Count    int        `json:"count"`
	}{
		Code:     200,
		Msg:      "ok",
		GameList: []GameInfo{},
		Count:    totalCount,
	}

	if gjson.GetBytes(body, "start").Exists() {
		start = int(gjson.GetBytes(body, "start").Int())
	}

	if gjson.GetBytes(body, "count").Exists() {
		count = int(gjson.GetBytes(body, "count").Int())
	}

	pageEnd := start + count
	if pageEnd > len(sendRobots) {
		pageEnd = len(sendRobots)
	}

	for i := start; i < pageEnd; i++ {
		client := sendRobots[i]
		infoList.GameList = append(infoList.GameList, GameInfo{
			GroupId:           client.GroupId,
			GameId:            client.GameId,
			GradeId:           client.GradeId,
			Name:              client.Name,
			GetRobotConfigUrl: client.GetRobotConfigUrl,
			PutRobotConfigUrl: client.PutRobotConfigUrl,
			CheckRobotUrl:     client.CheckRobotUrl,
			OfflineRobotUrl:   client.OfflineRobotUrl,
			Forceoffroboturl:  client.Forceoffroboturl,
			Forceonroboturl:   client.Forceonroboturl,
			RobotCount:        client.RobotCount,
		})
	}

	ReturnRsp(w, &infoList)
}

//公共接口，返回http请求数据
func ReturnRsp(w http.ResponseWriter, rsp interface{}) {
	re, err := json.Marshal(rsp)
	if err != nil {
		logs.Debug("ReturnRsp error:", err)
		w.Write([]byte(err.Error()))
	} else {
		w.Write(re)
	}
}
