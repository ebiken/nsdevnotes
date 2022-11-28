package main

import (
	"encoding/binary"
	"fmt"
	"syscall"

	"github.com/vishvananda/netlink/nl"
	"golang.org/x/sys/unix"
)

// Use syscall package to parse Netlink Messages
// https://pkg.go.dev/syscall#NlMsghdr
//	type NlMsghdr struct {
// 		Len   uint32
// 		Type  uint16
// 		Flags uint16
// 		Seq   uint32
// 		Pid   uint32
//	}
// https://pkg.go.dev/syscall#NetlinkMessage
//	type NetlinkMessage struct {
//		Header NlMsghdr
//		Data   []byte
//	}
//
//	type RtMsg struct {
//		Family   uint8
//		Dst_len  uint8
//		Src_len  uint8
//		Tos      uint8
//		Table    uint8
//		Protocol uint8
//		Scope    uint8
//		Type     uint8
//		Flags    uint32
//	}
//
//	type NetlinkRouteAttr struct {
//		Attr  RtAttr
//		Value []byte
//	}
//
//	type RtAttr struct {
//		Len  uint16
//		Type uint16
//	}
//
//// functions
//	func ParseNetlinkMessage(b []byte) ([]NetlinkMessage, error)
//	func ParseNetlinkRouteAttr(m *NetlinkMessage) ([]NetlinkRouteAttr, error)
//
//// RTNLGRP_* used to Subscribe
// RTNLGRP_NONE        = 0x0
// RTNLGRP_LINK        = 0x1
// RTNLGRP_NOTIFY      = 0x2
// RTNLGRP_NEIGH       = 0x3
// RTNLGRP_TC          = 0x4
// RTNLGRP_IPV4_IFADDR = 0x5
// RTNLGRP_IPV4_MROUTE = 0x6
// RTNLGRP_IPV4_ROUTE  = 0x7
// RTNLGRP_IPV4_RULE   = 0x8
// RTNLGRP_IPV6_IFADDR = 0x9
// RTNLGRP_IPV6_MROUTE = 0xa
// RTNLGRP_IPV6_ROUTE  = 0xb
// RTNLGRP_IPV6_IFINFO = 0xc
// RTNLGRP_IPV6_PREFIX = 0x12
// RTNLGRP_IPV6_RULE   = 0x13
// RTNLGRP_ND_USEROPT  = 0x14

