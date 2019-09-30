/**
* 玩家操作
**/
package main

import (
	"encoding/json"
)

type FsmUserOperate struct {
	UpMark   int
	Mark     int
	RC       *ExtRobotClient
	Overtime int
}

func (this *FsmUserOperate) InitFSM(mark int, rc *ExtRobotClient) {
	this.Mark = mark
	this.RC = rc
}

func (this *FsmUserOperate) GetMark() int {
	return this.Mark
}

func (this *FsmUserOperate) Run(upMark int, overtime int) {
	DebugLog("进入游戏状态：玩家操作")
	this.UpMark = upMark
	this.Overtime = overtime
	this.addListener() // 添加监听
}

func (this *FsmUserOperate) Leave() {
	this.removeListener()
}

func (this *FsmUserOperate) onEvent(interface{}) {

}

// 添加网络监听
func (this *FsmUserOperate) addListener() {
	this.RC.Handle[MSG_GAME_NGameOperate] = this.onGameOperate
	this.RC.Handle[MSG_GAME_RGameOperate] = this.onRetGameOperate
}

// 删除网络监听
func (this *FsmUserOperate) removeListener() {
	delete(this.RC.Handle, MSG_GAME_NGameOperate)
}

type UserOperateData struct {
	Sid         int
	OperateAuth int
	MinCoin     int64
	MaxCoin     int64
}

func (this *FsmUserOperate) onGameOperate(str string) {
	data := UserOperateData{}
	json.Unmarshal([]byte(str), &data)

	if data.Sid != this.RC.Sid {
		return
	}
	DebugLog(str)

	var time = this.Overtime / 1000
	if time < 1 {
		this.onRepOperate(0, data)
	} else {
		t, _ := GetRandomNum(1, int(time))
		this.RC.TimeTicker.AddTimer(t, this.onRepOperate, data)
	}
}

func (this *FsmUserOperate) onRepOperate(id int, d interface{}) {
	data := d.(UserOperateData)
	type OperateOdds struct {
		Mark int
		Min  int
		Max  int
	}
	operates := []OperateOdds{}

	if (data.OperateAuth & OperateAuthQP) != 0 { //弃牌
		operates = append(operates, OperateOdds{
			Mark: OperateAuthQP,
			Min:  0,
			Max:  15,
		})
	}
	if (data.OperateAuth & OperateAuthSH) != 0 { //allin
		operates = append(operates, OperateOdds{
			Mark: OperateAuthSH,
			Min:  16,
			Max:  30,
		})
	}
	if (data.OperateAuth & OperateAuthKP) != 0 { //开牌
		operates = append(operates, OperateOdds{
			Mark: OperateAuthKP,
			Min:  31,
			Max:  50,
		})
	}
	if (data.OperateAuth & OperateAuthGZ) != 0 { //跟注
		operates = append(operates, OperateOdds{
			Mark: OperateAuthGZ,
			Min:  51,
			Max:  80,
		})
	}
	if (data.OperateAuth & OperateAuthJZ) != 0 { //加注
		operates = append(operates, OperateOdds{
			Mark: OperateAuthJZ,
			Min:  81,
			Max:  99,
		})
	}

	if len(operates) != 0 {
		curOperate := -1
		for {
			odds, _ := GetRandomNum(0, 100)
			for _, v := range operates {
				if odds >= v.Min && odds <= v.Max {
					curOperate = v.Mark
					break
				}
			}
			if curOperate != -1 {
				break
			}
		}

		// 发送操作
		var value int64 = 0
		//加注的值
		if curOperate == OperateAuthJZ {
			add := (data.MaxCoin - data.MinCoin) / 11
			if add < 2 {
				value = data.MaxCoin
			} else {
				r, _ := GetRandomNum64(1, add)
				value = data.MinCoin + int64(r*11)
			}

			if value == data.MinCoin {
				value = value + 1
			}
		}

		if this.RC.CarryCoin < value {
			value = 0
			curOperate = OperateAuthSH
		} else if value > data.MaxCoin {
			value = data.MaxCoin
		}

		this.RC.AddMsgNative(MSG_GAME_QGameOperate, struct {
			Id          int32
			OperateAuth int
			RaiseValue  int64
		}{
			Id:          MSG_GAME_QGameOperate,
			OperateAuth: curOperate,
			RaiseValue:  value,
		})
	}
}

func (this *FsmUserOperate) onRetGameOperate(str string) {
	DebugLog("接收到操作结果 %s", str)
}
