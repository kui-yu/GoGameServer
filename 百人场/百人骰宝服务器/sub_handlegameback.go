package main

//请求返回大厅
func (this *ExtDesk)HandleGameBack(p *ExtPlayer,d *DkInMsg){
	obj := new(GSGameBack)
	obj.Id = MSG_GAME_INFO_BACK_REPLAY
	if p.IsBet{
		obj.Err = "已投注不能返回大厅"
		obj.Result = false
		p.SendNativeMsg(MSG_GAME_INFO_BACK_REPLAY,obj)
		return
	}else{
		obj.Err = ""
		obj.Result = true
		p.SendNativeMsg(MSG_GAME_INFO_BACK_REPLAY,obj)
		//调用底层
		this.LeaveByForce(p)
	}
}
