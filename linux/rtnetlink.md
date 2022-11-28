# rtnetlink: NETLINK_ROUTE

This is WIP draft.

Refer [/examples/gonlsub/](/examples/gonlsub/) for RTNETLINK explanation & sample code.


## memo

- NETLINK TYPE for NETLINK_ROUTE is 0
    - https://sites.uclouvain.be/SystInfo/usr/include/linux/netlink.h
    - `#define NETLINK_ROUTE		0	/* Routing/device hook				*/`
- RTM_NEWROUTE: 24 (0x18)
    - tshark: `Message type: Add network route (24)`
- xxx

### RTM_NEWROUTE

dump using gonlsub.go

```
> msgs, from, err := nlSock.Receive()

ebiken@dcig171:~/sandbox/nsdevnotes/examples/gonlsub$ go run gonlsub.go
Starting gonlsub.go
msgs: [{{60 24 1536 1669519338 130693} [2 32 0 0 254 3 0 1 0 0 0 0 8 0 15 0 254 0 0 0 8 0 1 0 10 11 11 99 8 0 5 0 172 20 104 1 8 0 4 0 5 0 0 0]}]
from: &{16 0 0 64 {0 0 0 0}}
err:  <nil>
```

strace

```
>>> root@dcig171:/home/ebiken# strace ip route add 10.11.11.99/32 via 172.20.104.1 dev eno1
```

```json
sendmsg(3, {msg_name={sa_family=AF_NETLINK, nl_pid=0, nl_groups=00000000}, msg_namelen=12,
    msg_iov=[{iov_base={
        {len=52, type=RTM_NEWROUTE, flags=NLM_F_REQUEST|NLM_F_ACK|NLM_F_EXCL|NLM_F_CREATE, seq=1669520597, pid=0},
        {rtm_family=AF_INET, rtm_dst_len=32, rtm_src_len=0, rtm_tos=0, rtm_table=RT_TABLE_MAIN, rtm_protocol=RTPROT_BOOT, rtm_scope=RT_SCOPE_UNIVERSE, rtm_type=RTN_UNICAST, rtm_flags=0},
        [
            {{nla_len=8, nla_type=RTA_DST}, inet_addr("10.11.11.99")},
            {{nla_len=8, nla_type=RTA_GATEWAY}, inet_addr("172.20.104.1")},
            {{nla_len=8, nla_type=RTA_OIF}, if_nametoindex("eno1")}
        ]
        }, iov_len=52
    }],
    msg_iovlen=1, msg_controllen=0, msg_flags=0}, 0) = 52

recvmsg(3, {msg_name={sa_family=AF_NETLINK, nl_pid=0, nl_groups=00000000}, msg_namelen=12,
    msg_iov=[{iov_base=NULL, iov_len=0}],
    msg_iovlen=1, msg_controllen=0, msg_flags=MSG_TRUNC}, MSG_PEEK|MSG_TRUNC) = 36

recvmsg(3, {msg_name={sa_family=AF_NETLINK, nl_pid=0, nl_groups=00000000}, msg_namelen=12,
    msg_iov=[{iov_base={
        {len=36, type=NLMSG_ERROR, flags=NLM_F_CAPPED, seq=1669520597, pid=132894},
        {error=0, msg={len=52, type=RTM_NEWROUTE, flags=NLM_F_REQUEST|NLM_F_ACK|NLM_F_EXCL|NLM_F_CREATE, seq=1669520597, pid=0}}
        }, iov_len=32768
    }], msg_iovlen=1, msg_controllen=0, msg_flags=0}, 0) = 36
```

tshark example of RTM_NEWROUTE

```
> ebiken@dcig171:~$ sudo ip route add 10.11.11.99/32 via 172.20.104.1 dev eno1
> 10.11.11.99 => hex: 0a0b0b63

Frame 15: 68 bytes on wire (544 bits), 68 bytes captured (544 bits) on interface nlmon0, id 0
Linux netlink (cooked header)
    Link-layer address type: Netlink (824)
    Family: Route (0x0000)
Linux rtnetlink (route netlink) protocol
    Netlink message header (type: Add network route)
        Length: 52
        Message type: Add network route (24)
        Flags: 0x0605
            .... .... .... ...1 = Request: 1
            .... .... .... ..0. = Multipart message: 0
            .... .... .... .1.. = Ack: 1
            .... .... .... 0... = Echo: 0
            .... .... ...0 .... = Dump inconsistent: 0
            .... .... ..0. .... = Dump filtered: 0
            .... ...0 .... .... = Specify tree root: 0
            .... ..1. .... .... = Return all matching: 1
            .... .1.. .... .... = Atomic: 1
        Flags: 0x0605
            .... .... .... ...1 = Request: 1
            .... .... .... ..0. = Multipart message: 0
            .... .... .... .1.. = Ack: 1
            .... .... .... 0... = Echo: 0
            .... .... ...0 .... = Dump inconsistent: 0
            .... .... ..0. .... = Dump filtered: 0
            .... ...0 .... .... = Replace: 0
            .... ..1. .... .... = Excl: 1
            .... .1.. .... .... = Create: 1
            .... 0... .... .... = Append: 0
        Sequence: 1669518030
        Port ID: 0
    Address family: AF_INET (2)
    Length of destination: 32
    Length of source: 0
    TOS filter: 0x00
    Routing table ID: 254
    Routing protocol: boot (0x03)
    Route origin: global route (0x00)
    Route type: Gateway or direct route (0x01)
    Route flags: 0x00000000
    Attribute: Route destination address
        Len: 8
        Type: 0x0001, Route destination address (1)
            0... .... .... .... = Nested: 0
            .0.. .... .... .... = Network byte order: 0
            Attribute type: Route destination address (1)
        Data: 0a0b0b63
    Attribute: Gateway of the route
        Len: 8
        Type: 0x0005, Gateway of the route (5)
            0... .... .... .... = Nested: 0
            .0.. .... .... .... = Network byte order: 0
            Attribute type: Gateway of the route (5)
        Data: ac146801
    Attribute: Output interface index: 5
        Len: 8
        Type: 0x0004, Output interface index (4)
            0... .... .... .... = Nested: 0
            .0.. .... .... .... = Network byte order: 0
            Attribute type: Output interface index (4)
        Output interface index: 5
```

## Const definitions

https://sites.uclouvain.be/SystInfo/usr/include/linux/rtnetlink.h.html

```c
enum {
        RTM_BASE        = 16,
#define RTM_BASE        RTM_BASE

        RTM_NEWLINK        = 16,
#define RTM_NEWLINK        RTM_NEWLINK
        RTM_DELLINK,
#define RTM_DELLINK        RTM_DELLINK
        RTM_GETLINK,
#define RTM_GETLINK        RTM_GETLINK
        RTM_SETLINK,
#define RTM_SETLINK        RTM_SETLINK

        RTM_NEWADDR        = 20,
#define RTM_NEWADDR        RTM_NEWADDR
        RTM_DELADDR,
#define RTM_DELADDR        RTM_DELADDR
        RTM_GETADDR,
#define RTM_GETADDR        RTM_GETADDR

        RTM_NEWROUTE        = 24,
#define RTM_NEWROUTE        RTM_NEWROUTE
        RTM_DELROUTE,
#define RTM_DELROUTE        RTM_DELROUTE
        RTM_GETROUTE,
#define RTM_GETROUTE        RTM_GETROUTE
... snip ...
}
```