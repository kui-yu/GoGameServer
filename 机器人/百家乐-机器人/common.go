package main

//此文件不可修改
import (
	"crypto/rand"
	"fmt"
	"logs"
	"math/big"
)

// 随机数生成器
func GetRandomNum(min, max int) (int, error) {
	maxBigInt := big.NewInt(int64(max - min))
	i, err := rand.Int(rand.Reader, maxBigInt)
	if i.Int64() < 0 {
		return 0, err
	}
	return int(i.Int64()) + min, err
}

// 删除数组里面的值，并返回数组
func DeleteIntArrayFromValue(arr []int, value int) []int {
	for i, v := range arr {
		if v == value {
			return append(arr[:i], arr[i+1:]...)
		}
	}
	return arr
}

func DebugLog(format string, a ...interface{}) {
	if gameConfig.GCGameServer.LogPrint {
		logs.Debug(format, a)
	}
}

func ErrorLog(format string, a ...interface{}) {
	if gameConfig.GCGameServer.ErrorPrint {
		logs.Error(format, a)
	}
}

func TestLog(format string, a ...interface{}) {
	logs.Debug(fmt.Sprintf(format, a))
}
