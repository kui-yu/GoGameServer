package main

import (
	"crypto/rand"
	"math/big"
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

//数组随机
func ListShuffleByInt32(list []int32) []int32 {

	var sourceList []int32

	MVCard := append([]int32{}, list...)
	// 随机打乱牌型
	for i := 0; i < len(list); i++ {
		//打乱
		randIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(MVCard))))
		//添加
		sourceList = append(sourceList, MVCard[randIndex.Int64()])
		// 移除已经添加的牌
		MVCard = ListDelOneByInt32(MVCard, int32(randIndex.Int64()))
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

//删除某个元素
func ListDelOneByInt32(list []int32, num int32) []int32 {

	sourceList := append([]int32{}, list...)
	sourceList = append(sourceList[:num], sourceList[num+1:]...)

	return sourceList
}

//排序
func ListSortDesc(cards []int) []int {
	cs := append([]int{}, cards...)
	for i := 0; i < len(cs)-1; i++ {
		for j := i + 1; j < len(cs); j++ {
			vi := GetCardValue(cs[i])
			vj := GetCardValue(cs[j])
			if vi < vj || ((vi == vj) && (GetCardColor(cs[i]) > GetCardColor(cs[j]))) {
				vt := cs[i]
				cs[i] = cs[j]
				cs[j] = vt
			}
		}
	}
	return cs
}

//排序
func ListSortAsc(cards []int) []int {
	cs := append([]int{}, cards...)
	for i := 0; i < len(cs)-1; i++ {
		for j := i + 1; j < len(cs); j++ {
			vi := GetCardValue(cs[i])
			vj := GetCardValue(cs[j])
			if vi > vj || ((vi == vj) && (GetCardColor(cs[i]) > GetCardColor(cs[j]))) {
				vt := cs[i]
				cs[i] = cs[j]
				cs[j] = vt
			}
		}
	}
	return cs
}

//数组新增元素
func ListAdd(list *[]int, element int) {
	(*list) = append((*list), element)
}
