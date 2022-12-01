# gonlsub

An elementry code written just to practice netlink in Golang which does below:

- subscribe to group `RTNLGRP_IPV4_ROUTE`
- receive broadcasted message
- parse and print netlink / rtnetlink messages
- parse array of RTA_NEXTHOP/RTA_GATEWAY inside RTA_MULTIPATH

## Next Hop and Next Hop Group

- rtmsg type (in NlMsghdr) for Next Hop is `RTM_NEWNEXTHOP` (104).
- So far it looks like this is NOT included in `RTNLGRP_*` thus cannot receive update.
- You can still receive `RTM_NEWROUTE` and `RTM_DELROUTE` using nexthop id
- This message will includes Attributes `RTA_NH_ID` and  `RTA_MULTIPATH`

### route using Next Hop Object

```
>> below messages will not show up even if you Subscribe to `NETLINK_ROUTE`.
> ip nexthop add id 11 via 172.20.105.173 dev eno1

>> below messages will be received via Subscribe.
> ip route add 10.11.12.13/32 nhid 11

-----------------------------------
NlMsghdr | Len:68, Type:RTM_NEWROUTE, Flags:600, Seq:1669595081, Pid:153947
rtmsg: {2 32 0 0 254 3 0 1 0}
rtmsg: RtMsg |
  Family:   AF_INET (2)
  Dst_len:  32
  Src_len:  0
  Tos:      0
  Table:    254
  Protocol: RTPROT_BOOT (3)
  Scope:    RT_SCOPE_UNIVERSE (0)
  Type:     RTN_UNICAST (1)
  Flags:    0
RtAttr | Len:8, Type:RTA_TABLE, Value:[254 0 0 0]
RtAttr | Len:8, Type:RTA_DST, Value:[10 11 12 13]
RtAttr | Len:8, Type:RTA_NH_ID, Value:[11 0 0 0]
RtAttr | Len:8, Type:RTA_GATEWAY, Value:[172 20 105 173]
RtAttr | Len:8, Type:RTA_OIF, Value:[5 0 0 0]
```

### route using Next Hop Group

`RTA_MULTIPATH` will be included when `RTA_NH_ID` is pointing to Next Hop Group. (nexthop id 3 in below example)

Note that you do NOT need to set `RTA_MULTIPATH` in sendmsg to configure route using nexthop group.
But when the route update is annouced, it will include `RTA_MULTIPATH`. (For backword compatibility?)

