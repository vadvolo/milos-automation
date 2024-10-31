#!/bin/bash

git clone -b release https://github.com/netbox-community/netbox-docker.git
cp ./docker-compose.yml ./netbox-docker/docker-compose.override.yml

git clone https://github.com/EvilFreelancer/docker-routeros.git
