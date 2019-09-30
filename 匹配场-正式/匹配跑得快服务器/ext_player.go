package main

type ExtPlayer struct {
	Player
	HandCards    []byte        //手牌
	TuoGuan      bool          //是否托管
	Pass         bool          //玩家是否过牌
	IsDan        bool          //是否报单
	IsBaoPei     bool          //包赔
	IsQuanGUan   bool          //是否全关
	BoomBalance  int           //炸弹结算金额
	BaoPeiCoins  int64         //包赔金币
	Booms        int           //炸弹数量
	BeBooms      int           //被炸数
	BeBoomPlayer []int32       //炸你的玩家椅子id列表
	GetCoins     int           //盈利
	WaterProft   float64       //抽水
	WinForMap    map[int32]int //从哪些玩家赢取金币
	LoseForMap   map[int32]int //从输给了哪些玩家金币
}