```
>> below messages will not show up even if you Subscribe to `NETLINK_ROUTE`.
> ip nexthop add id 1 via 172.20.105.172 dev eno1
> ip nexthop add id 2 via 172.20.105.173 dev eno1
> ip nexthop add id 3 group 1/2

>> below messages will be received via Subscribe.
> ip route add 10.11.12.13/32 nhid 3
> ip route del 10.11.12.13/32

-----------------------------------
NlMsghdr | Len:88, Type:RTM_NEWROUTE, Flags:600, Seq:1669611617, Pid:167576
rtmsg: {2 32 0 0 254 3 0 1 0}
rtmsg: RtMsg |
  Family:   AF_INET (2)
  Dst_len:  32
  Src_len:  0
  Tos:      0
  Table:    254
  Protocol: RTPROT_BOOT (3)
  Scope:    RT_SCOPE_UNIVERSE (0)
  Type:     RTN_UNICAST (1)
  Flags:    0
RtAttr | Len:8, Type:RTA_TABLE, Value:254
RtAttr | Len:8, Type:RTA_DST, IPv4:10.11.12.13
RtAttr | Len:8, Type:RTA_NH_ID, Value:3
RtAttr | Len:36, Type:RTA_MULTIPATH
  | rtnexthop: Len:16, Flags:0, Hops:0, Ifindex:5
  | RTA: Len:8, Type:RTA_GATEWAY, IPv4:172.20.105.172
  | rtnexthop: Len:16, Flags:0, Hops:0, Ifindex:5
  | RTA: Len:8, Type:RTA_GATEWAY, IPv4:172.20.105.173
-----------------------------------
NlMsghdr | Len:88, Type:RTM_DELROUTE, Flags:0, Seq:1669611617, Pid:167577
rtmsg: {2 32 0 0 254 3 0 1 0}
rtmsg: RtMsg |
  Family:   AF_INET (2)
  Dst_len:  32
  Src_len:  0
  Tos:      0
  Table:    254
  Protocol: RTPROT_BOOT (3)
  Scope:    RT_SCOPE_UNIVERSE (0)
  Type:     RTN_UNICAST (1)
  Flags:    0
RtAttr | Len:8, Type:RTA_TABLE, Value:254
RtAttr | Len:8, Type:RTA_DST, IPv4:10.11.12.13
RtAttr | Len:8, Type:RTA_NH_ID, Value:3
RtAttr | Len:36, Type:RTA_MULTIPATH
  | rtnexthop: Len:16, Flags:0, Hops:0, Ifindex:5
  | RTA: Len:8, Type:RTA_GATEWAY, IPv4:172.20.105.172
  | rtnexthop: Len:16, Flags:0, Hops:0, Ifindex:5
  | RTA: Len:8, Type:RTA_GATEWAY, IPv4:172.20.105.173

sendmsg(3, {msg_name={sa_family=AF_NETLINK, nl_pid=0, nl_groups=00000000}, msg_namelen=12, msg_iov=[{iov_base=[{nlmsg_len=52, nlmsg_type=RTM_NEWNEXTHOP, nlmsg_flags=NLM_F_REQUEST|NLM_F_ACK|NLM_F_EXCL|NLM_F_CREATE, nlmsg_seq=1669711532, nlmsg_pid=0}, {nh_family=AF_UNSPEC, nh_scope=RT_SCOPE_UNIVERSE, nh_protocol=RTPROT_UNSPEC, nh_flags=0}, [[{nla_len=8, nla_type=NHA_ID}, 3], [{nla_len=20, nla_type=NHA_GROUP}, [{id=1, weight=0}, {id=2, weight=0}]]]], iov_len=52}], msg_iovlen=1, msg_controllen=0, msg_flags=0}, 0) = 52
```

## IPv4 ROUTE (no next hop object)

```
> $ sudo ip route del 10.11.11.99/32
> $ sudo ip route add 10.11.11.99/32 via 172.20.104.1 dev eno1

nsdevnotes/examples/gonlsub$ go run gonlsub.go
Starting gonlsub.go
-----------------------------------
msg.Header {60 24 1536 1669553033 149079}
NlMsghdr | Len:60, Type:RTM_NEWROUTE, Flags:600, Seq:1669553033, Pid:149079
rtmsg: {2 32 0 0 254 3 0 1 0}
rtmsg: RtMsg |
  Family:   AF_INET (2)
  Dst_len:  32
  Src_len:  0
  Tos:      0
  Table:    254
  Protocol: RTPROT_BOOT (3)
  Scope:    RT_SCOPE_UNIVERSE (0)
  Type:     RTN_UNICAST (1)
  Flags:    0
RtAttr | Len:8, Type:RTA_TABLE, Value:[254 0 0 0]
RtAttr | Len:8, Type:RTA_DST, Value:[10 11 11 99]
RtAttr | Len:8, Type:RTA_GATEWAY, Value:[172 20 104 1]
RtAttr | Len:8, Type:RTA_OIF, Value:[5 0 0 0]
-----------------------------------
msg.Header {60 25 0 1669553060 149086}
NlMsghdr | Len:60, Type:RTM_DELROUTE, Flags:0, Seq:1669553060, Pid:149086
rtmsg: {2 32 0 0 254 3 0 1 0}
rtmsg: RtMsg |
  Family:   AF_INET (2)
  Dst_len:  32
  Src_len:  0
  Tos:      0
  Table:    254
  Protocol: RTPROT_BOOT (3)
  Scope:    RT_SCOPE_UNIVERSE (0)
  Type:     RTN_UNICAST (1)
  Flags:    0
RtAttr | Len:8, Type:RTA_TABLE, Value:[254 0 0 0]
RtAttr | Len:8, Type:RTA_DST, Value:[10 11 11 99]
RtAttr | Len:8, Type:RTA_GATEWAY, Value:[172 20 104 1]
RtAttr | Len:8, Type:RTA_OIF, Value:[5 0 0 0]
```

## IPv4 MULTIPATH ROUTE (no next hop object)


