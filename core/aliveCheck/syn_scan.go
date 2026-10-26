//go:build !nosyn

package aliveCheck

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/jackpal/gateway"
	"github.com/libp2p/go-netroute"
	"io"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"
)

// 全局缓存：只获取一次路由和网卡信息（单网卡场景，程序运行期间不会变）
var (
	globalSrcIp   net.IP
	globalSrcIp6  net.IP
	globalSrcMac  net.HardwareAddr
	globalGw      net.IP
	globalDevName string
	globalErr     error
	once          sync.Once
)

// SynScanner SYN 扫描器
type SynScanner struct {
	srcMac, gwMac net.HardwareAddr // MAC 地址
	devName       string           // 网卡设备名（pcap 格式: \Device\NPF_{...}）

	// 源 IP 地址
	srcIp, srcIp6 net.IP

	// pcap 句柄
	handle *pcap.Handle

	// 序列化选项和缓冲区复用
	opts    gopacket.SerializeOptions
	bufPool *sync.Pool

	// 存活端口结果通道
	openPortChan chan string
	ctx          context.Context
	cancel       context.CancelFunc

	// MAC 地址缓存（避免重复 ARP）
	macCache     map[string]net.HardwareAddr
	macCacheLock sync.RWMutex

	// 速率控制
	rateLimiter *time.Ticker

	// 状态
	isDone bool
}

// initOnce 只获取一次路由和网卡信息（程序运行期间不会变）
func initOnce(firstIp net.IP) {
	once.Do(func() {
		globalSrcIp, globalSrcIp6, globalSrcMac, globalGw, globalDevName, globalErr = getRouter(firstIp)
	})
}

// NewSynScanner 创建 SYN 扫描器
// firstIp: 用于路由选择的第一个目标 IP（仅第一次有效，后续复用缓存）
// rate: 每秒发包数，推荐 20000-50000，SYN 扫描可以跑很高的速率
func NewSynScanner(firstIp net.IP, rate int) (*SynScanner, error) {
	if rate <= 0 {
		rate = 30000 // 默认 30000 pps，高速扫描
	}

	// 只获取一次路由和网卡信息（单网卡场景，程序运行期间不会变）
	initOnce(firstIp)
	if globalErr != nil {
		return nil, globalErr
	}
	if globalDevName == "" {
		return nil, errors.New("failed to get network device")
	}

	srcIp := globalSrcIp
	srcIp6 := globalSrcIp6
	srcMac := globalSrcMac
	gw := globalGw
	devName := globalDevName

	ss := &SynScanner{
		opts: gopacket.SerializeOptions{
			FixLengths:       true,
			ComputeChecksums: true,
		},
		srcIp:   srcIp,
		srcIp6:  srcIp6,
		srcMac:  srcMac,
		devName: devName,
		bufPool: &sync.Pool{
			New: func() interface{} {
				return gopacket.NewSerializeBuffer()
			},
		},
		openPortChan: make(chan string, 1000),
		macCache:     make(map[string]net.HardwareAddr),
		rateLimiter:  time.NewTicker(time.Second / time.Duration(rate)),
	}

	// 创建上下文
	ss.ctx, ss.cancel = context.WithCancel(context.Background())

	// 打开网卡（每个包最大 1024 字节，不开启混杂模式，永久阻塞）
	// devName 是通过 pcap.FindAllDevs() 得到的正确格式: \Device\NPF_{GUID}
	handle, err := pcap.OpenLive(devName, 1024, false, pcap.BlockForever)
	if err != nil {
		return nil, fmt.Errorf("pcap.OpenLive failed: %w", err)
	}
	ss.handle = handle

	// 设置 BPF 过滤，只抓我们关心的包，大幅降低性能开销
	err = handle.SetBPFFilter(fmt.Sprintf(
		"ether dst %s && (arp || tcp[tcpflags] == tcp-syn|tcp-ack)",
		srcMac.String(),
	))
	if err != nil {
		handle.Close()
		return nil, err
	}

	// 启动接收协程
	go ss.recv()

	// 获取网关 MAC 地址
	if gw != nil && !gw.Equal(net.IPv4zero) {
		dstMac, err := ss.getHwAddr(gw)
		if err != nil {
			handle.Close()
			return nil, err
		}
		ss.gwMac = dstMac
	}

	return ss, nil
}

