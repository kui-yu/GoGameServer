package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

const (
	Porsche     = 0 //保时捷
	BMW         = 1 //宝马
	Benz        = 2 //奔驰
	VW          = 3 //大众
	Min_Porsche = 4 //保时捷
	Min_BMW     = 5 //宝马
	Min_Benz    = 6 //奔驰
	Min_VW      = 7 //大众
)

var CarTypeMultiple map[int]float32

func InitMultiple() {
	CarTypeMultiple = make(map[int]float32)

	CarTypeMultiple[Porsche] = 36.8
	CarTypeMultiple[BMW] = 27.6
	CarTypeMultiple[Benz] = 18.4
	CarTypeMultiple[VW] = 9.2
	CarTypeMultiple[Min_Porsche] = 4.6
	CarTypeMultiple[Min_BMW] = 4.6
	CarTypeMultiple[Min_Benz] = 4.6
	CarTypeMultiple[Min_VW] = 4.6
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
