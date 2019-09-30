package main

import (
	"crypto/rand"
	"math/big"
)

//数组随机
func ListShuffle(list []int32) []int32 {

	var sourceList []int32

	MVCard := append([]int32{}, list...)
	// 随机打乱牌型
	for i := 0; i < len(list); i++ {
		//打乱
		randIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(MVCard))))
		//添加
		sourceList = append(sourceList, MVCard[randIndex.Int64()])
		// 移除已经添加的牌
		MVCard = ListDelOne(MVCard, int32(randIndex.Int64()))
	}
	return sourceList
}

func ListShuffleByInt(list []int) []int {
	list32 := ListIntToInt32(list)
	rsList32 := ListShuffle(list32)
	return ListInt32ToInt(rsList32)
}

func ListByteToInt32(list []byte) []int32 {
	var sourceList []int32
	for _, v := range list {
		sourceList = append(sourceList, int32(v))
	}
	return sourceList
}

func ListIntToInt32(list []int) []int32 {
	var sourceList []int32
	for _, v := range list {
		sourceList = append(sourceList, int32(v))
	}
	return sourceList
}

//删除某个元素
func ListDelOne(list []int32, num int32) []int32 {

	sourceList := append([]int32{}, list...)
	sourceList = append(sourceList[:num], sourceList[num+1:]...)

	return sourceList
}

func ListInt32ToInt(list []int32) []int {
	var sourceList []int
	for _, v := range list {
		sourceList = append(sourceList, int(v))
	}
	return sourceList
}
