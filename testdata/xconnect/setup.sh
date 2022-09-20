#!/bin/bash

sudo ip link add dev xconnect-1 type veth peer name xc-1
sudo ip link add dev xconnect-2 type veth peer name xc-2
sudo ip link add dev xconnect-3 type veth peer name xc-3
sudo ip netns add ns1
sudo ip netns add ns2
sudo ip netns add ns3

sudo ip link set xc-1 netns ns1
sudo ip link set xc-2 netns ns2
sudo ip link set xc-3 netns ns3
sudo ip netns exec ns1 ip addr add 169.254.1.10/24 dev xc-1
sudo ip netns exec ns2 ip addr add 169.254.1.20/24 dev xc-2
sudo ip netns exec ns3 ip addr add 169.254.1.30/24 dev xc-3

sudo ip link set dev xconnect-1 up
sudo ip link set dev xconnect-2 up
sudo ip link set dev xconnect-3 up
sudo ip netns exec ns1 ip link set dev xc-1 up
sudo ip netns exec ns2 ip link set dev xc-2 up
sudo ip netns exec ns3 ip link set dev xc-3 up
