package hydra

import (
	"context"

	"github.com/yqcs/prismx/core/models"
	"github.com/yqcs/prismx/utils/netUtils"
)

func VncWeakPass(res any) any {

	var (
		t   = res.(models.HydraTask)
		msg = models.MSG{
			Name: "VNC WeakPassword",
			Type: "WeakPassword",
			Payload: models.Dict{
				User:     t.Dict.User,
				Password: t.Dict.Password,
			},
			Target: t.Target,
		}
	)

	conn, err := netUtils.SendDialTimeout("tcp", t.Target, t.Config.Timeout)
	if err != nil {
		return nil
	}
	vc, err := vnc.Connect(context.Background(), conn, vnc.NewClientConfig(t.Dict.Password))
	if err != nil {
		return nil
	}
	vc.Close()
	return msg
}
