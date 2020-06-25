package main

import (
	"fmt"
	"time"
)

// 演示协程
func go_worker(name string) {
	for i := 0; i < 10; i++ {
		fmt.Println("我是一个go协程, 我的名字是 ", name, "----")
		time.Sleep(1 * time.Second)
	}
	fmt.Println(name, " 执行行行完毕!")
}

// 演示channel
func rec_worker(rec chan string) {

	str := <-rec

	fmt.Println("i am rec_worker")
	fmt.Println("from main func data --- ", str)
}

func main() {

	// 创建channel
	c := make(chan string)

	go go_worker("tom")

	go go_worker("jack")

	go rec_worker(c)

	c <- "hello rec_worker"

	fmt.Println("i am main")

	//防止止main函数执行行行完毕,程序退出
	for {
		time.Sleep(1 * time.Second)
	}
}
