package main

import (
	"math/rand"
)

//数组头添加数据
func PushArr(arr []GameInfo, obj GameInfo) []GameInfo {
	res := make([]GameInfo, 0)
	res = append(res, obj)
	res = append(res, arr...)
	if len(res) > 8 {
		return res[:8]
	}
	return res
}
//获取骰子数
func GetCount() []int64 {
	arr := []int64{}
	for i := 1; i <= 3; i++ {
		count := rand.Intn(6) + 1
		arr = append(arr, int64(count))
	}
	return arr
}
