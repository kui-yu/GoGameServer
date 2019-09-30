package main

//个人战绩回复
func (this *ExtDesk) HandleRecord(p *ExtPlayer, d *DkInMsg) {

	info := GSRecordInfos{
		Id:    MSG_GAME_INFO_RECORD_INFO_REPLY,
		Infos: p.RecordInfos,
	}
	p.SendNativeMsg(MSG_GAME_INFO_RECORD_INFO_REPLY, &info)
}