```
> $ ip route add 10.11.11.11/32 \
nexthop via 172.20.105.174 dev eno1 \
nexthop via 172.20.105.175 dev eno1

nsdevnotes/examples/gonlsub$ go run gonlsub.go
Starting gonlsub.go
-----------------------------------
NlMsghdr | Len:80, Type:RTM_NEWROUTE, Flags:600, Seq:1669864316, Pid:224911
rtmsg: {2 32 0 0 254 3 0 1 0}
rtmsg: RtMsg |
  Family:   AF_INET (2)
  Dst_len:  32
  Src_len:  0
  Tos:      0
  Table:    254
  Protocol: RTPROT_BOOT (3)
  Scope:    RT_SCOPE_UNIVERSE (0)
  Type:     RTN_UNICAST (1)
  Flags:    0
RtAttr | Len:8, Type:RTA_TABLE, Value:254
RtAttr | Len:8, Type:RTA_DST, IPv4:10.11.11.11
RtAttr | Len:36, Type:RTA_MULTIPATH
  | rtnexthop: Len:16, Flags:0, Hops:0, Ifindex:5
  | RTA: Len:8, Type:RTA_GATEWAY, IPv4:172.20.105.174
  | rtnexthop: Len:16, Flags:0, Hops:0, Ifindex:5
  | RTA: Len:8, Type:RTA_GATEWAY, IPv4:172.20.105.175

# strace ip route add 10.11.11.11/32 nexthop via 172.20.105.174 dev eno1 nexthop via 172.20.105.175 dev eno1

sendmsg(3, {msg_name={sa_family=AF_NETLINK, nl_pid=0, nl_groups=00000000}, msg_namelen=12, msg_iov=[{iov_base=[{nlmsg_len=72, nlmsg_type=RTM_NEWROUTE, nlmsg_flags=NLM_F_REQUEST|NLM_F_ACK|NLM_F_EXCL|NLM_F_CREATE, nlmsg_seq=1669864316, nlmsg_pid=0}, {rtm_family=AF_INET, rtm_dst_len=32, rtm_src_len=0, rtm_tos=0, rtm_table=RT_TABLE_MAIN, rtm_protocol=RTPROT_BOOT, rtm_scope=RT_SCOPE_UNIVERSE, rtm_type=RTN_UNICAST, rtm_flags=0}, [[{nla_len=8, nla_type=RTA_DST}, inet_addr("10.11.11.11")], [{nla_len=36, nla_type=RTA_MULTIPATH}, [[{rtnh_len=16, rtnh_flags=0, rtnh_hops=0, rtnh_ifindex=if_nametoindex("eno1")}, [{nla_len=8, nla_type=RTA_GATEWAY}, inet_addr("172.20.105.174")]], [{rtnh_len=16, rtnh_flags=0, rtnh_hops=0, rtnh_ifindex=if_nametoindex("eno1")}, [{nla_len=8, nla_type=RTA_GATEWAY}, inet_addr("172.20.105.175")]]]]]], iov_len=72}], msg_iovlen=1, msg_controllen=0, msg_flags=0}, 0) = 72

Netlink Message Type: RTM_NEWROUTE
RT Message:
    rtm_family=AF_INET
    rtm_dst_len=32
    rtm_src_len=0
    rtm_tos=0
    rtm_table=RT_TABLE_MAIN
    rtm_protocol=RTPROT_BOOT
    rtm_scope=RT_SCOPE_UNIVERSE
    rtm_type=RTN_UNICAST
    rtm_flags=0
Netlink Attribute:
    {nla_len=8, nla_type=RTA_DST}, inet_addr("10.11.11.11")
    {nla_len=36, nla_type=RTA_MULTIPATH}
        {rtnh_len=16, rtnh_flags=0, rtnh_hops=0, rtnh_ifindex=if_nametoindex("eno1")}
            {nla_len=8, nla_type=RTA_GATEWAY}, inet_addr("172.20.105.174")]
        {rtnh_len=16, rtnh_flags=0, rtnh_hops=0, rtnh_ifindex=if_nametoindex("eno1")}
            {nla_len=8, nla_type=RTA_GATEWAY}, inet_addr("172.20.105.175")
```


## Prefix Len of RTA_DST (rtm_dst_len)

prefix len will be set in rtmsg: `rtm_dst_len`.

For example, to set route to `10.11.12.0/24` set 24 to `rtm_dst_len`.