// Scan 扫描单个 IP 和端口
func (ss *SynScanner) Scan(dstIp net.IP, dstPort int) error {
	if ss.isDone {
		return io.EOF
	}

	// 速率控制
	select {
	case <-ss.rateLimiter.C:
	case <-ss.ctx.Done():
		return ss.ctx.Err()
	}

	// 获取目标 MAC 地址
	var dstMac net.HardwareAddr
	if ss.gwMac != nil {
		// 通过网关，直接用网关 MAC
		dstMac = ss.gwMac
	} else {
		// 内网 IP，查缓存或发 ARP
		ipStr := dstIp.String()
		ss.macCacheLock.RLock()
		cachedMac, ok := ss.macCache[ipStr]
		ss.macCacheLock.RUnlock()
		if ok {
			dstMac = cachedMac
		} else {
			var err error
			dstMac, err = ss.getHwAddr(dstIp)
			if err != nil {
				return err
			}
			ss.macCacheLock.Lock()
			ss.macCache[ipStr] = dstMac
			ss.macCacheLock.Unlock()
		}
	}

	// 构建以太网层
	eth := layers.Ethernet{
		SrcMAC:       ss.srcMac,
		DstMAC:       dstMac,
		EthernetType: layers.EthernetTypeIPv4,
	}

	// 构建 IP 层（支持 IPv4/IPv6）
	var ip4 *layers.IPv4
	var ip6 *layers.IPv6
	if dstIp.To4() != nil {
		ip4 = &layers.IPv4{
			SrcIP:    ss.srcIp,
			DstIP:    dstIp,
			Version:  4,
			TTL:      128,
			Id:       uint16(40000 + rand.Intn(10000)),
			Flags:    layers.IPv4DontFragment,
			Protocol: layers.IPProtocolTCP,
		}
	} else {
		eth.EthernetType = layers.EthernetTypeIPv6
		ip6 = &layers.IPv6{
			Version:    6,
			NextHeader: layers.IPProtocolTCP,
			HopLimit:   64,
			SrcIP:      ss.srcIp6,
			DstIP:      dstIp,
		}
	}

	// 构建 TCP SYN 包
	tcp := layers.TCP{
		SrcPort: layers.TCPPort(49000 + rand.Intn(10000)), // 随机源端口，范围 49000-58999
		DstPort: layers.TCPPort(dstPort),
		SYN:     true,
		Window:  65280,
		Seq:     uint32(500000 + rand.Intn(10000)),
		Options: []layers.TCPOption{
			{
				OptionType:   layers.TCPOptionKindMSS,
				OptionLength: 4,
				OptionData:   []byte{0x05, 0x50}, // MSS = 1360
			},
			{OptionType: layers.TCPOptionKindNop},
			{
				OptionType:   layers.TCPOptionKindWindowScale,
				OptionLength: 3,
				OptionData:   []byte{0x08},
			},
			{OptionType: layers.TCPOptionKindNop},
			{OptionType: layers.TCPOptionKindNop},
			{
				OptionType:   layers.TCPOptionKindSACKPermitted,
				OptionLength: 2,
			},
		},
	}

	// 发送数据包
	if ip4 != nil {
		tcp.SetNetworkLayerForChecksum(ip4)
		return ss.send(&eth, ip4, &tcp)
	} else {
		tcp.SetNetworkLayerForChecksum(ip6)
		return ss.send(&eth, ip6, &tcp)
	}
}

// GetResults 获取存活端口结果通道
func (ss *SynScanner) GetResults() <-chan string {
	return ss.openPortChan
}

