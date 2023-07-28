# Performance measurement

- [L2 bridge between two VM (Ubutnu Multipass)](#l2-bridge-between-two-vm-ubutnu-multipass)
  - [Result](#result)
  - [VM Config](#vm-config)
  - [Host Interface (bridge)](#host-interface-bridge)
- [WIP: L3 Routing among 3 VMs (Ubutnu Multipass)](#wip-l3-routing-among-3-vms-ubutnu-multipass)
- [L3 routing between netns](#l3-routing-between-netns)
- [Specs](#specs)
  - [Host 171](#host-171)


## L2 bridge between two VM (Ubutnu Multipass)

Create two VMs

```
multipass launch 22.04 --name vm1 --cpus 16 --disk 20G --memory 16G
multipass launch 22.04 --name vm2 --cpus 16 --disk 20G --memory 16G
```

Login and install iperf3

```
sudo apt update
sudo apt install -y iperf3
```

Run iperf

```
> vm1: server

iperf3 -s

> vm2: client

iperf3 -c 10.28.65.177 -p 5201
```

### Result

```
ubuntu@vm2:~$ iperf3 -c 10.28.65.177 -p 5201
Connecting to host 10.28.65.177, port 5201
[  5] local 10.28.65.38 port 42840 connected to 10.28.65.177 port 5201
[ ID] Interval           Transfer     Bitrate         Retr  Cwnd
[  5]   0.00-1.00   sec  1.48 GBytes  12.7 Gbits/sec    0   3.13 MBytes
[  5]   1.00-2.00   sec  1.61 GBytes  13.8 Gbits/sec    0   3.13 MBytes
[  5]   2.00-3.00   sec  1.71 GBytes  14.6 Gbits/sec    0   3.13 MBytes
[  5]   3.00-4.00   sec  1.57 GBytes  13.5 Gbits/sec    0   3.13 MBytes
[  5]   4.00-5.00   sec  1.46 GBytes  12.5 Gbits/sec    0   3.13 MBytes
[  5]   5.00-6.00   sec  1.50 GBytes  12.9 Gbits/sec    0   3.13 MBytes
[  5]   6.00-7.00   sec  1.55 GBytes  13.3 Gbits/sec    0   3.13 MBytes
[  5]   7.00-8.00   sec  1.51 GBytes  12.9 Gbits/sec    0   3.13 MBytes
[  5]   8.00-9.00   sec  1.42 GBytes  12.2 Gbits/sec    0   3.13 MBytes
[  5]   9.00-10.00  sec  1.33 GBytes  11.4 Gbits/sec    0   3.13 MBytes
- - - - - - - - - - - - - - - - - - - - - - - - -
[ ID] Interval           Transfer     Bitrate         Retr
[  5]   0.00-10.00  sec  15.1 GBytes  13.0 Gbits/sec    0             sender
[  5]   0.00-10.04  sec  15.1 GBytes  12.9 Gbits/sec                  receiver
```

### VM Config

```
ubuntu@vm1:~$ ip a show ens3
2: ens3: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc fq_codel state UP group default qlen 1000
    link/ether 52:54:00:70:c4:c6 brd ff:ff:ff:ff:ff:ff
    altname enp0s3
    inet 10.28.65.177/24 metric 100 brd 10.28.65.255 scope global ens3
       valid_lft forever preferred_lft forever
    inet6 fe80::5054:ff:fe70:c4c6/64 scope link
       valid_lft forever preferred_lft forever

ubuntu@vm1:~$ ethtool -i ens3
driver: virtio_net
version: 1.0.0
firmware-version:
expansion-rom-version:
bus-info: 0000:00:03.0
supports-statistics: yes
supports-test: no
supports-eeprom-access: no
supports-register-dump: no
supports-priv-flags: no
```

```
ubuntu@vm2:~$ ip a show ens3
2: ens3: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc fq_codel state UP group default qlen 1000
    link/ether 52:54:00:4d:4f:79 brd ff:ff:ff:ff:ff:ff
    altname enp0s3
    inet 10.28.65.38/24 metric 100 brd 10.28.65.255 scope global ens3
       valid_lft forever preferred_lft forever
    inet6 fe80::5054:ff:fe4d:4f79/64 scope link
       valid_lft forever preferred_lft forever

ubuntu@vm2:~$ ethtool -i ens3
driver: virtio_net
version: 1.0.0
firmware-version:
expansion-rom-version:
bus-info: 0000:00:03.0
supports-statistics: yes
supports-test: no
supports-eeprom-access: no
supports-register-dump: no
supports-priv-flags: no
```

### Host Interface (bridge)

```
12: mpqemubr0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default qlen 1000
    link/ether 52:54:00:2e:d7:c8 brd ff:ff:ff:ff:ff:ff
    inet 10.28.65.1/24 brd 10.28.65.255 scope global mpqemubr0
       valid_lft forever preferred_lft forever
    inet6 fe80::5054:ff:fe2e:d7c8/64 scope link
       valid_lft forever preferred_lft forever
```

## WIP: L3 Routing among 3 VMs (Ubutnu Multipass)

> how to change multipass backend to lxd https://qiita.com/ynott/items/01e2913539c664b6559d
> TODO: virsh でやった方が楽？
 
Create bridge

```
sudo ip link add br0 type bridge
sudo ip link add br1 type bridge
```
Create two VMs

```
multipass launch 22.04 --name vm13 --cpus 16 --disk 20G --memory 16G --network br0
multipass launch 22.04 --name vm2 --cpus 16 --disk 20G --memory 16G
```


## L3 routing between netns

Create two netns.
Run iperf

```
TBD
```


## Specs

### Host 171

```
$ lscpu
Architecture:                    x86_64
CPU op-mode(s):                  32-bit, 64-bit
Byte Order:                      Little Endian
Address sizes:                   46 bits physical, 48 bits virtual
CPU(s):                          64
On-line CPU(s) list:             0-63
Thread(s) per core:              2
Core(s) per socket:              16
Socket(s):                       2
NUMA node(s):                    2
Vendor ID:                       GenuineIntel
CPU family:                      6
Model:                           85
Model name:                      Intel(R) Xeon(R) Gold 5218 CPU @ 2.30GHz
Stepping:                        7
CPU MHz:                         1000.420
CPU max MHz:                     3900.0000
CPU min MHz:                     1000.0000
BogoMIPS:                        4600.00
```
