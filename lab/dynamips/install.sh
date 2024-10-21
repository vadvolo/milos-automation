#!/bin/bash

apt update
apt install -y sudo ssh dynamips dynagen telnet telnetd net-tools bridge-utils iproute2 uml-utilities openvpn inetutils-ping wget

mkdir /home/ubuntu/ios-7200
mkdir /home/ubuntu/next-hop-lab
wget -P /home/ubuntu/ios-7200 ftp://ftp.radio-portal.su/Cisco/\!cisco_ios/7200/c7200-jk9s-mz.124-13a.bin

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

ifconfig tap0 up
ifconfig tap1 up
ifconfig tap2 up

ifconfig br0 172.20.0.20/16

cd /home/ubuntu/
dynamips -H 7200 &
dynagen lab.net

tail -f /dev/null
