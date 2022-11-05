// Based on gwind/go-netlink-socket-monitor.go
// https://gist.github.com/gwind/05f5f649d93e6015cf47ffa2b2fd9713

package main

import (
	"fmt"
	"net"
	"syscall"
	"unsafe"

	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink/nl"
)

const TCPF_ALL = 0xFFF

const (
	INET_DIAG_NONE = iota
	INET_DIAG_MEMINFO
	INET_DIAG_INFO
	INET_DIAG_VEGASINFO
	INET_DIAG_CONG
	INET_DIAG_TOS
	INET_DIAG_TCLASS
	INET_DIAG_SKMEMINFO
	INET_DIAG_SHUTDOWN
	INET_DIAG_DCTCPINFO
	INET_DIAG_PROTOCOL
	INET_DIAG_SKV6ONLY
)

// PoolingではなくSubscribeしてPushできるか調べてみます。

// some comments added from: https://man7.org/linux/man-pages/man7/sock_diag.7.html
func main() {
	fmt.Println("Starting gonlmon.go")

	// The request starts with a struct nlmsghdr header described in
	// netlink(7) with nlmsg_type field set to SOCK_DIAG_BY_FAMILY.
	//
	// If the nlmsg_flags field of the struct nlmsghdr header has the
	// NLM_F_DUMP flag set, it means that a list of sockets is being
	// requested; otherwise it is a query about an individual socket.
	req := nl.NewNetlinkRequest(SOCK_DIAG_BY_FAMILY, syscall.NLM_F_DUMP)

	// It is followed by a header specific to the address family that
	// starts with a common part shared by all address families:
	// struct sock_diag_req {
	// 	__u8 sdiag_family;
	// 	__u8 sdiag_protocol;
	// };
	msg := NewInetDiagReqV2(
		// For IPv4 and IPv6 sockets, the request is represented in the
    	// following structure:
        //   struct inet_diag_req_v2 {
        //       __u8    sdiag_family;
        //       __u8    sdiag_protocol;
        //       __u8    idiag_ext;
        //       __u8    pad;
        //       __u32   idiag_states;
        //       struct inet_diag_sockid id;
        //   };

		// sdiag_family: This should be set to either AF_INET or AF_INET6
		syscall.AF_INET,		// SDiagFamily
		// sdiag_protocol: This should be set to one of IPPROTO_TCP, IPPROTO_UDP, or IPPROTO_UDPLITE.
		syscall.IPPROTO_TCP,	// SDiagProtocol
		
		// idiag_states (IDiagStates)
		// This is a bit mask that defines a filter of socket states.
		// Only those sockets whose states are in this mask will be
		// reported. Ignored when querying for an individual socket.
		// See TcpStatesMap later in this code for list of flags.
		//// Everything
		//TCPF_ALL,
		//// Ignore TCP_SYN_RECV, TCP_TIME_WAIT, TCP_CLOSE, TCP_CLOSE_WAIT
		TCPF_ALL & ^((1<<TCP_SYN_RECV)|(1<<TCP_TIME_WAIT)|(1<<TCP_CLOSE)|(1<<TCP_CLOSE_WAIT)),
		//// Ignore TCP_SYN_RECV
		//TCPF_ALL & ^(1<<TCP_SYN_RECV),
		//// Only TCP_TIME_WAIT
		//(1<<TCP_TIME_WAIT),
	)
	//// Set more InetDiagReqV2 params if required
	// idiag_ext: This is a set of flags defining what kind of extended
	// information to report.  Each requested kind of information
	// is reported back as a netlink attribute as described below:
	// INET_DIAG_TOS, INET_DIAG_TCLASS, INET_DIAG_MEMINFO, INET_DIAG_SKMEMINFO, INET_DIAG_CONG
	// INET_DIAG_INFO
	//   The payload associated with this attribute is
	//   specific to the address family.  For TCP sockets,
	//   it is an object of type struct tcp_info.
	// The values are defined in inet_diag.h (the INET_DIAG_*-constants).
	msg.IDiagExt |= (1 << (INET_DIAG_INFO - 1))
	req.AddData(msg)

	res, err := req.Execute(syscall.NETLINK_INET_DIAG, 0)
	if err != nil {
		logrus.Error("req.Execute error: ", err)
		return
	}

	fmt.Println(IndexInetDiagMsg())
	for _, data := range res {
		m := ParseInetDiagMsg(data)
		fmt.Println(m)

		//m := (*InetDiagMsg)(unsafe.Pointer(&data[0]))
		//fmt.Println("\n\n", i, data)
		//fmt.Println(i, m)
		//fmt.Printf(
		//	"%3d %-6d %-5d %-5d %-8d %-8d %-10d %-10d %-5d %-10d %s\n",
		//	i,
		//	m.idiag_family,
		//	m.idiag_state,
		//	m.idiag_timer,
		//	m.idiag_retrans,
		//	m.idiag_expires,
		//	m.idiag_rqueue,
		//	m.idiag_wqueue,
		//	m.idiag_uid,
		//	m.idiag_inode,
		//	showID(m.id),
		//)
	}
}

// inetdiag.go --------------------------------------------//

const (
	SizeofInetDiagReqV2 = 0x38
)

