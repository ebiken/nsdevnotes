#!/usr/bin/bash

# This script will create(remove) namespace.
#   namespace: ns1, ns2, ns3, ns4
#   interface: vnet1, vnet2, vnet3, vnet4 (created by sonic.xml)

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

    # Create netns
    run ip netns add ns1
    run ip netns add ns2
    run ip netns add ns3
    run ip netns add ns4

    ### IPv4 addressing schema
    # 192.168.X.Y/24
    #   X : first digit of vlan (vlan1001: X=1, vlan1002: X=2)
    #   Y : always 1 for vlan interface
    #       100 + nsN for netns interface

    # ns1: vlan1001, untagged
    run ip link set vnet1 netns ns1
    run ip netns exec ns1 ip link set dev lo up
    run ip netns exec ns1 ethtool --offload vnet1 rx off tx off
    run ip netns exec ns1 ip addr add 192.168.1.101/24 dev vnet1
    run ip netns exec ns1 ip link set dev vnet1 up
    run ip netns exec ns1 ip route add default via 192.168.1.1

    # ns2: vlan1001, untagged
    run ip link set vnet2 netns ns2
    run ip netns exec ns2 ip link set dev lo up
    run ip netns exec ns2 ethtool --offload vnet2 rx off tx off
    run ip netns exec ns2 ip addr add 192.168.1.102/24 dev vnet2
    run ip netns exec ns2 ip link set dev vnet2 up
    run ip netns exec ns2 ip route add default via 192.168.1.1

    # ns3: vlan1002, untagged
    run ip link set vnet3 netns ns3
    run ip netns exec ns3 ip link set dev lo up
    run ip netns exec ns3 ethtool --offload vnet3 rx off tx off
    run ip netns exec ns3 ip addr add 192.168.2.103/24 dev vnet3
    run ip netns exec ns3 ip link set dev vnet3 up
    run ip netns exec ns3 ip route add default via 192.168.2.1

    # ns4: vlan1002, tagged
    run ip link set vnet4 netns ns4
    run ip netns exec ns4 ip link set dev lo up
    run ip netns exec ns4 ethtool --offload vnet4 rx off tx off
    run ip netns exec ns4 ip link add link vnet4 name vnet4.1002 type vlan id 1002
    run ip netns exec ns4 ip addr add 192.168.2.104/24 dev vnet4.1002
    run ip netns exec ns4 ip link set dev vnet4 up
    # run ip netns exec ns4 ip link set dev vnet4.1002 up
    run ip netns exec ns4 ip route add default via 192.168.2.1

    exit 1
}

destroy_netns () {
    echo "destroy_network"
    silent ip netns del ns1
    silent ip netns del ns2
    silent ip netns del ns3
    silent ip netns del ns4
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