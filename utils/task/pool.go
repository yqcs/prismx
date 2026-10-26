package task

import "sync"

// Pool 极简工作池 - 使用 Go 原生原语实现并发控制
// 替代 ants 框架，仅保留核心功能：限制最大并发数 + 等待所有任务完成
type Pool struct {
	sema chan struct{} // 信号量，控制最大并发数
	wg   sync.WaitGroup
}

// NewPool 创建工作池，maxWorkers 为最大并发数
func NewPool(maxWorkers int) *Pool {
	if maxWorkers <= 0 {
		maxWorkers = 100 // 默认值保护
	}
	return &Pool{
		sema: make(chan struct{}, maxWorkers),
	}
}

// Go 提交一个无参数任务
func (p *Pool) Go(task func()) {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		// 获取信号量令牌
		p.sema <- struct{}{}
		defer func() {
			// 释放令牌
			<-p.sema
			// panic 恢复，避免单个任务崩溃整个扫描
			_ = recover()
		}()
		task()
	}()
}

// GoWith 提交一个带参数的任务，替代 ants.PoolWithFunc.Invoke
func (p *Pool) GoWith(arg any, task func(any)) {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		p.sema <- struct{}{}
		defer func() {
			<-p.sema
			_ = recover()
		}()
		task(arg)
	}()
}

// Wait 等待所有任务完成
func (p *Pool) Wait() {
	p.wg.Wait()
}

// Add 手动增加计数（兼容原有逻辑）
func (p *Pool) Add(delta int) {
	p.wg.Add(delta)
}

// Done 手动完成计数（兼容原有逻辑）
func (p *Pool) Done() {
	p.wg.Done()
}
