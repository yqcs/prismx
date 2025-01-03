package hydra

import (
	"github.com/hirochachacha/go-smb2"
	"prismx_cli/core/models"
	"prismx_cli/utils/netUtils"
)

func SMBWeakPass(res any) any {
	var (
		t   = res.(models.HydraTask)
		msg = models.MSG{
			Name: "SMB WeakPassword",
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
	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     t.Dict.User,
			Password: t.Dict.Password,
		},
	}
	s, err := d.Dial(conn)
	if err != nil {
		return nil
	}
	s.Logoff()
	conn.Close()
	return msg
}
