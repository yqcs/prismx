package hydra

import (
	"github.com/jlaffaye/ftp"
	"net"
	"prismx_cli/core/models"
	"prismx_cli/utils/netUtils"
)

func FtpWeakPass(res any) any {

	var (
		t   = res.(models.HydraTask)
		msg = models.MSG{
			Name: "FTP WeakPassword",
			Type: "WeakPassword",
			Payload: models.Dict{
				User:     t.Dict.User,
				Password: t.Dict.Password,
			},
			Target: t.Target,
		}
	)
	conn, err := ftp.Dial(t.Target, ftp.DialWithDialFunc(func(network, address string) (net.Conn, error) {
		return netUtils.SendDialTimeout("tcp", t.Target, t.Config.Timeout)
	}))
	if err != nil {
		return nil
	}
	if err = conn.Login(t.Dict.User, t.Dict.Password); err != nil {
		return nil
	}
	conn.Logout()
	return msg
}
