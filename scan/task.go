package scan

import (
	"github.com/panjf2000/ants/v2"
	"prismx_cli/core/models"
	"prismx_cli/utils/logger"
	"prismx_cli/utils/task"
	"time"
)

type TaskPool struct {
	Scan      *task.Pool
	Params    models.ScanParams
	HydraTask *models.HydraTask
}

func (t *TaskPool) NewPoolWithFunc(pool *task.Pool, invoke func(), function func(any)) {
	//任务函数
	pool.PoolWithFunc, _ = ants.NewPoolWithFunc(t.Params.Thread, func(i interface{}) {
		function(i)
		t.Scan.Wg.Done()
	})
	//任务下发函数
	invoke()
	//实体队列堵塞
	pool.Wg.Wait()
	//释放锁
	//清除任务
	pool.PoolWithFunc.Release()
}

// Start 全部存活端口
func (t *TaskPool) Start() {
	start := time.Now()
	//捕捉启动日志
	logger.Info(logger.Global.Color().Yellow("Start running scan task"))
	//任务堵塞流
	t.NewPoolWithFunc(t.Scan, t.TaskInChan, t.TaskFunc)
	//捕捉结束日志
	logger.Info(logger.Global.Color().Yellow("The task has ended, taking - " + time.Since(start).String()))
}