func main() {
	fmt.Println("Starting gonlsub.go")

	// nl_linux.go: func Subscribe(protocol int, groups ...uint) (*NetlinkSocket, error) {
	// List of unix consts: https://pkg.go.dev/golang.org/x/sys/unix#pkg-constants

	//nlSock, err := nl.Subscribe(unix.NETLINK_ROUTE, unix.RTNLGRP_IPV4_ROUTE)
	nlSock, err := nl.Subscribe(
		syscall.NETLINK_ROUTE, 
		syscall.RTNLGRP_NONE,        // 0x0
		syscall.RTNLGRP_LINK,        // 0x1
		syscall.RTNLGRP_NOTIFY,      // 0x2
		syscall.RTNLGRP_NEIGH,       // 0x3
		syscall.RTNLGRP_TC,          // 0x4
		syscall.RTNLGRP_IPV4_IFADDR, // 0x5
		syscall.RTNLGRP_IPV4_MROUTE, // 0x6
		syscall.RTNLGRP_IPV4_ROUTE,  // 0x7
		syscall.RTNLGRP_IPV4_RULE,   // 0x8
		syscall.RTNLGRP_IPV6_IFADDR, // 0x9
		syscall.RTNLGRP_IPV6_MROUTE, // 0xa
		syscall.RTNLGRP_IPV6_ROUTE,  // 0xb
		syscall.RTNLGRP_IPV6_IFINFO, // 0xc
		syscall.RTNLGRP_IPV6_PREFIX, // 0x12
		syscall.RTNLGRP_IPV6_RULE,   // 0x13
		syscall.RTNLGRP_ND_USEROPT,  // 0x14
	)
	if err != nil {
		fmt.Println("Error on creating the socket: %v", err)
	}

	nlSock.SetReceiveTimeout(&unix.Timeval{Sec: 1, Usec: 0})
	for {
		//msgs, from, err := nlSock.Receive()
		msgs, _, _ := nlSock.Receive()
		if msgs != nil {
			fmt.Println("-----------------------------------")

			for _, msg := range msgs { // msg => NetlinkMessage
				//fmt.Printf("msg.Header %v\n", msg.Header)
				fmt.Printf("NlMsghdr | Len:%d, Type:%s, Flags:%x, Seq:%v, Pid:%v\n",
					msg.Header.Len,
					RtmMap[msg.Header.Type],
					msg.Header.Flags,
					msg.Header.Seq,
					msg.Header.Pid,
				)

				myrtm := getMyRtMsg(msg.Data[0:syscall.SizeofRtMsg])
				fmt.Printf("rtmsg: %v\n", myrtm)
				fmt.Printf("rtmsg: %s\n", myrtm.String())

				nras, _ := syscall.ParseNetlinkRouteAttr(&msg)
				for _, nra := range nras {
					//fmt.Printf("nra: %v\n", nra)
					s := fmt.Sprintf("RtAttr | Len:%v, Type:%s", nra.Attr.Len, RtaMap[nra.Attr.Type])
					switch nra.Attr.Type {
					case syscall.RTA_MULTIPATH:
						// RTA_MULTIPATH is array of [RTA_NEXTHOP + RTA_GATEWAY]
						parseRtNexthop := func(v []byte) ([]byte) {
							l := binary.LittleEndian.Uint16(v[0:2])
							i := array2int32(v[4:8])
							// https://pkg.go.dev/syscall#RtNexthop
							rtnexthop := syscall.RtNexthop {
								Len:     l, // uint16
								Flags:   uint8(v[2]), // uint8
								Hops:    uint8(v[3]), // uint8
								Ifindex: i, // int32
							}
							s += fmt.Sprintf("\n  | rtnexthop: Len:%v, Flags:%v, Hops:%v, Ifindex:%v",
								rtnexthop.Len,
								rtnexthop.Flags,
								rtnexthop.Hops,
								rtnexthop.Ifindex,
							)
							// Parse RTAs
							v = v[unix.SizeofRtNexthop:] // unix.SizeofRtNexthop(8)
							//rtalen := rtnexthop.Len - unix.SizeofRtNexthop
							rtalen := binary.LittleEndian.Uint16(v[0:2])
							rtatype := binary.LittleEndian.Uint16(v[2:4])
							s += fmt.Sprintf("\n  | RTA: Len:%v, Type:%s",
								rtalen,
								RtaMap[rtatype],
							)
							if rtatype == syscall.RTA_GATEWAY {
								if (rtalen-4) == 4 {
									s += fmt.Sprintf(", IPv4:%d.%d.%d.%d", v[4], v[5], v[6], v[7])
								} else if (rtalen-4) == 16 {
									s += fmt.Sprintf(", IPv6:%x", v[4:rtalen])
								} else {
									s += fmt.Sprintf(", UNKOWN:%x", v[4:rtalen])
								}
							} else {
								s += fmt.Sprintf(", Value:%x", v[4:rtalen])
							}
							v = v[rtalen:]
							return v
						}
						rest := nra.Value
						for len(rest) > 0 {
							rest = parseRtNexthop(rest)
						}
					case syscall.RTA_DST:
						v := nra.Value
						l := len(v)
						if l == 4 {
							s += fmt.Sprintf(", IPv4:%d.%d.%d.%d", v[0], v[1], v[2], v[3])
						} else if l == 16 {
							s += fmt.Sprintf(", IPv6:%x", v)
						} else {
							s += fmt.Sprintf(", UNKOWN:%x", v)
						}
					case syscall.RTA_TABLE:
						i := array2int32(nra.Value[0:4])
						s += fmt.Sprintf(", Value:%v", i)
					case 30: // RTA_NH_ID
						i := array2int32(nra.Value[0:4])
						s += fmt.Sprintf(", Value:%v", i)
					default:
						s += fmt.Sprintf(", Value:%v", nra.Value)
					}
					s += "\n"
					fmt.Print(s)
				}
			}
		}
	}

	//time.Sleep(4 * time.Second)
	// Close the socket
	nlSock.Close()
}

