// package main

// import (
// 	"fmt"
// 	"time"
// )

// func running() {
// 	var times int
// 	//构建一个无限循环
// 	for true {
// 		times++
// 		fmt.Println("tick", times)
// 		//延时1s
// 		time.Sleep(time.Second)
// 	}
// }
// func main() {
// 	//并发执行程序
// 	go running()
// 	//接受命令行输入，不做任何事情
// 	var input string
// 	fmt.Scan(&input)
// 	fmt.Print("你输入的是:", input)
// }

// package main

// import (
// 	"fmt"
// 	"time"
// )

// func main() {
// 	go func() {
// 		var times int
// 		for true {
// 			times++
// 			fmt.Println("tick", times)
// 			time.Sleep(time.Second)
// 		}
// 	}()
// 	var input string
// 	fmt.Scan(&input)
// 	fmt.Print("你输入的是:", input)
// }

// package main

// // "fmt"

// func main() {
// 	//创建一个空接口通道
// 	ch := make(chan interface{})
// 	//将0放入通道中
// 	ch <- 0
// 	//将hello 字符串放入通道中
// 	ch <- "hello"
// }

// package main

// func main() {
// 	// 创建一个整型通道
// 	ch := make(chan int)
// 	// 尝试将0通过通道发送
// 	ch <- 0
// }

// package main

// import (
// 	"fmt"
// )

// func main() {
// 	//构建一个通道
// 	// ch := make(chan int)
// 	//开启一个并发匿名函数
// 	go func() {
// 		fmt.Println("goroutine start")
// 		//通过通道通知main的goroutine
// 		// ch <- 0
// 		fmt.Println("goroutine end")
// 	}()
// 	fmt.Println("wait goroutine")
// 	//等待匿名 goroutine
// 	// <-ch
// 	fmt.Println("all done")
// }

// package main

// import (
// 	"fmt"
// 	"time"
// )

// func main() {
// 	//构建一个通道
// 	ch := make(chan int)
// 	//开启一个并发匿名函数
// 	go func() {
// 		//从3循环到0
// 		for i := 3; i >= 0; i-- {
// 			//发送3-0之间的数值
// 			ch <- i
// 			//每次发送完时等待
// 			time.Sleep(time.Second)
// 		}
// 	}()
// 	//遍历接收数据
// 	for data := range ch {
// 		//打印通道数据
// 		fmt.Println(data)
// 		//当遇到数据0时，退出接收循环
// 		if data == 0 {
// 			break
// 		}
// 	}
// }

package main

import (
	"fmt"
)

func printer(c chan int) {
	//开始无限循环等待数据
	for true {
		//从 channel中获取一个数据
		data := <-c
		//将 0视为数据结束
		if data == 0 {
			break
		}
		//打印数据
		fmt.Println(data)
	}
	//通知main已经结束循环
	c <- 0
}

func main() {
	//创建一个channel
	c := make(chan int)
	//并发执行 printer ，传入channel
	go printer(c)

	for i := 1; i <= 10; i++ {
		//将数据通过channel 投给printer
		c <- i
	}
	//通知并发的printer 结束循环（没数据啦）
	c <- 0
	//等待printer 结束 （搞定喊我)
	<-c
}
