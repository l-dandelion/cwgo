package pool

import (
	"sync"
)

type Pool interface {
	Add()
	Done()
	Wait()
	Init(int)
	Num() int
	Cap() int
}

type myPool struct {
	pool  chan bool      //控制数量
	group sync.WaitGroup //用于等待所有爬虫执行完毕
}

//增加爬虫
func (pool *myPool) Add() {
	pool.pool <- true
	pool.group.Add(1)
}

//爬虫结束
func (pool *myPool) Done() {
	<-pool.pool
	pool.group.Done()
}

//等待所有爬虫结束
func (pool *myPool) Wait() {
	pool.group.Wait()
}

func (pool *myPool) Init(maxNum int) {
	pool.pool = make(chan bool, maxNum)
}

//进程池中已分配的数量
func (pool *myPool) Num() int {
	return len(pool.pool)
}

//进程池容量
func (pool *myPool) Cap() int {
	return cap(pool.pool)
}

func New(maxNum int) Pool {
	return &myPool{
		pool: make(chan bool, maxNum),
	}
}
