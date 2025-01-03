package aliveCheck

import (
	"golang.org/x/net/icmp"
	"net"
	"os/exec"
	"os/user"
	"prismx_cli/utils/netUtils"
	"runtime"
	"strings"
	"time"
)

// HostAliveCheck 主机存活检测
func HostAliveCheck(ip string, ping bool, timeout time.Duration, fuzz bool) (endFlag bool) {
	u, err := user.Current()
	if err != nil || u == nil || ping || (runtime.GOOS != "windows" && (u.Gid != "0" || u.Uid != "0")) {
		endFlag = execCommandPing(ip)
	} else {
		//优先尝试监听本地icmp探测
		conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
		if err == nil {
			endFlag = RunIcmpLocal(ip, conn)
		} else {
			//尝试无监听icmp探测
			endFlag = icmpAlive(ip, timeout)
		}
	}
	//启动模糊检测
	if !endFlag && fuzz {
		endFlag = tcpScanPortCheck(ip, timeout)
	}
	return endFlag
}

func RunIcmpLocal(ip string, conn *icmp.PacketConn) bool {
	defer conn.Close()
	dst, err := net.ResolveIPAddr("ip", ip)
	if err != nil {
		return false
	}

	result := make(chan bool)
	go func() {
		for {
			buff := make([]byte, 100)
			_, sourceIP, _ := conn.ReadFrom(buff)
			if sourceIP != nil && sourceIP.String() == ip {
				result <- true
				return
			}
		}
	}()

	_, err = conn.WriteTo(makeMsg(ip), dst)
	if err != nil {
		return false
	}

	select {
	case <-result:
		return true
	case <-time.After(time.Second * 5):
		return false
	}
}

// icmpAlive 存活检测
func icmpAlive(ip string, timeout time.Duration) bool {

	startTime := time.Now()
	conn, err := netUtils.SendDialTimeout("ip4:icmp", ip, timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	if err = conn.SetDeadline(startTime.Add(timeout)); err != nil {
		return false
	}
	msg := makeMsg(ip)
	if _, err = conn.Write(msg); err != nil {
		return false
	}
	receive := make([]byte, 60)
	if _, err = conn.Read(receive); err != nil {
		return false
	}
	return true
}

func makeMsg(host string) []byte {
	msg := make([]byte, 40)
	id0, id1 := genIdentifier(host)
	msg[0] = 8
	msg[1] = 0
	msg[2] = 0
	msg[3] = 0
	msg[4], msg[5] = id0, id1
	msg[6], msg[7] = genSequence(1)
	check := checkSum(msg[0:40])
	msg[2] = byte(check >> 8)
	msg[3] = byte(check & 255)
	return msg
}

func checkSum(msg []byte) uint16 {
	sum := 0
	length := len(msg)
	for i := 0; i < length-1; i += 2 {
		sum += int(msg[i])*256 + int(msg[i+1])
	}
	if length%2 == 1 {
		sum += int(msg[length-1]) * 256
	}
	sum = (sum >> 16) + (sum & 0xffff)
	sum = sum + (sum >> 16)
	answer := uint16(^sum)
	return answer
}

func genSequence(v int16) (byte, byte) {
	ret1 := byte(v >> 8)
	ret2 := byte(v & 255)
	return ret1, ret2
}

func genIdentifier(host string) (byte, byte) {
	return host[0], host[1]
}
func execCommandPing(ip string) bool {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("ping", ip, "-n", "1", "-w", "200")
	case "linux":
		cmd = exec.Command("/bin/sh", "-c", "ping -c 1 "+ip)
	case "darwin":
		cmd = exec.Command("ping", ip, "-c", "1", "-W", "200")
	case "freebsd":
		cmd = exec.Command("ping", "-c", "1", "-W", "200", ip)
	case "openbsd":
		cmd = exec.Command("ping", "-c", "1", "-w", "200", ip)
	case "netbsd":
		cmd = exec.Command("ping", "-c", "1", "-w", "2", ip)
	default:
		cmd = exec.Command("ping", "-c", "1", ip)
	}

	if output, err := cmd.Output(); err == nil && strings.Contains(strings.ToLower(string(output)), "ttl=") {
		return true
	}
	return false
}
