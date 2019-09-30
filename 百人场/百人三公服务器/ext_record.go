package main

import (
	"encoding/json"

	"time"
)

var G_AllRecord AllRecord

func init() {
	G_AllRecord.Data = make(map[int64]*[]RecordData)
}

//type RecordData struct {
//	Grade     int     //场次类型
//	YaZhu     int     //押注金额
//	YingKui   int     //盈亏
//	YzQuYu    []int   //押注区域
//	XianPai   []int32 //闲牌
//	ZhuangPai []int32 //庄牌
//	Date      int64   //时间
//	DaXiao    int     //0无此值，1大，2小
//}
type RecordData struct {
	GradeType int    //游戏场次
	AllBet    int64  //总下注
	WinCoins  int64  //输赢金币
	BetArea   []Area //押注区域
	CardType  int    //牌型
	Time      string
	Date      int64 //阶段时间
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
	rsp.Data = G_AllRecord.GetRecordByRange(p.Uid, req.Start, req.End)
	p.SendNativeMsg(MSG_GAME_INFO_GET_RECORD_REPLY, &rsp)
}
