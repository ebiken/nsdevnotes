# Running SONiC on containerlab

SONiC を [containerlab](https://containerlab.dev/) で動かす方法のメモ。  
複数のインスタンスを上げたい時に利用可能だが、sonic-vs は機能的な制約（sFlow動作しない等）があるため注意。  
詳細は本家の解説ページ https://containerlab.dev/manual/kinds/sonic-vs/ を参照。

Table of Contents
- [Quic Start using demo01.clab.yml](#quic-start-using-demo01clabyml)
  - [Install containerlab on Ubuntu 20.04](#install-containerlab-on-ubuntu-2004)
  - [demo02 : run two linux containers + SONiC](#demo02--run-two-linux-containers--sonic)
  - [demo01 : run two linux containers](#demo01--run-two-linux-containers)

## Quic Start using demo01.clab.yml

- https://containerlab.dev/quickstart/
- Try VM based? https://netdevops.me/2021/containerlab-your-network-centric-labs-with-a-docker-ux/

### Install containerlab on Ubuntu 20.04

```
Run: bash -c "$(curl -sL https://get.containerlab.dev)"

### Log ---------------------------------------------------
ebiken@xxx:~/containerlab$ lsb_release -a
No LSB modules are available.
Distributor ID: Ubuntu
Description:    Ubuntu 18.04.6 LTS
Release:        18.04
Codename:       bionic

ebiken@xxx:~/containerlab$ bash -c "$(curl -sL https://get.containerlab.dev)"
Downloading https://github.com/srl-labs/containerlab/releases/download/v0.25.1/containerlab_0.25.1_linux_amd64.deb
Preparing to install containerlab 0.25.1 from package
[sudo] password for ebiken:
Selecting previously unselected package containerlab.
(Reading database ... 226241 files and directories currently installed.)
Preparing to unpack .../containerlab_0.25.1_linux_amd64.deb ...
Unpacking containerlab (0.25.1) ...
Setting up containerlab (0.25.1) ...

                           _                   _       _
                 _        (_)                 | |     | |
 ____ ___  ____ | |_  ____ _ ____   ____  ____| | ____| | _
/ ___) _ \|  _ \|  _)/ _  | |  _ \ / _  )/ ___) |/ _  | || \
( (__| |_|| | | | |_( ( | | | | | ( (/ /| |   | ( ( | | |_) )
\____)___/|_| |_|\___)_||_|_|_| |_|\____)_|   |_|\_||_|____/

    version: 0.25.1
     commit: 2f68fb3d
       date: 2022-03-22T09:49:14Z
     source: https://github.com/srl-labs/containerlab
 rel. notes: https://containerlab.dev/rn/0.25/#0251

ebiken@xxx:~/containerlab$ ls /etc/containerlab/*
/etc/containerlab/lab-examples:
br01    clos01  clos03  cvx02  sonic01  srl02  srlceos01  srlfrr01     templated02  vr02  vr04  vxlan01
cert01  clos02  cvx01   frr01  srl01    srl03  srlcrpd01  templated01  vr01         vr03  vr05

/etc/containerlab/templates:
base__srl.tmpl  base__vr-sros.tmpl  graph  srl-ifaces__srl.tmpl
```

### demo02 : run two linux containers + SONiC

- https://containerlab.dev/manual/kinds/sonic-vs/

```
sudo containerlab deploy --topo demo02.clab.yml
sudo containerlab destroy --topo demo02.clab.yml

> 3. Setup IP in the virtual switch docker
$ docker exec -it clab-demo02-sonic bash
root@2e9b5c2dc2a2:/# 
config interface ip add Ethernet0 10.0.0.0/31
config interface ip add Ethernet4 10.0.0.2/31
config interface startup Ethernet0
config interface startup Ethernet4

### 
#    "INTERFACE": {
#        "Ethernet0": {},
#        "Ethernet0|10.0.0.0/31": {},
#        "Ethernet4": {},
#        "Ethernet4|10.0.0.2/31": {}
#    },
#    "PORT": {
#        "Ethernet0": {
#            "admin_status": "up",
#            "alias": "fortyGigE0/0",
#            "index": "0",
#            "lanes": "25,26,27,28",
#            "speed": "100000"
#        },
#        "Ethernet4": {
#            "admin_status": "up",
#            "alias": "fortyGigE0/4",
#            "index": "1",
#            "lanes": "29,30,31,32",
#            "speed": "100000"
#        },

root@f495dd5429c2:/# show interface status

> Setup n1,n2 hosts

sudo ip netns exec clab-demo02-n1 ifconfig eth1 10.0.0.1/31
sudo ip netns exec clab-demo02-n1 ip route del default via 172.20.20.1
sudo ip netns exec clab-demo02-n1 ip route add default via 10.0.0.0
sudo ip netns exec clab-demo02-n2 ifconfig eth1 10.0.0.3/31
sudo ip netns exec clab-demo02-n2 ip route del default via 172.20.20.1
sudo ip netns exec clab-demo02-n2 ip route add default via 10.0.0.2

docker exec -it clab-demo01-n1 sh
sudo containerlab destroy --topo demo01.clab.yml

### Log ---------------------------------------------------
ebiken@xxx:~/containerlab$ sudo containerlab deploy --topo demo02.clab.yml
INFO[0000] Containerlab v0.25.1 started
INFO[0000] Parsing & checking topology file: demo02.clab.yml
INFO[0000] Creating lab directory: /home/ebiken/containerlab/clab-demo02
INFO[0000] Creating docker network: Name="clab", IPv4Subnet="172.20.20.0/24", IPv6Subnet="2001:172:20:20::/64", MTU="1500"
INFO[0000] Creating container: "n2"
INFO[0000] Creating container: "sonic"
INFO[0000] Creating container: "n1"
INFO[0002] Creating virtual wire: n1:eth1 <--> sonic:eth1
INFO[0003] Creating virtual wire: n2:eth1 <--> sonic:eth4
INFO[0003] Adding containerlab host entries to /etc/hosts file
+---+-------------------+--------------+-----------------+----------+---------+----------------+----------------------+
| # |       Name        | Container ID |      Image      |   Kind   |  State  |  IPv4 Address  |     IPv6 Address     |
+---+-------------------+--------------+-----------------+----------+---------+----------------+----------------------+
| 1 | clab-demo02-n1    | d94d47e79cc3 | alpine:latest   | linux    | running | 172.20.20.2/24 | 2001:172:20:20::2/64 |
| 2 | clab-demo02-n2    | 87d5b4017cca | alpine:latest   | linux    | running | 172.20.20.3/24 | 2001:172:20:20::3/64 |
| 3 | clab-demo02-sonic | db7fd0156303 | docker-sonic-vs | sonic-vs | running | 172.20.20.4/24 | 2001:172:20:20::4/64 |
+---+-------------------+--------------+-----------------+----------+---------+----------------+----------------------+

```

### demo01 : run two linux containers

https://containerlab.dev/manual/kinds/linux/

```
> create and run lab, connect to container, destroy lab.

sudo containerlab deploy --topo demo01.clab.yml
docker exec -it clab-demo01-n1 sh
sudo containerlab destroy --topo demo01.clab.yml

### Log ---------------------------------------------------

ebiken@xxx:~/containerlab$ cat demo01.clab.yml
# a simple topo of two alpine containers connected with each other
name: demo01

topology:
  nodes:
    n1:
      kind: linux
      image: alpine:latest
    n2:
      kind: linux
      image: alpine:latest
  links:
    - endpoints: ["n1:eth1","n2:eth1"]

> Run: sudo containerlab deploy --topo demo01.clab.yml
ebiken@xxx:~/containerlab$ sudo containerlab deploy --topo demo01.clab.yml
INFO[0000] Containerlab v0.25.1 started
INFO[0000] Parsing & checking topology file: demo01.clab.yml
INFO[0000] Creating lab directory: /home/ebiken/containerlab/clab-demo01
INFO[0000] Creating docker network: Name="clab", IPv4Subnet="172.20.20.0/24", IPv6Subnet="2001:172:20:20::/64", MTU="1500"
INFO[0000] Creating container: "n1"
INFO[0000] Creating container: "n2"
INFO[0002] Creating virtual wire: n1:eth1 <--> n2:eth1
INFO[0002] Adding containerlab host entries to /etc/hosts file
+---+----------------+--------------+---------------+-------+---------+----------------+----------------------+
| # |      Name      | Container ID |     Image     | Kind  |  State  |  IPv4 Address  |     IPv6 Address     |
+---+----------------+--------------+---------------+-------+---------+----------------+----------------------+
| 1 | clab-demo01-n1 | d5d33136b62e | alpine:latest | linux | running | 172.20.20.3/24 | 2001:172:20:20::3/64 |
| 2 | clab-demo01-n2 | 695d61e6e36d | alpine:latest | linux | running | 172.20.20.2/24 | 2001:172:20:20::2/64 |
+---+----------------+--------------+---------------+-------+---------+----------------+----------------------+

> Enter container
ebiken@xxx:~/containerlab$ docker exec -it clab-demo01-n1 sh
/ # ip a
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
    inet6 ::1/128 scope host
       valid_lft forever preferred_lft forever
46: eth0@if47: <BROADCAST,MULTICAST,UP,LOWER_UP,M-DOWN> mtu 1500 qdisc noqueue state UP
    link/ether 02:42:ac:14:14:03 brd ff:ff:ff:ff:ff:ff
    inet 172.20.20.3/24 brd 172.20.20.255 scope global eth0
       valid_lft forever preferred_lft forever
    inet6 2001:172:20:20::3/64 scope global flags 02
       valid_lft forever preferred_lft forever
    inet6 fe80::42:acff:fe14:1403/64 scope link
       valid_lft forever preferred_lft forever
49: eth1@if48: <BROADCAST,MULTICAST,UP,LOWER_UP,M-DOWN> mtu 9500 qdisc noqueue state UP
    link/ether aa:c1:ab:9e:a5:0d brd ff:ff:ff:ff:ff:ff
    inet6 fe80::a8c1:abff:fe9e:a50d/64 scope link
       valid_lft forever preferred_lft forever

> Destroy lab
ebiken@xxx:~/containerlab$ sudo containerlab destroy --topo demo01.clab.yml
INFO[0000] Parsing & checking topology file: demo01.clab.yml
INFO[0000] Destroying lab: demo01
INFO[0000] Removed container: clab-demo01-n1
INFO[0000] Removed container: clab-demo01-n2
INFO[0000] Removing containerlab host entries from /etc/hosts file
```
