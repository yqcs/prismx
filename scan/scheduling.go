package scan

import (
	"time"

	"github.com/panjf2000/ants/v2"
	fingerprint "github.com/yqcs/fingerscan"
	"github.com/yqcs/prismx/core/models"
	"github.com/yqcs/prismx/utils/logger"
	"github.com/yqcs/prismx/utils/task"
)

var (
	CaptureHost   func(id, host string)
	CaptureFinger func(id string, host string, port int, finger fingerprint.AppFinger)
	CaptureVul    func(id string, host string, port int, vul models.MSG)
	CaptureHydra  func(id string, target, app string, msg models.MSG)
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
