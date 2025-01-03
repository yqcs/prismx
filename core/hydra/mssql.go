package hydra

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	mssql "github.com/denisenkom/go-mssqldb"
	"net"
	"prismx_cli/core/models"
)

func MSSQLWeakPass(res any) any {
	var (
		t   = res.(models.HydraTask)
		msg = models.MSG{
			Name: "MSSQL WeakPassword",
			Type: "WeakPassword",
			Payload: models.Dict{
				User:     t.Dict.User,
				Password: t.Dict.Password,
			},
			Target: t.Target,
		}
	)
	host, port, err := net.SplitHostPort(t.Target)
	if err != nil {
		return nil
	}
	conn, err := mssql.NewConnector(fmt.Sprintf("server=%s;user id=%s;password=%s;port=%v;encrypt=disable;timeout=%v", host, t.Dict.User, t.Dict.Password, port, t.Config.Timeout))
	if err != nil {
		return nil
	}
	conn.Dialer = &proxyDialer{
		timeout: t.Config.Timeout,
	}
	db := sql.OpenDB(conn)
	if err = db.Ping(); err != nil {
		return nil
	}
	db.Close()
	return msg
}
