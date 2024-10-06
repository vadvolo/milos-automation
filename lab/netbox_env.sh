#!/bin/bash

PORT=":8080"
IP=`echo $(docker exec -it netbox-docker-netbox-1 bash -c "ifconfig eth0 | grep -oP '(?<=inet\s)\d+(\.\d+){3}'")`

export NETBOX_IP=$IP
export NETBOX_URL="http://"${IP:0:${#IP}-1}":8080"
export NETBOX_PORT="8080"

echo $NETBOX_IP
echo $NETBOX_URL
echo $NETBOX_PORT