// Wait 等待剩余响应
func (ss *SynScanner) Wait() {
	// 等待最后一个包的响应（2 秒）
	for i := 0; i < 20; i++ {
		time.Sleep(100 * time.Millisecond)
	}
}

// Close 关闭扫描器
func (ss *SynScanner) Close() {
	ss.isDone = true
	ss.rateLimiter.Stop()

	if ss.handle != nil {
		// Linux 下 pcap BlockForever 模式无法自动退出，需要发个包唤醒
		eth := layers.Ethernet{
			SrcMAC:       ss.srcMac,
			DstMAC:       ss.srcMac,
			EthernetType: layers.EthernetTypeARP,
		}
		arp := layers.ARP{
			AddrType:          layers.LinkTypeEthernet,
			Protocol:          layers.EthernetTypeIPv4,
			HwAddressSize:     6,
			ProtAddressSize:   4,
			Operation:         layers.ARPReply,
			SourceHwAddress:   []byte(ss.srcMac),
			SourceProtAddress: []byte(ss.srcIp),
			DstHwAddress:      []byte(ss.srcMac),
			DstProtAddress:    []byte(ss.srcIp),
		}
		handle, _ := pcap.OpenLive(ss.devName, 1024, false, time.Second)
		if handle != nil {
			buf := ss.bufPool.Get().(gopacket.SerializeBuffer)
			gopacket.SerializeLayers(buf, ss.opts, &eth, &arp)
			handle.WritePacketData(buf.Bytes())
			handle.Close()
			buf.Clear()
			ss.bufPool.Put(buf)
		}
		ss.handle.Close()
	}

	ss.cancel()
	close(ss.openPortChan)
}

// send 发送数据包
func (ss *SynScanner) send(l ...gopacket.SerializableLayer) error {
	buf := ss.bufPool.Get().(gopacket.SerializeBuffer)
	defer func() {
		buf.Clear()
		ss.bufPool.Put(buf)
	}()
	if err := gopacket.SerializeLayers(buf, ss.opts, l...); err != nil {
		return err
	}
	return ss.handle.WritePacketData(buf.Bytes())
}

// getHwAddr 获取目标 MAC 地址（自动识别 IPv4/IPv6）
func (ss *SynScanner) getHwAddr(ip net.IP) (net.HardwareAddr, error) {
	if ip.To4() != nil {
		return ss.getHwAddrV4(ip)
	}
	return ss.getHwAddrV6(ip)
}

