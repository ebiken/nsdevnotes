# pytools: SendPacket

- read YAML file with packet information
- use scapy to send packet based on the information

Table of Contents
- [How to Use](#how-to-use)
- [prerequisits](#prerequisits)
- [YAML format example](#yaml-format-example)
- [Reference](#reference)
- [Example tshark output](#example-tshark-output)

## How to Use

```
$ sudo ./sendpacket.py --help
Usage: sendpacket.py [OPTIONS] FILE_YAML

Options:
  -c, --count INTEGER    number of packets to send
  -i, --send_iface TEXT  interface name to send packet
  -s, --show             show packet details
  --debug                show debug messages
  --help                 Show this message and exit.
```
> Note: You can also specify send_iface in yaml file.

Examples

```
> send 3 packets to interface ens1f0

sudo ./sendpacket.py pkt_srv6_ipv4_01.yaml -i ens1f0 -c 3

> send 1 packet. to interface ens1f0 or lo
> show packet details (-s) 
> show debug messages like parsed YAML (--debug)

sudo ./sendpacket.py pkt_srv6_ipv4_01.yaml -i ens1f0 -s --debug
sudo ./sendpacket.py pkt_srv6_ipv6_01.yaml -i lo -s --debug
```

## prerequisits

```
pip3 install pyyaml
```

## YAML format example

- Field names are ones listed with `ls(<scapy_header>)` like `ls(TCP)`.
- Header names (e.g. `ether`, `ipv6`) are lower case of the scapy header name with some exceptions for long ones. (e.g. `srh` for `IPv6ExtHdrSegmentRouting`)
- Copy [pkt_template.yaml](pkt_template.yaml) and remove unnessesary headers/fields.

```
packet:
  send_iface: "lo"
  ether:
    src: "02:03:04:05:06:01"
    dst: "02:03:04:05:06:02"
  ipv6:
    src: "2001:db8::100"
    dst: "2001:db8::1"
  srh: #IPv6ExtHdrSegmentRouting
    segleft: 1
    addresses:
      - "2001:db8:4::100"
      - "2001:db8::1"
  inner_packet:
    ip:
      src: "10.0.0.100"
      dst: "10.10.0.100"
    udp:
      sport: 1234
      dport: 4321
```

## Reference

- SRH: https://datatracker.ietf.org/doc/rfc8754/

## Example tshark output

```
$ sudo ./sendpacket.py pkt_srv6_ipv6_01.yaml -i lo -s --debug

# tshark -i lo -f "not tcp" -O ipv6
Running as user "root" and group "root". This could be dangerous.
Capturing on 'Loopback: lo'
Frame 1: 142 bytes on wire (1136 bits), 142 bytes captured (1136 bits) on interface lo, id 0
Ethernet II, Src: MS-NLB-PhysServer-03_04:05:06:01 (02:03:04:05:06:01), Dst: MS-NLB-PhysServer-03_04:05:06:02 (02:03:04:05:06:02)
Internet Protocol Version 6, Src: 2001:db8::100, Dst: 2001:db8::1
    0110 .... = Version: 6
    .... 0000 0000 .... .... .... .... .... = Traffic Class: 0x00 (DSCP: CS0, ECN: Not-ECT)
        .... 0000 00.. .... .... .... .... .... = Differentiated Services Codepoint: Default (0)
        .... .... ..00 .... .... .... .... .... = Explicit Congestion Notification: Not ECN-Capable Transport (0)
    .... .... .... 0000 0000 0000 0000 0000 = Flow Label: 0x00000
    Payload Length: 88
    Next Header: Routing Header for IPv6 (43)
    Hop Limit: 64
    Source: 2001:db8::100
    Destination: 2001:db8::1
    Routing Header for IPv6 (Segment Routing)
        Next Header: IPv6 (41)
        Length: 4
        [Length: 40 bytes]
        Type: Segment Routing (4)
        Segments Left: 1
        First segment: 1
        Flags: 0x00
            0... .... = Unused: 0x0
            .0.. .... = Protected: False
            ..0. .... = OAM: False
            ...0 .... = Alert: Not Present
            .... 0... = HMAC: Not Present
            .... .000 = Unused: 0x0
            [Expert Info (Note/Undecoded): Dissection for SRH TLVs not yet implemented]
                [Dissection for SRH TLVs not yet implemented]
                [Severity level: Note]
                [Group: Undecoded]
        Reserved: 0000
        Address[0]: 2001:db8:4::100 [next segment]
        Address[1]: 2001:db8::1
        [Segments in Traversal Order]
            Address[1]: 2001:db8::1
            Address[0]: 2001:db8:4::100 [next segment]
Internet Protocol Version 6, Src: 2001:db8:aa::100, Dst: 2001:db8:aa::1
    0110 .... = Version: 6
    .... 0000 0000 .... .... .... .... .... = Traffic Class: 0x00 (DSCP: CS0, ECN: Not-ECT)
        .... 0000 00.. .... .... .... .... .... = Differentiated Services Codepoint: Default (0)
        .... .... ..00 .... .... .... .... .... = Explicit Congestion Notification: Not ECN-Capable Transport (0)
    .... .... .... 0000 0000 0000 0000 0000 = Flow Label: 0x00000
    Payload Length: 8
    Next Header: UDP (17)
    Hop Limit: 64
    Source: 2001:db8:aa::100
    Destination: 2001:db8:aa::1
User Datagram Protocol, Src Port: 1234, Dst Port: 4321
```