package hydra

import (
	"prismx_cli/core/models"
	"prismx_cli/utils/netUtils"
	"strings"
	"time"
)

func MemcachedWeakPass(res any) any {
	t := res.(models.HydraTask)
	client, err := netUtils.SendDialTimeout("tcp", t.Target, t.Config.Timeout)
	if err != nil {
		return nil
	}
	defer client.Close()
	if client.SetDeadline(time.Now().Add(t.Config.Timeout)) != nil {
		return nil
	}
	_, err = client.Write([]byte("stats\n"))
	if err != nil {
		return nil
	}
	rev := make([]byte, 1024)
	n, err := client.Read(rev)
	if err != nil {
		return nil
	}
	if !strings.Contains(string(rev[:n]), "STAT") {
		return nil
	}
	return models.MSG{
		Name:    "Memcached unauthorized",
		Type:    "Unauthorized",
		Payload: models.Dict{},
		Target:  t.Target,
	}
}
