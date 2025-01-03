package hydra

import (
	"golang.org/x/crypto/ssh"
	"prismx_cli/core/models"
	"prismx_cli/utils/netUtils"
)

func SSHWeakPass(res any) any {
	var (
		t   = res.(models.HydraTask)
		msg = models.MSG{
			Name: "SSH WeakPassword",
			Type: "WeakPassword",
			Payload: models.Dict{
				User:     t.Dict.User,
				Password: t.Dict.Password,
			},
			Target: t.Target,
		}
	)

	proxy, err := netUtils.SendDialTimeout("tcp", t.Target, t.Config.Timeout)
	if err != nil {
		return nil
	}
	defer proxy.Close()
	config := &ssh.ClientConfig{
		User:            t.Dict.User,
		Timeout:         t.Config.Timeout,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth:            []ssh.AuthMethod{ssh.Password(t.Dict.Password)},
	}
	conn, _, _, err := ssh.NewClientConn(proxy, t.Target, config)
	if err != nil {
		return nil
	}
	defer conn.Close()
	return msg
}
