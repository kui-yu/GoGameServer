// log project log.go
package logs

import (
	"bufio"
	"fmt"
	"os"

	"github.com/astaxie/beego/logs"
)

var LogError, LogDebug *logs.BeeLogger
var FuncCallDepth = 3

func init() {
	CreateFileLog()
	//
	LogError = logs.NewLogger(10000)
	LogError.EnableFuncCallDepth(true)
	LogError.SetLevel(8)
	LogError.SetLogger("maxdays", "30")
	LogError.SetLogger("file", `{"filename":"./log/server_err.log"}`)
	//
	LogDebug = logs.NewLogger(10000)
	LogDebug.EnableFuncCallDepth(true)
	LogDebug.SetLevel(8)
	LogDebug.SetLogger("maxdays", "30")
	LogDebug.SetLogger("file", `{"filename":"./log/server_debug.log"}`)

}

func CreateFile(filepath string) error {
	file, err := os.OpenFile(filepath, os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
	}
	//
	defer file.Close()
	//
	wFile := bufio.NewWriter(file)
	wFile.Flush()
	return err
}

func CreateFileLog() {
	_, err := os.Stat("./log")
	if err != nil {
		err = os.Mkdir("./log", 0700)
		if err != nil {
			panic(err)
		}
	}
	os.Chdir("./log")
	filepath := "server_err.log"
	filepath2 := "server_debug.log"
	CreateFile(filepath)
	CreateFile(filepath2)
	os.Chdir("./../")
}

func Debug(format string, v ...interface{}) {
	fmt.Println(fmt.Sprintf(format, v...))
	LogDebug.SetLogFuncCallDepth(FuncCallDepth)
	LogDebug.Debug(format, v...)
}

func Error(format string, v ...interface{}) {
	LogError.SetLogFuncCallDepth(FuncCallDepth)
	LogError.Error(format, v...)
}
