package hydra

import (
	"errors"
	"sync"

	"github.com/yqcs/prismx/core/models"
	"github.com/yqcs/prismx/utils/go-rdp/core"
	"github.com/yqcs/prismx/utils/go-rdp/protocol/nla"
	"github.com/yqcs/prismx/utils/go-rdp/protocol/pdu"
	"github.com/yqcs/prismx/utils/go-rdp/protocol/rfb"
	"github.com/yqcs/prismx/utils/go-rdp/protocol/sec"
	"github.com/yqcs/prismx/utils/go-rdp/protocol/t125"
	"github.com/yqcs/prismx/utils/go-rdp/protocol/tpkt"
	"github.com/yqcs/prismx/utils/go-rdp/protocol/x224"
	"github.com/yqcs/prismx/utils/netUtils"
)

type rdpClient struct {
	Host string
	tpkt *tpkt.TPKT
	x224 *x224.X224
	mcs  *t125.MCSClient
	sec  *sec.Client
	pdu  *pdu.Client
	vnc  *rfb.RFB
}

func RdpWeakPass(res any) any {

	var (
		t   = res.(models.HydraTask)
		msg = models.MSG{
			Name: "RDP WeakPassword",
			Type: "WeakPassword",
			Payload: models.Dict{
				User:     t.Dict.User,
				Password: t.Dict.Password,
			},
			Target: t.Target,
		}
	)

	g := rdpClient{Host: t.Target}

	conn, err := netUtils.SendDialTimeout("tcp", g.Host, t.Config.Timeout)
	if err != nil {
		return nil
	}
	defer conn.Close()
	g.tpkt = tpkt.New(core.NewSocketLayer(conn), nla.NewNTLMv2("", t.Dict.User, t.Dict.Password))
	g.x224 = x224.New(g.tpkt)
	g.mcs = t125.NewMCSClient(g.x224)
	g.sec = sec.NewClient(g.mcs)
	g.pdu = pdu.NewClient(g.sec)

	g.sec.SetUser(t.Dict.User)
	g.sec.SetPwd(t.Dict.Password)
	g.sec.SetDomain("")

	g.tpkt.SetFastPathListener(g.sec)
	g.sec.SetFastPathListener(g.pdu)
	g.pdu.SetFastPathSender(g.tpkt)

	//g.x224.SetRequestedProtocol(x224.PROTOCOL_SSL)
	//g.x224.SetRequestedProtocol(x224.PROTOCOL_RDP)

	err = g.x224.Connect()
	if err != nil {
		return nil
	}
	wg := &sync.WaitGroup{}
	breakFlag := false
	wg.Add(1)

	g.pdu.On("error", func(e error) {
		err = e
		g.pdu.Emit("done")
	})
	g.pdu.On("close", func() {
		err = errors.New("close")
		g.pdu.Emit("done")
	})
	g.pdu.On("success", func() {
		err = nil
		g.pdu.Emit("done")
	})
	g.pdu.On("ready", func() {
		g.pdu.Emit("done")
	})
	g.pdu.On("update", func(rectangles []pdu.BitmapData) {
	})
	g.pdu.On("done", func() {
		if breakFlag == false {
			breakFlag = true
			wg.Done()
		}
	})
	wg.Wait()
	if err != nil {
		return nil
	}
	return msg
}
