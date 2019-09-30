package main

import (
	"fmt"
	"math/rand"
	"time"
)

//数组随机
func ListShuffle(list []int) []int {
	var sourceList []int
	fmt.Println(len(list), "牌数")
	MVCard := append([]int{}, list...)
	// 随机打乱牌型
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	perm := r.Perm(len(MVCard))
	for i := 0; i < len(list); i++ {
		sourceList = append(sourceList, list[perm[i]])
	}
	return sourceList
}
func ListShuffleByInt(list []int) []int {
	return ListShuffle(list)
}

//删除原列表某个元素 ListDelMore(list,删除的元素,1)，返回原列表
func ListDelMore(list []int, from int, num int) []int {
	sourceList := append([]int{}, list...)
	sourceList = append(sourceList[:from], sourceList[from+num:]...)
	return sourceList
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

func SortPlayerCard(handCard [][]int) (p int) {
	max := handCard[0]
	p = 0
	for i := 1; i < len(handCard); i++ {
		maxc, maxco := SortHandCard(max)
		hcc, hcco := SortHandCard(handCard[i])
		result := GetResult(maxc, maxco, hcc, hcco)
		if result == 1 {
			max = handCard[i]
			p = i
		}
	}
	return p
}

func IsCoinEnough(coin int64, paycoin []int64, bscore int64, mincoin int64, cardtype int) int {
	sum := int64(0)
	for i := 0; i < len(paycoin); i++ {
		sum += paycoin[i]
	}
	if cardtype == 0 {
		if (coin-sum*bscore)/bscore > mincoin {
			return 1
		} else {
			return 0
		}
	} else {
		if (coin-sum*bscore)/bscore > mincoin*2 {
			return 1
		} else {
			return 0
		}
	}
}

func SortCard(allCard []int) []int {
	for n := len(allCard); n > 0; n-- { //冒泡排序
		for i := 0; i < n-1; i++ {
			if allCard[i] > allCard[i+1] {
				allCard[i], allCard[i+1] = allCard[i+1], allCard[i]
			}
		}
	}
	return allCard
}

func DuplicateRemoval(list []int) []int {
	sign := make([]int, len(list))
	for i := 0; i < len(list); i++ {
		for j := i + 1; j < len(list); j++ {
			if list[i] == list[j] {
				sign[j] = 1
			}
		}
	}

	newlist := make([]int, 0, len(list))

	for i := 0; i < len(list); i++ {
		if sign[i] != 1 {
			newlist = append(newlist, list[i])
		}
	}
	return newlist
}
