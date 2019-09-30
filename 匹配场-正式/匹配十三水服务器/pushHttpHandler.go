package main

import (
	"encoding/json"
	"io/ioutil"
	"logs"
	"net/http"
)

/////////////////////////////////////////////////////
//初始化http请求
func InitHttpHandle() {
	http.HandleFunc("/msgtoclient", Handler_MsgToClient)
}

func Handler_MsgToClient(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//参数获取
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		ReturnRsp(w, &PBaseMsg{
			Code: 1,
			Err:  "参数读取错误",
		})
		return
	}
	//
	arg := PMsgToHallWebMsg{}
	json.Unmarshal(body, &arg)

	var uid int64
	for k, v := range arg.Server {
		if k == "uid" {
			uid = int64(v.(float64))
		}
	}

	G_MgrPushMsg.PushPlayerInfoChange(uid, arg.Server)
	err = G_MgrPushMsg.PushMsg(uid, &arg.Client)
	if err != nil {
		ReturnRsp(w, &PBaseMsg{
			Code: 2,
			Err:  err.Error(),
		})
		return
	}

	ReturnRsp(w, &PBaseMsg{
		Code: 0,
	})
}

////////////////////////////////////////////////////////
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
