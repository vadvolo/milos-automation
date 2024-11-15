#!/bin/bash

echo "netstat -tlpn" >> /root/.bashrc
apt update
apt install -y sudo ssh dynamips dynagen telnet telnetd net-tools bridge-utils iproute2 uml-utilities openvpn inetutils-ping wget

mkdir /home/ubuntu/next-hop-lab

ip tuntap add tap0 mode tap
ip tuntap add tap1 mode tap
ip tuntap add tap2 mode tap
ip tuntap add tap3 mode tap
ip tuntap add tap4 mode tap
ip tuntap add tap5 mode tap

brctl addbr br0
ip addr flush dev eth0
brctl addif br0 eth0
brctl addif br0 tap0
brctl addif br0 tap1

brctl addbr br1
ip addr flush dev eth1
brctl addif br1 eth1
brctl addif br1 tap2

brctl addbr br2
ip addr flush dev eth2
brctl addif br2 eth2
brctl addif br2 tap3

brctl addbr br3
ip addr flush dev eth3
brctl addif br3 eth3
brctl addif br3 tap4

brctl addbr br4
ip addr flush dev eth4
brctl addif br4 eth4
brctl addif br4 tap5

ifconfig tap0 up
ifconfig tap1 up
ifconfig tap2 up
ifconfig tap3 up
ifconfig tap4 up
ifconfig tap5 up

ifconfig br0 172.20.0.20/24
ifconfig br1 10.1.1.20/24
ifconfig br2 10.1.2.20/24
ifconfig br3 10.1.3.20/24
ifconfig br4 10.1.4.20/24

cd /home/ubuntu/
dynamips -H 7200 &
sleep 3s
dynagen lab.net

tail -f /dev/null