var RtmMap = map[uint16]string{
	syscall.RTM_NEWLINK:  "RTM_NEWLINK",
	syscall.RTM_DELLINK:  "RTM_DELLINK",
	syscall.RTM_GETLINK:  "RTM_GETLINK",
	syscall.RTM_SETLINK:  "RTM_SETLINK",
	syscall.RTM_NEWADDR:  "RTM_NEWADDR",
	syscall.RTM_DELADDR:  "RTM_DELADDR",
	syscall.RTM_GETADDR:  "RTM_GETADDR",
	syscall.RTM_NEWROUTE: "RTM_NEWROUTE",
	syscall.RTM_DELROUTE: "RTM_DELROUTE",
	syscall.RTM_GETROUTE: "RTM_GETROUTE",
	// RTM_*NEXTHOP below are not defined in syscall
	104: "RTM_NEWNEXTHOP",
	105: "RTM_DELNEXTHOP",
	106: "RTM_GETNEXTHOP",
}

// Routing message attributes (enum rtattr_type_t)
// include/uapi/linux/rtnetlink.h

var RtaMap = map[uint16]string{
	syscall.RTA_UNSPEC:    "RTA_UNSPEC",
	syscall.RTA_DST:       "RTA_DST",
	syscall.RTA_SRC:       "RTA_SRC",
	syscall.RTA_IIF:       "RTA_IIF",
	syscall.RTA_OIF:       "RTA_OIF",
	syscall.RTA_GATEWAY:   "RTA_GATEWAY",
	syscall.RTA_PRIORITY:  "RTA_PRIORITY",
	syscall.RTA_PREFSRC:   "RTA_PREFSRC",
	syscall.RTA_METRICS:   "RTA_METRICS",
	syscall.RTA_MULTIPATH: "RTA_MULTIPATH",
	//RTA_PROTOINFO, /* no longer used */
	syscall.RTA_FLOW:      "RTA_FLOW",
	syscall.RTA_CACHEINFO: "RTA_CACHEINFO",
	//RTA_SESSION, /* no longer used */
	//RTA_MP_ALGO, /* no longer used */
	syscall.RTA_TABLE: "RTA_TABLE",
	// RTA_* below are not defined in syscall
	16: "RTA_MARK",
    17: "RTA_MFC_STATS",
    18: "RTA_VIA",
    19: "RTA_NEWDST",
    20: "RTA_PREF",
    21: "RTA_ENCAP_TYPE",
    22: "RTA_ENCAP",
    23: "RTA_EXPIRES",
    24: "RTA_PAD",
    25: "RTA_UID",
    26: "RTA_TTL_PROPAGATE",
    27: "RTA_IP_PROTO",
    28: "RTA_SPORT",
    29: "RTA_DPORT",
    30: "RTA_NH_ID",
	//__RTA_MAX
}

var RtMsgFamilyMap = map[uint8]string{
	syscall.AF_INET:   "AF_INET",   // 0x2
	syscall.AF_INET6:  "AF_INET6",  // 0xa
	syscall.AF_PACKET: "AF_PACKET", // 0x11
	syscall.AF_ROUTE:  "AF_ROUTE",  // 0x10
	syscall.AF_UNIX:   "AF_UNIX",   // 0x1
	syscall.AF_UNSPEC: "AF_UNSPEC", // 0x0
}

var RtMsgProtoMap = map[uint8]string{
	syscall.RTPROT_UNSPEC:   "RTPROT_UNSPEC",   // 0x0
	syscall.RTPROT_REDIRECT: "RTPROT_REDIRECT", // 0x1
	syscall.RTPROT_KERNEL:   "RTPROT_KERNEL",   // 0x2
	syscall.RTPROT_BOOT:     "RTPROT_BOOT",     // 0x3
	syscall.RTPROT_STATIC:   "RTPROT_STATIC",   // 0x4
	syscall.RTPROT_GATED:    "RTPROT_GATED",    // 0x8
	syscall.RTPROT_RA:       "RTPROT_RA",       // 0x9
	syscall.RTPROT_MRT:      "RTPROT_MRT",      // 0xa
	syscall.RTPROT_ZEBRA:    "RTPROT_ZEBRA",    // 0xb
	syscall.RTPROT_BIRD:     "RTPROT_BIRD",     // 0xc
	syscall.RTPROT_DNROUTED: "RTPROT_DNROUTED", // 0xd
	syscall.RTPROT_XORP:     "RTPROT_XORP",     // 0xe
	syscall.RTPROT_NTK:      "RTPROT_NTK",      // 0xf
	syscall.RTPROT_DHCP:     "RTPROT_DHCP",     // 0x10
}

