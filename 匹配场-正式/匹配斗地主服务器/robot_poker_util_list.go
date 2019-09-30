package main

import (
	"crypto/rand"
	"math/big"
	"strconv"
	"strings"
)

//数组随机
func ListShuffle(list []int) []int {

	var sourceList []int

	MVCard := append([]int{}, list...)
	// 随机打乱牌型
	for i := 0; i < len(list); i++ {
		//打乱
		randIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(MVCard))))
		//添加
		sourceList = append(sourceList, MVCard[randIndex.Int64()])
		// 移除已经添加的牌
		MVCard = ListDelMore(MVCard, int(randIndex.Int64()), 1)
	}
	return sourceList
}

func ListShuffleByInt(list []int) []int {
	return ListShuffle(list)
}

func ListIntToByte(list []int) []byte {
	var sourceList []byte
	for _, v := range list {
		sourceList = append(sourceList, byte(v))
	}
	return sourceList
}

func ListByteToInt(list []byte) []int {
	var sourceList []int
	for _, v := range list {
		sourceList = append(sourceList, int(v))
	}
	return sourceList
}

func ListInt32ToInt(list []int32) []int {
	var sourceList []int
	for _, v := range list {
		sourceList = append(sourceList, int(v))
	}
	return sourceList
}

//删除原列表某个元素 ListDelMore(list,删除的元素,1)，返回原列表
func ListDelMore(list []int, from int, num int) []int {
	sourceList := append([]int{}, list...)
	sourceList = append(sourceList[:from], sourceList[from+num:]...)
	return sourceList
}

//获取想要的参数，移除对应列表的元素，返回想要的参数列表
func ListDelGet(list []int, from int, num int) []int {

	sourceList := append([]int{}, list[from:from+num]...)
	//删除第0个元素开始
	list = ListDelMore(list, from, num)
	//返回
	return sourceList
}

//返回想要的参数列表
func ListGet(list []int, from int, num int) []int {
	//返回
	end := len(list)
	if end > from+num {
		end = from + num
	}
	return append([]int{}, list[from:end]...)
}

func ListGetByByte(list []byte, from int, num int) []byte {
	//返回
	end := len(list)
	if end > from+num {
		end = from + num
	}
	return append([]byte{}, list[from:end]...)
}

//list 删除 list
func ListDelList(list []int, delList []int) []int {
	for _, v := range delList {
		for index, n := range list {
			if v == n {
				list = ListDelMore(list, index, 1)
				break
			}
		}
	}
	return list
}

func ListDelListByByte(list []byte, delList []byte) []byte {
	tempList := ListByteToInt(list)
	tempDelList := ListByteToInt(delList)

	rsList := ListDelList(tempList, tempDelList)

	return ListIntToByte(rsList)
}

func StrSplitToCards(str string) []byte {

	var cards []byte

	strCards := strings.Split(str, " ")
	for _, c := range strCards {
		ic, _ := strconv.Atoi(c)
		cards = append(cards, byte(ic))
	}

	return cards
}
