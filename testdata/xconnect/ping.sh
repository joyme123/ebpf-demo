#!/bin/bash

sudo ip netns exec ns1 ping 169.254.1.20 &
sudo ip netns exec ns1 ping 169.254.1.30 &
