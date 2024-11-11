#!/bin/bash

echo "netstat -tlpn" >> /root/.bashrc
apt update
apt install -y sudo ssh dynamips dynagen telnet telnetd net-tools bridge-utils iproute2 uml-utilities openvpn inetutils-ping wget

mkdir /home/ubuntu/next-hop-lab

ip tuntap add tap0 mode tap
ip tuntap add tap1 mode tap
ip tuntap add tap2 mode tap
ip tuntap add tap3 mode tap

brctl addbr br0
ip addr flush dev eth0
brctl addif br0 eth0
brctl addif br0 tap0
brctl addif br0 tap1
brctl addif br0 tap2
brctl addif br0 tap3

ifconfig tap0 up
ifconfig tap1 up
ifconfig tap2 up
ifconfig tap3 up

ifconfig br0 172.20.0.20/24

cd /home/ubuntu/
dynamips -H 7200 &
sleep 3s
dynagen lab.net

tail -f /dev/null
