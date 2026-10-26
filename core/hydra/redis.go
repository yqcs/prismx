package hydra

import (
	"net"

	"github.com/go-redis/redis"
	"github.com/yqcs/prismx/core/models"
	"github.com/yqcs/prismx/utils/netUtils"
)

func RedisWeakPass(res any) any {

	var (
		t   = res.(models.HydraTask)
		msg = models.MSG{
			Name: "Redis WeakPassword",
			Type: "WeakPassword",
			Payload: models.Dict{
				User:     t.Dict.User,
				Password: t.Dict.Password,
			},
			Target: t.Target,
		}
	)
	redisShell := redis.NewClient(&redis.Options{
		Addr:     t.Target,
		DB:       0,
		Password: t.Dict.Password,
		Dialer: func() (net.Conn, error) {
			return netUtils.SendDialTimeout("tcp", t.Target, t.Config.Timeout)
		},
	})

	pong, err := redisShell.Ping().Result()
	if err != nil {
		return nil
	}
	defer redisShell.Close()

	//如果不等于pong那么就是蜜罐，任务直接停止
	if pong != "PONG" {
		return nil
	}
	if t.Dict.Password == "" {
		msg.Name = "Redis Unauthorized"
		msg.Type = "Unauthorized"
	}
	return msg
}