const (
	TCPDIAG_GETSOCK     = 18 // linux/inet_diag.h
	SOCK_DIAG_BY_FAMILY = 20 // linux/sock_diag.h
)

// netinet/tcp.h
const (
	_               = iota
	TCP_ESTABLISHED = iota
	TCP_SYN_SENT
	TCP_SYN_RECV
	TCP_FIN_WAIT1
	TCP_FIN_WAIT2
	TCP_TIME_WAIT
	TCP_CLOSE
	TCP_CLOSE_WAIT
	TCP_LAST_ACK
	TCP_LISTEN
	TCP_CLOSING
)

const (
	TCP_ALL = 0xFFF
)

var TcpStatesMap = map[uint8]string{
	TCP_ESTABLISHED: "established",
	TCP_SYN_SENT:    "syn_sent",
	TCP_SYN_RECV:    "syn_recv",
	TCP_FIN_WAIT1:   "fin_wait1",
	TCP_FIN_WAIT2:   "fin_wait2",
	TCP_TIME_WAIT:   "time_wait",
	TCP_CLOSE:       "close",
	TCP_CLOSE_WAIT:  "close_wait",
	TCP_LAST_ACK:    "last_ack",
	TCP_LISTEN:      "listen",
	TCP_CLOSING:     "closing",
}

var DiagFamilyMap = map[uint8]string{
	syscall.AF_INET:  "tcp",
	syscall.AF_INET6: "tcp6",
}

type be16 [2]byte

func (v be16) Int() int {
	// (*(*[SizeofInetDiagReqV2]byte)(unsafe.Pointer(req)))[:]
	v2 := (*(*uint16)(unsafe.Pointer(&v)))
	return int(nl.Swap16(v2))
}

type be32 [4]byte

// linux/inet_diag.h
type InetDiagSockId struct {
	IDiagSPort  be16
	IDiagDPort  be16
	IDiagSrc    [4]be32
	IDiagDst    [4]be32
	IDiagIf     uint32
	IDiagCookie [2]uint32
}

func (id *InetDiagSockId) SrcIPv4() net.IP {
	return ipv4(id.IDiagSrc[0])
}

func (id *InetDiagSockId) DstIPv4() net.IP {
	return ipv4(id.IDiagDst[0])
}

func (id *InetDiagSockId) SrcIPv6() net.IP {
	return ipv6(id.IDiagSrc)
}

func (id *InetDiagSockId) DstIPv6() net.IP {
	return ipv6(id.IDiagDst)
}

func (id *InetDiagSockId) SrcIP() net.IP {
	return ip(id.IDiagSrc)
}

func (id *InetDiagSockId) DstIP() net.IP {
	return ip(id.IDiagDst)
}

func ip(bytes [4]be32) net.IP {
	if isIpv6(bytes) {
		return ipv6(bytes)
	} else {
		return ipv4(bytes[0])
	}
}

func isIpv6(original [4]be32) bool {
	for i := 1; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if original[i][j] != 0 {
				return true
			}
		}
	}
	return false
}

func ipv4(original be32) net.IP {
	return net.IPv4(original[0], original[1], original[2], original[3])
}

func ipv6(original [4]be32) net.IP {
	ip := make(net.IP, net.IPv6len)
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			ip[4*i+j] = original[i][j]
		}
	}
	return ip
}

func (id *InetDiagSockId) String() string {
	return fmt.Sprintf("%s:%d -> %s:%d", id.SrcIP().String(), id.IDiagSPort.Int(), id.DstIP().String(), id.IDiagDPort.Int())
}

type InetDiagReqV2 struct {
	SDiagFamily   uint8
	SDiagProtocol uint8
	IDiagExt      uint8
	Pad           uint8
	IDiagStates   uint32
	Id            InetDiagSockId
}

func (req *InetDiagReqV2) Serialize() []byte {
	return (*(*[SizeofInetDiagReqV2]byte)(unsafe.Pointer(req)))[:]
}

func (req *InetDiagReqV2) Len() int {
	return SizeofInetDiagReqV2
}

func NewInetDiagReqV2(family, protocol uint8, states uint32) *InetDiagReqV2 {
	return &InetDiagReqV2{
		SDiagFamily:   family,
		SDiagProtocol: protocol,
		IDiagStates:   states,
	}
}

type InetDiagMsg struct {
	IDiagFamily  uint8
	IDiagState   uint8
	IDiagTimer   uint8
	IDiagRetrans uint8
	Id           InetDiagSockId
	IDiagExpires uint32
	IDiagRqueue  uint32
	IDiagWqueue  uint32
	IDiagUid     uint32
	IDiagInode   uint32
}

func IndexInetDiagMsg() string {
	s := "Family,       State, InetDiagSockId\n"
	s += "--------------------------------------------------------"
	return s
}

func (msg *InetDiagMsg) String() string {
	return fmt.Sprintf("%6s, %11s, %s",
		DiagFamilyMap[msg.IDiagFamily],
		TcpStatesMap[msg.IDiagState],
		msg.Id.String(),
	)
}

func ParseInetDiagMsg(data []byte) *InetDiagMsg {
	return (*InetDiagMsg)(unsafe.Pointer(&data[0]))
}
