package task

import (
	"github.com/panjf2000/ants/v2"
	"sync"
)

// Pool 任务池
type Pool struct {
	//PoolWithFunc 队列
	PoolWithFunc *ants.PoolWithFunc
	//堵塞器
	Wg *sync.WaitGroup
}

// NewPool 实例化工作池使用
func NewPool() *Pool {
	return &Pool{
		Wg: &sync.WaitGroup{},
	}
}
