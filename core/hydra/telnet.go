package hydra

import (
	"prismx_cli/core/models"
	"prismx_cli/utils/go-telnet"
)

func TelnetWeakPass(res any) any {
	var (
		t   = res.(models.HydraTask)
		msg = models.MSG{
			Name: "Telnet WeakPassword",
			Type: "WeakPassword",
			Payload: models.Dict{
				User:     t.Dict.User,
				Password: t.Dict.Password,
			},
			Target: t.Target,
		}
	)

	c := &telnet.Client{
		UserName:     t.Dict.User,
		Password:     t.Dict.Password,
		LastResponse: "",
		ServerType:   telnet.UsernameAndPassword,
	}

	if err := c.Connect(t.Target, t.Config.Timeout); err != nil {
		return nil
	}

	if err := c.Login(); err != nil {
		return nil
	}
	return msg
}
