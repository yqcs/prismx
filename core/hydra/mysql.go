package hydra

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"net"
	"prismx_cli/core/models"
	"prismx_cli/utils/netUtils"
)

type noLog struct{}

func (noLog) Print(v ...interface{}) {}
func init() {
	mysql.SetLogger(noLog{})
}

func MySQLWeakPass(res any) any {

	var (
		t   = res.(models.HydraTask)
		msg = models.MSG{
			Name: "MySQL WeakPassword",
			Type: "WeakPassword",
			Payload: models.Dict{
				User:     t.Dict.User,
				Password: t.Dict.Password,
			},
			Target: t.Target,
		}
	)

	db, err := sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%s)/information_schema?charset=utf8&timeout=%v", t.Dict.User, t.Dict.Password, t.Target, t.Config.Timeout))
	if err != nil {
		return nil
	}

	//设置代理
	mysql.RegisterDialContext("tcp", func(ctx context.Context, addr string) (net.Conn, error) {
		return netUtils.SendDialTimeout("tcp", addr, t.Config.Timeout)
	})

	defer db.Close()

	if err = db.Ping(); err != nil {
		return nil
	}
	return msg
}
