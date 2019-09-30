package main

import (
	"encoding/json"
	"fmt"

	// "logs"
	"logs"
	"math/rand"
)

func (this *ExtDesk) HandleGameAuto(p *ExtPlayer, d *DkInMsg) {
	logs.Debug("所有玩家：", this.Players)
	logs.Debug("房间主人ID", this.FkOwner)
	//如果这个桌子已经开通了，那么用户可以直接进去，无需判断金币是否足够
	//如果该桌子不是空闲状态，则提示不能加入 并发送错误原因
	if this.GameState != GAME_STATUS_FREE {
		p.SendNativeMsg(MSG_GAME_FK_JOIN_REPLY, &GFkJoinReply{
			Id:     MSG_GAME_FK_JOIN_REPLY, //玩家加入游戏响应
			Result: 1,                      //0 为成功，其余失败。
			Err:    "房间暂时无法进入，可能已经开始游戏",
		})
		this.DeskMgr.LeaveDo(p.Uid) //从管理器中删除该玩家
		return                      //结束该方法
	}
	//超过人不让进
	if len(this.Players) >= this.TableConfig.PlayerNum && this.TableConfig.PlayerNum != 0 {
		p.SendNativeMsg(MSG_GAME_FK_JOIN_REPLY, &GFkJoinReply{
			Id:     MSG_GAME_FK_JOIN_REPLY,
			Result: 1,
			Err:    "房间人数已经满啦",
		})
		this.DeskMgr.LeaveDo(p.Uid) //从管理器中删除该玩家
		return                      //结束该方法
	}
	//按顺序设置 chairid(0,1,2)  （椅子ID）
	for i := 0; i < this.TableConfig.PlayerNum; i++ {
		exist := false
		for _, v := range this.Players {
			if v.ChairId == int32(i) {
				exist = true
				break
			}
		}
		if !exist { //如果玩家列表没有存在与i重复的 chairid ，就给这位玩家添加上去
			//加入玩家赋予椅子
			p.ChairId = int32(i)
			break
		}
	}
	//添加该玩家 到桌子的所有玩家切片中
	this.Players = append(this.Players, p)
	//当ChairId==0 的时候，证明进入这个房间的是第一个人，也就是房主，房主进入的时候，房间才开始初始化
	if p.ChairId == 0 {
		logs.Debug("房间号码:", this.FkNo)
		//房主进入，房间配置
		if p.Uid == this.FkOwner {
			//底分，局数配置，携带金币，牌局人护士，玩法模式，金币消耗
			this.TableConfig = GATableConfig{} //初始化
			err := json.Unmarshal([]byte(this.FkInfo), &this.TableConfig)
			logs.Debug("炸弹配置：", this.TableConfig.Boom)
			fmt.Println("conselect:", this.TableConfig.CanSelect)
			fmt.Println("局数配置:", this.TableConfig.MatchCount)
			if err != nil {
				//解析失败
				logs.Debug("创建房间 解析失败")
				// this.ResetTable()
				// this.TimerOver(nil)
				this.ReSet()
				return
			}
			//如果是积分专场，那么初始化 其BaseCore
			this.TableConfig.BaseScore = 1
			//判断房主带的房卡是否足够
			needCards := this.getPayMoney()
			if needCards > p.RoomCard {
				p.SendNativeMsg(MSG_GAME_FK_JOIN_REPLY, &GFkJoinReply{
					Id:     MSG_GAME_FK_JOIN_REPLY,
					Result: 1,
					Err:    "携带房卡不足" + string(needCards) + ",请充值!",
				})
				// this.ResetTable()
				// this.TimerOver(nil)
				this.ReSet()
				return
			}
		}
	}
	//如果 以上条件通过，那么发送房卡匹配成功
	p.SendNativeMsg(MSG_GAME_FK_JOIN_REPLY, &GFkJoinReply{
		Id:     MSG_GAME_FK_JOIN_REPLY,
		Result: 0,
	})
	//发送房间信息
	if this.JuHao == "" {
		this.JuHao = GetJuHao()
	}
	gtablejinfoReplay := GTableInfoReply{
		Id:      MSG_GAME_INFO_ROOM_NOTIFY,
		TableId: this.FkNo,
		Config:  this.TableConfig,
	}
	p.SendNativeMsg(MSG_GAME_INFO_ROOM_NOTIFY, &gtablejinfoReplay)
	stage := GStageInfo{
		Id:        MSG_GAME_INFO_STAGE,   //状态Id
		Stage:     int32(this.GameState), //状态
		StageTime: 0,                     //状态时间
	}
	p.SendNativeMsg(MSG_GAME_INFO_STAGE, &stage)
	//底分赋值
	this.Bscore = int32(this.TableConfig.BaseScore)
	fmt.Println("桌子底分:", this.Bscore)
	this.MaxDouble = int32(G_DbGetGameServerData.MaxTimes)
	this.ClearTimer()
	//发送玩家信息
	for _, v := range this.Players {
		gameReply := GInfoAutoGameReply{
			Id: MSG_GAME_INFO_AUTO_REPLY,
		}
		p.Coins = 0
		for _, p := range this.Players {
			seat := GSeatInfo{
				Uid:  p.Uid,
				Nick: p.Nick,
				Cid:  p.ChairId,
				Sex:  p.Sex,
				Head: p.Head,
				Lv:   p.Lv,
				Coin: p.Coins,
			}
			if p.isReady == 0 {
				seat.Ready = false
			} else {
				seat.Ready = true
				// fmt.Println("玩家:", v.Nick, "已经准备就绪！！！！！！！！！！！")
			}
			if p.Uid != v.Uid { //不是自己的玩家 名字隐藏
				seat.Nick = "***" + seat.Nick[len(seat.Nick)-4:]
			}
			gameReply.Seat = append(gameReply.Seat, seat)
		}
		v.SendNativeMsg(MSG_GAME_INFO_AUTO_REPLY, &gameReply) //返回玩家椅子信息
	}
	// //判断人员是否已满，开启游戏
	// if len(this.Players) >= GCONFIG.PlayerNum {
	// 	logs.Debug("第", this.Round, "把游戏开始！！！！")
	// }
}

