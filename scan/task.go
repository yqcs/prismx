package scan

import (
	"prismx_cli/core/models"
	"prismx_cli/utils/logger"
	"time"
)

type TaskPool struct {
	Params    models.ScanParams
	HydraTask *models.HydraTask
}

// Start 启动四阶段扫描流水线
func (t *TaskPool) Start() {
	start := time.Now()
	//捕捉启动日志
	logger.Info(logger.Global.Color().Yellow("Start running scan task"))
	//四阶段流水线直接执行
	t.TaskInChan()
	//捕捉结束日志
	logger.Info(logger.Global.Color().Yellow("The task has ended, taking - " + time.Since(start).String()))
}
