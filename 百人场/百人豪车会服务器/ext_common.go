package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

const (
	Ferrari      = iota //法拉利
	Maserati            //玛莎拉蒂
	Porsche             //保时捷
	Benz                //奔驰
	Min_Ferrari         //小法拉利
	Min_Maserati        //小玛莎拉蒂
	Min_Porsche         //小保时捷
	Min_Benz            //小奔驰
)

var CarTypeMultiple map[int]float32

func InitMultiple() {
	CarTypeMultiple = make(map[int]float32)
	CarTypeMultiple[Porsche] = gameConfig.Multiple.Porsche
	CarTypeMultiple[Benz] = gameConfig.Multiple.Benz
	CarTypeMultiple[Ferrari] = gameConfig.Multiple.Ferrari
	CarTypeMultiple[Maserati] = gameConfig.Multiple.Maserati
	CarTypeMultiple[Min_Porsche] = gameConfig.Multiple.Min_Porsche
	CarTypeMultiple[Min_Benz] = gameConfig.Multiple.Min_Benz
	CarTypeMultiple[Min_Ferrari] = gameConfig.Multiple.Min_Ferrari
	CarTypeMultiple[Min_Maserati] = gameConfig.Multiple.Min_Maserati

	fmt.Println("初始化:", CarTypeMultiple)
}

func FormatDeskId(deskId int, grade int) string {
	first := ""
	switch grade {
	case 1:
		first = "C"
	case 2:
		first = "Z"
	case 3:
		first = "G"
	}

	return fmt.Sprintf("%s%04d", first, deskId+1)
}

func DebugLog(format string, args ...interface{}) {
	// return
	fmt.Println(fmt.Sprintf(format, args))
}

// 获得当前时间的毫秒
func GetTimeMS() int64 {
	return time.Now().UnixNano() / 1e6
}

// 随机数生成器
func GetRandomNum(min, max int) (int, error) {
	maxBigInt := big.NewInt(int64(max - min))
	i, err := rand.Int(rand.Reader, maxBigInt)
	if i.Int64() < 0 {
		return 0, err
	}
	return int(i.Int64()) + min, err
}
