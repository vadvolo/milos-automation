# Check if host system is Linux, if sudo is installed, and if the user is not root
# Use "sudo" if all conditions are true
SUDO := $(shell if [ "$$(uname -s)" = "Linux" ] && [ -x "$$(command -v sudo)" ] && [ "$$(id -u)" -ne 0 ]; then echo "sudo"; fi)

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
	@cd ./netbox-docker && docker compose build --no-cache && cd ..

build: build_annet build_netbox

rebuild:
	rm -rf netbox-docker || true && \
	$(MAKE) build

services_stop:
	@cd ./netbox-docker && docker compose down -v && cd ..

check_cisco_image:
	@if !(ls ./vm_images/*.bin 1> /dev/null 2>&1); then echo "\n\n !!! ./vm_images doesn't have any images !!!\n" && exit 1; fi

unpack_cisco_image:
	@cd vm_images && find . -name '*.bin' -exec sh -c 'f=$$(echo $$1 | cut -d "." -f 2,3); unzip -p ".$$f.bin" > ".$$f.image"' _ {} \; && cd ..

# Build and destroy architecture

lab00: services_stop check_cisco_image unpack_cisco_image
	@cp ./topologies/lab00_basic_cisco/docker-compose.override.yml ./netbox-docker/docker-compose.override.yml && \
	cd netbox-docker/ && docker compose up -d && cd .. && \
	docker exec --user root netbox /opt/netbox/netbox/manage.py loaddata /home/ubuntu/src/netbox-dump.json && \
	docker exec -u root -t -i annet netsshsetup -a 172.20.0.101 -b ios -l annet -p annet -P telnet -v cisco --ipdomain nh.com && \
	docker exec -u root -t -i annet netsshsetup -a 172.20.0.102 -b ios -l annet -p annet -P telnet -v cisco --ipdomain nh.com && \
	docker exec -u root -t -i annet netsshsetup -a 172.20.0.103 -b ios -l annet -p annet -P telnet -v cisco --ipdomain nh.com && \
	eval "docker exec -u root -t -i annet /bin/bash"

lab00_test:
	@python3 ./topologies/lab00_basic_cisco/lab_test.py

lab01: services_stop
	@cp ./topologies/lab01_basic_frr/docker-compose.override.yml ./netbox-docker/docker-compose.override.yml && \
	cd topologies/lab01_basic_frr && \
	$(SUDO) find operational_configs -mindepth 1 -not -name '.gitkeep' -delete || true && \
	cp -r default_configs/* operational_configs/ && \
	cd ../../netbox-docker/ && docker compose up -d --build && cd .. && \
	docker exec --user root netbox /opt/netbox/netbox/manage.py loaddata /home/ubuntu/src/netbox-dump.json && \
	eval "docker exec -u root -t -i annet /bin/bash"

lab10: services_stop check_cisco_image unpack_cisco_image
	@cp ./topologies/lab10_dc_cisco/docker-compose.override.yml ./netbox-docker/docker-compose.override.yml && \
	cd netbox-docker/ && docker compose up -d && cd .. && \
	docker exec --user root netbox /opt/netbox/netbox/manage.py loaddata /home/ubuntu/src/netbox-dump.json && \
	docker exec -u root -t -i annet netsshsetup -a 172.20.0.101 -b ios -l annet -p annet -P telnet -v cisco --ipdomain nh.com && \
	docker exec -u root -t -i annet netsshsetup -a 172.20.0.102 -b ios -l annet -p annet -P telnet -v cisco --ipdomain nh.com && \
	docker exec -u root -t -i annet netsshsetup -a 172.20.0.103 -b ios -l annet -p annet -P telnet -v cisco --ipdomain nh.com && \
	docker exec -u root -t -i annet netsshsetup -a 172.20.0.104 -b ios -l annet -p annet -P telnet -v cisco --ipdomain nh.com && \
	docker exec -u root -t -i annet netsshsetup -a 172.20.0.105 -b ios -l annet -p annet -P telnet -v cisco --ipdomain nh.com && \
	eval "docker exec -u root -t -i annet /bin/bash"

lab11: services_stop
	@cp ./topologies/lab11_dc_frr/docker-compose.override.yml ./netbox-docker/docker-compose.override.yml && \
	cd topologies/lab11_dc_frr && \
	$(SUDO) find operational_configs -mindepth 1 -not -name '.gitkeep' -delete || true && \
	cp -r default_configs/* operational_configs/ && \
	cd ../../netbox-docker/ && docker compose up -d --build && cd .. && \
	docker exec --user root netbox /opt/netbox/netbox/manage.py loaddata /home/ubuntu/src/netbox-dump.json && \
	eval "docker exec -u root -t -i annet /bin/bash"

lab12: services_stop check_cisco_image unpack_cisco_image
	@cp ./topologies/lab12_dc_arista_cisco_frr/docker-compose.override.yml ./netbox-docker/docker-compose.override.yml && \
	cd topologies/lab12_dc_arista_cisco_frr && \
	$(SUDO) find operational_configs -mindepth 1 -not -name '.gitkeep' -delete || true && \
	cp -r default_configs/* operational_configs/ && \
	cd ../../netbox-docker/ && docker compose up -d --build && cd .. && \
	docker exec --user root netbox /opt/netbox/netbox/manage.py loaddata /home/ubuntu/src/netbox-dump.json && \
	docker exec -u root -t -i annet netsshsetup -a 172.20.0.103 -b ios -l annet -p annet -P telnet -v cisco --ipdomain nh.com && \
	docker exec -u root -t -i annet netsshsetup -a 172.20.0.105 -b ios -l annet -p annet -P telnet -v cisco --ipdomain nh.com && \
	eval "docker exec -u root -t -i annet /bin/bash"

# Netbox utilities

netbox_export:
	@docker exec --user root netbox /opt/netbox/netbox/manage.py dumpdata -a -o /home/ubuntu/src/netbox-dump.json --exclude extras.cachedvalue --exclude extras.objectchange --exclude extras.report --exclude extras.script --exclude django_rq.queue --exclude auth.permission --exclude contenttypes.contenttype

netbox_import:
	@docker exec --user root netbox /opt/netbox/netbox/manage.py loaddata /home/ubuntu/src/netbox-dump.json

cc:
	isort topologies

style:
	pylint -j8 --score n --reports=no topologies || true && flake8 topologies || true && mypy topologies
