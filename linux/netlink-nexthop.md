# netlink/rtnetlink - Next Hop Object & Next Hop Group

> 特に言及がない場合は Linux v6.0 をベースに解説しています

本ページでは ip route を設定する際の netlink (rtnetlink) の動作を、 Linux v5.3 で導入された Next Hop Object を利用しない従来の方法と、利用した場合を比較を中心に解説しています。

- netlink に関する基本的な説明は [Linux Netlink](./netlink.md) を参照
- Linux における IP Routing （Fib や nexthop を含む）に関しては [Linux IP Routing](./iprouting.md) を参照

目次

- [strace と RTM\_NEWNEXTHOP](#strace-と-rtm_newnexthop)
- [ip route 追加時の netlink/rtnetlink の具体例](#ip-route-追加時の-netlinkrtnetlink-の具体例)
  - [nexthop 利用無し（従来）](#nexthop-利用無し従来)
  - [nexthop 利用無し（従来） Multipath](#nexthop-利用無し従来-multipath)
  - [nexthop を利用](#nexthop-を利用)
  - [nexthop group を利用 Multipath](#nexthop-group-を利用-multipath)
- [route add ipv6](#route-add-ipv6)
- [reference](#reference)

## strace と RTM_NEWNEXTHOP

strace コマンドと Kernel 間でやり取りされる netlink message をモニタ可能な便利なツールです。

以降の解説では、ip コマンドの前に `strace` を付けて実行しています。

Next Hop Object に関するメッセージである `RTM_NEWNEXTHOP` には [strace v5.15 (2021-10-14) から対応](https://fossies.org/linux/strace/ChangeLog) していますので、もしそれ以前のバージョンの場合は以下のように v5.15 以上にアップデートが必要です。

yum/apt コマンドによるアップデートができればベストですが、もし yum/apt で v5.15 以降にアップデートされない場合は以下のように Source Code からのビルド＆インストールが必要になります。

Ubuntu 20.04.4 で strace v6.0 をビルド＆インストールした際の手順は以下の通りです。

```
> https://github.com/strace/strace/releases/tag/v6.0
> download strace-6.0.tar.xz

$ tar xf strace-6.0.tar.xz
$ cd strace-6.0
$ ./configure --disable-mpers
$ make
$ sudo make install

$ which strace
/usr/local/bin/strace

$ strace --version
strace -- version 6.0
Copyright (c) 1991-2022 The strace developers <https://strace.io>.
This is free software; see the source for copying conditions.  There is NO
warranty; not even for MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.

Optional features enabled: stack-trace=libunwind no-m32-mpers no-mx32-mpers
```

## ip route 追加時の netlink/rtnetlink の具体例

ip route を設定する際の netlink/rtnetlink の動作を確認しましょう。

具体例として、以下３パターンを比較します。
"Multipath" は２つ（以上）の nexthop が存在する route を意味します。

- nexthop 利用無し（従来）
  - `ip route add 10.11.11.99/32 via 172.20.104.1 dev eno1`
- nexthop 利用無し Multipath（従来）
  - `ip route add 10.11.11.11/32 nexthop via 172.20.105.174 dev eno1 nexthop via 172.20.105.175 dev eno1`
- nexthop を利用
  - `ip nexthop add id 11 via 172.20.105.173 dev eno1`
  - `ip route add 10.11.12.13/32 nhid 11`
- nexthop group を利用 Multipath
  - `ip nexthop add id 1 via 172.20.105.172 dev eno1`
  - `ip nexthop add id 2 via 172.20.105.173 dev eno1`
  - `ip nexthop add id 3 group 1/2`
  - `ip route add 10.11.12.13/32 nhid 3`

それぞれの解説では設定のための rtnetlink message である sendmsg だけを抜粋しています。
デバイス名の解決など、その前後でもメッセージがやりとりされる場合がありますので、詳細は以下 strace のログを参照してください。（ご自身の環境で strace コマンドを入力してみる事をお勧めします）

- strace logs
  - [nexthop 利用無し（従来）](logs/strace-ip-route-add-no-nexthop.log)
  - [nexthop 利用無し（従来）Multipath](logs/strace-ip-route-add-no-nexthop-multipath.log)
  - [nexthop を利用](logs/strace-ip-route-add-nexthop.log)
  - [nexthop group を利用 Multipath](logs/strace-ip-route-add-nexthop-group.log)

共通の解説

- `RTA_OIF` の値である `if_nametoindex(<name>)` は、デバイス名から dev index を求める関数です。
  - 実際の netlink message ではデバイス名（e.g. `eno1`）ではなく、数字（ID）が送信される事に留意してください。
- `rtmsg rtm_table` や `RTA_TABLE` で利用される route table の一覧は `cat /etc/iproute2/rt_tables` から取得可能
  - `RT_TABLE_MAIN` (Table ID: 254)  は Table ID を指定せず route を追加した場合追加されるテーブル
- 宛先アドレス（PREFIX）について
  - route entry の宛先アドレス情報は `RTA_DST` に保持されますが、PREFIX Length は `rtm_dst_len` に保持されます。
  - 宛先 prefix/length に関する情報が RT Message と Attribute の異なる場所に格納されるので注意が必要です。

### nexthop 利用無し（従来）

nexthop object を用いない従来の方法では、 `RTM_NEWROUTE` に nexthop に関する情報が含まれます。

```
# ip route add 10.11.11.99/32 via 172.20.104.1 dev eno1

sendmsg(3, {msg_name={sa_family=AF_NETLINK, nl_pid=0, nl_groups=00000000}, msg_namelen=12, msg_iov=[{iov_base={{len=52, type=RTM_NEWROUTE, flags=NLM_F_REQUEST|NLM_F_ACK|NLM_F_EXCL|NLM_F_CREATE, seq=1669690078, pid=0}, {rtm_family=AF_INET, rtm_dst_len=32, rtm_src_len=0, rtm_tos=0, rtm_table=RT_TABLE_MAIN, rtm_protocol=RTPROT_BOOT, rtm_scope=RT_SCOPE_UNIVERSE, rtm_type=RTN_UNICAST, rtm_flags=0}, [{{nla_len=8, nla_type=RTA_DST}, inet_addr("10.11.11.99")}, {{nla_len=8, nla_type=RTA_GATEWAY}, inet_addr("172.20.104.1")}, {{nla_len=8, nla_type=RTA_OIF}, if_nametoindex("eno1")}]}, iov_len=52}], msg_iovlen=1, msg_controllen=0, msg_flags=0}, 0) = 52

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
    {nla_len=8, nla_type=RTA_DST}, inet_addr("10.11.11.99")
    {nla_len=8, nla_type=RTA_GATEWAY}, inet_addr("172.20.104.1")
    {nla_len=8, nla_type=RTA_OIF}, if_nametoindex("eno1")
```

### nexthop 利用無し（従来） Multipath

nexthop object を用いない従来の方法でも、複数の nexthop 設定（Multipath）は可能です。
`RTA_MULTIPATH` の値として、 `rtnexthop` 構造体を nexthop の数だけ利用し、 `rtnexthop` の中の `rtnh_ifindex` に Output Interface ID（`RTA_OIF` 相当）を、`rtnexthop` の値に `RTA_GATEWAY` として gateway アドレスをセットします。


```c
// include/uapi/linux/rtnetlink.h
struct rtnexthop {
    unsigned short      rtnh_len;
    unsigned char       rtnh_flags;
    unsigned char       rtnh_hops;
    int         rtnh_ifindex;
};
```

```
# ip route add 10.11.11.11/32 nexthop via 172.20.105.174 dev eno1 nexthop via 172.20.105.175 dev eno1

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
            {nla_len=8, nla_type=RTA_GATEWAY}, inet_addr("172.20.105.174")
        {rtnh_len=16, rtnh_flags=0, rtnh_hops=0, rtnh_ifindex=if_nametoindex("eno1")}
            {nla_len=8, nla_type=RTA_GATEWAY}, inet_addr("172.20.105.175")
```


### nexthop を利用

Next Hop Object を利用する場合は、まず `RTM_NEWNEXTHOP` メッセージを送信して nexthop を作成し、そのIDを（`RTA_GATEWAY` や `RTA_OIF` の代わりに） `RTM_NEWROUTE` の Attribute である `RTA_NH_ID` に指定して送信します。


```
> ip nexthop add id 11 via 172.20.105.173 dev eno1

sendmsg(3, {msg_name={sa_family=AF_NETLINK, nl_pid=0, nl_groups=00000000}, msg_namelen=12, msg_iov=[{iov_base=[{nlmsg_len=48, nlmsg_type=RTM_NEWNEXTHOP, nlmsg_flags=NLM_F_REQUEST|NLM_F_ACK|NLM_F_EXCL|NLM_F_CREATE, nlmsg_seq=1669695871, nlmsg_pid=0}, {nh_family=AF_INET, nh_scope=RT_SCOPE_UNIVERSE, nh_protocol=RTPROT_UNSPEC, nh_flags=0}, [[{nla_len=8, nla_type=NHA_ID}, 11], [{nla_len=8, nla_type=NHA_GATEWAY}, inet_addr("172.20.105.173")], [{nla_len=8, nla_type=NHA_OIF}, if_nametoindex("eno1")]]], iov_len=48}], msg_iovlen=1, msg_controllen=0, msg_flags=0}, 0) = 48

Netlink Message Type: RTM_NEWNEXTHOP
Next Hop Message:
    nh_family=AF_INET,
    nh_scope=RT_SCOPE_UNIVERSE,
    nh_protocol=RTPROT_UNSPEC,
    nh_flags=0
Netlink Attribute:
    {nla_len=8, nla_type=NHA_ID}, 11
    {nla_len=8, nla_type=NHA_GATEWAY}, inet_addr("172.20.105.173")
    {nla_len=8, nla_type=NHA_OIF}, if_nametoindex("eno1")
```

```
> ip route add 10.11.12.13/32 nhid 11

sendmsg(3, {msg_name={sa_family=AF_NETLINK, nl_pid=0, nl_groups=00000000}, msg_namelen=12, msg_iov=[{iov_base=[{nlmsg_len=44, nlmsg_type=RTM_NEWROUTE, nlmsg_flags=NLM_F_REQUEST|NLM_F_ACK|NLM_F_EXCL|NLM_F_CREATE, nlmsg_seq=1669710575, nlmsg_pid=0}, {rtm_family=AF_INET, rtm_dst_len=32, rtm_src_len=0, rtm_tos=0, rtm_table=RT_TABLE_MAIN, rtm_protocol=RTPROT_BOOT, rtm_scope=RT_SCOPE_UNIVERSE, rtm_type=RTN_UNICAST, rtm_flags=0}, [[{nla_len=8, nla_type=RTA_DST}, inet_addr("10.11.12.13")], [{nla_len=8, nla_type=RTA_NH_ID}, "\x0b\x00\x00\x00"]]], iov_len=44}], msg_iovlen=1, msg_controllen=0, msg_flags=0}, 0) = 44

Netlink Message Type: RTM_NEWROUTE
Next Hop Message:
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
    {nla_len=8, nla_type=RTA_DST}, inet_addr("10.11.12.13")
    {nla_len=8, nla_type=RTA_NH_ID}, "\x0b\x00\x00\x00"
```

### nexthop group を利用 Multipath


Next Hop Object を利用して Multipath を設置する場合は、以下３ステップを辿ります。

1. `RTM_NEWNEXTHOP` メッセージを送信し nexthop を作成（２個以上）
2. `RTM_NEWNEXTHOP` メッセージを送信し `NHA_GROUP` に "1." で作成した nexthop の ID を指定し nexthop group を作成
3. nexthop group の ID を `RTA_NH_ID` に指定して送信


```
> ip nexthop add id 1 via 172.20.105.172 dev eno1
> ip nexthop add id 2 via 172.20.105.173 dev eno1

sendmsg(3, {msg_name={sa_family=AF_NETLINK, nl_pid=0, nl_groups=00000000}, msg_namelen=12, msg_iov=[{iov_base=[{nlmsg_len=48, nlmsg_type=RTM_NEWNEXTHOP, nlmsg_flags=NLM_F_REQUEST|NLM_F_ACK|NLM_F_EXCL|NLM_F_CREATE, nlmsg_seq=1669711458, nlmsg_pid=0}, {nh_family=AF_INET, nh_scope=RT_SCOPE_UNIVERSE, nh_protocol=RTPROT_UNSPEC, nh_flags=0}, [[{nla_len=8, nla_type=NHA_ID}, 1], [{nla_len=8, nla_type=NHA_GATEWAY}, inet_addr("172.20.105.172")], [{nla_len=8, nla_type=NHA_OIF}, if_nametoindex("eno1")]]], iov_len=48}], msg_iovlen=1, msg_controllen=0, msg_flags=0}, 0) = 48

sendmsg(3, {msg_name={sa_family=AF_NETLINK, nl_pid=0, nl_groups=00000000}, msg_namelen=12, msg_iov=[{iov_base=[{nlmsg_len=48, nlmsg_type=RTM_NEWNEXTHOP, nlmsg_flags=NLM_F_REQUEST|NLM_F_ACK|NLM_F_EXCL|NLM_F_CREATE, nlmsg_seq=1669711492, nlmsg_pid=0}, {nh_family=AF_INET, nh_scope=RT_SCOPE_UNIVERSE, nh_protocol=RTPROT_UNSPEC, nh_flags=0}, [[{nla_len=8, nla_type=NHA_ID}, 2], [{nla_len=8, nla_type=NHA_GATEWAY}, inet_addr("172.20.105.173")], [{nla_len=8, nla_type=NHA_OIF}, if_nametoindex("eno1")]]], iov_len=48}], msg_iovlen=1, msg_controllen=0, msg_flags=0}, 0) = 48

Netlink Message Type: RTM_NEWNEXTHOP
Next Hop Message:
    nh_family=AF_INET,
    nh_scope=RT_SCOPE_UNIVERSE,
    nh_protocol=RTPROT_UNSPEC,
    nh_flags=0
Netlink Attribute:
    {nla_len=8, nla_type=NHA_ID}, 1
    {nla_len=8, nla_type=NHA_GATEWAY}, inet_addr("172.20.105.172")
    {nla_len=8, nla_type=NHA_OIF}, if_nametoindex("eno1")

Netlink Message Type: RTM_NEWNEXTHOP
Next Hop Message:
    nh_family=AF_INET,
    nh_scope=RT_SCOPE_UNIVERSE,
    nh_protocol=RTPROT_UNSPEC,
    nh_flags=0
Netlink Attribute:
    {nla_len=8, nla_type=NHA_ID}, 2
    {nla_len=8, nla_type=NHA_GATEWAY}, inet_addr("172.20.105.173")
    {nla_len=8, nla_type=NHA_OIF}, if_nametoindex("eno1")

> ip nexthop add id 3 group 1/2

sendmsg(3, {msg_name={sa_family=AF_NETLINK, nl_pid=0, nl_groups=00000000}, msg_namelen=12, msg_iov=[{iov_base=[{nlmsg_len=52, nlmsg_type=RTM_NEWNEXTHOP, nlmsg_flags=NLM_F_REQUEST|NLM_F_ACK|NLM_F_EXCL|NLM_F_CREATE, nlmsg_seq=1669711532, nlmsg_pid=0}, {nh_family=AF_UNSPEC, nh_scope=RT_SCOPE_UNIVERSE, nh_protocol=RTPROT_UNSPEC, nh_flags=0}, [[{nla_len=8, nla_type=NHA_ID}, 3], [{nla_len=20, nla_type=NHA_GROUP}, [{id=1, weight=0}, {id=2, weight=0}]]]], iov_len=52}], msg_iovlen=1, msg_controllen=0, msg_flags=0}, 0) = 52

Netlink Message Type: RTM_NEWNEXTHOP
Next Hop Message:
    nh_family=AF_UNSPEC,
    nh_scope=RT_SCOPE_UNIVERSE,
    nh_protocol=RTPROT_UNSPEC,
    nh_flags=0
Netlink Attribute:
    {nla_len=8, nla_type=NHA_ID}, 3
    {nla_len=20, nla_type=NHA_GROUP}, [ {id=1, weight=0}, {id=2, weight=0} ]


> ip route add 10.11.12.13/32 nhid 3

sendmsg(3, {msg_name={sa_family=AF_NETLINK, nl_pid=0, nl_groups=00000000}, msg_namelen=12, msg_iov=[{iov_base=[{nlmsg_len=44, nlmsg_type=RTM_NEWROUTE, nlmsg_flags=NLM_F_REQUEST|NLM_F_ACK|NLM_F_EXCL|NLM_F_CREATE, nlmsg_seq=1669711569, nlmsg_pid=0}, {rtm_family=AF_INET, rtm_dst_len=32, rtm_src_len=0, rtm_tos=0, rtm_table=RT_TABLE_MAIN, rtm_protocol=RTPROT_BOOT, rtm_scope=RT_SCOPE_UNIVERSE, rtm_type=RTN_UNICAST, rtm_flags=0}, [[{nla_len=8, nla_type=RTA_DST}, inet_addr("10.11.12.13")], [{nla_len=8, nla_type=RTA_NH_ID}, "\x03\x00\x00\x00"]]], iov_len=44}], msg_iovlen=1, msg_controllen=0, msg_flags=0}, 0) = 44

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
    {nla_len=8, nla_type=RTA_DST}, inet_addr("10.11.12.13")
    {nla_len=8, nla_type=RTA_NH_ID}, "\x03\x00\x00\x00"
```

## route add ipv6

IPv6 を設定する際も同様となります。
なお、 `rtm_family=AF_INET6` の場合は `RTA_DST` の長さが `AF_INET` の場合と異なる事に注意してください。

```
> ip route add 2001:db8:ffff::/64 dev veth103

sendmsg(3, {msg_name={sa_family=AF_NETLINK, nl_pid=0, nl_groups=00000000}, msg_namelen=12, msg_iov=[{iov_base={{len=56, type=RTM_NEWROUTE, flags=NLM_F_REQUEST|NLM_F_ACK|NLM_F_EXCL|NLM_F_CREATE, seq=1669814439, pid=0}, {rtm_family=AF_INET6, rtm_dst_len=64, rtm_src_len=0, rtm_tos=0, rtm_table=RT_TABLE_MAIN, rtm_protocol=RTPROT_BOOT, rtm_scope=RT_SCOPE_UNIVERSE, rtm_type=RTN_UNICAST, rtm_flags=0}, [{{nla_len=20, nla_type=RTA_DST}, 2001:db8:ffff::}, {{nla_len=8, nla_type=RTA_OIF}, if_nametoindex("veth103")}]}, iov_len=56}], msg_iovlen=1, msg_controllen=0, msg_flags=0}, 0) = 56

Netlink Message Type: RTM_NEWROUTE
RT Message:
    rtm_family=AF_INET6
    rtm_dst_len=64
    rtm_src_len=0
    rtm_tos=0
    rtm_table=RT_TABLE_MAIN
    rtm_protocol=RTPROT_BOOT
    rtm_scope=RT_SCOPE_UNIVERSE
    rtm_type=RTN_UNICAST
    rtm_flags=0
Netlink Attribute:
    {nla_len=20, nla_type=RTA_DST}, 2001:db8:ffff::
    {nla_len=8, nla_type=RTA_OIF}, if_nametoindex("veth103")
```


## reference

- https://wiki.slank.dev/book/types.html
  - NHA_* (Next Hop Attribute) 含む netlink attribute が１ページにまとまってる
- 2022-11-23 [Netlinkと友達になろう](https://eniyo0.hatenablog.com/entry/2022/11/23/180135)
  - 日本語の平易な解説BLOG
- 