```
> ip route add 10.11.12.0/24 nhid 11

-----------------------------------
NlMsghdr | Len:68, Type:RTM_NEWROUTE, Flags:600, Seq:1669596915, Pid:154076
rtmsg: {2 24 0 0 254 3 0 1 0}
rtmsg: RtMsg |
  Family:   AF_INET (2)
  Dst_len:  24
  Src_len:  0
  Tos:      0
  Table:    254
  Protocol: RTPROT_BOOT (3)
  Scope:    RT_SCOPE_UNIVERSE (0)
  Type:     RTN_UNICAST (1)
  Flags:    0
RtAttr | Len:8, Type:RTA_TABLE, Value:[254 0 0 0]
RtAttr | Len:8, Type:RTA_DST, Value:[10 11 12 0]
RtAttr | Len:8, Type:RTA_NH_ID, Value:[11 0 0 0]
RtAttr | Len:8, Type:RTA_GATEWAY, Value:[172 20 105 173]
RtAttr | Len:8, Type:RTA_OIF, Value:[5 0 0 0]
```

## tshark decode output of RTA_MULTIPATH

Looks like it doesn't support decoding data part of `RTA_MULTIPATH`.


```
$ tshark -v
TShark (Wireshark) 3.2.3 (Git v3.2.3 packaged as 3.2.3-1)


Frame 102: 104 bytes on wire (832 bits), 104 bytes captured (832 bits) on interface nlmon0, id 0
    Interface id: 0 (nlmon0)
        Interface name: nlmon0
    Encapsulation type: Linux Netlink (158)
    Arrival Time: Nov 28, 2022 02:02:32.821352526 UTC
    [Time shift for this packet: 0.000000000 seconds]
    Epoch Time: 1669600952.821352526 seconds
    [Time delta from previous captured frame: 0.000019221 seconds]
    [Time delta from previous displayed frame: 0.000019221 seconds]
    [Time since reference or first frame: 46.296049573 seconds]
    Frame Number: 102
    Frame Length: 104 bytes (832 bits)
    Capture Length: 104 bytes (832 bits)
    [Frame is marked: False]
    [Frame is ignored: False]
    [Protocols in frame: netlink:netlink-route]
Linux netlink (cooked header)
    Link-layer address type: Netlink (824)
    Family: Route (0x0000)
Linux rtnetlink (route netlink) protocol
    Netlink message header (type: Add network route)
        Length: 88
        Message type: Add network route (24)
        Flags: 0x0600
            .... .... .... ...0 = Request: 0
            .... .... .... ..0. = Multipart message: 0
            .... .... .... .0.. = Ack: 0
            .... .... .... 0... = Echo: 0
            .... .... ...0 .... = Dump inconsistent: 0
            .... .... ..0. .... = Dump filtered: 0
        Sequence: 1669600953
        Port ID: 155117
    Address family: AF_INET (2)
    Length of destination: 24
    Length of source: 0
    TOS filter: 0x00
    Routing table ID: 254
    Routing protocol: boot (0x03)
    Route origin: global route (0x00)
    Route type: Gateway or direct route (0x01)
    Route flags: 0x00000000
    Attribute: RTA_TABLE
        Len: 8
        Type: 0x000f, RTA_TABLE (15)
            0... .... .... .... = Nested: 0
            .0.. .... .... .... = Network byte order: 0
            Attribute type: RTA_TABLE (15)
        Data: fe000000
    Attribute: Route destination address
        Len: 8
        Type: 0x0001, Route destination address (1)
            0... .... .... .... = Nested: 0
            .0.. .... .... .... = Network byte order: 0
            Attribute type: Route destination address (1)
        Data: 0a0b0c00
    Attribute
        Len: 8
        Type: 0x001e
            0... .... .... .... = Nested: 0
            .0.. .... .... .... = Network byte order: 0
            Attribute type: Unknown (30)
        Data: 03000000
    Attribute: RTA_MULTIPATH
        Len: 36
        Type: 0x0009, RTA_MULTIPATH (9)
            0... .... .... .... = Nested: 0
            .0.. .... .... .... = Network byte order: 0
            Attribute type: RTA_MULTIPATH (9)
        Data: 100000000500000008000500ac1469ac1000000005000000…
```