package hydra

import (
	"github.com/yqcs/prismx/core/models"
	"github.com/yqcs/prismx/utils/go-snmp"
	"github.com/yqcs/prismx/utils/netUtils"
)

func SNMPWeakPass(res any) any {
	var (
		t   = res.(models.HydraTask)
		msg = models.MSG{
			Name:    "SNMP Unauthorized",
			Type:    "Unauthorized",
			Payload: models.Dict{},
			Target:  t.Target,
		}
	)

	// Open a UDP connection to the target
	conn, err := netUtils.SendDialTimeout("udp", t.Target, t.Config.Timeout)
	if err != nil {
		return nil
	}

	snmp := &go_snmp.GoSNMP{t.Target, "public", go_snmp.Version2c, t.Config.Timeout, conn}

	resp, err := snmp.Get(".1.3.6.1.2.1.1.1.0")
	if err != nil {
		return nil
	}
	for _, v := range resp.Variables {
		switch v.Type {
		case go_snmp.OctetString:
			return msg
		}
	}
	return nil
}
