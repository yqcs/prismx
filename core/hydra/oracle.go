package hydra

import (
	"net"
	"prismx_cli/core/models"
	"prismx_cli/utils/go-ora"
	"strconv"
)

func OracleWeakPass(res any) any {
	var serviceName = []string{
		"orcl",
		"xe",
		"oracle",
	}

	var (
		t   = res.(models.HydraTask)
		msg = models.MSG{
			Name: "Oracle WeakPassword",
			Type: "WeakPassword",
			Payload: models.Dict{
				User:     t.Dict.User,
				Password: t.Dict.Password,
			},
			Target: t.Target,
		}
	)
	host, por, err := net.SplitHostPort(t.Target)
	if err != nil {
		return nil
	}
	atoi, err := strconv.Atoi(por)
	if err != nil {
		return nil
	}
	for _, service := range serviceName {
		connection, err := go_ora.NewConnection(go_ora.BuildUrl(host, atoi, service, t.Dict.User, t.Dict.Password, nil))
		if err != nil {
			continue
		}
		if err := connection.Open(); err != nil {
			continue
		}
		connection.Close()
		return msg
	}
	return nil
}
