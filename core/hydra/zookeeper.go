package hydra

import (
	"bytes"

	"github.com/yqcs/prismx/core/models"
	"github.com/yqcs/prismx/utils/netUtils"
)

func ZookeeperWeakPass(res any) any {
	var (
		t   = res.(models.HydraTask)
		msg = models.MSG{
			Name:    "Zookeeper Unauthorized",
			Type:    "Unauthorized",
			Payload: models.Dict{},
			Target:  t.Target,
		}
	)
	conn, err := netUtils.SendDialTimeout("tcp", t.Target, t.Config.Timeout)
	if err != nil {
		return nil
	}
	_, err = conn.Write([]byte("envi"))
	if err != nil {
		return nil
	}
	reply := make([]byte, 1024)
	n, err := conn.Read(reply)
	if err != nil {
		return nil
	}
	if bytes.Contains(reply[:n], []byte("Environment")) {
		return msg
	}
	return nil
}
