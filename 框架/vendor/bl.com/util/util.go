package util

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"time"
)

const (
	Card_A byte = iota + 1
	Card_2
	Card_3
	Card_4
	Card_5
	Card_6
	Card_7
	Card_8
	Card_9
	Card_10
	Card_J
	Card_Q
	Card_K
)

func GetRandomNum(min, max int) (int, error) {
	maxBigInt := big.NewInt(int64(max - min))
	i, err := rand.Int(rand.Reader, maxBigInt)
	if i.Int64() < 0 {
		return 0, err
	}
	return int(i.Int64()) + min, err
}

func Shuffle(v []byte) []byte {
	ret := append([]byte{}, v...)

	for i := range ret {
		r, _ := GetRandomNum(i, len(ret))
		ret[i], ret[r] = ret[r], ret[i]
	}

	return ret
}

func BuildGameId(gameType int) string {
	t := time.Now().UnixNano()
	s := t / 1e9
	ns := t % 1e9

	return strconv.Itoa(gameType) + "-" + strconv.Itoa(int(s)) + "-" + strconv.Itoa(int(ns))
}

func BuildRoomId(gradeType, id int) string {
	if gradeType == 1 {
		return fmt.Sprintf("T%04d", id)
	}

	if gradeType == 2 {
		return fmt.Sprintf("C%04d", id)
	}

	if gradeType == 3 {
		return fmt.Sprintf("M%04d", id)
	}

	if gradeType == 4 {
		return fmt.Sprintf("H%04d", id)
	}

	return fmt.Sprintf("N%04d", id)
}

func LessInt32List(sub, min []int32) []int32 {
	var ret []int32

	for i := range sub {
		ret = append(ret, sub[i]-min[i])
	}

	return ret
}

func LessInt64List(sub, min []int64) []int64 {
	var ret []int64

	for i := range sub {
		if min[i] < 0 {
			ret = append(ret, sub[i])
			continue
		}
		ret = append(ret, sub[i]-min[i])
	}

	return ret
}
