# Check if sudo is installed and if the user is not root
# Use "sudo" if both are true
SUDO := $(shell if [ -x "$$(command -v sudo)" ] && [ "$$(id -u)" -ne 0 ]; then echo "sudo"; fi)

### Services (annet and netbox) control

build_annet:
	@cd ./annet && \
	cp -r ../scripts/netsshsetup ./netsshsetup && \
	docker build --no-cache -t annet . && \
	rm -rf ./netsshsetup

build_netbox:
	@git clone -b release https://github.com/netbox-community/netbox-docker.git
	@cp ./netbox-docker/docker-compose.yml ./netbox-docker/docker-compose.yml.bak
	@sed -i.bak 's/VERSION-v4.1-3.0.2/VERSION-v3.7/g' ./netbox-docker/docker-compose.yml
	@sed -i.bak 's/0s/00s/g' ./netbox-docker/docker-compose.yml
	@cd ./netbox-docker && docker-compose build --no-cache && cd ..

build: build_annet build_netbox

rebuild:
	rm -rf netbox-docker || true && \
	$(MAKE) build

services_stop:
	@cd ./netbox-docker && docker-compose down -v && cd ..

check_cisco_image:
	if [ ! -f "./vm_images/c7200-jk9s-mz.124-13a.bin" ]; then echo "\n\n !!! ./vm_images doesn't have c7200-jk9s-mz.124-13a.bin !!!\n" && exit 1; fi 

try_unpack_cisco_image:
	if [ ! -f "./vm_images/c7200-jk9s-mz.124-13a.image" ]; then unzip -p ./vm_images/c7200-jk9s-mz.124-13a.bin > ./vm_images/c7200-jk9s-mz.124-13a.image; fi

# Build and destroy architecture


lab00: services_stop check_cisco_image try_unpack_cisco_image
	@echo "Lab00 Cisco Basic Scenario is launching..." && \
	cp ./topologies/lab00_basic_cisco/docker-compose.override.yml ./netbox-docker/docker-compose.override.yml && \
	cd netbox-docker/ && docker-compose up -d && cd .. && \
	docker exec --user root netbox /opt/netbox/netbox/manage.py loaddata /home/ubuntu/src/netbox-dump.json && \
	echo "!!! Don't forget add Cisco IOS image into ../vm_images folder !!!" && \
	eval "docker exec -u root -t -i annet /bin/bash"

lab01: services_stop try_unpack_cisco_image
	@cp ./topologies/lab01_basic_frr/docker-compose.override.yml ./netbox-docker/docker-compose.override.yml && \
	cd topologies/lab01_basic_frr && \
	$(SUDO) find operational_configs -mindepth 1 -not -name '.gitkeep' -delete || true && \
	cp -r default_configs/* operational_configs/ && \
	cd ../../netbox-docker/ && docker-compose up -d --build && cd .. && \
	docker exec --user root netbox /opt/netbox/netbox/manage.py loaddata /home/ubuntu/src/netbox-dump.json && \
	eval "docker exec -u root -t -i annet /bin/bash"

lab10: services_stop
	cp ./topologies/lab10_dc_cisco/docker-compose.override.yml ./netbox-docker/docker-compose.override.yml && \
	cd netbox-docker/ && docker-compose up -d && cd .. && \
	docker exec --user root netbox /opt/netbox/netbox/manage.py loaddata /home/ubuntu/src/netbox-dump.json && \
	echo "!!! Don't forget add iol image into ./lab00_basic_cisco/src/ios-7200 folder !!!" && \
	eval "docker exec -u root -t -i annet /bin/bash"

lab12: services_stop check_cisco_image try_unpack_cisco_image
	cp ./topologies/lab12_dc_arista_cisco_frr/docker-compose.override.yml ./netbox-docker/docker-compose.override.yml && \
	cd topologies/lab12_dc_arista_cisco_frr && \
	$(SUDO) find operational_configs -mindepth 1 -not -name '.gitkeep' -delete || true && \
	cp -r default_configs/* operational_configs/ && \
	cd ../../netbox-docker/ && docker-compose up -d --build && cd .. && \
	docker exec --user root netbox /opt/netbox/netbox/manage.py loaddata /home/ubuntu/src/netbox-dump.json && \
	eval "docker exec -u root -t -i annet /bin/bash"

# Netbox utilities

netbox_export:
	@docker exec --user root netbox /opt/netbox/netbox/manage.py dumpdata -a -o /home/ubuntu/src/netbox-dump.json --exclude extras.cachedvalue --exclude extras.objectchange --exclude extras.report --exclude extras.script --exclude django_rq.queue --exclude auth.permission --exclude contenttypes.contenttype

netbox_import:
	@docker exec --user root netbox /opt/netbox/netbox/manage.py loaddata /home/ubuntu/src/netbox-dump.json

cc:
	isort topologies

style:
	pylint -j8 --score n --reports=no topologies && \
    flake8 topologies