// getHwAddrV4 IPv4 ARP 请求
func (ss *SynScanner) getHwAddrV4(arpDst net.IP) (net.HardwareAddr, error) {
	// 构建 ARP 请求
	eth := layers.Ethernet{
		SrcMAC:       ss.srcMac,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeARP,
	}
	arp := layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         layers.ARPRequest,
		SourceHwAddress:   []byte(ss.srcMac),
		SourceProtAddress: []byte(ss.srcIp),
		DstHwAddress:      []byte{0, 0, 0, 0, 0, 0},
		DstProtAddress:    []byte(arpDst),
	}

	if err := ss.send(&eth, &arp); err != nil {
		return nil, err
	}

	// 等待 ARP 回复，最多 600ms
	start := time.Now()
	for {
		ss.macCacheLock.RLock()
		mac, ok := ss.macCache[arpDst.String()]
		ss.macCacheLock.RUnlock()
		if ok {
			return mac, nil
		}
		if time.Since(start) > 600*time.Millisecond {
			return nil, errors.New("arp timeout")
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// getHwAddrV6 IPv6 NDP 请求（简化版，暂不实现）
func (ss *SynScanner) getHwAddrV6(ip net.IP) (net.HardwareAddr, error) {
	return nil, errors.New("ipv6 not supported yet")
}

// recv 接收数据包
func (ss *SynScanner) recv() {
	var ethLayer layers.Ethernet
	var ipLayer layers.IPv4
	var ipv6Layer layers.IPv6
	var tcpLayer layers.TCP
	var arpLayer layers.ARP

	parser := gopacket.NewDecodingLayerParser(
		layers.LayerTypeEthernet,
		&ethLayer,
		&ipLayer,
		&ipv6Layer,
		&tcpLayer,
		&arpLayer,
	)

	var foundLayerTypes []gopacket.LayerType
	var data []byte
	var err error

	for {
		data, _, err = ss.handle.ReadPacketData()
		if err != nil {
			if err == io.EOF || ss.isDone {
				return
			}
			continue
		}

		err = parser.DecodeLayers(data, &foundLayerTypes)
		if len(foundLayerTypes) == 0 {
			continue
		}

		// 处理 ARP 回复
		if arpLayer.SourceProtAddress != nil {
			ipStr := net.IP(arpLayer.SourceProtAddress).String()
			ss.macCacheLock.Lock()
			ss.macCache[ipStr] = arpLayer.SourceHwAddress
			ss.macCacheLock.Unlock()
			arpLayer.SourceProtAddress = nil // 重置标记
			continue
		}

		// 处理 TCP 包 - 只抓源端口在 49000-58999 范围的（我们发出的 SYN 包的源端口）
		if tcpLayer.DstPort >= 49000 && tcpLayer.DstPort <= 58999 {
			// 只关心 SYN-ACK（端口开放）
			if tcpLayer.SYN && tcpLayer.ACK {
				var srcIp net.IP
				if ethLayer.EthernetType == layers.EthernetTypeIPv6 {
					srcIp = ipv6Layer.SrcIP
				} else {
					srcIp = ipLayer.SrcIP
				}
				result := fmt.Sprintf("%s:%d", srcIp.String(), tcpLayer.SrcPort)

				// 发送结果（非阻塞）
				select {
				case ss.openPortChan <- result:
				default:
				}

				// 回复 RST，主动断开连接，避免半开连接堆积
				ss.sendRST(&ethLayer, &ipLayer, &tcpLayer, &ipv6Layer)
			}
		}
	}
}

// sendRST 发送 RST 包断开连接
func (ss *SynScanner) sendRST(
	ethRecv *layers.Ethernet,
	ipRecv *layers.IPv4,
	tcpRecv *layers.TCP,
	ip6Recv *layers.IPv6,
) {
	eth := layers.Ethernet{
		SrcMAC:       ss.srcMac,
		DstMAC:       ethRecv.SrcMAC,
		EthernetType: ethRecv.EthernetType,
	}

	var ip4 *layers.IPv4
	var ip6 *layers.IPv6

	if ethRecv.EthernetType == layers.EthernetTypeIPv6 {
		ip6 = &layers.IPv6{
			Version:    6,
			NextHeader: layers.IPProtocolTCP,
			HopLimit:   64,
			SrcIP:      ss.srcIp6,
			DstIP:      ip6Recv.SrcIP,
		}
	} else {
		ip4 = &layers.IPv4{
			SrcIP:    ss.srcIp,
			DstIP:    ipRecv.SrcIP,
			Version:  4,
			TTL:      64,
			Protocol: layers.IPProtocolTCP,
		}
	}

	// 发送 RST + ACK
	tcp := layers.TCP{
		SrcPort: tcpRecv.DstPort,
		DstPort: tcpRecv.SrcPort,
		RST:     true,
		ACK:     true,
		Seq:     tcpRecv.Ack,
		Ack:     tcpRecv.Seq + 1,
		Window:  0,
	}

	if ip4 != nil {
		tcp.SetNetworkLayerForChecksum(ip4)
		ss.send(&eth, ip4, &tcp)
	} else if ip6 != nil {
		tcp.SetNetworkLayerForChecksum(ip6)
		ss.send(&eth, ip6, &tcp)
	}
}

// getRouter 根据目标 IP 获取路由信息（源 IP、源 MAC、网关、网卡名）
// getDevByIp 根据 IP 获取 pcap 设备名（关键：返回 \Device\NPF_{...} 格式）
func getDevByIp(ip net.IP) (devName string, err error) {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		return "", err
	}
	for _, d := range devices {
		for _, address := range d.Addresses {
			_ip := address.IP
			if _ip != nil && _ip.IsGlobalUnicast() && _ip.Equal(ip) {
				return d.Name, nil // 这里返回的是 pcap 的正确设备名格式: \Device\NPF_{GUID}
			}
		}
	}
	return "", errors.New("can not find pcap dev")
}

// getIfaceMac 根据网卡 IP 获取源 IP、源 IPv6、MAC 地址
func getIfaceMac(ifaceAddr net.IP) (src net.IP, src6 net.IP, mac net.HardwareAddr) {
	interfaces, _ := net.Interfaces()
	var s4 = ifaceAddr.To4() != nil
	for _, iface := range interfaces {
		var ip, ip6 net.IP
		if addrs, err := iface.Addrs(); err == nil {
			for _, addr := range addrs {
				var ipNet = addr.(*net.IPNet)
				if !ipNet.IP.IsGlobalUnicast() {
					continue
				}
				if ipNet.IP.To4() != nil {
					if !s4 {
						ip = ipNet.IP.To4()
					}
				} else {
					if s4 {
						ip6 = ipNet.IP
					}
				}
				if ipNet.Contains(ifaceAddr) {
					if s4 {
						ip = ipNet.IP.To4()
					} else {
						ip6 = ipNet.IP
					}
					mac = iface.HardwareAddr
				}
			}
			if mac != nil {
				src = ip
				src6 = ip6
				return
			}
		}
	}
	return
}

// getRouter 根据目标 IP 获取路由信息
// 核心：通过 pcap.FindAllDevs() 获取正确的 \Device\NPF_{...} 格式网卡名
func getRouter(dst net.IP) (srcIp, srcIp6 net.IP, srcMac net.HardwareAddr, gw net.IP, devName string, err error) {
	// 同网段：直接获取源 IP 的网卡信息
	srcIp, srcIp6, srcMac = getIfaceMac(dst)
	if srcIp == nil {
		// 不同网段：通过路由表找网关
		r, errRoute := netroute.New()
		if errRoute == nil {
			var sip net.IP
			_, gw, sip, errRoute = r.Route(dst)
			if errRoute == nil {
				srcIp, srcIp6, srcMac = getIfaceMac(sip)
			}
		}
		// 如果路由表查找失败，尝试获取默认网关
		if errRoute != nil || srcMac == nil {
			gw, errRoute = gateway.DiscoverGateway()
			if errRoute == nil {
				srcIp, srcIp6, srcMac = getIfaceMac(gw)
			}
		}
	}

	if gw.To4() != nil {
		gw = gw.To4()
	}
	if srcIp.To4() != nil {
		srcIp = srcIp.To4()
	}

	// 关键：通过 pcap.FindAllDevs() 获取正确的 pcap 网卡名格式
	devName, err = getDevByIp(srcIp)
	if err != nil && srcIp6 != nil {
		// 如果 IPv4 失败，尝试用 IPv6
		devName, err = getDevByIp(srcIp6)
	}

	if (srcIp == nil && srcIp6 == nil) || err != nil || srcMac == nil {
		if err == nil {
			err = errors.New("router not found")
		}
		return nil, nil, nil, nil, "", fmt.Errorf("no router found: %w", err)
	}

	return srcIp, srcIp6, srcMac, gw, devName, nil
}

// CheckSYNAvailable 检测当前环境是否支持 SYN 扫描
func CheckSYNAvailable() bool {
	// 尝试枚举 pcap 设备
	ifaces, err := pcap.FindAllDevs()
	if err != nil || len(ifaces) == 0 {
		return false
	}
	// 至少有一个非回环设备
	for _, iface := range ifaces {
		if strings.Contains(iface.Name, "NPF") { // Windows pcap 设备
			return true
		}
		if len(iface.Addresses) > 0 && !iface.Addresses[0].IP.IsLoopback() {
			return true
		}
	}
	return len(ifaces) > 0
}
