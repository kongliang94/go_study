package main

import (
	"fmt"
	"time"
)

/* 有关Task任务相关定义及操作 */
//定义任务Task类型,每一个任务Task都可以抽象成一个函数
type Task struct {
	f func() error //一个无参的函数类型
}

//通过NewTask来创建一个Task
func NewTask(f func() error) *Task {
	task := Task{
		f,
	}
	return &task
}

//执行行行Task任务的方方法
func (t *Task) Execute() {
	t.f() //调用用任务所绑定的函数
}

/* 有关协程池的定义及操作 */
//定义池类型
type Pool struct {
	// EntryChannel 对外接收Task的入口
	EntryChannel chan *Task
	// worker_num 协程池最大worker数量,限定Goroutine的个数
	worker_num int
	// JobsChannel 协程池内部的任务就绪队列列
	JobsChannel chan *Task
}

// NewPool 创建一个协程池
func NewPool(cap int) *Pool {
	pool := Pool{
		EntryChannel: make(chan *Task),
		worker_num:   cap,
		JobsChannel:  make(chan *Task),
	}
	return &pool
}

//协程池创建一个worker并且开始工工作
func (p *Pool) worker(work_ID int) {
	//worker不不断的从JobsChannel内部任务队列列中拿任务
	for task := range p.JobsChannel {
		//如果拿到任务,则执行行行task任务
		task.Execute()
		fmt.Println("worker ID ", work_ID, " 执行行行完毕任务")
	}
}

//让协程池Pool开始工工作
func (p *Pool) Run() {
	//1,首首先根据协程池的worker数量量限定,开启固定数量量的Worker,
	//每一一个Worker用用一一个Goroutine承载
	for i := 0; i < p.worker_num; i++ {
		go p.worker(i)
	}
	//2, 从EntryChannel协程池入入口口取外界传递过来的任务
	//并且将任务送进JobsChannel中
	for task := range p.EntryChannel {
		p.JobsChannel <- task
	}
	//3, 执行行行完毕需要关闭JobsChannel
	close(p.JobsChannel)
	//4, 执行行行完毕需要关闭EntryChannel
	close(p.EntryChannel)
}

func main() {

	//创建一一个Task
	t := NewTask(func() error {
		time.Sleep(1 * time.Second)
		fmt.Println(time.Now())
		return nil
	})

	t1 := NewTask(func() error {
		time.Sleep(2 * time.Second)
		fmt.Println(time.Now())
		return nil
	})
	//创建一个协程池,最大大开启3个协程worker
	p := NewPool(5)
	//开一个协程 不不断的向 Pool 输送打印一条时间的task任务
	go func() {
		for {
			time.Sleep(5 * time.Second)
			p.EntryChannel <- t
			p.EntryChannel <- t1
		}
	}()
	//启动协程池p
	p.Run()

}