var RtMsgScopeMap = map[uint8]string{
	syscall.RT_SCOPE_UNIVERSE: "RT_SCOPE_UNIVERSE", // 0x0
	syscall.RT_SCOPE_SITE:     "RT_SCOPE_SITE",     // 0xc8
	syscall.RT_SCOPE_LINK:     "RT_SCOPE_LINK",     // 0xfd
	syscall.RT_SCOPE_HOST:     "RT_SCOPE_HOST",     // 0xfe
	syscall.RT_SCOPE_NOWHERE:  "RT_SCOPE_NOWHERE",  // 0xff
}

var RtMsgTypeMap = map[uint8]string{
	syscall.RTN_UNSPEC:      "RTN_UNSPEC",      // 0x0
	syscall.RTN_UNICAST:     "RTN_UNICAST",     // 0x1
	syscall.RTN_LOCAL:       "RTN_LOCAL",       // 0x2
	syscall.RTN_BROADCAST:   "RTN_BROADCAST",   // 0x3
	syscall.RTN_ANYCAST:     "RTN_ANYCAST",     // 0x4
	syscall.RTN_MULTICAST:   "RTN_MULTICAST",   // 0x5
	syscall.RTN_BLACKHOLE:   "RTN_BLACKHOLE",   // 0x6
	syscall.RTN_UNREACHABLE: "RTN_UNREACHABLE", // 0x7
	syscall.RTN_PROHIBIT:    "RTN_PROHIBIT",    // 0x8
	syscall.RTN_THROW:       "RTN_THROW",       // 0x9
	syscall.RTN_NAT:         "RTN_NAT",         // 0xa
	syscall.RTN_XRESOLVE:    "RTN_XRESOLVE",    // 0xb
}

type myRtMsg struct {
	Family   uint8
	Dst_len  uint8
	Src_len  uint8
	Tos      uint8
	Table    uint8
	Protocol uint8
	Scope    uint8
	Type     uint8
	Flags    uint32
}

// func getMyRtMsg(rtm syscall.RtMsg) myRtMsg {
func getMyRtMsg(data []byte) myRtMsg {
	var r myRtMsg
	r.Family = data[0]
	r.Dst_len = data[1]
	r.Src_len = data[2]
	r.Tos = data[3]
	r.Table = data[4]
	r.Protocol = data[5]
	r.Scope = data[6]
	r.Type = data[7]
	r.Flags = binary.LittleEndian.Uint32(data[8:12])

	return r
}

func (m *myRtMsg) String() string {
	s := "RtMsg |\n"
	s += fmt.Sprintf("  Family:   %s (%d)\n", RtMsgFamilyMap[m.Family], m.Family)
	s += fmt.Sprintf("  Dst_len:  %v\n", m.Dst_len)
	s += fmt.Sprintf("  Src_len:  %v\n", m.Src_len)
	s += fmt.Sprintf("  Tos:      %v\n", m.Tos)
	s += fmt.Sprintf("  Table:    %v\n", m.Table)
	s += fmt.Sprintf("  Protocol: %s (%d)\n", RtMsgProtoMap[m.Protocol], m.Protocol)
	s += fmt.Sprintf("  Scope:    %s (%d)\n", RtMsgScopeMap[m.Scope], m.Scope)
	s += fmt.Sprintf("  Type:     %s (%d)\n", RtMsgTypeMap[m.Type], m.Type)
	s += fmt.Sprintf("  Flags:    %v", m.Flags)
	return s
}

func array2int32(v []byte) int32 {
	var r int32
	r |= int32(v[0])
	r |= int32(v[1]) << 8
	r |= int32(v[2]) << 16
	r |= int32(v[3]) << 24
	return r
}
