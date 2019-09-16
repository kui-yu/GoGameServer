// ModelGame project main.go
package main

import (
	"logs"
	"net/http"

	"strconv"

	"math/rand"
	_ "net/http/pprof"
	"runtime"
	"time"

	"code.google.com/p/go.net/websocket"
)

func main() {
	logs.Debug("服务器运行开始")
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().UnixNano())
	http.Handle("/robotserver", websocket.Server{Handler: DoHandler})

	http.Handle("/", http.FileServer(http.Dir("home/")))
	//
	InitHttpHandle()
	//
	err := http.ListenAndServe(":"+strconv.Itoa(GCONFIG.Port), nil)
	if err != nil {
		logs.Debug("main Serve:", err)
	}
}
