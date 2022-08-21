#!/usr/bin/bash

# This script will create(remove) namespace.
#   namespace: ns1, ns2, ns3, ns4
#   brXX created manually
#   interface: vnet1, vnet2, vnet3, vnet4 (created by sonic.xml and assigned to brXX.
#   interface: veth1, veth2, veth3, veth4 would be pair to veth-brxx
#   interface: veth-brxx will belong to bridge

if [[ $(id -u) -ne 0 ]] ; then echo "Please run with sudo" ; exit 1 ; fi

set -e

if [ -n "$SUDO_UID" ]; then
    uid=$SUDO_UID
else
    uid=$UID
fi

run () {
    echo "$@"
    "$@" || exit 1
}

silent () {
    "$@" 2> /dev/null || true
}

create_netns () {
    echo "create_netns for sonic-demo01"

    # Create bridge
    run ip link add br01 type bridge
    run ip link add br02 type bridge
    run ip link add br03 type bridge
    run ip link add br04 type bridge

    run ip link set br01 up
    run ip link set br02 up
    run ip link set br03 up
    run ip link set br04 up

    # Create netns
    run ip netns add ns1
    run ip netns add ns2
    run ip netns add ns3
    run ip netns add ns4

    run ip link add name veth1 type veth peer name veth-br01
    run ip link add name veth2 type veth peer name veth-br02
    run ip link add name veth3 type veth peer name veth-br03
    run ip link add name veth4 type veth peer name veth-br04

    ### IPv4 addressing schema
    # 192.168.X.Y/24
    #   X : first digit of vlan (vlan1001: X=1, vlan1002: X=2)
    #   Y : always 1 for vlan interface
    #       100 + nsN for netns interface

    # ns1: vlan1001, untagged
    run ip link set veth1 netns ns1
    run ip netns exec ns1 ip link set dev lo up
    run ip netns exec ns1 ethtool --offload veth1 rx off tx off
    run ip netns exec ns1 ip addr add 192.168.1.101/24 dev veth1
    run ip netns exec ns1 ip link set dev veth1 up
    run ip netns exec ns1 ip route add default via 192.168.1.1
    run ip link set veth-br01 master br01 up

    # ns2: vlan1001, untagged
    run ip link set veth2 netns ns2
    run ip netns exec ns2 ip link set dev lo up
    run ip netns exec ns2 ethtool --offload veth2 rx off tx off
    run ip netns exec ns2 ip addr add 192.168.1.102/24 dev veth2
    run ip netns exec ns2 ip link set dev veth2 up
    run ip netns exec ns2 ip route add default via 192.168.1.1
    run ip link set veth-br02 master br02 up

    # ns3: vlan1002, untagged
    run ip link set veth3 netns ns3
    run ip netns exec ns3 ip link set dev lo up
    run ip netns exec ns3 ethtool --offload veth3 rx off tx off
    run ip netns exec ns3 ip addr add 192.168.2.103/24 dev veth3
    run ip netns exec ns3 ip link set dev veth3 up
    run ip netns exec ns3 ip route add default via 192.168.2.1
    run ip link set veth-br03 master br03 up

    # ns4: vlan1002, tagged
    run ip link set veth4 netns ns4
    run ip netns exec ns4 ip link set dev lo up
    run ip netns exec ns4 ethtool --offload veth4 rx off tx off
    run ip netns exec ns4 ip link add link veth4 name veth4.1002 type vlan id 1002
    run ip netns exec ns4 ip addr add 192.168.2.104/24 dev veth4.1002
    run ip netns exec ns4 ip link set dev veth4 up
    # run ip netns exec ns4 ip link set dev veth4.1002 up
    run ip netns exec ns4 ip route add default via 192.168.2.1
    run ip link set veth-br04 master br04 up

    exit 1
}

destroy_netns () {
    echo "destroy_network"
    silent ip netns del ns1
    silent ip netns del ns2
    silent ip netns del ns3
    silent ip netns del ns4

    silent ip link del name veth1
    silent ip link del name veth2
    silent ip link del name veth3
    silent ip link del name veth4

    silent ip link set br01 down
    silent ip link set br02 down
    silent ip link set br03 down
    silent ip link set br04 down

    silent brctl delbr br01
    silent brctl delbr br02
    silent brctl delbr br03
    silent brctl delbr br04

    #exit 1
}

while getopts "cd" ARGS;
do
    case $ARGS in
    c )
        #NUM=$OPTARG
        destroy_netns
        create_netns
        exit 1;;
    d )
        #NUM=$OPTARG
        destroy_netns
        exit 1;;
    esac
done

cat << EOF
usage: sudo ./$(basename $BASH_SOURCE) <option>
option:
    -c : destroy_netns and then create_netns
    -d : destroy_netns
EOF
