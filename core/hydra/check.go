package hydra

import (
	"context"
	"fmt"
	"net"
	"prismx_cli/core/models"
	"prismx_cli/utils/logger"
	"prismx_cli/utils/netUtils"
	"prismx_cli/utils/task"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

type proxyDialer struct {
	timeout time.Duration
}

// DialContext adheres to the mssql.Dialer interface.
func (c *proxyDialer) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	return netUtils.SendDialTimeout(network, addr, c.timeout)
}

func RunCheck(t *models.HydraTask, app models.HydraAppFunc) {

	if len(models.UserDict[app.App]) != 0 && len(models.PassDict[app.App]) == 0 {
		//有用户名无密码的情况
		for _, item := range models.UserDict[app.App] {
			t.DictList = append(t.DictList, models.Dict{User: item})
		}
	} else if len(models.UserDict[app.App]) == 0 && len(models.PassDict[app.App]) != 0 {
		//无用户名有密码的情况
		for _, item := range models.UserDict[app.App] {
			t.DictList = append(t.DictList, models.Dict{Password: item})
		}
	} else if len(models.UserDict[app.App]) == 0 && len(models.PassDict[app.App]) == 0 {
		//无账户无密码情况
		t.DictList = append(t.DictList, models.Dict{})
	} else {
		//根据服务名称获取对应字典列表
		for _, user := range models.UserDict[app.App] {
			for _, pass := range models.PassDict[app.App] {
				t.DictList = append(t.DictList, models.Dict{User: user, Password: strings.ReplaceAll(pass, "{user}", user)})
			}
		}
	}

	// 创建弱密码爆破池，固定 30 并发
	t.Scan = task.NewPool(100)

	// 找到弱密码的终止标记，找到第一个后立即终止剩余任务
	var found uint32

	for _, item := range t.DictList {
		// 先检查是否已找到正确密码，找到就跳过剩余任务
		if atomic.LoadUint32(&found) == 1 {
			break
		}

		dict := item
		t.Scan.Go(func() {
			// 每个任务开始前再检查一次，避免已找到后还继续执行
			if atomic.LoadUint32(&found) == 1 {
				return
			}

			// 创建任务副本，避免并发修改 t.Dict 导致数据竞争
			task := *t
			task.Dict = dict

			ctx, cancel := context.WithTimeout(context.Background(), t.Config.Timeout)
			defer cancel()

			done := make(chan any, 1) // 带缓冲，确保发送不阻塞

			go func() {
				defer func() { _ = recover() }() // 内层 goroutine 也加 panic 保护
				done <- app.Func(task)
			}()

			select {
			case <-ctx.Done():
				return
			case out := <-done:
				if out != nil {
					// 标记已找到，让后续任务快速跳过
					atomic.StoreUint32(&found, 1)

					data := out.(models.MSG)
					account := data.Payload.(models.Dict)
					if data.Type == "Unauthorized" {
						logger.ScanMessage(logger.Global.Color().Green(fmt.Sprintf("%-"+strconv.Itoa(20)+"s", t.Target)) + fmt.Sprintf(" [%s] Detect Unauthorized access", logger.Global.Color().YellowBg(app.App)))
					} else {
						logger.ScanMessage(logger.Global.Color().Green(fmt.Sprintf("%-"+strconv.Itoa(20)+"s", t.Target)) + fmt.Sprintf(" [%s] Detect Weak Password. Login: %s, Password: %s", logger.Global.Color().YellowBg(app.App), logger.Global.Color().Red(account.User), logger.Global.Color().Red(account.Password)))
					}
				}
			}
		})
	}

	t.Scan.Wait()
}
