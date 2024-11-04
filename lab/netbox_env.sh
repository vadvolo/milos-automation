#!/bin/bash

PORT=":8080"
# IP=`echo $(docker exec -u root -it netbox-docker-netbox-1 bash -c "apt upgrade && apt install net-tools && echo "\n\n" && ifconfig eth0 | grep inet | grep -oP '(?<=inet\s)\d+(\.\d+){3}'")`
IP=`echo $(docker inspect netbox-docker-netbox-1 | grep "IPAddress" | tail -n 1 | awk '{print $2}' | tr -d '",')`

export NETBOX_TOKEN="a630dcefcb191982869e7576190e79bfd569d33c"
export NETBOX_IP=$IP
export NETBOX_URL="http://"${IP}":8080"
export NETBOX_PORT="8080"
