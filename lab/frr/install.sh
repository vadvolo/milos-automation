#!/usr/bin/env bash

for DAEMON in $FRR_DAEMONS
do
    if [ -s "/usr/lib/frr/${DAEMON}" ]
    then
        touch /etc/frr/${DAEMON}.conf
        chown frr:frr /etc/frr/${DAEMON}.conf
        sed -i "s/${DAEMON}=no/${DAEMON}=yes/" /etc/frr/daemons
    else
        echo "No such daemon: ${DAEMON}"
    fi
done

cp /${HOSTNAME}.conf /etc/frr/frr.conf
chown frr:frr /etc/frr/frr.conf

service frr restart > /dev/null 2>&1

tail -f /dev/null
