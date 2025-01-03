package hydra

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"net"
	"prismx_cli/core/models"
	"prismx_cli/utils/netUtils"
)

func PGSQLWeakPass(res any) any {
	var (
		t   = res.(models.HydraTask)
		msg = models.MSG{
			Name: "Postgres WeakPassword",
			Type: "WeakPassword",
			Payload: models.Dict{
				User:     t.Dict.User,
				Password: t.Dict.Password,
			},
			Target: t.Target,
		}
	)
	config, err := pgx.ParseConfig(fmt.Sprintf("postgres://%v:%v@%s/postgres?sslmode=disable", t.Dict.User, t.Dict.Password, t.Target))
	if err != nil {
		return nil
	}
	config.DialFunc = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return netUtils.SendDialTimeout(network, addr, t.Config.Timeout)
	}
	ctx := context.Background()

	conn, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		return nil
	}
	defer conn.Close(ctx)
	if err = conn.Ping(ctx); err != nil {
		return nil
	}
	return msg
}
