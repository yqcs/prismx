package hydra

import (
	"context"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"net"
	"prismx_cli/core/models"
	"prismx_cli/utils/logger"
	"prismx_cli/utils/netUtils"
	"prismx_cli/utils/task"
	"strconv"
	"strings"
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
	t.Scan = task.NewPool()
	t.Scan.PoolWithFunc, _ = ants.NewPoolWithFunc(30, func(i interface{}) {
		defer t.Scan.Wg.Done()

		ctx, cancel := context.WithTimeout(context.Background(), t.Config.Timeout)
		defer cancel()

		done := make(chan any)

		go func() {
			done <- app.Func(i)
		}()

		select {
		case <-ctx.Done():
			return
		case out := <-done:
			if out != nil {
				data := out.(models.MSG)
				account := data.Payload.(models.Dict)
				if data.Type == "Unauthorized" {
					logger.ScanMessage(logger.Global.Color().Green(fmt.Sprintf("%-"+strconv.Itoa(20)+"s", t.Target)) + fmt.Sprintf(" [%s] Detect Unauthorized access", logger.Global.Color().YellowBg(app.App)))
				} else {
					logger.ScanMessage(logger.Global.Color().Green(fmt.Sprintf("%-"+strconv.Itoa(20)+"s", t.Target)) + fmt.Sprintf(" [%s] Detect Weak Password. Login: %s, Password: %s", logger.Global.Color().YellowBg(app.App), logger.Global.Color().Red(account.User), logger.Global.Color().Red(account.Password)))
				}
			}
			close(done)
		}
	})

	for _, item := range t.DictList {
		t.Scan.Wg.Add(1)
		t.Dict = item
		t.Scan.PoolWithFunc.Invoke(*t)
	}
	t.Scan.Wg.Wait()
	t.Scan.PoolWithFunc.Release()
}
