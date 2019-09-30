// kclient project main.go
package main

import (
	"time"

	"math/rand"
	"strconv"
)

var G_ALLROBOTER []*Roboter

func main() {
	rand.Seed(time.Now().UnixNano())
	uuid := time.Now().UnixNano() * 10000
	uuid = 1
	for i := 0; i < GCONFIG.Num; i++ {
		var p Roboter
		p.SInit()
		G_ALLROBOTER = append(G_ALLROBOTER, &p)
		go p.HandleMsg()
		xuhao := uuid + int64(i)
		go p.Start("robot_"+strconv.FormatInt(xuhao, 10), xuhao, uuid)
	}

	t := time.NewTicker(time.Second)
	for _ = range t.C {
		for _, v := range G_ALLROBOTER {
			v.MTimerChan <- true
		}
	}
}
