package main

import (
	"encoding/json"
	"logs"
	"time"
)

var G_AllRecord AllRecord

func init() {
	G_AllRecord.Data = make(map[int64]*[]RecordData)
}

type RecordData struct {
	RoomName    string  //房间号
	AllBet      int64   //总投注
	WinOrLost   int64   //总盈利
	BetArea     []int64 //玩家押注区域
	TBankerCard int32   // 庄牌型
	TypeList    []int32 // 牌型（东西南北）
	EndTime     string  //结束时间
	Date        int64   //与客户端无关，可忽略
}

type AllRecord struct {
	Data     map[int64]*[]RecordData
	Interval int64
}

//添加玩家记录集
func (this *AllRecord) AddRecord(uid int64, d *RecordData) {
	v, ok := this.Data[uid]
	if !ok {
		v = &[]RecordData{}
		this.Data[uid] = v
	}
	rd := []RecordData{*d}
	*v = append(rd, (*v)...)
}

//删除玩家记录
func (this *AllRecord) DelRecordByUid(uid int64) {
	delete(this.Data, uid)
}

//获取区间记录列表
func (this *AllRecord) GetRecordByRange(uid int64, start, end int) []RecordData {
	logs.Debug("接收到用户请求游戏记录!!!!")
	this.TimerClearRecord()
	result := []RecordData{}
	if end <= start {
		return result
	}
	v, ok := this.Data[uid]
	if ok {
		if start >= len(*v) {
			return result
		}
		if end >= len(*v) {
			end = len(*v)
		}
		result = append([]RecordData{}, (*v)[start:end]...)
	}
	logs.Debug("...........222", result)
	return result
}

func (this *AllRecord) TimerClearRecord() {
	t := time.Now().Unix()
	for _, v := range this.Data {
		nrecord := []RecordData{}
		for _, v2 := range *v {
			if t-v2.Date > 86400*3 {
				continue
			}
			nrecord = append(nrecord, v2)
		}
		*v = nrecord
	}
}

///////////////////////////////////
type GGetRecordReq struct {
	Id    int32 //协议号
	Start int   //开始区间,从0开始，0表示第一个
	End   int   //结束区间
}

type GGetRecordRsp struct {
	Id   int32
	Data []RecordData
}

func (this *ExtDesk) HandleGetRecord(p *ExtPlayer, d *DkInMsg) {
	req := GGetRecordReq{}
	err := json.Unmarshal([]byte(d.Data), &req)
	if err != nil {
		return
	}
	//
	rsp := GGetRecordRsp{
		Id: MSG_GAME_INFO_GET_RECORD_REPLY,
	}
	logs.Debug("游戏记录：", rsp.Data)
	rsp.Data = G_AllRecord.GetRecordByRange(p.Uid, req.Start, req.End)
	logs.Debug(".............................1", rsp.Data)
	p.SendNativeMsg(MSG_GAME_INFO_GET_RECORD_REPLY, &rsp)
}