func (this *ExtDesk) TimerSendCard(d interface{}) {
	this.CallOrGet = 1
	logs.Debug("正在发牌，该玩家选择的游戏模式是：", this.TableConfig.GameType)
	this.CurCid = int32(rand.Intn(10) % len(this.Players)) //随机一个玩家抢地主阶段是第一个
	//判断房间中是什么模式
	if this.TableConfig.GameType == 4 {
		//洗牌
		this.Lz_CardMgr.shuffle()
		//发牌17张
		for _, v := range this.Players {
			v.HandCard = []byte{}
			hd := this.Lz_CardMgr.sendCard(17)
			v.SetHandCard(Sort(hd))
		}
		//底牌
		this.DiPai = []byte{}
		this.DiPai = this.Lz_CardMgr.sendCard(3)
		//发送牌消息
		sd := GGameSendCardNotify{
			Id:    MSG_GAME_INFO_SEND_NOTIFY,
			Rount: this.Round,
			Cid:   this.CurCid,
			Lz:    this.Lz_CardMgr.Lz_Lz,
		}
		for _, v := range this.Players {
			for _, v1 := range v.HandCard {
				sd.HandsCards = append(sd.HandsCards, int(v1))
			}
			sd.HandsCards = []int{}
			v.SendNativeMsg(MSG_GAME_INFO_SEND_NOTIFY, &sd)
		}
	} else if this.TableConfig.GameType == 1 {
		//洗牌
		this.CardMgr.Shuffle()
		// 发牌 17 张
		for _, v := range this.Players {
			v.HandCard = []byte{}
			hd := this.CardMgr.SendHandCard(17)
			v.SetHandCard(Sort(hd))
		}
		// 底牌
		this.DiPai = []byte{}
		this.DiPai = this.CardMgr.SendHandCard(3)
		//发送牌消息
		sd := GGameSendCardNotify{
			Id:    MSG_GAME_INFO_SEND_NOTIFY,
			Rount: this.Round,
			Cid:   this.CurCid,
		}
		for _, v := range this.Players {
			for _, v1 := range v.HandCard {
				sd.HandsCards = append(sd.HandsCards, int(v1))
			}
			fmt.Println(v.Nick, "的手牌是：", sd.HandsCards)
			v.SendNativeMsg(MSG_GAME_INFO_SEND_NOTIFY, &sd)
			sd.HandsCards = []int{}
		}
	}
	//发牌进入抢地主(叫地主)阶段

	this.AddTimer(TIMER_START, TIMER_START_NUM, this.TimerDealPoker, nil)
}

//发牌动画
func (this *ExtDesk) TimerDealPoker(d interface{}) {
	logs.Debug("到发牌动画了")
	this.GameState = GAME_STATUS_CALL
	fmt.Println("刚开始的时候Calltype", this.TableConfig.CallType)
	if this.TableConfig.CallType == 1 {
		this.BroadStageTime(TIMER_CALL_NUM)
		//发牌后进入叫分阶段，开启叫分阶段的定时器
		this.AddTimer(TIMER_CALL, TIMER_CALL_NUM, this.TimerCall, nil)
	} else {
		this.BroadStageTime(TIMER_GETGMS_NUM)
		//发牌后进入抢地主阶段，开启抢地主阶段定时器
		this.AddTimer(TIMER_GETGMS, TIMER_GETGMS_NUM, this.TimerGet, nil)
	}
}

//创建房卡失败重置
func (this *ExtDesk) ReSet() {
	this.ClearTimer()

	this.GameState = GAME_STATUS_END
	this.BroadStageTime(0)
	this.GameOverLeave()

	//归还桌子
	this.GameState = GAME_STATUS_FREE
	this.ResetTable()
	this.DeskMgr.BackDesk(this)
}
