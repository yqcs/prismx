package hydra

import (
	"context"
	"github.com/masterzen/winrm"
	"net"
	"os"
	"prismx_cli/core/models"
	"prismx_cli/utils/netUtils"
	"strconv"
)

func WinRMWeakPass(res any) any {
	var (
		t   = res.(models.HydraTask)
		msg = models.MSG{
			Name:    "WinRM WeakPassword",
			Type:    "WeakPassword",
			Payload: models.Dict{},
			Target:  t.Target,
		}
	)

	params := winrm.DefaultParameters
	params.Dial = func(network, addr string) (net.Conn, error) {
		return netUtils.SendDialTimeout("tcp", t.Target, t.Config.Timeout)
	}
	host, port, _ := net.SplitHostPort(t.Target)
	intPort, _ := strconv.Atoi(port)
	client, err := winrm.NewClientWithParameters(winrm.NewEndpoint(host, intPort, false, false, nil, nil, nil, t.Config.Timeout), t.Dict.User, t.Dict.Password, params)
	if err != nil {
		return nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err = client.RunWithContext(ctx, "echo ok > nul", os.Stdout, os.Stderr)
	if err != nil {
		return nil
	}
	return msg
}
